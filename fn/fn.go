package fn

import (
	"golang.org/x/exp/constraints"
)

func Map[XS ~[]X, YS ~[]Y, X, Y any](xs XS, f func(X) Y) YS {
	ys := make(YS, len(xs))
	for i, x := range xs {
		ys[i] = f(x)
	}
	return ys
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

func Identity[X any](x X) X {
	return x
}

func IsMultipleOf[X constraints.Integer](a X) func(X) bool {
	return func(x X) bool { return x%a == 0 }
}