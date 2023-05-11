package parser

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/bonnou-shounen/bakusai"
)

var rePathParam = regexp.MustCompile(`([^=/]+)=(\d+)`)

func CanonizeThreadURI(uri string) (string, error) {
	param := map[string]string{}

	matches := rePathParam.FindAllStringSubmatch(uri, -1)
	for _, m := range matches {
		param[m[1]] = m[2]
	}

	ids := map[string]int{}

	for _, key := range []string{"acode", "ctgid", "bid", "tid"} {
		id, err := strconv.Atoi(param[key])
		if err != nil {
			return "", fmt.Errorf(`on strconv.Atoi("%s"): %w`, param[key], err)
		}

		ids[key] = id
	}

	return fmt.Sprintf(
		"%sthr_res/acode=%d/ctgid=%d/bid=%d/tid=%d/",
		bakusai.RootURI,
		ids["acode"],
		ids["ctgid"],
		ids["bid"],
		ids["tid"],
	), nil
}
