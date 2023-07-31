package fn

import (
	"github.com/sw965/omw"
	"golang.org/x/exp/constraints"
)

func Map[YS ~[]Y, XS ~[]X, X, Y any](xs XS, f func(X) Y) YS {
	return omw.MapFunc[XS, YS](xs, f)
}

func MapIndex[YS ~[]Y, XS ~[]X, X, Y any](xs XS, f func(int, X) Y) YS {
	y := make(YS, len(xs))
	for i, x := range xs {
		y[i] = f(i, x)
	}
	return y
}

func MapArg2[YS ~[]Y, XS1 ~[]X1, XS2 ~[]X2, X1, X2, Y any](xs1 XS1, xs2 XS2, f func(X1, X2) Y) YS {
	y := make(YS, len(xs1))
	for i, x1 := range xs1 {
		y[i] = f(x1, xs2[i])
	}
	return y
}

func MapNest2[YS ~[]Y, XS1 ~[]X1, XS2 ~[]X2, X1, X2, Y any](xs1 XS1, xs2 XS2, f func(X1, X2) Y) YS {
	y := make(YS, 0, len(xs1) * len(xs2))
	for _, x1 := range xs1 {
		for _, x2 := range xs2 {
			y = append(y, f(x1, x2))
		}
	}
	return y
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

func FilterIndex[XS ~[]X, X any](xs XS, f func(int, X) bool) XS {
	ys := make(XS, 0, len(xs))
	for i, x := range xs {
		if f(i, x) {
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
