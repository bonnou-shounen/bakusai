package bakusai

import (
	"fmt"
	"regexp"
	"strconv"
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

func (r *Res) ParseURI(uri string) error {
	keys := []string{"acode", "ctgid", "bid", "tid", "rrid"}

	param := map[string]int{}
	if err := parseURI(uri, keys, param); err != nil {
		return err
	}

	r.AreaCode = param["acode"]
	r.CategoryID = param["ctgid"]
	r.BoardID = param["bid"]
	r.ThreadID = param["tid"]
	r.RRID = param["rrid"]

	return nil
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
	PrevTID       int
	NextTID       int
	PageNum       int
}

func (t *Thread) URI() string {
	if t.PageNum == 0 {
		return fmt.Sprintf(
			`%s/thr_res/acode=%d/ctgid=%d/bid=%d/tid=%d/`,
			RootURI, t.AreaCode, t.CategoryID, t.BoardID, t.ThreadID,
		)
	}

	if t.PageNum < 0 {
		// 最新からの相対
		return fmt.Sprintf(
			`%s/thr_res/acode=%d/ctgid=%d/bid=%d/tid=%d/p=%d/`,
			RootURI, t.AreaCode, t.CategoryID, t.BoardID, t.ThreadID, -t.PageNum,
		)
	}

	// 最初から
	return fmt.Sprintf(
		`%s/thr_res/acode=%d/ctgid=%d/bid=%d/tid=%d/p=%d/rw=1/`,
		RootURI, t.AreaCode, t.CategoryID, t.BoardID, t.ThreadID, t.PageNum,
	)
}

func (t *Thread) ParseURI(uri string) error {
	keys := []string{"acode", "ctgid", "bid", "tid"}

	param := map[string]int{}
	if err := parseURI(uri, keys, param); err != nil {
		return err
	}

	t.AreaCode = param["acode"]
	t.CategoryID = param["ctgid"]
	t.BoardID = param["bid"]
	t.ThreadID = param["tid"]

	return nil
}

var uriParamsRe = regexp.MustCompile(`([^=/]+)=(\d+)`)

func parseURI(uri string, keys []string, param map[string]int) error {
	matches := uriParamsRe.FindAllStringSubmatch(uri, -1)
	for _, match := range matches {
		param[match[1]], _ = strconv.Atoi(match[2])
	}

	for _, key := range keys {
		if _, ok := param[key]; !ok {
			return fmt.Errorf("missing: %s", key)
		}
	}

	return nil
}
