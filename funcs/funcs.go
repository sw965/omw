package funcs

import (
	"fmt"
)

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

func MapI[X, Y any](xs []X, f func(X, int) Y) []Y {
	ys := make([]Y, len(xs))
	for i, x := range xs {
		ys[i] = f(x, i)
	}
	return ys
}

func MapErrI[X, Y any](xs []X, f func(X, int) (Y, error)) ([]Y, error) {
	ys := make([]Y, len(xs))
	for i, x := range xs {
		y, err := f(x, i)
		if err != nil {
			return nil, err
		}
		ys[i] = y
	}
	return ys, nil
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

func FilterKeepIndices[X any](xs []X, f func(X) bool) []int {
	idxs := make([]int, 0, len(xs))
	for i, x := range xs {
		if f(x) {
			idxs = append(idxs, i)
		}
	}
	return idxs
}

func FilterErrKeepIndices[X any](xs []X, f func(X) (bool, error)) ([]int, error) {
	idxs := make([]int, 0, len(xs))
	for i, x := range xs {
		keep, err := f(x)
		if err != nil {
			return nil, err
		}
		if keep {
			idxs = append(idxs, i)
		}
	}
	return idxs, nil
}

func Fold[X, Y any](xs []X, init Y, f func(Y, X) Y) Y {
    acc := init
    for _, x := range xs {
        acc = f(acc, x)
    }
    return acc
}

func FoldErr[X, Y any](xs []X, init Y, f func(Y, X) (Y, error)) (Y, error) {
	var err error
    acc := init
    for _, x := range xs {
        acc, err = f(acc, x)
		if err != nil {
			var y Y
			return y, err
		}
    }
    return acc, nil
}

func Scan[X, Y any](xs []X, init Y, f func(Y, X) Y) []Y {
    ys := make([]Y, 0, len(xs)+1)
    acc := init
    ys = append(ys, acc)
    for _, x := range xs {
        acc = f(acc, x)
        ys = append(ys, acc)
    }
    return ys
}

func ScanErr[X, Y any](xs []X, init Y, f func(Y, X) (Y, error)) ([]Y, error) {
    ys := make([]Y, 0, len(xs)+1)
    acc := init
    ys = append(ys, acc)
    for _, x := range xs {
        var err error
        acc, err = f(acc, x)
        if err != nil {
            return nil, err
        }
        ys = append(ys, acc)
    }
    return ys, nil
}

func ZipWith[X1, X2, Y any](xs1 []X1, xs2 []X2, f func(X1, X2) Y) ([]Y, error) {
    n := len(xs1)
	if n != len(xs2) {
		return nil, fmt.Errorf("len(xs1) != len(xs2)")
	}

    z := make([]Y, n)
    for i := 0; i < n; i++ {
        z[i] = f(xs1[i], xs2[i])
    }
    return z, nil
}

func ZipWith3[X1, X2, X3, Y any](xs1 []X1, xs2 []X2, xs3 []X3, f func(X1, X2, X3) Y) ([]Y, error) {
	n := len(xs1)
	if n != len(xs2) || n != len(xs3) {
		return nil, fmt.Errorf("n != len(xs2) || n != len(xs3)")
	}

	z := make([]Y, n)
	for i := 0; i < n; i++ {
		z[i] = f(xs1[i], xs2[i], xs3[i])
	}
	return z, nil
}

func Juxt[Y, X any](x X, fs []func(X) Y) []Y {
    ys := make([]Y, len(fs))
    for i, f := range fs {
        ys[i] = f(x)
    }
    return ys
}

func Identity[X any](x X) X {
	return x
}