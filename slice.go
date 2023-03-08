package omw

import (
	"golang.org/x/exp/constraints"
)

func MakeSliceFunc[Y any](n int, f func(int) Y) []Y {
	ys := make([]Y, n)
	for i := 0; i < n; i++ {
		ys[i] = f(i)
	}
	return ys
}

func MakeSliceRange[Y constraints.Integer | constraints.Float](start, end, step Y) []Y {
	n := int((end - start) / step)
	ys := make([]Y, n)
	for i := 0; i < n; i++ {
		ys[i] = start + (step * Y(i))
	}
	return ys
}

func MapFunc[X, Y any](xs []X, f func(X) Y) []Y {
	ys := make([]Y, len(xs))
	for i, x := range xs {
		ys[i] = f(x)
	}
	return ys
}

func Filter[X any](xs []X, f func(X) bool) []X {
	ys := make([]X, 0, len(xs))
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

func Contains[X comparable](xs []X, a X) bool {
	for _, x := range xs {
		if x == a {
			return true
		}
	}
	return false
}

func AllContains[X comparable](xs, as []X) bool {
	for _, a := range as {
		if !Contains(xs, a) {
			return false
		}
	}
	return true
}

func Count[X comparable](xs []X, a X) int {
	y := 0
	for _, x := range xs {
		if x == a {
			y += 1
		}
	}
	return y
}

func Equals[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func Sort[X any](xs []X, f func(int, int) bool) []X {
	n := len(xs)
	ys := MapFunc(xs, Identity[X])
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if f(j, j+1) {
				ys[j], ys[j+1] = ys[j+1], ys[j]
			}
		}
	}
	return ys
}

func Reverse[X any](xs []X) []X {
	n := len(xs)
	ys := make([]X, 0, n)
	for i := n - 1; i > -1; i-- {
		ys = append(ys, xs[i])
	}
	return ys
}

func IsUnique[X comparable](xs []X) bool {
	for _, x := range xs {
		if Count(xs, x) != 1 {
			return false
		}
	}
	return true
}

func IndicesAccess[X any](xs []X, indices []int) []X {
	ys := make([]X, len(indices))
	for i, idx := range indices {
		ys[i] = xs[idx]
	}
	return ys
}

func Index[X comparable](xs []X, a X) int {
	for i, x := range xs {
		if x == a {
			return i
		}
	}
	return -1
}

func Indices[X comparable](xs []X, a X) []int {
	ys := make([]int, 0, len(xs))
	for i, x := range xs {
		if x == a {
			ys = append(ys, i)
		}
	}
	return ys
}

func PointersToValues[X any](xs []*X) []X {
	ys := make([]X, len(xs))
	for i, x := range xs {
		ys[i] = *x
	}
	return ys
}

func ValuesToPointers[X any](xs []X) []*X {
	ys := make([]*X, len(xs))
	for i, x := range xs {
		ys[i] = &x
	}
	return ys
}

func Permutation[X any](xs []X, r int) [][]X {
	n := len(xs)
	numss := PermutationNumberss(n, r)
	access := func(nums []int) []X { return IndicesAccess(xs, nums) }
	return MapFunc(numss, access)
}

func Combination[X any](xs []X, r int) [][]X {
	n := len(xs)
	numss := CombinationNumberss(n, r)
	access := func(nums []int) []X { return IndicesAccess(xs, nums) }
	return MapFunc(numss, access)
}