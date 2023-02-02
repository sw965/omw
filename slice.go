package omw

import (
	"golang.org/x/exp/constraints"
)

func MakeSliceFunc[T any](length int, f func(int) T) []T {
	y := make([]T, length)
	for i := 0; i < length; i++ {
		y[i] = f(i)
	}
	return y
}

func MakeSliceRange[T constraints.Integer | constraints.Float](start, end, step T) []T {
	yLen := int((end - start) / step)
	y := make([]T, yLen)
	for i := 0; i < yLen; i++ {
		y[i] = start + (step * T(i))
	}
	return y
}

func MapFunc[X, Y any](xs []X, f func(X) Y) []Y {
	ys := make([]Y, len(xs))
	for i, x := range xs {
		ys[i] = f(x)
	}
	return ys
}

func Filter[X any](xs []X, f func(X) bool) []X {
	ys := make([]X, len(xs))
	for _, x := range xs {
		if f(x) {
			ys = append(ys, x)
		}
	}
	return ys
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

func Sort[T any](xs []T, f func(int, int) bool) []T {
	y := MapFunc(xs, Identity[T])
	yLen := len(y)

	for i := 0; i < yLen-1; i++ {
		for j := 0; j < yLen-i-1; j++ {
			if f(j, j+1) {
				y[j], y[j+1] = y[j+1], y[j]
			}
		}
	}
	return y
}

func Reverse[T any](xs []T) []T {
	yLen := len(xs)
	ys := make([]T, 0, yLen)
	for i := yLen - 1; i > -1; i-- {
		ys = append(ys, xs[i])
	}
	return ys
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
	for i, index := range indices {
		y[i] = xs[index]
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
