package parser_test

import (
	"encoding/json"
	"fmt"
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

func AreEqualJSON(t *testing.T, got interface{}, want string) (bool, error) {
	t.Helper()

	b, err := json.Marshal(got)
	if err != nil {
		t.Errorf("on json.Marshal(got): %s", err.Error())
	}

	var o1 interface{}

	if err := json.Unmarshal(b, &o1); err != nil {
		return false, fmt.Errorf("on json.Unmarshal(got): %s", err.Error())
	}

	var o2 interface{}

	if err := json.Unmarshal([]byte(want), &o2); err != nil {
		return false, fmt.Errorf("on json.Unmarshal(want): %s", err.Error())
	}

	return reflect.DeepEqual(o1, o2), nil
}

func TestParseThread(t *testing.T) {
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
							<div itemprop="commentText">コメント1</div>
						</div>
						<div class="article" id="res2">
							<div itemprop="commentText">コメント2</div>
						</div>
					</dl>
					<div class="paging">
						<span class="paging_number">1</span>
					</div>
				</td>
			</table>`,
			want: `{"AreaCode":0, "CategoryID":0, "BoardID":0, "ThreadID":0, "PrevTID":0, "NextTID":0,
				"PageNum":1, "DatePublished":"2022-02-03T12:34:00+09:00",
				"Author":"",
				"Title":"スレタイ",
				"ResList":[
					{"AreaCode":0, "CategoryID":0, "BoardID":0, "ThreadID":0,
						"RRID":1,
						"CommentTime":"0001-01-01T00:00:00Z",
						"Name":"",
						"CommentText":"コメント1"
					},
					{"AreaCode":0, "CategoryID":0, "BoardID":0, "ThreadID":0,
						"RRID":2,
						"CommentTime":"0001-01-01T00:00:00Z",
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

			equals, err := AreEqualJSON(t, got, tt.want)
			if err != nil {
				t.Errorf("AreEqualJSON() error = %v", err)
			}

			if !equals {
				t.Errorf("ParseThread() = %v, want %v", got, tt.want)
			}
		})
	}
}
