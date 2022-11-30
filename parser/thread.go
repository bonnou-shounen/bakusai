package parser

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/bonnou-shounen/bakusai"
)

func ParseThread(tdResList *goquery.Selection, thread *bakusai.Thread) error {
	if err := parseHeader(tdResList, thread); err != nil {
		return fmt.Errorf(`on parseHeader(): %w`, err)
	}

	resList, err := ParseResList(tdResList, &bakusai.Res{URI: thread.URI})
	if err != nil {
		return fmt.Errorf(`on ParseResList(): %w`, err)
	}

	thread.ResList = resList

	return nil
}

func parseHeader(tdResList *goquery.Selection, thread *bakusai.Thread) error {
	// ヘッダ
	dtHeaeder := tdResList.Find(`dt.titleheader`)
	if dtHeaeder.Length() == 0 {
		return fmt.Errorf(`missing: dt.titleheader`)
	}

	// スレッドID
	strNO := dtHeaeder.Find(`span.thr_id a`).Text()
	if strNO == "" {
		return fmt.Errorf(`empty: span.thr_id a`)
	}

	strTID := strings.TrimPrefix(strNO, "NO.")

	threadID, err := strconv.Atoi(strTID)
	if err != nil {
		return fmt.Errorf(`on strconv.Atoi("%s"): %w`, strTID, err)
	}

	thread.ThreadID = threadID

	// タイトル
	thread.Title = dtHeaeder.Find(`div#title_thr`).Text()

	// スレ立て日時
	strDatePublished := dtHeaeder.Find(`span[itemprop="datePublished"]`).Text()
	if strDatePublished == "" {
		return fmt.Errorf(`empty: span[itemprop="datePublished"]`)
	}

	strDatePublished += ":00 +09:00"

	datePublished, err := time.Parse("2006/01/02 15:04:05 -07:00", strDatePublished)
	if err != nil {
		return fmt.Errorf(`on time.Parse("%s"): %w`, strDatePublished, err)
	}

	thread.DatePublished = datePublished

	// 前スレ・次スレ
	divPager := tdResList.Find(`div#thr_pager`)
	thread.PrevURI = divPager.Find(`div.sre_mae a`).AttrOr("href", "")
	thread.NextURI = divPager.Find(`div.sre_tugi a`).AttrOr("href", "")

	return nil
}
