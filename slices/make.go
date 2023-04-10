package slices

import (
	"github.com/sw965/omw"
	"golang.org/x/exp/constraints"
)

func MakeFunc[XS ~[]X, X any](n int, f func(int) X) XS {
	y := make(XS, n)
	for i := 0; i < n; i++ {
		y[i] = f(i)
	}
	return y
}

func MakeIntegerRange[XS ~[]X, X constraints.Integer](start, end, step X) XS {
	return omw.Slices_MakeIntegerRange[XS](start, end, step)
}