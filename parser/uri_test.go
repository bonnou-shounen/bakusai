package parser_test

import (
	"reflect"
	"testing"

	"github.com/bonnou-shounen/bakusai"
	"github.com/bonnou-shounen/bakusai/parser"
)

func TestParseThreadURI(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		uri     string
		want    *bakusai.Thread
		wantErr bool
	}{
		{
			name: "success",
			uri:  "https://bakusai.com/thr_res_show/acode=1/ctgid=22/bid=333/tid=44444444/",
			want: &bakusai.Thread{AreaCode: 1, CategoryID: 22, BoardID: 333, ThreadID: 44444444},
		},
		{
			name: "ignore page number",
			uri:  "https://bakusai.com/thr_res_show/acode=1/ctgid=22/bid=333/tid=44444444/p=55/tp=1/rw=1",
			want: &bakusai.Thread{AreaCode: 1, CategoryID: 22, BoardID: 333, ThreadID: 44444444},
		},
		{
			name:    "wrong uri",
			uri:     "http://jbbs.shitaraba.net/bbs/read.cgi/computer/12345/67890/123/",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := parser.ParseThreadURI(tt.uri)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseThreadURI() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseResURI() = %v, want %v", got, tt.want)
			}
		})
	}
}
