package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/bonnou-shounen/bakusai/scraper"
	"github.com/bonnou-shounen/bakusai/util"
)

func main() {
	var uri string

	if len(os.Args) > 1 {
		uri = os.Args[1]
	} else {
		uri = readURL()
	}

	if err := run(uri); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}
}

func run(uri string) error {
	ctx := context.Background()

	thread, err := scraper.ScrapeThread(ctx, uri)
	if err != nil {
		return fmt.Errorf(`on scraper.GetThread("%s"): %w`, uri, err)
	}

	fmt.Fprintf(os.Stdout,
		"# T: %s\n# C: %s\n# A: %s\n# U: %s\n# P: %s\n# N: %s\n",
		thread.Title,
		thread.DatePublished.Format("2006/01/02 15:04"),
		thread.Author,
		thread.URI(),
		util.ThreadPrevURI(thread),
		util.ThreadNextURI(thread),
	)

	for _, res := range thread.ResList {
		fmt.Fprintf(os.Stdout, "==== %d %s\n", res.RRID, res.CommentTime.Format("2006/01/02 15:04"))
		fmt.Fprintf(os.Stdout, "%s\n\n", res.CommentText)
	}

	return nil
}

func readURL() string {
	fmt.Fprint(os.Stderr, "paste URL: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	return scanner.Text()
}
