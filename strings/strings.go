package strings

import (
	"strings"
)

func Replace(old, new string, n int) func(string) string {
	return func(s string) string {
		return strings.Replace(s, old, new, n)
	}
}