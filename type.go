package bakusai

import (
	"fmt"
	"time"
)

const RootURI = "https://bakusai.com"

type Res struct {
	URI         string
	ResID       int
	Name        string
	CommentTime time.Time
	CommentText string
}

type Thread struct {
	URI           string
	ThreadID      int
	DatePublished time.Time
	Author        string
	Title         string
	ResList       []*Res
	PrevURI       string
	NextURI       string
}

func (t *Thread) PageURI(page int) string {
	if page < 0 {
		return fmt.Sprintf("%sp=%d/", t.URI, -page)
	}

	if page == 0 {
		return t.URI
	}

	if page == 1 {
		return fmt.Sprintf("%srw=1/", t.URI)
	}

	return fmt.Sprintf("%sp=%d/rw=1/", t.URI, page)
}
