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

func GetThread(ctx context.Context, argURI string) (*bakusai.Thread, error) {
	uriMaker, err := parser.ParseThreadURI(argURI)
	if err != nil {
		return nil, fmt.Errorf(`on thread.ParseURI("%s"): %w`, argURI, err)
	}

	uriMaker.PageNum = 0
	uri := uriMaker.URI()

	lastThread, err := getThreadOnPage(ctx, uri)
	if err != nil {
		return nil, fmt.Errorf(`on getThread([page=last]): %w`, err)
	}

	lastRRID := lastThread.ResList[len(lastThread.ResList)-1].RRID
	lastPage := (lastRRID + 49) / 50

	var resList []*bakusai.Res

	if lastPage == 1 {
		resList = lastThread.ResList
	} else {
		pageThreads, err := getPageThreads(ctx, uriMaker, lastPage)
		if err != nil {
			return nil, fmt.Errorf(`on getPageThreads(): %w`, err)
		}

		for _, pageThread := range pageThreads {
			resList = append(resList, pageThread.ResList...)
		}
	}

	return &bakusai.Thread{
		AreaCode:      uriMaker.AreaCode,
		CategoryID:    uriMaker.CategoryID,
		BoardID:       uriMaker.BoardID,
		ThreadID:      uriMaker.ThreadID,
		DatePublished: lastThread.DatePublished,
		Author:        lastThread.Author,
		Title:         lastThread.Title,
		ResList:       resList,
		PrevTID:       lastThread.PrevTID,
		NextTID:       lastThread.NextTID,
	}, nil
}

func getPageThreads(ctx context.Context, uriMaker *bakusai.Thread, lastPage int) ([]*bakusai.Thread, error) {
	pageThreads := make([]*bakusai.Thread, lastPage)

	eg, egCtx := errgroup.WithContext(ctx, 5)

	for page := 1; page < lastPage; page++ {
		page := page

		uriMaker.PageNum = page
		uri := uriMaker.URI()

		eg.Go(func() error {
			pageThread, err := getThreadOnPage(egCtx, uri)
			if err != nil {
				return fmt.Errorf(`on getThread([page=%d]): %w`, page, err)
			}

			pageThreads[page] = pageThread

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("on goroutine: %w", err)
	}

	return pageThreads, nil
}

func getThreadOnPage(ctx context.Context, uri string) (*bakusai.Thread, error) {
	resp, err := resty.New().
		SetRetryCount(10).SetRetryWaitTime(1 * time.Second).
		AddRetryCondition(func(r *resty.Response, err error) bool {
			if !strings.Contains(string(r.Body()), `<td class="reslist_td">`) {
				fmt.Fprintf(os.Stderr, "retry: %s\n", r.Request.URL)

				return true
			}

			return false
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
