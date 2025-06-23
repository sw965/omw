package fn

import (
	"golang.org/x/exp/constraints"
)

func Tabulate[Y any](n int, f func(int) Y) []Y {
	ys := make([]Y, n)
	for i := 0; i < n; i++ {
		ys[i] = f(i)
	}
	return ys
}

func Map[X, Y any](xs []X, f func(X) Y) []Y {
	ys := make([]Y, len(xs))
	for i, x := range xs {
		ys[i] = f(x)
	}
	return ys
}

func MapWithError[X, Y any](xs []X, f func(X) (Y, error)) ([]Y, error) {
	ys := make([]Y, len(xs))
	for i, x := range xs {
		y, err := f(x)
		if err != nil {
			return ys, err
		}
		ys[i] = y
	}
	return ys, nil
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

func Accumulate[X constraints.Ordered](xs []X) []X {
	a := make([]X, len(xs))
	var current X
	for i, x := range xs {
		current += x
		a[i] = current
	}
	return a
}

func Identity[X any](x X) X {
	return x
}
