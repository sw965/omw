package slices

import (
	"golang.org/x/exp/slices"
	"golang.org/x/exp/constraints"
	"github.com/sw965/omw"
)

func MakeFunc[XS ~[]X, X any](n int, f func(int) X) XS {
	y := make(XS, n)
	for i := 0; i < n; i++ {
		y[i] = f(i)
	}
	return y
}

type IntegerRange[NS ~[]N, N constraints.Integer] struct {
	Start N
	End N
	Step N
}

func(rng *IntegerRange[NS, N]) Make() NS {
	n := int((rng.End - 1 - rng.Start) / rng.Step) + 1
	y := make(NS, n)
	for i := 0; i < n; i++ {
		y[i] = rng.Start + (rng.Step * N(i))
	}
	return y
}

func IsSubset[XS ~[]X, X comparable](xs, subs XS) bool {
	for _, sub := range subs {
		if !slices.Contains(xs, sub) {
			return false
		}
	}
	return true
}

func Access[XS ~[]X, X any](xs XS, indices ...int) XS {
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
	return omw.Reverse(xs)
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

func ToUnique[XS ~[]X, X comparable](xs XS) XS {
	y := make(XS, 0, len(xs))
	for _, x := range xs {
		if !slices.Contains(y, x) {
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

func Permutation[XSS ~[]XS, XS ~[]X, X any](xs XS, r int) XSS {
	n := len(xs)
	idxss := omw.GetPermutation(n, r)
	access := func(idxs []int) XS {return Access(xs, idxs...) }
	return omw.MapFunc[[][]int, XSS](idxss, access)
}

func Combination[XSS ~[]XS, XS ~[]X, X any](xs XS, r int) XSS {
	n := len(xs)
	idxss := omw.GetCombination(n, r)
	access := func(idxs []int) XS { return Access(xs, idxs...) }
	return omw.MapFunc[[][]int, XSS](idxss, access)
}

func Pop[XS ~[]X, X any](xs XS, idx int) (XS, X) {
	result := make(XS, 0, len(xs) - 1)
	var y X
	for i, x := range xs {
		if i == idx {
			y = x
		} else {
			result = append(result, x)
		}
	}
	return result, y
}

func Sorted[XS ~[]X, X constraints.Ordered](xs XS) XS {
	clone := slices.Clone(xs)
	slices.Sort(clone)
	return clone
}