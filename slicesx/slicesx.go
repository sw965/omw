/*
Package slicesx provides generic utility functions for slices, leveraging Go 1.23+ iterators.
It primarily focuses on combinatorial algorithms such as permutations, combinations, and Cartesian products, designed to be memory efficient by yielding elements sequentially.

Package slicesx は、Go 1.23+ のイテレータを活用したスライスのためのジェネリックなユーティリティ関数を提供します。
主に、順列、組み合わせ、直積などの組み合わせアルゴリズムに焦点を当てており、要素を順次生成（yield）することでメモリ効率良くなるように設計されています。

# Usage

Most functions in this package return an `iter.Seq[S]`, allowing them to be used directly in `for-range` loops.

このパッケージの多くの関数は `iter.Seq[S]` を返すため、`for-range` ループで直接使用することができます。

    s := []int{1, 2, 3}
    // Generate permutations / 順列を生成
    for p := range slicesx.Permutations(s, 2) {
        fmt.Println(p)
    }

# Combinatorial Functions / 組み合わせ関数

The package supports the following operations:
このパッケージは以下の操作をサポートしています:

  - Permutations: Generates all permutations of length r. (nPr)
    順列: 長さ r のすべての順列を生成します。

  - Sequences: Generates all permutations with repetition of length r. (n^r)
    重複順列: 長さ r のすべての重複順列を生成します。

  - Combinations: Generates all combinations of length r. (nCr)
    組み合わせ: 長さ r のすべての組み合わせを生成します。

  - CartesianProducts: Generates the Cartesian product of multiple slices.
    直積: 複数のスライスの直積を生成します。

  - Argsort: Returns the indices that would sort the slice.
    Argsort: スライスをソートするインデックスを返します。
*/
package slicesx

import (
	"cmp"
	"iter"
	"slices"
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

// Argsort returns the indices that would sort the slice in ascending order.
// The original slice is not modified.
//
// スライスを昇順にソートした場合のインデックスの並びを返します。
// 元のスライスは変更されません。
func Argsort[S ~[]E, E cmp.Ordered](s S) []int {
	return ArgsortFunc(s, func(a, b E) int {
		return cmp.Compare(a, b)
	})
}

// ArgsortFunc returns the indices that would sort the slice using the provided comparison function.
// The original slice is not modified.
//
// The function f must return a negative number when a < b, a positive number when a > b,
// and zero when a == b.
//
// 提供された比較関数を使用してスライスをソートした場合のインデックスの並びを返します。
// 元のスライスは変更されません。
//
// 比較関数 f は、a < b の場合に負の値、a > b の場合に正の値、
// a == b の場合に 0 を返す必要があります（cmp.Compare と同様の仕様です）。
func ArgsortFunc[S ~[]E, E any](s S, f func(a, b E) int) []int {
    idxs := make([]int, len(s))
    for i := range idxs {
        idxs[i] = i
    }
    slices.SortFunc(idxs, func(i, j int) int {
        return f(s[i], s[j])
    })
    return idxs
}