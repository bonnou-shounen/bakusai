package parser

import (
	"fmt"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/bonnou-shounen/bakusai"
)

func ParseThread(tdResList *goquery.Selection) (*bakusai.Thread, error) {
	thread, err := parseHeader(tdResList)
	if err != nil {
		return nil, fmt.Errorf(`on parseHeader(): %w`, err)
	}

	resList, err := ParseResList(tdResList)
	if err != nil {
		return nil, fmt.Errorf(`on ParseResList(): %w`, err)
	}

	thread.ResList = resList

	return thread, nil
}

func parseHeader(tdResList *goquery.Selection) (*bakusai.Thread, error) {
	thread := bakusai.Thread{}

	// タイトル
	dtTitle := tdResList.Find(`dt.titleheader`)
	if dtTitle.Length() == 0 {
		return nil, fmt.Errorf(`missing: dt.titleheader`)
	}

	thread.Title = dtTitle.Find(`div#title_thr`).Text()

	// スレ立て日時
	strDatePublished := dtTitle.Find(`span[itemprop="datePublished"]`).Text()
	if strDatePublished == "" {
		return nil, fmt.Errorf(`empty: span[itemprop="datePublished"]`)
	}

	strDatePublished += ":00 +09:00"

	datePublished, err := time.Parse("2006/01/02 15:04:05 -07:00", strDatePublished)
	if err != nil {
		return nil, fmt.Errorf(`on time.Parse("%s"): %w`, strDatePublished, err)
	}

	thread.DatePublished = datePublished

	// ページ番号
	strPageNum := tdResList.Find(`div.paging span.paging_number`).Text()

	pageNum, err := strconv.Atoi(strPageNum)
	if err != nil {
		return nil, fmt.Errorf(`on strconv.Atoi("%s"): %w`, strPageNum, err)
	}

	thread.PageNum = pageNum

	// 前スレ・次スレ
	divPager := tdResList.Find(`div#thr_pager`)

	uriPrev := divPager.Find(`div.sre_mae a`).AttrOr("href", "")
	if uriPrev != "" {
		t, err := ParseThreadURI(uriPrev)
		if err != nil {
			return nil, fmt.Errorf(`on ParseThreadURI("%s"): %w`, uriPrev, err)
		}

		thread.PrevTID = t.ThreadID
	}

	uriNext := divPager.Find(`div.sre_tugi a`).AttrOr("href", "")
	if uriNext != "" {
		t, err := ParseThreadURI(uriNext)
		if err != nil {
			return nil, fmt.Errorf(`on ParseThreadURI("%s"): %w`, uriNext, err)
		}

		thread.NextTID = t.ThreadID
	}

	return &thread, nil
}
