package scraper

import (
	"context"
	"fmt"
	"log"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/bonnou-shounen/bakusai"
	"github.com/bonnou-shounen/bakusai/parser"
)

func ScrapeThread(ctx context.Context, uri string) (*bakusai.Thread, error) {
	log.Println("scraping last page...")

	lastThread, lastResID, err := scrapeLastThread(ctx, uri)
	if err != nil {
		return nil, fmt.Errorf(`on scrapeLastThread(): %w`, err)
	}

	lastPage := (lastResID + 49) / 50

	log.Printf("was page %d", lastPage)

	if lastPage == 1 {
		return lastThread, nil
	}

	thread := bakusai.Thread{
		URI:     lastThread.URI,
		ResList: make([]*bakusai.Res, 0, lastResID),
	}

	partThreads, err := scrapeThreads(ctx, &thread, lastPage-1)
	if err != nil {
		return nil, fmt.Errorf(`on scrapeThreads(): %w`, err)
	}

	for page := 1; page < lastPage; page++ {
		thread.ResList = append(thread.ResList, partThreads[page].ResList...)
	}

	thread.ResList = append(thread.ResList, lastThread.ResList[lastPage*50-lastResID:]...)

	return &thread, nil
}

func scrapeLastThread(ctx context.Context, argURI string) (*bakusai.Thread, int, error) {
	uri, err := parser.CanonizeThreadURI(argURI)
	if err != nil {
		return nil, 0, fmt.Errorf(`on parser.CanonizeThreadURI("%s"): %w`, argURI, err)
	}

	thread, err := ScrapePartThread(ctx, uri)
	if err != nil {
		return nil, 0, fmt.Errorf(`on ScrapePartThread(last): %w`, err)
	}

	thread.URI = uri

	lastResID := thread.ResList[len(thread.ResList)-1].ResID

	return thread, lastResID, nil
}

func scrapeThreads(ctx context.Context, thread *bakusai.Thread, lastPage int) ([]*bakusai.Thread, error) {
	threads := make([]*bakusai.Thread, lastPage+1)

	// 1リクエスト/秒のアクセス頻度制限があるようなので控えめに
	concurrency := 2

	eg, egCtx := errgroup.WithContext(ctx)
	eg.SetLimit(concurrency)

	for page := 1; page <= lastPage; page++ {
		page := page

		time.Sleep(time.Duration(1.0/concurrency) * time.Second)
		eg.Go(func() error {
			log.Printf("scraping page %d...", page)

			partThread, err := ScrapePartThread(egCtx, thread.PageURI(page))
			if err != nil {
				return fmt.Errorf(`on Scrape([page=%d]): %w`, page, err)
			}

			threads[page] = partThread

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf(`from eg.Go(): %w`, err)
	}

	return threads, nil
}
