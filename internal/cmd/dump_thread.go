package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/bonnou-shounen/bakusai"
	"github.com/bonnou-shounen/bakusai/scraper"
)

type DumpThread struct {
	URI    string `arg:"" optional:"" help:"the thread URI"`
	OptURI string `name:"uri" optional:"" hidden:""`
}

func (c *DumpThread) Run() error {
	ctx := context.Background()

	uri := c.getURI()
	if uri == "" {
		return fmt.Errorf(`on getURI(): missing URI`)
	}

	thread, err := scraper.ScrapeThread(ctx, uri)
	if err != nil {
		return fmt.Errorf(`on scraper.ScrapeThread("%s"): %w`, uri, err)
	}

	c.dump(thread)

	return nil
}

func (c *DumpThread) getURI() string {
	if c.URI != "" {
		return c.URI
	}

	if c.OptURI != "" {
		return c.OptURI
	}

	return c.readURI()
}

func (c *DumpThread) readURI() string {
	fmt.Fprint(os.Stderr, "paste thread URI: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	return scanner.Text()
}

func (c *DumpThread) dump(thread *bakusai.Thread) {
	fmt.Fprintf(os.Stdout,
		"# T: %s\n# C: %s\n# A: %s\n# U: %s\n# P: %s\n# N: %s\n",
		thread.Title,
		thread.DatePublished.Format("2006/01/02 15:04"),
		thread.Author,
		thread.URI,
		thread.PrevURI,
		thread.NextURI,
	)

	for _, res := range thread.ResList {
		fmt.Fprintf(os.Stdout,
			"==== %d %s %s\n%s\n\n",
			res.ResID, res.CommentTime.Format("2006/01/02 15:04"), res.Name,
			res.CommentText,
		)
	}
}
