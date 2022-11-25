package parser

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/bonnou-shounen/bakusai"
)

func ParseResURI(uri string) (*bakusai.Res, error) {
	res := bakusai.Res{}

	keys := []string{"acode", "ctgid", "bid", "tid", "rrid"}

	param, err := parseURI(uri, keys)
	if err != nil {
		return nil, err
	}

	res.AreaCode = param["acode"]
	res.CategoryID = param["ctgid"]
	res.BoardID = param["bid"]
	res.ThreadID = param["tid"]
	res.ResID = param["rrid"]

	return &res, nil
}

func ParseThreadURI(uri string) (*bakusai.Thread, error) {
	thread := bakusai.Thread{}

	keys := []string{"acode", "ctgid", "bid", "tid"}

	param, err := parseURI(uri, keys)
	if err != nil {
		return nil, err
	}

	thread.AreaCode = param["acode"]
	thread.CategoryID = param["ctgid"]
	thread.BoardID = param["bid"]
	thread.ThreadID = param["tid"]

	return &thread, nil
}

var uriParamsRe = regexp.MustCompile(`([^=/]+)=(\d+)`)

func parseURI(uri string, keys []string) (map[string]int, error) {
	param := map[string]int{}

	matches := uriParamsRe.FindAllStringSubmatch(uri, -1)
	for _, match := range matches {
		param[match[1]], _ = strconv.Atoi(match[2])
	}

	for _, key := range keys {
		if _, ok := param[key]; !ok {
			return nil, fmt.Errorf("missing: %s", key)
		}
	}

	return param, nil
}
