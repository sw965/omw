package fn

import (
	"golang.org/x/exp/constraints"
	"github.com/sw965/omw"
)

func Map[XS ~[]X, YS ~[]Y, X, Y any](xs XS, f func(X) Y) YS {
	return omw.MapFunc[XS, YS](xs, f)
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

func Any[XS ~[]X, X any](xs XS, f func(X) bool) bool {
	for _, x := range xs {
		if f(x) {
			return true
		}
	}
	return false
}

func Identity[X any](x X) X {
	return x
}

func IsRemainderZero[X constraints.Integer](x X) func(X) bool {
	return func(a X) bool { return a%x == 0 }
}

func ToStrTilde[X, Y ~string](x X) Y {
	return Y(x)
}