package slices

import (
	"github.com/sw965/omw"
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
	"fmt"
)

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

func CountFunc[XS ~[]X, X any](xs XS, f func(x X) bool) int {
	y := 0
	for _, x := range xs {
		if f(x) {
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

func IndicesFunc[XS ~[]X, X any](xs XS, f func(X) bool) []int {
	y := make([]int, 0, len(xs))
	for i, x := range xs {
		if f(x) {
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
	access := func(idxs []int) XS { return Access(xs, idxs...) }
	return omw.MapFunc[[][]int, XSS](idxss, access)
}

func Combination[XSS ~[]XS, XS ~[]X, X any](xs XS, r int) XSS {
	n := len(xs)
	idxss := omw.GetCombination(n, r)
	access := func(idxs []int) XS { return Access(xs, idxs...) }
	return omw.MapFunc[[][]int, XSS](idxss, access)
}

func Sorted[XS ~[]X, X constraints.Ordered](xs XS) XS {
	clone := slices.Clone(xs)
	slices.Sort(clone)
	return clone
}

func SortedFunc[XS ~[]X, X any](xs XS, f func(a, b X) bool) XS {
	clone := slices.Clone(xs)
	slices.SortFunc(clone, f)
	return clone
}

func Product[XSS ~[]XS, XS ~[]X, X any](xss ...XS) XSS {
	n := len(xss)
	m := 1
	for _, xs := range xss {
		m *= len(xs)
	}
	result := make(XSS, 0, m)

	var f func(nest int, xs XS)
	f = func(nest int, ys XS) {
		if nest == n {
			result = append(result, ys)
			return
		}

		for _, x := range xss[nest] {
			clone := slices.Clone(ys)
			f(nest+1, append(clone, x))
		}
	}
	f(0, make(XS, 0, n))
	return result
}

func All(bs []bool) bool {
	for _, b := range bs {
		if !b {
			return false
		}
	}
	return true
}

func Any(bs []bool) bool {
	for _, b := range bs {
		if b {
			return true
		}
	}
	return false
}

func GetEnd[XS ~[]X, X any](xs XS) (X, error) {
	if len(xs) == 0 {
		var x X
		return x, fmt.Errorf("len(xs) == 0")
	}
	return xs[len(xs)-1], nil
}
