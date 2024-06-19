package strings

import (
	"strings"
	"unicode/utf8"
)

func Len(s string) int {
	return utf8.RuneCountInString(s)
}

func EmptyPadding(s string, max int) string {
	n := Len(s)
	if n >= max {
		return s
	}

	ret := make([]string, max)
	for i, e := range strings.Split(s, "") {
		ret[i] = e
	}
	return strings.Join(ret, "")
}