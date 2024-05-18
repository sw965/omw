package fn

func Map[YS ~[]Y, XS ~[]X, X, Y any](xs XS, f func(X) Y) YS {
	ys := make(YS, len(xs))
	for i, x := range xs {
		ys[i] = f(x)
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

func Product[YS ~[]Y, Y any, XS1 ~[]X1, XS2 ~[]X2, X1, X2 any](xs1 XS1, xs2 XS2, f func(X1, X2) Y) YS {
	ys := make(YS, 0, len(xs1) * len(xs2))
	for _, x1 := range xs1 {
		for _, x2 := range xs2 {
			ys = append(ys, f(x1, x2))
		}
	}
	return ys
}

func Identity[X any](x X) X {
	return x
}