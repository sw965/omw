package fn

import (
	"fmt"
	"github.com/sw965/omw"
)

func Map[YS ~[]Y, XS ~[]X, X, Y any](xs XS, f func(X) Y) YS {
	return omw.MapFunc[YS, XS](xs, f)
}

func Zip[YS ~[]Y, XS1 ~[]X1, XS2 ~[]X2, X1, X2, Y any](xs1 XS1, xs2 XS2, f func(X1, X2) Y) (YS, error) {
	if len(xs1) != len(xs2) {
		return YS{}, fmt.Errorf("スライスの長さが、一致しないので、Zip関数を適用出来ない")
	}

	ys := make(YS, len(xs1))
	for i := range xs1 {
		ys[i] = f(xs1[i], xs2[i])
	}
	return ys, nil
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

func Product2[YS ~[]Y, XS1 ~[]X1, XS2 ~[]X2, X1, X2, Y any](xs1 XS1, xs2 XS2, f func(X1, X2) Y) YS {
	ys := make(YS, 0, len(xs1) * len(xs2))
	for _, x1 := range xs1 {
		for _, x2 := range xs2 {
			y := f(x1, x2)
			ys = append(ys, y)
		}
	}
	return ys
}

func Product3[YS ~[]Y, XS1 ~[]X1, XS2 ~[]X2, XS3 ~[]X3, X1, X2, X3, Y any](xs1 XS1, xs2 XS2, xs3 XS3, f func(X1, X2, X3) Y) YS {
	ys := make(YS, 0, len(xs1) * len(xs2) * len(xs3))
	for _, x1 := range xs1 {
		for _, x2 := range xs2 {
			for _, x3 := range xs3 {
				y := f(x1, x2, x3)
				ys = append(ys, y)
			}
		}
	}
	return ys
}

func Product4[YS ~[]Y, XS1 ~[]X1, XS2 ~[]X2, XS3 ~[]X3, XS4 ~[]X4, X1, X2, X3, X4, Y any](xs1 XS1, xs2 XS2, xs3 XS3, xs4 XS4, f func(X1, X2, X3, X4) Y) YS {
	ys := make(YS, 0, len(xs1) * len(xs2) * len(xs3) * len(xs4))
	for _, x1 := range xs1 {
		for _, x2 := range xs2 {
			for _, x3 := range xs3 {
				for _, x4 := range xs4 {
					y := f(x1, x2, x3, x4)
					ys = append(ys, y)
				}
			}
		}
	}
	return ys
}

func Memo[M ~map[K]V, KS ~[]K, K comparable, V any](ks KS, f func(K)V) M {
	m := M{}
	for _, k := range ks {
		m[k] = f(k)
	}
	return m
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

func ToIntTilde[X, Y ~int](x X) Y {
	return Y(x)
}

func ToStrTilde[X, Y ~string](x X) Y {
	return Y(x)
}