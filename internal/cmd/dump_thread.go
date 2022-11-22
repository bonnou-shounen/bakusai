package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/bonnou-shounen/bakusai"
	"github.com/bonnou-shounen/bakusai/parser"
	"github.com/bonnou-shounen/bakusai/scraper"
)

type DumpThread struct {
	URI    string `arg:"" optional:"" help:"the thread URI"`
	OptURI string `name:"uri" optional:"" hidden:""`
}

func (d *DumpThread) Run(o *Option) error {
	ctx := context.Background()

	uri := d.getURI()
	if uri == "" {
		return fmt.Errorf(`on getURI(): missing URI`)
	}

	thread, err := scraper.ScrapeThread(ctx, uri)
	if err != nil {
		return fmt.Errorf(`on scraper.ScrapeThread("%s"): %w`, uri, err)
	}

	d.dump(thread)

	return nil
}

func (d *DumpThread) getURI() string {
	if d.URI != "" {
		return d.URI
	}

	if d.OptURI != "" {
		return d.OptURI
	}

	return d.readURI()
}

func (d *DumpThread) readURI() string {
	fmt.Fprint(os.Stderr, "paste thread URI: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	return scanner.Text()
}

func (d *DumpThread) dump(thread *bakusai.Thread) {
	fmt.Fprintf(os.Stdout,
		"# T: %s\n# C: %s\n# A: %s\n# U: %s\n# P: %s\n# N: %s\n",
		thread.Title,
		thread.DatePublished.Format("2006/01/02 15:04"),
		thread.Author,
		thread.URI(),
		canonicalURI(thread.PrevURI),
		canonicalURI(thread.NextURI),
	)

	for _, res := range thread.ResList {
		fmt.Fprintf(os.Stdout, "==== %d %s\n", res.RRID, res.CommentTime.Format("2006/01/02 15:04"))
		fmt.Fprintf(os.Stdout, "%s\n\n", res.CommentText)
	}
}

func canonicalURI(uri string) string {
	if uri == "" {
		return ""
	}

	thread, err := parser.ParseThreadURI(uri)
	if err != nil {
		return ""
	}

	return thread.URI()
}
