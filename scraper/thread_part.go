package scraper

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"

	"github.com/bonnou-shounen/bakusai"
	"github.com/bonnou-shounen/bakusai/parser"
)

func ScrapePartThread(ctx context.Context, uri string) (*bakusai.Thread, error) {
	resp, err := getPage(ctx, uri, func(r *resty.Response, _ error) bool {
		// アクセス頻度制限にかかると記事のないエラーページが返る
		return !strings.Contains(string(r.Body()), `<td class="reslist_td">`)
	})
	if err != nil {
		return nil, fmt.Errorf(`on getPage(): %w`, err)
	}

	tdResList, err := findTDResList(bytes.NewReader(resp.Body()))
	if err != nil {
		return nil, fmt.Errorf(`on findTDResList(): %w`, err)
	}

	thread := bakusai.Thread{URI: uri}
	if err := parser.ParseThread(tdResList, &thread); err != nil {
		return nil, fmt.Errorf(`on parser.ParseThread(): %w`, err)
	}

	return &thread, nil
}

func getPage(ctx context.Context, uri string,
	retryCondition func(*resty.Response, error) bool,
) (*resty.Response, error) {
	restyClient := resty.New().
		SetRetryCount(5).
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(5 * time.Second).
		AddRetryCondition(retryCondition).
		AddRetryHook(func(r *resty.Response, _ error) {
			log.Printf("retry: %s", r.Request.URL)
		})

	time.Sleep(1 * time.Second)

	resp, err := restyClient.R().SetContext(ctx).Get(uri)
	if err != nil {
		return nil, fmt.Errorf(`on resty.Get("%s"): %w`, uri, err)
	}

	if retryCondition(resp, nil) {
		return nil, fmt.Errorf(`on resty.Get("%s"): still need retry`, uri)
	}

	return resp, nil
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
