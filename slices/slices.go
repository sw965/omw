package slices

import (
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"

	omwmath "github.com/sw965/omw/math"
)

func MakeInteger[S ~[]I, I constraints.Integer](start, end I) S {
	n := end - start
	ret := make(S, int(n))
	for i := I(0); i < n; i++ {
		ret[i] = i
	}
	return ret
}

func Concat[S ~[]E, E any](ss ...S) S {
	n := 0
	for _, s := range ss {
		n += len(s)
	}

	ret := make(S, 0, n)
	for _, s := range ss {
		ret = append(ret, s...)
	}
	return ret
}

func Contains[S ~[]E, E comparable](s S) func(E) bool {
	return func(e E) bool {
		return slices.Contains(s, e)
	}
}

func Equal[S ~[]E, E comparable](s1 S) func(S)bool {
	return func(s2 S) bool {
		return slices.Equal(s1, s2)
	}
}

func Reverse[S ~[]E, E any](s S) S {
	n := len(s)
	ret := make(S, 0, n)
	for i := n - 1; i > -1; i-- {
		ret = append(ret, s[i])
	}
	return ret
}

func Count[S ~[]E, E comparable](s S, e E) int {
	ret := 0
	for _, si := range s {
		if si == e {
			ret += 1
		}
	}
	return ret
}

func Indices[S ~[]E, E comparable](s S, e E) []int {
	ret := make([]int, 0, len(s))
	for i, a := range s {
		if a == e {
			ret = append(ret, i)
		}
	}
	return ret
}

func IndicesFunc[S ~[]E, E any](s S, f func(E) bool) []int {
	ret := make([]int, 0, len(s))
	for i, e := range s {
		if f(e) {
			ret = append(ret, i)
		}
	}
	return ret
}

func CountFunc[S ~[]E, E any](s S, f func(x E) bool) int {
	ret := 0
	for _, si := range s {
		if f(si) {
			ret += 1
		}
	}
	return ret
}

func MinIndex[S ~[]E, E constraints.Ordered](s S) int {
	min := s[0]
	idx := 0
	for i, e := range s {
		if e < min {
			min = e
			idx = i
		}
	}
	return idx
}

func MinIndices[S ~[]E, E constraints.Ordered](s S) []int {
	min := omwmath.Min(s...)
	idxs := make([]int, 0, len(s))
	for i, e := range s {
		if e == min {
			idxs = append(idxs, i)
		}
	}
	return idxs
}

func MaxIndex[S ~[]E, E constraints.Ordered](s S) int {
	max := s[0]
	idx := 0
	for i, e := range s {
		if e > max {
			max = e
			idx = i
		}
	}
	return idx
}

func MaxIndices[S ~[]E, E constraints.Ordered](s S) []int {
	max := omwmath.Max(s...)
	ret := make([]int, 0, len(s))
	for i, e := range s {
		if e == max {
			ret = append(ret, i)
		}
	}
	return ret
}

func IndicesAccess[S ~[]E, E any](s S, idxs ...int) S {
	ret := make(S, len(idxs))
	for i, idx := range idxs {
		ret[i] = s[idx]
	}
	return ret
}

func IntPermutation(n, r int) [][]int {
	c := omwmath.PermutationCount(n, r)
	ret := make([][]int, 0, c)
	if r == 0 {
		return ret
	}
	var f func(int, []int)
	f = func(nest int, nums []int) {
		if nest == r {
			ret = append(ret, nums)
			return
		}
		for i := 0; i < n; i++ {
			isContinue := false
			for _, num := range nums {
				if num == i {
					isContinue = true
					break
				}
			}
			if isContinue {
				continue
			}
			clone := slices.Clone(nums)
			f(nest+1, append(clone, i))
		}
	}
	f(0, make([]int, 0, r))
	return ret
}

func Permutation[SS ~[]S, S ~[]E, E any](s S, r int) SS {
	n := len(s)
	idxss := IntPermutation(n, r)
	ret := make(SS, len(idxss))
	for i, idxs := range idxss {
		ret[i] = IndicesAccess(s, idxs...)
	}
	return ret
}

func IntSequence(n, r int) [][]int {
	c := omwmath.SequenceCount(n, r)
	ret := make([][]int, 0, c)
	if r == 0 {
		return ret
	}
	var f func(int, []int)
	f = func(nest int, nums []int) {
		if nest == r {
			ret = append(ret, nums)
			return
		}
		for i := 0; i < n; i++ {
			clone := slices.Clone(nums)
			f(nest+1, append(clone, i))
		}
	}
	f(0, make([]int, 0, r))
	return ret
}

func Sequence[SS ~[]S, S ~[]E, E any](s S, r int) SS {
	n := len(s)
	idxss := IntSequence(n, r)
	ret := make(SS, len(idxss))
	for i, idxs := range idxss {
		ret[i] = IndicesAccess(s, idxs...)
	}
	return ret
}

func IntCombination(n, r int) [][]int {
	nums := make([]int, r)
	for i := 0; i < r; i++ {
		nums[i] = i
	}

	c := omwmath.CombinationCount(n, r)
	ret := make([][]int, 0, c)
	if r == 0 {
		return ret
	}

	end := r - 1
	for i := 0; i < c; i++ {
		clone := slices.Clone(nums)
		ret = append(ret, clone)
		max := omwmath.Max(nums...)
		if max == (n - 1) {
			reversed := Reverse(nums)
			count := omwmath.CountConsecutiveDecrease(reversed...)
			idx := end - count
			if idx < 0 {
				break
			}
			nums[idx] += 1
			for j := idx + 1; j < r; j++ {
				nums[j] = nums[idx] + j - (idx)
			}
		} else {
			nums[end] += 1
		}
	}
	return ret
}

func Combination[SS ~[]S, S ~[]E, E any](s S, r int) SS {
	n := len(s)
	idxss := IntCombination(n, r)
	ret := make(SS, len(idxss))
	for i, idxs := range idxss {
		ret[i] = IndicesAccess(s, idxs...)
	}
	return ret
}

func Any(bs []bool) bool {
	for _, b := range bs {
		if b {
			return true
		}
	}
	return false
}

func AnyFunc[S ~[]E, E any](s S, f func(E) bool) bool {
	for _, e := range s {
		if f(e) {
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

func AllFunc[S ~[]E, E any](s S, f func(E) bool) bool {
	for _, e := range s {
		if !f(e) {
			return false
		}
	}
	return true
}

func End[XS ~[]X, X any](xs XS) X {
	n := len(xs)
	return xs[n-1]
}