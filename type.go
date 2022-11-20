package bakusai

import (
	"fmt"
	"time"
)

const RootURI = "https://bakusai.com"

type Res struct {
	AreaCode    int
	CategoryID  int
	BoardID     int
	ThreadID    int
	RRID        int
	CommentTime time.Time
	Name        string
	CommentText string
}

func (r *Res) URI() string {
	return fmt.Sprintf(
		`%s/thr_res_show/acode=%d/ctgid=%d/bid=%d/tid=%d/rrid=%d/`,
		RootURI, r.AreaCode, r.CategoryID, r.BoardID, r.ThreadID, r.RRID,
	)
}

type Thread struct {
	AreaCode      int
	CategoryID    int
	BoardID       int
	ThreadID      int
	DatePublished time.Time
	Author        string
	Title         string
	ResList       []*Res
	PrevURI       string
	NextURI       string
}

func (t *Thread) URI() string {
	return fmt.Sprintf(
		`%s/thr_res/acode=%d/ctgid=%d/bid=%d/tid=%d/`,
		RootURI, t.AreaCode, t.CategoryID, t.BoardID, t.ThreadID,
	)
}

func (t *Thread) LastPageURI() string {
	return t.URI()
}

func (t *Thread) PageURI(page int) string {
	if page == 0 {
		return t.URI()
	}

	if page < 0 {
		// 最新からの相対
		return fmt.Sprintf(
			`%s/thr_res/acode=%d/ctgid=%d/bid=%d/tid=%d/p=%d/`,
			RootURI, t.AreaCode, t.CategoryID, t.BoardID, t.ThreadID, -page,
		)
	}

	// 最初から
	return fmt.Sprintf(
		`%s/thr_res/acode=%d/ctgid=%d/bid=%d/tid=%d/p=%d/rw=1/`,
		RootURI, t.AreaCode, t.CategoryID, t.BoardID, t.ThreadID, page,
	)
}
