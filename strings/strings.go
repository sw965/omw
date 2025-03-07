package strings

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func Len(s string) int {
	return utf8.RuneCountInString(s)
}

func PadLeft(s, pad string, size int) (string, error) {
	if Len(pad) != 1 {
		return "", fmt.Errorf("strings.Len(pad) != 1")
	}

	n := Len(s)
	if n >= size {
		return s, nil
	}

	ret := make([]string, size)

	for i := 0; i < size-n; i++ {
		ret[i] = pad
	}

	for i, e := range strings.Split(s, "") {
		ret[i+n] = e
	}

	return strings.Join(ret, ""), nil
}

func PadRight(s, pad string, size int) (string, error) {
	if Len(pad) != 1 {
		return "", fmt.Errorf("strings.Len(pad) != 1")
	}

	n := Len(s)
	if n >= size {
		return s, nil
	}

	ret := make([]string, size)
	for i, e := range strings.Split(s, "") {
		ret[i] = e
	}

	for i := n; i < size; i++ {
		ret[i] = pad
	}
	return strings.Join(ret, ""), nil
}
