package omw

import (
	"golang.org/x/exp/constraints"
)

func MakeRangeInteger[IS ~[]I, I constraints.Integer](start, end I) IS {
	n := end - start
	y := make(IS, int(n))
	for i := I(0); i < n; i++ {
		y[i] = i
	}
	return y
}

func Reverse[XS ~[]X, X any](xs XS) XS {
	n := len(xs)
	y := make(XS, 0, n)
	for i := n - 1; i > -1; i-- {
		y = append(y, xs[i])
	}
	return y
}

func Concat[XS ~[]X, X any](xs1 XS, xs2 XS) XS {
	y := make(XS, 0, len(xs1) + len(xs2))
	for i := range xs1 {
		y = append(y, xs1[i])
	}

	for i := range xs2 {
		y = append(y, xs2[i])
	}
	return y
}

func CountElement[XS ~[]X, X comparable](xs XS, e X) int {
	count := 0
	for i := range xs {
		if xs[i] == e {
			count += 1
		}
	}
	return count
}

func CountIf[XS ~[]X, X any](xs XS, f func(x X) bool) int {
	count := 0
	for i := range xs {
		if f(xs[i]) {
			count += 1
		}
	}
	return count
}

func MinIndex[XS ~[]X, X constraints.Ordered](xs XS) int {
	min := xs[0]
	idx := 0
	for i := range xs {
		xi := x[i]
		if xi < min {
			min = xi
			idx = i
		}
	}
	return idx
}

func MinIndices[XS ~[]X, X constraints.Ordered](xs XS) []int {
	min := Min(xs...)
	idxs := make([]int, 0, len(xs))
	for i := range xs {
		if xs[i] == min {
			idxs = append(idxs, i)
		}
	}
	return idxs
}

func MaxIndex[XS ~[]X, X constraints.Ordered](xs XS) int {
	max := xs[0]
	idx := 0
	for i := range xs {
		xi := x[i]
		if xi > max {
			max = xi
			idx = i
		}
	}
	return idx
}

func MaxIndices[XS ~[]X, X constraints.Ordered](xs XS) []int {
	max := Max(xs...)
	idxs := make([]int, 0, len(xs))
	for i := range xs {
		if xs[i] == max {
			idxs = append(idxs, i)
		}
	}
	return idxs
}

func ElementsAtIndices[XS ~[]X, X any](xs XS, idxs ...int) XS {
	y := make(XS, len(idxs))
	for i := range idxs {
		y[i] = xs[idxs[i]]
	}
	return y
}

func Permutations[XSS ~[]XS, XS ~[]X, X any](xs XS, r int) XSS {
	n := len(xs)
	idxss := IntPermutations(n, r)
	result := make(XSS, len(idxss))
	for i := range idxss {
		result[i] = ElementsAtIndices(xs, idxss[i]...)
	}
	return result
}

func Combinations[XSS ~[]XS, XS ~[]X, X any](xs XS, r int) XSS {
	n := len(xs)
	idxss := IntCombinations(n, r)
	result := make(XSS, len(idxss))
	for i := range idxss {
		result[i] = ElementsAtIndices(xs, idxss[i]...)
	}
	return result
}

func Any(bs []bool) bool {
	for _, b := range bs {
		if b {
			return true
		}
	}
	return false
}

func AnyMatch[XS ~[]X, X any](xs XS, f func(X) bool) bool {
	for _, x := range xs {
		if f(x) {
			return true
		}
	}
	return false
}

func All(bs []bool) bool {
	for _, b := range bs {
		if !b {
			return false
		}
	}
	return true
}

func AllMatch[XS ~[]X, X any](xs XS, f func(X) bool) bool {
	for _, x := range xs {
		if !f(x) {
			return false
		}
	}
	return true
}