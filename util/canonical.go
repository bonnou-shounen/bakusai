package util

import (
	"fmt"
	"regexp"
	"strconv"
)

var uriParamsRe = regexp.MustCompile(`([^=/]+)=(\d+)`)

func CanonicalThreadURI(uri string) (string, error) {
	m := map[string]int{}

	matches := uriParamsRe.FindAllStringSubmatch(uri, -1)
	for _, match := range matches {
		m[match[1]], _ = strconv.Atoi(match[2])
	}

	keys := []string{"acode", "ctgid", "bid", "tid"}

	for _, key := range keys {
		if _, ok := m[key]; !ok {
			return "", fmt.Errorf("missing key: %s", key)
		}
	}

	return fmt.Sprintf(
		"https://bakusai.com/thr_res/acode=%d/ctgid=%d/bid=%d/tid=%d",
		m["acode"], m["ctgid"], m["bid"], m["tid"],
	), nil
}
