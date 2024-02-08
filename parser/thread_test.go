package parser_test

import (
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"

	"github.com/bonnou-shounen/bakusai"
	"github.com/bonnou-shounen/bakusai/parser"
)

func setUp(t *testing.T, name string) (*goquery.Selection, func()) {
	t.Helper()

	f, err := os.Open("testdata/" + name + ".html")
	if err != nil {
		t.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		t.Fatal(err)
	}

	tdResList := doc.Find("td.reslist_td")
	if tdResList.Length() == 0 {
		t.Fatal("missing: td.reslist_td")
	}

	return tdResList, func() { f.Close() }
}

func TestParseThread(t *testing.T) {
	t.Parallel()

	tdResList, tearDown := setUp(t, "thread-1")
	defer tearDown()

	thread := bakusai.Thread{}

	err := parser.ParseThread(tdResList, &thread)
	if err != nil {
		t.Errorf(`parser.ParseThread(): error %v`, err)
	}

	if thread.ThreadID == 0 {
		t.Errorf(`thread.ThreadID is zero-value`)
	}

	if thread.Title == "" {
		t.Errorf(`thread.Title is empty`)
	}

	if thread.DatePublished.IsZero() {
		t.Errorf(`thread.DatePublished is zero-value`)
	}

	if len(thread.ResList) == 0 {
		t.Errorf(`thread.ResList is empty`)
	}

	res := thread.ResList[0]
	if res.ResID == 0 {
		t.Errorf(`res.ResID is zero-value`)
	}

	if res.CommentTime.IsZero() {
		t.Errorf(`res.CommentTime is zero-value`)
	}

	if res.CommentText == "" {
		t.Errorf(`res.CommentText is empty`)
	}
}
