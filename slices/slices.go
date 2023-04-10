package slices

import (
	"github.com/sw965/omw"
)

func Access[XS ~[]X, X any](xs XS , indices ...int) XS {
	y := make(XS, len(indices))
	for i, idx := range indices {
		y[i] = xs[idx]
	}
	return y
}

func Count[XS ~[]X, X comparable](xs XS, a X) int {
	y := 0
	for _, x := range xs {
		if x == a {
			y += 1
		}
	}
	return y
}

func Reverse[XS ~[]X, X any](xs XS) XS {
	return omw.Slices_Reverse(xs)
}

func Indices[XS ~[]X, X comparable](xs XS, a X) []int {
	y := make([]int, 0, len(xs))
	for i, x := range xs {
		if x == a {
			y = append(y, i)
		}
	}
	return y
}