package omw

import (
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

func MakeRangeInteger[S ~[]I, I constraints.Integer](start, end I) S {
	n := end - start
	ret := make(S, int(n))
	for i := I(0); i < n; i++ {
		ret[i] = i
	}
	return ret
}

func EqualSlice[S ~[]E, E comparable](s1 S) func(S)bool {
	return func(s2 S) bool {
		return slices.Equal(s1, s2)
	}
}

func ReverseSlice[S ~[]E, E any](s S) S {
	n := len(s)
	ret := make(S, 0, n)
	for i := n - 1; i > -1; i-- {
		ret = append(ret, s[i])
	}
	return ret
}

func CountElement[S ~[]E, E comparable](s S, e E) int {
	ret := 0
	for _, si := range s {
		if si == e {
			ret += 1
		}
	}
	return ret
}

func CountElementFunc[S ~[]E, E any](s S, f func(x E) bool) int {
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
	min := Min(s...)
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
	max := Max(s...)
	ret := make([]int, 0, len(s))
	for i, e := range s {
		if e == max {
			ret = append(ret, i)
		}
	}
	return ret
}

func ElementsAtIndices[S ~[]E, E any](s S, idxs ...int) S {
	ret := make(S, len(idxs))
	for i, idx := range idxs {
		ret[i] = s[idx]
	}
	return ret
}

func Permutations[SS ~[]S, S ~[]E, E any](s S, r int) SS {
	n := len(s)
	idxss := IntPermutations(n, r)
	ret := make(SS, len(idxss))
	for i, idxs := range idxss {
		ret[i] = ElementsAtIndices(s, idxs...)
	}
	return ret
}

func Combinations[SS ~[]S, S ~[]E, E any](s S, r int) SS {
	n := len(s)
	idxss := IntCombinations(n, r)
	ret := make(SS, len(idxss))
	for i, idxs := range idxss {
		ret[i] = ElementsAtIndices(s, idxs...)
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