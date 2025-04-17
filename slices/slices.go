package slices

import (
	omwmath "github.com/sw965/omw/math"
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

func MakeInteger[S ~[]I, I constraints.Integer](start, end I) S {
	n := end - start
	s := make(S, int(n))
	for i := I(0); i < n; i++ {
		s[i] = i
	}
	return s
}

func Concat[S ~[]E, E any](ss ...S) S {
	n := 0
	for _, s := range ss {
		n += len(s)
	}

	s := make(S, 0, n)
	for _, si := range ss {
		s = append(s, si...)
	}
	return s
}

func Contains[S ~[]E, E comparable](s S) func(E) bool {
	return func(e E) bool {
		return slices.Contains(s, e)
	}
}

func Equal[S ~[]E, E comparable](s1 S) func(S) bool {
	return func(s2 S) bool {
		return slices.Equal(s1, s2)
	}
}

func Reverse[S ~[]E, E any](s S) S {
	n := len(s)
	r := make(S, 0, n)
	for i := n - 1; i > -1; i-- {
		r = append(r, s[i])
	}
	return r
}

func Count[S ~[]E, E comparable](s S, e E) int {
	c := 0
	for _, si := range s {
		if si == e {
			c += 1
		}
	}
	return c
}

func Indices[S ~[]E, E comparable](s S, e E) []int {
	idxs := make([]int, 0, len(s))
	for i, si := range s {
		if si == e {
			idxs = append(idxs, i)
		}
	}
	return idxs
}

func IndicesFunc[S ~[]E, E any](s S, f func(E) bool) []int {
	idxs := make([]int, 0, len(s))
	for i, e := range s {
		if f(e) {
			idxs = append(idxs, i)
		}
	}
	return idxs
}

func CountFunc[S ~[]E, E any](s S, f func(x E) bool) int {
	c := 0
	for _, si := range s {
		if f(si) {
			c += 1
		}
	}
	return c
}

func MinIndex[S ~[]E, E constraints.Ordered](s S) int {
    min := s[0]
    idx := 0
    for i, e := range s[1:] {
        if e < min {
            min = e
            idx = i + 1
        }
    }
    return idx
}

func MinIndices[S ~[]E, E constraints.Ordered](s S) []int {
    n := len(s)
    if n == 0 {
        return nil
    }

    min := s[0]
    idxs := make([]int, 0, n)
    idxs = append(idxs, 0)

    for i := 1; i < n; i++ {
        e := s[i]
        switch {
        case e < min:
            min = e
            // capacityを残したままスライスを空にする。
            idxs = idxs[:0]
            idxs = append(idxs, i)
        case e == min:
            idxs = append(idxs, i)
        }
    }
    return idxs
}

func MaxIndex[S ~[]E, E constraints.Ordered](s S) int {
	max := s[0]
	idx := 0
	for i, e := range s[1:] {
		if e > max {
			max = e
			idx = i + 1
		}
	}
	return idx
}

func MaxIndices[S ~[]E, E constraints.Ordered](s S) []int {
	n := len(s)
	if n == 0 {
		return nil
	}

	max := s[0]
	idxs := make([]int, 0, n)
	idxs = append(idxs, 0)

	for i := 1; i < n; i++ {
		e := s[i]
		switch {
		case e > max:
			max = e
			//capacityを残したままスライスを空にする。
			idxs = idxs[:0]
			idxs = append(idxs, i)
		case e == max:
			idxs = append(idxs, i)
		}
	}
	return idxs
}

func Argsort[S ~[]E, E constraints.Ordered](s S) []int {
	idxs := make([]int, len(s))
	for i := range s {
		idxs[i] = i
	}

	slices.SortFunc(idxs, func(idx1, idx2 int) bool {
		return s[idx1] < s[idx2]
	})

	return idxs
}

func ElementsByIndices[S ~[]E, E any](s S, idxs ...int) S {
	es := make(S, len(idxs))
	for i, idx := range idxs {
		es[i] = s[idx]
	}
	return es
}

func IntPermutation(n, r int) [][]int {
	c := omwmath.PermutationCount(n, r)
	result := make([][]int, 0, c)
	if r == 0 {
		return result
	}
	var f func(int, []int)
	f = func(nest int, nums []int) {
		if nest == r {
			result = append(result, nums)
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
	return result
}

func Permutation[SS ~[]S, S ~[]E, E any](s S, r int) SS {
	n := len(s)
	idxss := IntPermutation(n, r)
	ss := make(SS, len(idxss))
	for i, idxs := range idxss {
		ss[i] = ElementsByIndices(s, idxs...)
	}
	return ss
}

// 重複順列
func IntSequence(n, r int) [][]int {
	c := omwmath.SequenceCount(n, r)
	result := make([][]int, 0, c)
	if r == 0 {
		return result
	}
	var f func(int, []int)
	f = func(nest int, nums []int) {
		if nest == r {
			result = append(result, nums)
			return
		}
		for i := 0; i < n; i++ {
			clone := slices.Clone(nums)
			f(nest+1, append(clone, i))
		}
	}
	f(0, make([]int, 0, r))
	return result
}

// 重複順列
func Sequence[SS ~[]S, S ~[]E, E any](s S, r int) SS {
	n := len(s)
	idxss := IntSequence(n, r)
	ss := make(SS, len(idxss))
	for i, idxs := range idxss {
		ss[i] = ElementsByIndices(s, idxs...)
	}
	return ss
}

func IntCombination(n, r int) [][]int {
	nums := make([]int, r)
	for i := 0; i < r; i++ {
		nums[i] = i
	}

	c := omwmath.CombinationCount(n, r)
	result := make([][]int, 0, c)
	if r == 0 {
		return result
	}

	end := r - 1
	for i := 0; i < c; i++ {
		clone := slices.Clone(nums)
		result = append(result, clone)
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
	return result
}

func Combination[SS ~[]S, S ~[]E, E any](s S, r int) SS {
	n := len(s)
	idxss := IntCombination(n, r)
	ss := make(SS, len(idxss))
	for i, idxs := range idxss {
		ss[i] = ElementsByIndices(s, idxs...)
	}
	return ss
}

func CartesianProduct[SS ~[]S, S ~[]E, E any](ss ...S) SS {
	if len(ss) == 0 {
		return SS{}
	}

	c := 1
	for _, s := range ss {
		c *= len(s)
	}

	result := make(SS, 0, c)

	var f func(nest int, nums S)
	f = func(nest int, nums S) {
		if nest == len(ss) {
			result = append(result, slices.Clone(nums))
			return
		}

		for _, e := range ss[nest] {
			f(nest+1, append(nums, e))
		}
	}
	f(0, make(S, 0, len(ss)))

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
