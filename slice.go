package omw

import (
	"golang.org/x/exp/constraints"
)

func MapFunc[X, Y any](xs []X, f func(X) Y) []Y {
	result := make([]Y, len(xs))
	for i, x := range xs {
		result[i] = f(x)
	}
	return result
}

func Filter[X any](xs []X, f func(X) bool) []X {
	result := make([]X, 0, len(xs))
	for _, x := range xs {
		if f(x) {
			result = append(result, x)
		}
	}
	return result
}

func Reduce[X, Y any](xs []X, f func(Y, X) Y, init Y) Y {
	result := init
	for _, x := range xs {
		result = f(result, x)
	}
	return result
}

func Equals[T comparable](xs1, xs2 []T) bool {
	for i, x1 := range xs1 {
		x2 := xs2[i]
		if x1 != x2 {
			return false
		}
	}
	return true
}

func Contains[T comparable](xs []T, x T) bool {
	for _, ele := range xs {
		if ele == x {
			return true
		}
	}
	return false
}

func AllContains[T comparable](xs, eles []T) bool {
	for _, ele := range eles {
		if !Contains(xs, ele) {
			return false
		}
	}
	return true
}

func Count[T comparable](xs []T, x T) int {
	result := 0
	for _, ele := range xs {
		if ele == x {
			result += 1
		}
	}
	return result
}

func Indices[T comparable](xs []T, x T) []int {
	result := make([]int, 0, len(xs))
	for i, ele := range xs {
		if ele == x {
			result = append(result, i)
		}
	}
	return result
}

func IndicesFunc[T any](xs []T, f func(T) bool) []int {
	result := make([]int, 0, len(xs))
	for i, x := range xs {
		if f(x) {
			result = append(result, i)
		}
	}
	return result
}

func IndicesAccess[T any](xs []T, indices []int) []T {
	result := make([]T, len(indices))
	for i, index := range indices {
		result[i] = xs[index]
	}
	return result
}

func MinIndices[T constraints.Ordered](xs []T) []int {
	min := Min(xs...)
	result := make([]int, 0, len(xs))
	for i, x := range xs {
		if x == min {
			result = append(result, i)
		}
	}
	return result
}

func MaxIndices[T constraints.Ordered](xs []T) []int {
	max := Max(xs...)
	result := make([]int, 0, len(xs))
	for i, x := range xs {
		if x == max {
			result = append(result, i)
		}
	}
	return result
}

func Sort[T any](xs []T, f func(int, int) bool) []T {
	length := len(xs)
	copyXs := make([]T, length)
	for i, x := range xs {
		copyXs[i] = x
	}

	for i := 0; i < length - 1; i++ {
		for j := 0; j < length - i - 1; j++ {
			if f(j, j + 1) {
				copyXs[j], copyXs[j + 1] = copyXs[j + 1], copyXs[j]
			}
		}
	}
	return copyXs
}

func DescSort[T constraints.Ordered](xs []T, f func(int, int) bool) []T {
	return Sort(xs, func(i, j int) bool { return xs[i] > xs[j] })
}

func AscSort[T constraints.Ordered](xs []T, f func(int, int) bool) []T {
	return Sort(xs, func(i, j int) bool { return xs[i] < xs[j] })
}

func All(xs []bool) bool {
	for _, x := range xs {
		if !x {
			return false
		}
	}
	return true
}

func Any(xs []bool) bool {
	for _, x := range xs {
		if x {
			return true
		}
	}
	return false
}

func Reverse[T any](xs []T) []T {
	length := len(xs)
	result := make([]T, 0, len(xs))
	for i := length - 1; i > -1; i-- {
		result = append(result, xs[i])
	}
	return result
}

func IsUnique[T comparable](xs []T) bool {
	for _, x := range xs {
		if Count(xs, x) != 1 {
			return false
		}
	}
	return true
}

func PointersToValues[T any](xs []*T) []T {
	result := make([]T, len(xs))
	for i, x := range xs {
		result[i] = *x
	}
	return result
}

func Permutation[T any](xs []T, n, r int) ([][]T, error) {
	totalNum, err := PermutationTotalNum(n, r)
	if err != nil {
		return [][]T{}, err
	}
	numberss, _ := PermutationNumberss(n, r)
	result := make([][]T, totalNum)
	for i, numbers := range numberss {
		result[i] = IndicesAccess(xs, numbers)
	}
	return result, nil
}

func Combination[T any](xs []T, n, r int) ([][]T, error) {
	totalNum, err := CombinationTotalNum(n, r)
	if err != nil {
		return [][]T{}, err
	}
	numberss, _ := CombinationNumberss(n, r)
	result := make([][]T, totalNum)
	for i, numbers := range numberss {
		result[i] = IndicesAccess(xs, numbers)
	}
	return result, nil
}