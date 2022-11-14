package scraper

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/bonnou-shounen/bakusai"
	"github.com/bonnou-shounen/bakusai/util"
	"github.com/jesse0michael/errgroup"
)

type PageInfo struct {
	bakusai.Thread
	MaxPage int
}

func GetThread(argURI string) (*bakusai.Thread, error) {
	baseURI, err := util.CanonicalThreadURI(argURI)
	if err != nil {
		return nil, fmt.Errorf(`on util.CanonicalThreadURI("%s"): %w`, argURI, err)
	}

	ctx := context.Background()
	client := bakusai.NewClient(nil)

	info := PageInfo{}

	articles, err := getArticlesOnPage(ctx, client, baseURI, 1, &info)
	if err != nil {
		return nil, fmt.Errorf(`on getFirstPage(): %w`, err)
	}

	thread := bakusai.Thread{
		URI:       baseURI,
		CreatedAt: info.CreatedAt,
		Title:     info.Title,
		PrevURI:   info.PrevURI,
		NextURI:   info.NextURI,
		Articles:  articles,
	}

	maxPage := info.MaxPage

	if maxPage >= 2 {
		eg, egCtx := errgroup.WithContext(ctx, 5)

		articlesOnPage := make([][]*bakusai.Article, maxPage+1)

		for page := 2; page <= maxPage; page++ {
			page := page

			eg.Go(func() error {
				articlesOnPage[page], err = getArticlesOnPage(egCtx, client, baseURI, page, nil)
				if err != nil {
					return fmt.Errorf(`on getRestPages(): %w`, err)
				}

				return nil
			})
		}

		if err := eg.Wait(); err != nil {
			return nil, fmt.Errorf("on goroutine: %w", err)
		}

		for page := 2; page <= maxPage; page++ {
			thread.Articles = append(thread.Articles, articlesOnPage[page]...)
		}
	}

	return &thread, nil
}

func getArticlesOnPage(ctx context.Context, bakusaiClient *bakusai.Client, baseURI string, page int, pageInfo *PageInfo) ([]*bakusai.Article, error) { //nolint:lll
	uri := fmt.Sprintf("%s/tp=1/rw=1/p=%d", baseURI, page)

	var articles []*bakusai.Article

	if err := try10(ctx, bakusaiClient, uri, page, func(reslistTD *goquery.Selection) (bool, error) {
		var err error

		articles, err = parseArticles(reslistTD)
		if err != nil {
			return false, fmt.Errorf(`on parseArticles(): %w`, err)
		}

		if pageInfo != nil {
			if err = parseInfo(reslistTD, pageInfo); err != nil {
				return false, fmt.Errorf(`on parseInfo(): %w`, err)
			}

			maxPage, err := parseMaxPage(reslistTD)
			if err != nil {
				if maxPage < 0 {
					return false, fmt.Errorf(`on parseMaxPage(): %w`, err)
				}
			}

			pageInfo.MaxPage = maxPage
		}

		return false, nil
	}); err != nil {
		return nil, err
	}

	if pageInfo != nil && pageInfo.MaxPage == 0 {
		maxPage, err := getMaxPage(ctx, bakusaiClient, baseURI)
		if err != nil {
			maxPage = 20
		}

		pageInfo.MaxPage = maxPage
	}

	return articles, nil
}

func try10(ctx context.Context, bakusaiClient *bakusai.Client, baseURI string, page int,
	callback func(selection *goquery.Selection) (bool, error),
) error {
	var lastError error

	uri := fmt.Sprintf("%s/tp=1/rw=1/p=%d", baseURI, page)

	for try := 0; try < 20; try++ {
		if try > 0 {
			fmt.Fprintf(os.Stderr, "warning: page %d: retry %d: %s\n", page, try, lastError.Error())
			time.Sleep(time.Duration(try) * time.Second)
		}

		retry, err := func() (bool, error) {
			resp, err := bakusaiClient.Get(ctx, uri)
			if err != nil {
				return false, fmt.Errorf(`on bakusai.Client.Get(..."%s"...): %w`, uri, err)
			}
			defer resp.Body.Close()

			doc, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				return true, fmt.Errorf(`on goquery.NewDocumentFromReader(): %w`, err)
			}

			reslistTD := doc.Find(`table#inner_container td.reslist_td`)
			if reslistTD.Length() == 0 {
				return true, fmt.Errorf(`on doc.Find("td.reslist_td"): not found`)
			}

			return callback(reslistTD)
		}()
		if err == nil {
			return nil
		}

		lastError = err

		if !retry {
			break
		}
	}

	return lastError
}

func parseInfo(reslistTD *goquery.Selection, info *PageInfo) error {
	titleDT := reslistTD.Find(`dt.titleheader`)

	createdAtStr := titleDT.Find(`span[itemprop="datePublished"]`).Text() + ":00 +09:00"

	createdAt, err := time.Parse("2006/01/02 15:04:05 -07:00", createdAtStr)
	if err != nil {
		return fmt.Errorf(`on time.Parse("%s"): %w`, createdAtStr, err)
	}

	info.CreatedAt = createdAt

	info.Title = titleDT.Find(`div#title_thr`).Text()

	pagerDiv := reslistTD.Find(`div#thr_pager`)

	maeA := pagerDiv.Find(`div.sre_mae a`)
	if maeA.Length() > 0 {
		prevURI, ok := maeA.Attr("href")
		if ok {
			info.PrevURI = prevURI
		}
	}

	tsugiA := pagerDiv.Find(`div.sre_tsugi a`)
	if tsugiA.Length() > 0 {
		nextURI, ok := tsugiA.Attr("href")
		if ok {
			info.NextURI = nextURI
		}
	}

	return nil
}

func parseMaxPage(reslistTD *goquery.Selection) (int, error) {
	pagingDiv := reslistTD.Find(`div.paging`)
	if strings.HasSuffix(pagingDiv.Text(), "...") {
		return 0, fmt.Errorf(`need jump`)
	}

	nextSpan := pagingDiv.Find(`span.paging_nextlink`)
	if nextSpan.Length() == 0 {
		nextSpan = pagingDiv.Find(`span.paging_next`)
	}

	maxPage, err := strconv.Atoi(nextSpan.Prev().Text())
	if err != nil {
		return -1, fmt.Errorf(`on strconv.Atoi(): %w`, err)
	}

	return maxPage, nil
}

func getMaxPage(ctx context.Context, bakusaiClient *bakusai.Client, baseURI string) (int, error) {
	var maxPage int

	if err := try10(ctx, bakusaiClient, baseURI, 21, func(reslistTD *goquery.Selection) (bool, error) {
		var err error

		maxPage, err = parseMaxPage(reslistTD)
		if err != nil {
			return false, fmt.Errorf(`on parseMaxPage(): %w`, err)
		}

		return false, nil
	}); err != nil {
		return 20, err
	}

	return maxPage, nil
}
