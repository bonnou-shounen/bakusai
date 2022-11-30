package parser

import (
	"fmt"
	"regexp"
	"strconv"
)

const URITop = "https://bakusai.com"

var uriParamsRe = regexp.MustCompile(`([^=/]+)=(\d+)`)

func CanonizeThreadURI(uri string) (string, error) {
	param := map[string]string{}

	matches := uriParamsRe.FindAllStringSubmatch(uri, -1)
	for _, m := range matches {
		param[m[1]] = m[2]
	}

	idOf := map[string]int{}

	for _, key := range []string{"acode", "ctgid", "bid", "tid"} {
		id, err := strconv.Atoi(param[key])
		if err != nil {
			return "", fmt.Errorf(`on strconv.Atoi("%s"): %w`, param[key], err)
		}

		idOf[key] = id
	}

	return fmt.Sprintf(
		"%s/thr_res/acode=%d/ctgid=%d/bid=%d/tid=%d/",
		URITop,
		idOf["acode"],
		idOf["ctgid"],
		idOf["bid"],
		idOf["tid"],
	), nil
}
