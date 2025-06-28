package funcs

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

func MapErr[X, Y any](xs []X, f func(X) (Y, error)) ([]Y, error) {
	ys := make([]Y, len(xs))
	for i, x := range xs {
		y, err := f(x)
		if err != nil {
			return nil, err
		}
		ys[i] = y
	}
	return ys, nil
}

func MapMemo[K comparable, V any](s []K, f func(K)V) M {
    m := M{}
    for _, k := range s {
        m[k] = f(k)
    }
    return m
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

func FilterErr[X any](xs []X, f func(X) (bool, error)) ([]X, error) {
	ys := make([]X, 0, len(xs))
	for _, x := range xs {
		keep, err := f(x)
		if err != nil {
			return nil, err
		}
		if keep {
			ys = append(ys, x)
		}
	}
	return ys, nil
}

func Accumulate[Y, X any](xs []X, f func(...X) Y) []Y {
	n := len(xs)
	ys := make([]Y, n)
	current := make([]X, 0, n)
	for i, x := range xs {
		current = append(current, x)
		ys[i] = f(current...)
	}
	return ys
}

func AccumulateErr[Y, X any](xs []X, f func(...X) (Y, error)) ([]Y, error) {
	n := len(xs)
	ys := make([]Y, n)
	current := make([]X, 0, n)
	for i, x := range xs {
		current = append(current, x)
		y, err := f(current...)
		if err != nil {
			return nil, err
		}
		ys[i] = y
	}
	return ys, nil
}

func Identity[X any](x X) X {
	return x
}
