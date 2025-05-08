package fn

import (
	"golang.org/x/exp/constraints"
)

func Tabulate[YS ~[]Y, Y any](n int, f func(int) Y) YS {
	ys := make(YS, n)
	for i := 0; i < n; i++ {
		ys[i] = f(i)
	}
	return ys
}

func Map[YS ~[]Y, XS ~[]X, X, Y any](xs XS, f func(X) Y) YS {
	ys := make(YS, len(xs))
	for i, x := range xs {
		ys[i] = f(x)
	}
	return ys
}

func MapWithError[YS ~[]Y, XS ~[]X, X, Y any](xs XS, f func(X) (Y, error)) (YS, error) {
	ys := make(YS, len(xs))
	for i, x := range xs {
		y, err := f(x)
		if err != nil {
			return ys, err
		}
		ys[i] = y
	}
	return ys, nil
}

func Filter[XS ~[]X, X any](xs XS, f func(X) bool) XS {
	ys := make(XS, 0, len(xs))
	for _, x := range xs {
		if f(x) {
			ys = append(ys, x)
		}
	}
	return ys
}

func Accumulate[XS ~[]X, X constraints.Ordered](xs XS) XS {
	a := make(XS, len(xs))
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
