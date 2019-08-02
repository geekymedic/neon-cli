package util

import (
	"strings"
)

func SplitTrimSpace(s, sep string) []string {
	var ret []string
	for _, item := range strings.Split(strings.TrimSpace(s), sep) {
		ret = append(ret, strings.TrimSpace(item))
	}
	return ret
}
