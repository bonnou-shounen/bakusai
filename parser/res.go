package parser

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/bonnou-shounen/bakusai"
)

func ParseResList(tdResList *goquery.Selection, tmpl *bakusai.Res) ([]*bakusai.Res, error) {
	divResList := tdResList.Find(`dl#res_list div.article`)
	if divResList.Length() == 0 {
		return nil, fmt.Errorf(`missing: dl#res_list div.article`)
	}

	resList := make([]*bakusai.Res, divResList.Length())

	var firstError error

	divResList.EachWithBreak(func(i int, articleDiv *goquery.Selection) bool {
		var res bakusai.Res

		if tmpl != nil {
			res = *tmpl // copy
		}

		if err := ParseRes(articleDiv, &res); err != nil {
			firstError = fmt.Errorf(`on ParseRes([%d]): %w`, i+1, err)

			return false
		}

		resList[i] = &res

		return true
	})

	if firstError != nil {
		return nil, firstError
	}

	return resList, nil
}

func ParseRes(divRes *goquery.Selection, res *bakusai.Res) error {
	strResnumb := divRes.Find(`span.resnumb`).Text()
	if strResnumb == "" {
		return fmt.Errorf(`empty: span.resnumb`)
	}

	strID := strings.TrimPrefix(strResnumb, "#")

	/*
		strID, exists := divRes.Attr("id")
		if !exists {
			return fmt.Errorf(`missing: divRes.Attr("id")`)
		}

		strID = strings.TrimPrefix(strID, "res")
	*/

	rrID, err := strconv.Atoi(strID)
	if err != nil {
		return fmt.Errorf(`on strconv.Atoi("%s"): %w`, strID, err)
	}

	res.ResID = rrID

	strCommentTime := divRes.Find(`span[itemprop="commentTime"]`).Text()
	if strCommentTime != "" {
		// 「削除済み」にはタイムスタンプがない
		strCommentTime += ":00 +09:00"

		commentTime, err := time.Parse("2006/01/02 15:04:05 -07:00", strCommentTime)
		if err != nil {
			return fmt.Errorf(`on time.Parse("%s"): %w`, strCommentTime, err)
		}

		res.CommentTime = commentTime
	}

	res.CommentText = divRes.Find(`div[itemprop="commentText"]`).Text()
	res.Name = divRes.Find(`dd.name`).Text()

	return nil
}
