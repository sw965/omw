package fn

import (
	"github.com/sw965/omw"
)

func Map[YS ~[]Y, XS ~[]X, X, Y any](xs XS, f func(X) Y) YS {
	return omw.MapFunc[XS, YS](xs, f)
}

func MapIndex[YS ~[]Y, XS ~[]X, X, Y any](xs XS, f func(int, X) Y) YS {
	ys := make(YS, len(xs))
	for i, x := range xs {
		ys[i] = f(i, x)
	}
	return ys
}

func MapError[YS ~[]Y, XS ~[]X, X, Y any](xs XS, f func(X) (Y, error)) (YS, error) {
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

func ToStrTilde[X, Y ~string](x X) Y {
	return Y(x)
}

func ToIntTilde[X, Y ~int](x X) Y {
	return Y(x)
}
