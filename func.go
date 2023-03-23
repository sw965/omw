package omw

import (
	"golang.org/x/exp/constraints"
)

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

func IsMultipleOf[X constraints.Integer](a X) func(X) bool {
	return func(x X) bool { return x%a == 0 }
}