package scraper

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/bonnou-shounen/bakusai"
	"github.com/bonnou-shounen/bakusai/parser"
)

func ScrapeThread(ctx context.Context, uri string) (*bakusai.Thread, error) {
	log.Println("scraping last page...")

	lastPart, lastResID, err := scrapeLastThreadPart(ctx, uri)
	if err != nil {
		return nil, fmt.Errorf(`on scrapeLastPart(): %w`, err)
	}

	lastPage := (lastResID + bakusai.MaxPageRes - 1) / bakusai.MaxPageRes
	log.Printf("was page %d", lastPage)

	if lastPage == 1 {
		return lastPart, nil
	}

	thread := *lastPart
	thread.ResList = make([]*bakusai.Res, 0, lastResID)

	parts, err := scrapeThreadParts(ctx, &thread, lastPage-1)
	if err != nil {
		return nil, fmt.Errorf(`on scrapeThreadPages(): %w`, err)
	}

	for page := 1; page < lastPage; page++ {
		thread.ResList = append(thread.ResList, parts[page].ResList...)
	}

	thread.ResList = append(thread.ResList, lastPart.ResList[lastPage*bakusai.MaxPageRes-lastResID:]...)

	return &thread, nil
}

func scrapeLastThreadPart(ctx context.Context, argURI string) (*bakusai.Thread, int, error) {
	uri, err := parser.CanonizeThreadURI(argURI)
	if err != nil {
		return nil, 0, fmt.Errorf(`on parser.CanonizeThreadURI("%s"): %w`, argURI, err)
	}

	part, err := ScrapeThreadPart(ctx, uri)
	if err != nil {
		return nil, 0, fmt.Errorf(`on ScrapePart(last): %w`, err)
	}

	part.URI = strings.Replace(part.URI, "/p=1/", "/", 1)
	lastResID := part.ResList[len(part.ResList)-1].ResID

	return part, lastResID, nil
}

func scrapeThreadParts(ctx context.Context, thread *bakusai.Thread, lastPage int) ([]*bakusai.Thread, error) {
	parts := make([]*bakusai.Thread, lastPage+1)

	for page := 1; page <= lastPage; page++ {
		log.Printf("scraping page %d...", page)

		part, err := ScrapeThreadPart(ctx, thread.PageURI(page))
		if err != nil {
			return nil, fmt.Errorf(`on Scrape([page=%d]): %w`, page, err)
		}

		parts[page] = part
	}

	return parts, nil
}
