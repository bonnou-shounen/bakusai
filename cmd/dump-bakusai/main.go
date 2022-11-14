package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/bonnou-shounen/bakusai/scraper"
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
	thread, err := scraper.GetThread(uri)
	if err != nil {
		return fmt.Errorf(`on scraper.GetThread("%s"): %w`, uri, err)
	}

	fmt.Fprintf(os.Stdout,
		"# T: %s\n# C: %s\n# U: %s\n# P: %s\n# N: %s\n",
		thread.Title,
		thread.CreatedAt.Format("2006/01/02 15:04"),
		thread.URI,
		thread.PrevURI,
		thread.NextURI,
	)

	for _, a := range thread.Articles {
		fmt.Fprintf(os.Stdout, "==== %d %s\n", a.ID, a.PostAt.Format("2006/01/02 15:04"))
		fmt.Fprintf(os.Stdout, "%s\n\n", a.BodyText)
	}

	return nil
}

func readURL() string {
	fmt.Fprint(os.Stderr, "paste URL: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	return scanner.Text()
}
