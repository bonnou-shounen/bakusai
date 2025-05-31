package scraper

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"

	"github.com/bonnou-shounen/bakusai"
	"github.com/bonnou-shounen/bakusai/parser"
)

func ScrapeThreadPart(ctx context.Context, uri string) (*bakusai.Thread, error) {
	respBody, finalURI, err := getPage(ctx, uri)
	if err != nil {
		return nil, err
	}
	defer respBody.Close()

	tdResList, err := findTDResList(respBody)
	if err != nil {
		return nil, fmt.Errorf(`on findTDResList(): %w`, err)
	}

	part := &bakusai.Thread{URI: finalURI}
	if err := parser.ParseThread(tdResList, part); err != nil {
		return nil, fmt.Errorf(`on parser.ParseThread(): %w`, err)
	}

	return part, nil
}

func getPage(ctx context.Context, uri string) (io.ReadCloser, string, error) {
	finalURI := uri

	retryCondition := func(r *resty.Response, _ error) bool {
		// アクセス頻度制限にかかると記事のないエラーページが返る
		return !strings.Contains(string(r.Body()), `<td class="reslist_td">`)
	}

	restyClient := resty.New().
		SetRetryCount(10).
		SetRetryWaitTime(5 * time.Second).
		SetRetryMaxWaitTime(30 * time.Second).
		AddRetryCondition(retryCondition).
		AddRetryHook(func(r *resty.Response, _ error) {
			log.Printf("retry: %s", r.Request.URL)
		}).
		SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, _ []*http.Request) error {
			finalURI = req.URL.String()
			return nil
		}))

	time.Sleep(1 * time.Second)

	resp, err := restyClient.R().SetContext(ctx).Get(uri)
	if err != nil {
		return nil, "", fmt.Errorf(`on resty.Get("%s"): %w`, uri, err)
	}

	if retryCondition(resp, nil) {
		return nil, "", fmt.Errorf(`on resty.Get("%s"): still need retry`, uri)
	}

	// http.Response.Body互換にする
	return io.NopCloser(bytes.NewReader(resp.Body())), finalURI, nil
}

func findTDResList(reader io.Reader) (*goquery.Selection, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf(`on goquery.NewDocumentFromReader(): %w`, err)
	}

	tdResList := doc.Find(`table#inner_container td.reslist_td`)
	if tdResList.Length() == 0 {
		return nil, fmt.Errorf(`missing: td.reslist_td`)
	}

	return tdResList, nil
}
