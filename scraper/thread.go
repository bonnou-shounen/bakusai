package scraper

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/bonnou-shounen/bakusai"
	"github.com/bonnou-shounen/bakusai/parser"
	"github.com/go-resty/resty/v2"
	"github.com/jesse0michael/errgroup"
)

func ScrapeThread(ctx context.Context, argURI string) (*bakusai.Thread, error) {
	thread, err := parser.ParseThreadURI(argURI)
	if err != nil {
		return nil, fmt.Errorf(`on thread.ParseURI("%s"): %w`, argURI, err)
	}

	lastThread, err := scrapeThreadOnPage(ctx, thread.URI())
	if err != nil {
		return nil, fmt.Errorf(`on getThread(last): %w`, err)
	}

	thread.DatePublished = lastThread.DatePublished
	thread.Author = lastThread.Author
	thread.Title = lastThread.Title
	thread.PrevURI = lastThread.PrevURI
	thread.NextURI = lastThread.NextURI

	lastRRID := lastThread.ResList[len(lastThread.ResList)-1].ResID
	lastPage := (lastRRID + 49) / 50

	if lastPage == 1 {
		thread.ResList = lastThread.ResList

		return thread, nil
	}

	pageThreads, err := scrapeThreads(ctx, thread, lastPage)
	if err != nil {
		return nil, fmt.Errorf(`on getPageThreads(): %w`, err)
	}

	for page := 1; page < lastPage; page++ {
		thread.ResList = append(thread.ResList, pageThreads[page].ResList...)
	}

	thread.ResList = append(thread.ResList, lastThread.ResList[lastPage*50-lastRRID:]...)

	return thread, nil
}

func scrapeThreads(ctx context.Context, thread *bakusai.Thread, lastPage int) ([]*bakusai.Thread, error) {
	pageThreads := make([]*bakusai.Thread, lastPage)

	/*
		1リクエスト/秒のアクセス頻度制限があるようで
		並列化や秒未満のリトライはエラーが頻発するため
		実質シングルスレッドにしておく
	*/
	eg, egCtx := errgroup.WithContext(ctx, 1)

	for page := 1; page < lastPage; page++ {
		page := page

		time.Sleep(1 * time.Second)

		eg.Go(func() error {
			threadOnPage, err := scrapeThreadOnPage(egCtx, thread.PageURI(page))
			if err != nil {
				return fmt.Errorf(`on getThread([page=%d]): %w`, page, err)
			}

			pageThreads[page] = threadOnPage

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("on goroutine: %w", err)
	}

	return pageThreads, nil
}

func scrapeThreadOnPage(ctx context.Context, uri string) (*bakusai.Thread, error) {
	resp, err := resty.New().
		SetRetryCount(10).SetRetryWaitTime(1 * time.Second).SetRetryMaxWaitTime(3 * time.Second).
		AddRetryCondition(func(r *resty.Response, err error) bool {
			// アクセス頻度制限にかかると記事のないエラーページが返る
			return !strings.Contains(string(r.Body()), `<td class="reslist_td">`)
		}).
		AddRetryHook(func(r *resty.Response, err error) {
			fmt.Fprintf(os.Stderr, "retry: %s\n", r.Request.URL)
		}).
		R().SetContext(ctx).
		Get(uri)
	if err != nil {
		return nil, fmt.Errorf(`on resty.Get("%s"): %w`, uri, err)
	}

	tdResList, err := findTDResList(bytes.NewReader(resp.Body()))
	if err != nil {
		return nil, fmt.Errorf(`on parser.FindTDResList(): %w`, err)
	}

	thread, err := parser.ParseThread(tdResList)
	if err != nil {
		return nil, fmt.Errorf(`on parser.ParseThread(): %w`, err)
	}

	// fmt.Fprintf(os.Stderr, "done: %s\n", uri)

	return thread, nil
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
