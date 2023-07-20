package fn

import (
	"github.com/sw965/omw"
	"golang.org/x/exp/constraints"
)

func Map[YS ~[]Y, XS ~[]X, X, Y any](xs XS, f func(X) Y) YS {
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

func All[XS ~[]X, X any](xs XS, f func(X) bool) bool {
	for _, x := range xs {
		if !f(x) {
			return false
		}
	}
	return true
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

func IdentityWithNilError[X any](x X) (X, error) {
	return x, nil
}

func IsRemainderZero[X constraints.Integer](x X) func(X) bool {
	return func(a X) bool { return a%x == 0 }
}

func ToStrTilde[X, Y ~string](x X) Y {
	return Y(x)
}
