package slicesx

import (
	"cmp"
	"fmt"
	"iter"
	"slices"
)

// 順列
func Permutations[S ~[]E, E any](s S, r int) iter.Seq[S] {
	return func(yield func(S) bool) {
		n := len(s)

		if r < 0 {
			return
		}

		if r == 0 {
			yield(make(S, 0))
			return
		}

		if r > n {
			return
		}

		idxs := make([]int, n)
		for i := range idxs {
			idxs[i] = i
		}
		cycles := make([]int, r)
		for i := range r {
			cycles[i] = n - i
		}

		emit := func() bool {
			out := make(S, r)
			for i := range r {
				out[i] = s[idxs[i]]
			}
			return yield(out)
		}

		if !emit() {
			return
		}

		for {
			advanced := false
			for i := r - 1; i >= 0; i-- {
				cycles[i]--
				if cycles[i] == 0 {
					tmp := idxs[i]
					copy(idxs[i:], idxs[i+1:])
					idxs[n-1] = tmp
					cycles[i] = n - i
					if i == 0 {
						return
					}
					continue
				}

				j := n - cycles[i]
				idxs[i], idxs[j] = idxs[j], idxs[i]

				if !emit() {
					return
				}
				advanced = true
				break
			}
			if !advanced {
				return
			}
		}
	}
}

// 重複順列
func Sequences[S ~[]E, E any](s S, r int) iter.Seq[S] {
	return func(yield func(S) bool) {
		n := len(s)

		if r < 0 {
			return
		}

		if r == 0 {
			yield(make(S, 0))
			return
		}

		if n == 0 {
			return
		}

		idxs := make([]int, r)

		for {
			out := make(S, r)
			for i := range r {
				out[i] = s[idxs[i]]
			}
			if !yield(out) {
				return
			}

			k := r - 1
			for ; k >= 0; k-- {
				idxs[k]++
				if idxs[k] < n {
					break
				}
				idxs[k] = 0
			}
			if k < 0 {
				return
			}
		}
	}
}

// 組合せ
func Combinations[S ~[]E, E any](s S, r int) iter.Seq[S] {
	return func(yield func(S) bool) {
		n := len(s)

		if r < 0 {
			return
		}

		if r == 0 {
			yield(make(S, 0))
			return
		}
		if r > n {
			return
		}

		idxs := make([]int, r)
		for i := range r {
			idxs[i] = i
		}

		for {
			out := make(S, r)
			for i := range r {
				out[i] = s[idxs[i]]
			}
			if !yield(out) {
				return
			}

			i := r - 1
			for ; i >= 0; i-- {
				if idxs[i] != i+n-r {
					break
				}
			}
			if i < 0 {
				return
			}
			idxs[i]++
			for j := i + 1; j < r; j++ {
				idxs[j] = idxs[j-1] + 1
			}
		}
	}
}

// 直積
func CartesianProducts[S ~[]E, E any](ss ...S) iter.Seq[S] {
	return func(yield func(S) bool) {
		k := len(ss)

		if k == 0 {
			yield(make(S, 0))
			return
		}

		for _, s := range ss {
			if len(s) == 0 {
				return
			}
		}

		idxs := make([]int, k)

		for {
			out := make(S, k)
			for i := range k {
				out[i] = ss[i][idxs[i]]
			}
			if !yield(out) {
				return
			}

			p := k - 1
			for ; p >= 0; p-- {
				idxs[p]++
				if idxs[p] < len(ss[p]) {
					break
				}
				idxs[p] = 0
			}
			if p < 0 {
				return
			}
		}
	}
}

func Counts[S ~[]E, E comparable](s S) map[E]int {
	c := make(map[E]int, len(s))
	for _, e := range s {
		c[e]++
	}
	return c
}

func Argsort[S ~[]E, E cmp.Ordered](s S) []int {
	return ArgsortFunc(s, cmp.Compare)
}

func ArgsortFunc[S ~[]E, E any](s S, f func(a, b E) int) []int {
	idxs := make([]int, len(s))
	for i := range idxs {
		idxs[i] = i
	}
	slices.SortStableFunc(idxs, func(i, j int) int {
		return f(s[i], s[j])
	})
	return idxs
}

func IsUnique[S ~[]E, E comparable](s S) bool {
	if len(s) <= 1 {
		return true
	}
	seen := make(map[E]struct{}, len(s))
	for _, e := range s {
		if _, ok := seen[e]; ok {
			return false
		}
		seen[e] = struct{}{}
	}
	return true
}

func ElementsByIndices[S ~[]E, E any](s S, idxs ...int) (S, error) {
	result := make(S, len(idxs))
	n := len(s)

	for i, idx := range idxs {
		if idx < 0 || idx >= n {
			return nil, fmt.Errorf("0 <= idx < %d であるべき: idx = %d", n, idx)
		}
		result[i] = s[idx]
	}
	return result, nil
}
