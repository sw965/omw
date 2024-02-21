package slices

import (
	"github.com/sw965/omw"
	"golang.org/x/exp/slices"
	"fmt"
	"golang.org/x/exp/constraints"
)

func NewSequentialInteger[IS ~[]I, I constraints.Integer](start, end I) IS {
	n := end - start
	y := make(IS, int(n))

	for i := I(0); i < n; i++ {
		y[i] = i
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

func IndexAccess[XS ~[]X, X any](xs XS) func(int) X {
	return func(idx int) X {
		return xs[idx]
	}
}

func IndicesAccess[XS ~[]X, X any](xs XS) func([]int) XS {
	return func(idxs []int) XS {
		y := make(XS, len(idxs))
		for i, idx := range idxs {
			y[i] = xs[idx]
		}
		return y
	}
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
	return omw.MapFunc[XSS](idxss, IndicesAccess(xs))
}

func Combination[XSS ~[]XS, XS ~[]X, X any](xs XS, r int) XSS {
	n := len(xs)
	idxss := omw.GetCombination(n, r)
	return omw.MapFunc[XSS](idxss, IndicesAccess(xs))
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

func End[XS ~[]X, X any](xs XS) (X, error) {
	if len(xs) == 0 {
		var x X
		return x, fmt.Errorf("len(xs) == 0")
	}
	return xs[len(xs)-1], nil
}

func Sorted[XS ~[]X, X constraints.Ordered](xs XS) XS {
	clone := slices.Clone(xs)
	slices.Sort(clone)
	return clone
}

func Delete[XS ~[]X, X any](xs XS, idxs ...int) (XS, XS) {
	n := len(xs)
	d := make(XS, 0, len(idxs))
	y := make(XS, 0, n)
	for i := 0; i < n; i++ {
		x := xs[i]
		if slices.Contains(idxs, i) {
			d = append(d, x)
		} else {
			y = append(y, x)
		}
	}
	return y, d
}

func Fill[XS ~[]X, X any](xs XS, x X) XS {
	y := make(XS, len(xs))
	for i := range xs {
		y[i] = x
	}
	return y
}

func Concat[XS ~[]X, X any](xs1 XS, xs2 XS) XS {
	y := make(XS, 0, len(xs1) + len(xs2))
	for _, x1 := range xs1 {
		y = append(y, x1)
	}

	for _, x2 := range xs2 {
		y = append(y, x2)
	}
	return y
}

func Binary(n int) []int {
    y := make([]int, 0, (n/2)+1)
    q := n
    for {
        m := q%2
        y = append(y, m)
        q = q>>1
        if q < 2 {
            y = append(y, q)
            break
        }
    }
    return omw.Reverse(y)
}

func Product2[XS1 ~[]X1, XS2 ~[]X2, YS ~[]Y, X1, X2, Y any](xs1 XS1, xs2 XS2, f func(X1, X2) Y) YS {
	ys := make(YS, 0, len(xs1) * len(xs2))
	for _, x1 := range xs1 {
		for _, x2 := range xs2 {
			ys = append(ys, f(x1, x2))
		}
	}
	return ys
}

func Product3[XS1 ~[]X1, XS2 ~[]X2, XS3 ~[]X3, YS ~[]Y, X1, X2, X3, Y any](xs1 XS1, xs2 XS2, xs3 XS3, f func(X1, X2, X3) Y) YS {
	ys := make(YS, 0, len(xs1) * len(xs2) * len(xs3))
	for _, x1 := range xs1 {
		for _, x2 := range xs2 {
			for _, x3 := range xs3 {
				ys = append(ys, f(x1, x2, x3))
			}
		}
	}
	return ys
}