package slicesx

import (
	"iter"
)

// Permutations returns a sequence of all permutations of r elements from s.
//
// s から r 個の要素を選ぶ順列のシーケンスを返します。
func Permutations[S ~[]E, E any](s S, r int) iter.Seq[S] {
	return func(yield func(S) bool) {
		n := len(s)

		if r < 0 {
			return
		}

		if r == 0 {
			if !yield(make(S, 0)) {
				return
			}
			return
		}

		if r > n {
			return
		}

		indices := make([]int, n)
		for i := range indices {
			indices[i] = i
		}
		cycles := make([]int, r)
		for i := 0; i < r; i++ {
			cycles[i] = n - i
		}

		emit := func() bool {
			out := make(S, r)
			for i := 0; i < r; i++ {
				out[i] = s[indices[i]]
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
					tmp := indices[i]
					copy(indices[i:], indices[i+1:])
					indices[n-1] = tmp
					cycles[i] = n - i
					if i == 0 {
						return
					}
					continue
				}

				j := n - cycles[i]
				indices[i], indices[j] = indices[j], indices[i]

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

// Sequences returns a sequence of all tuples of length r from s (permutations with repetition).
//
// s から長さ r の要素の列（重複順列）をすべて返します。
func Sequences[S ~[]E, E any](s S, r int) iter.Seq[S] {
	return func(yield func(S) bool) {
		n := len(s)

		if r < 0 {
			return
		}

		if r == 0 {
			if !yield(make(S, 0)) {
				return
			}
			return
		}

		if n == 0 {
			return
		}

		idx := make([]int, r)

		for {
			out := make(S, r)
			for i := 0; i < r; i++ {
				out[i] = s[idx[i]]
			}
			if !yield(out) {
				return
			}

			k := r - 1
			for ; k >= 0; k-- {
				idx[k]++
				if idx[k] < n {
					break
				}
				idx[k] = 0
			}
			if k < 0 {
				return
			}
		}
	}
}

// Combinations returns a sequence of all combinations of r elements from s.
//
// s から r 個の要素を選ぶ組合せのシーケンスを返します。
func Combinations[S ~[]E, E any](s S, r int) iter.Seq[S] {
	return func(yield func(S) bool) {
		n := len(s)

		if r < 0 {
			return
		}

		if r == 0 {
			if !yield(make(S, 0)) {
				return
			}
			return
		}
		if r > n {
			return
		}

		idx := make([]int, r)
		for i := 0; i < r; i++ {
			idx[i] = i
		}

		for {
			out := make(S, r)
			for i := 0; i < r; i++ {
				out[i] = s[idx[i]]
			}
			if !yield(out) {
				return
			}

			i := r - 1
			for ; i >= 0; i-- {
				if idx[i] != i+n-r {
					break
				}
			}
			if i < 0 {
				return
			}
			idx[i]++
			for j := i + 1; j < r; j++ {
				idx[j] = idx[j-1] + 1
			}
		}
	}
}

// CartesianProducts returns the Cartesian product of the input slices.
//
// 入力された複数のスライスの直積を返します。
func CartesianProducts[S ~[]E, E any](ss ...S) iter.Seq[S] {
	return func(yield func(S) bool) {
		k := len(ss)

		if k == 0 {
			if !yield(make(S, 0)) {
				return
			}
			return
		}

		for _, s := range ss {
			if len(s) == 0 {
				return
			}
		}

		idx := make([]int, k)

		for {
			out := make(S, k)
			for i := 0; i < k; i++ {
				out[i] = ss[i][idx[i]]
			}
			if !yield(out) {
				return
			}

			p := k - 1
			for ; p >= 0; p-- {
				idx[p]++
				if idx[p] < len(ss[p]) {
					break
				}
				idx[p] = 0
			}
			if p < 0 {
				return
			}
		}
	}
}

// Counts returns a map containing the counts of each unique element in the slice.
// It iterates through the provided slice and increments the counter for each element.
//
// Counts は、スライス内の各ユニークな要素の出現回数を含むマップを返します。
// 提供されたスライスを反復処理し、各要素のカウンタをインクリメントします。
func Counts[S ~[]E, E comparable](s S) map[E]int {
	c := make(map[E]int, len(s))
	for _, e := range s {
		c[e]++
	}
	return c
}
