package slices

import (
	"golang.org/x/exp/slices"
	"golang.org/x/exp/constraints"
)

func Get[E any](es []E, indices ...int) []E {
	y := make([]E, len(indices))
	for i, idx := range indices {
		y[i] = es[idx]
	}
	return y
}

func Make[E any](n int, f func(int) E) []E {
	y := make([]E, n)
	for i := 0; i < n; i++ {
		y[i] = f(i)
	}
	return y
}

func Range[E constraints.Integer | constraints.Float](start, end, step E) []E {
	n := int((end - start) / step)
	y := make([]E, n)
	for i := 0; i < n; i++ {
		y[i] = start + (step * E(i))
	}
	return y
}

func Count[E comparable](es []E, a E) int {
	y := 0
	for _, e := range es {
		if e == a {
			y += 1
		}
	}
	return y
}

func Reverse[E any](es []E) []E {
	n := len(es)
	y := make([]E, 0, n)
	for i := n - 1; i > -1; i-- {
		y = append(y, es[i])
	}
	return y
}

func AllContains[E comparable](es, as []E) bool {
	for _, a := range as {
		if !slices.Contains(es, a) {
			return false
		}
	}
	return true
}

func Indices[E comparable](es []E, a E) []int {
	y := make([]int, 0, len(es))
	for i, e := range es {
		if e == a {
			y = append(y, i)
		}
	}
	return y
}

func ToUnique[E comparable](es []E) []E {
	y := make([]E, 0, len(es))
	for _, e := range es {
		if !slices.Contains(es, e) {
			y = append(y, e)
		}
	}
	return y
}

func IsUnique[E comparable](es []E) bool {
	for _, e := range es {
		if Count(es, e) != 1 {
			return false
		}
	}
	return true
}