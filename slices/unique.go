package slices

import (
	"golang.org/x/exp/slices"
)

func ToUnique[XS ~[]X, X comparable](xs XS) XS {
	y := make(XS, 0, len(xs))
	for _, x := range xs {
		if !slices.Contains(xs, x) {
			y = append(y, x)
		}
	}
	return y
}

func IsUnique[XS ~[]X, X comparable](xs XS) bool {
	for _, x := range xs {
		if Count(xs, x) != 1 {
			return false
		}
	}
	return true
}