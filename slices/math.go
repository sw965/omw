package slices

import (
	"github.com/sw965/omw"
	"golang.org/x/exp/slices"
)

func Permutation[XSS ~[]XS, XS ~[]X, X any](xs XS, r int) XSS {
	n := len(xs)
	numss := omw.Math_Permutation(n, r)
	access := func(nums []int) XS {return Access(xs, nums...) }
	return omw.Fn_Map[[][]int, XSS](numss, access)
}

func Combination[XSS ~[]XS, XS ~[]X, X any](xs XS, r int) XSS {
	n := len(xs)
	numss := omw.Math_Combination(n, r)
	access := func(nums []int) XS { return Access(xs, nums...) }
	return omw.Fn_Map[[][]int, XSS](numss, access)
}

func IsSubset[XS ~[]X, X comparable](xs, subs XS) bool {
	for _, sub := range subs {
		if !slices.Contains(xs, sub) {
			return false
		}
	}
	return true
}