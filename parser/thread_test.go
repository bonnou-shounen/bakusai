package parser_test

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/bonnou-shounen/bakusai/parser"
)

func GetTDResList(t *testing.T, html string) *goquery.Selection {
	t.Helper()

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatal(err)
	}

	tdResList := doc.Find("td.reslist_td")
	if tdResList.Length() == 0 {
		t.Fatal("missing: td.reslist_td")
	}

	return tdResList
}

func AreEqualJSON(t *testing.T, got interface{}, want string) bool {
	t.Helper()

	b, err := json.Marshal(got)
	if err != nil {
		t.Fatal(err)
	}

	var o1 interface{}

	if err := json.Unmarshal(b, &o1); err != nil {
		t.Fatal(err)
	}

	var o2 interface{}

	if err := json.Unmarshal([]byte(want), &o2); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(o1, o2) {
		t.Errorf("got: %s\n", string(b))

		return false
	}

	return true
}

func TestParseThread(t *testing.T) { //nolint:funlen
	t.Parallel()

	tests := []struct {
		name    string
		html    string
		want    string
		wantErr bool
	}{
		{
			name: "normal",
			html: `<table>
				<td class="reslist_td">
					<dt class="titleheader">
						<div id="title_thr">スレタイ</div>
						<span itemprop="datePublished">2022/02/03 12:34</span>
					</dt>
					<dl id="res_list">
						<div class="article" id="res1">
							<span itemprop="commentTime">2022/02/04 23:45</span>
							<div itemprop="commentText">コメント1</div>
						</div>
						<div class="article" id="res2">
							<span itemprop="commentTime">2022/02/05 06:17</span>
							<div itemprop="commentText">コメント2</div>
						</div>
					</dl>
					<div class="paging">
						<span class="paging_number">1</span>
					</div>
				</td>
			</table>`,
			want: `{"AreaCode":0, "CategoryID":0, "BoardID":0, "ThreadID":0, "PrevURI":"", "NextURI":"",
				"DatePublished":"2022-02-03T12:34:00+09:00",
				"Author":"",
				"Title":"スレタイ",
				"ResList":[
					{"AreaCode":0, "CategoryID":0, "BoardID":0, "ThreadID":0,
						"RRID":1,
						"CommentTime":"2022-02-04T23:45:00+09:00",
						"Name":"",
						"CommentText":"コメント1"
					},
					{"AreaCode":0, "CategoryID":0, "BoardID":0, "ThreadID":0,
						"RRID":2,
						"CommentTime":"2022-02-05T06:17:00+09:00",
						"Name":"",
						"CommentText":"コメント2"
					}
				]
			}`,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := parser.ParseThread(GetTDResList(t, tt.html))
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseThread() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !AreEqualJSON(t, got, tt.want) {
				t.Errorf("ParseThread() = %v, want %v", got, tt.want)
			}
		})
	}
}
