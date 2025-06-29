package slices

import (
	"slices"
	"golang.org/x/exp/constraints"
	omath "github.com/sw965/omw/math"
)

func MakeInteger[S ~[]I, I constraints.Integer](start, end I) S {
	n := end - start
	s := make(S, int(n))
	for i := I(0); i < n; i++ {
		s[i] = start + i
	}
	return s
}

func Reversed[S ~[]E, E any](s S) S {
	s = slices.Clone(s)
	slices.Reverse(s)
	return s
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

func CountFunc[S ~[]E, E any](s S, f func(E) bool) int {
	c := 0
	for _, e := range s {
		if f(e) {
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
    slices.SortFunc(idxs, func(i, j int) int {
        switch {
        case s[i] < s[j]:
            return -1
        case s[i] > s[j]:
            return 1
        default:
            return 0
        }
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

func IntPermutations(n, r int) [][]int {
	c := omath.PermutationsCount(n, r)
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

func Permutations[S ~[]E, E any](s S, r int) []S {
	n := len(s)
	idxss := IntPermutations(n, r)
	ss := make([]S, len(idxss))
	for i, idxs := range idxss {
		ss[i] = ElementsByIndices(s, idxs...)
	}
	return ss
}

// 重複順列
func IntSequences(n, r int) [][]int {
	c := omath.SequencesCount(n, r)
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

func Sequences[S ~[]E, E any](s S, r int) []S {
	n := len(s)
	idxss := IntSequences(n, r)
	ss := make([]S, len(idxss))
	for i, idxs := range idxss {
		ss[i] = ElementsByIndices(s, idxs...)
	}
	return ss
}

func IntCombinations(n, r int) [][]int {
	nums := make([]int, r)
	for i := 0; i < r; i++ {
		nums[i] = i
	}

	c := omath.CombinationsCount(n, r)
	result := make([][]int, 0, c)
	if r == 0 {
		return result
	}

	end := r - 1
	for i := 0; i < c; i++ {
		clone := slices.Clone(nums)
		result = append(result, clone)
		max := omath.Max(nums...)
		if max == (n - 1) {
			reversed := slices.Clone(nums)
			slices.Reverse(reversed)
			count := omath.CountConsecutiveDecrease(reversed...)
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

func Combinations[S ~[]E, E any](s S, r int) []S {
	n := len(s)
	idxss := IntCombinations(n, r)
	ss := make([]S, len(idxss))
	for i, idxs := range idxss {
		ss[i] = ElementsByIndices(s, idxs...)
	}
	return ss
}

func CartesianProducts[S ~[]E, E any](ss ...S) []S {
	if len(ss) == 0 {
		return nil
	}

	c := 1
	for _, s := range ss {
		c *= len(s)
	}
	result := make([]S, 0, c)
	
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

func ToUnique[S ~[]E, E comparable](s S) S {
	u := make(S, 0, len(s))
	for _, e := range s {
		if Count(u, e) == 0 {
			u = append(u, e)
		}
	}
	return u
}

func ToUniqueFunc[S ~[]E, E comparable](s S, f func(E) bool) S {
	u := make(S, 0, len(s))
	for _, e := range s {
		if CountFunc(s, f) == 0 {
			u = append(u, e)
		}
	}
	return u
}