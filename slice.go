package omw

import (
	"golang.org/x/exp/constraints"
)

func MakeSliceFunc[T any](n int, f func(int) T) []T {
	y := make([]T, n)
	for i := 0; i < n; i++ {
		y[i] = f(i)
	}
	return y
}

func MakeSliceRange[T constraints.Integer | constraints.Float](start, end, step T) []T {
	yn := int((end - start) / step)
	y := make([]T, yn)
	for i := 0; i < yn; i++ {
		y[i] = start + (step * T(i))
	}
	return y
}

func MapFunc[X, Y any](xs []X, f func(X) Y) []Y {
	y := make([]Y, len(xs))
	for i, x := range xs {
		y[i] = f(x)
	}
	return y
}

func Filter[X any](xs []X, f func(X) bool) []X {
	y := make([]X, 0, len(xs))
	for _, x := range xs {
		if f(x) {
			y = append(y, x)
		}
	}
	return y
}

func Reduce[X, Y any](xs []X, f func(X, Y) Y, init Y) Y {
	y := init
	for _, x := range xs {
		y = f(x, y)
	}
	return y
}

func Contains[T comparable](xs []T, v T) bool {
	for _, x := range xs {
		if x == v {
			return true
		}
	}
	return false
}

func AllContains[T comparable](xs, vs []T) bool {
	for _, v := range vs {
		if !Contains(xs, v) {
			return false
		}
	}
	return true
}

func Count[T comparable](xs []T, v T) int {
	y := 0
	for _, x := range xs {
		if x == v {
			y += 1
		}
	}
	return y
}

func Equals[T comparable](xs1, xs2 []T) bool {
	if len(xs1) != len(xs2) {
		return false
	}
	for i, x1 := range xs1 {
		x2 := xs2[i]
		if x1 != x2 {
			return false
		}
	}
	return true
}

func Sort[T any](xs []T, f func(int, int) bool) []T {
	n := len(xs)
	y := MapFunc(xs, Identity[T])
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if f(j, j+1) {
				y[j], y[j+1] = y[j+1], y[j]
			}
		}
	}
	return y
}

func Reverse[T any](xs []T) []T {
	yn := len(xs)
	y := make([]T, 0, yn)
	for i := yn - 1; i > -1; i-- {
		y = append(y, xs[i])
	}
	return y
}

func IsUnique[T comparable](xs []T) bool {
	for _, x := range xs {
		if Count(xs, x) != 1 {
			return false
		}
	}
	return true
}

func IndicesAccess[T any](xs []T, indices []int) []T {
	y := make([]T, len(indices))
	for i, idx := range indices {
		y[i] = xs[idx]
	}
	return y
}

func Index[T comparable](xs []T, v T) int {
	for i, x := range xs {
		if x == v {
			return i
		}
	}
	return -1
}

func Indices[T comparable](xs []T, v T) []int {
	y := make([]int, 0, len(xs))
	for i, x := range xs {
		if x == v {
			y = append(y, i)
		}
	}
	return y
}

func PointersToValues[T any](xs []*T) []T {
	y := make([]T, len(xs))
	for i, x := range xs {
		y[i] = *x
	}
	return y
}

func ValuesToPointers[T any](xs []T) []*T {
	y := make([]*T, len(xs))
	for i, x := range xs {
		y[i] = &x
	}
	return y
}

func Permutation[T any](xs []T, n, r int) [][]T {
	numss := PermutationNumberss(n, r)
	access := func(nums []int) []T { return IndicesAccess(xs, nums) }
	return MapFunc(numss, access)
}

func Combination[T any](xs []T, n, r int) [][]T {
	numss := CombinationNumberss(n, r)
	access := func(nums []int) []T { return IndicesAccess(xs, nums) }
	return MapFunc(numss, access)
}
