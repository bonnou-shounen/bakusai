package scraper

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/bonnou-shounen/bakusai"
)

func parseArticle(articleDiv *goquery.Selection) (*bakusai.Article, error) {
	idStr, exists := articleDiv.Attr("id")
	if !exists {
		return nil, fmt.Errorf(`div.Attr("id"): not found`)
	}

	article := &bakusai.Article{}

	var err error

	article.ID, err = strconv.Atoi(strings.TrimPrefix(idStr, "res"))
	if err != nil {
		return nil, fmt.Errorf(`strconv.Atoi("%s"): %w`, idStr, err)
	}

	timeSpan := articleDiv.Find(`span[itemprop="commentTime"]`)
	if timeSpan.Length() > 0 {
		postAtStr := timeSpan.Text() + ":00 +09:00"

		article.PostAt, err = time.Parse("2006/01/02 15:04:05 -07:00", postAtStr)
		if err != nil {
			return nil, fmt.Errorf(`time.Parse("%s"): %w`, postAtStr, err)
		}
	}

	article.BodyText = articleDiv.Find(`div[itemprop="commentText"]`).Text()

	article.AuthorName = articleDiv.Find(`dd.name`).Text()

	return article, nil
}

func parseArticles(reslistTD *goquery.Selection) ([]*bakusai.Article, error) {
	articleDivs := reslistTD.Find(`dl#res_list div.article`)
	if articleDivs.Length() == 0 {
		return nil, fmt.Errorf(`on td.Find("dl#res_list div.article"): not found`)
	}

	var articles []*bakusai.Article

	var firstError error

	articleDivs.EachWithBreak(func(n int, articleDiv *goquery.Selection) bool {
		article, err := parseArticle(articleDiv)
		if err != nil {
			firstError = fmt.Errorf(`on getArticle(%d): %w`, n, err)

			return false
		}

		articles = append(articles, article)

		return true
	})

	if firstError != nil {
		return nil, firstError
	}

	return articles, nil
}
