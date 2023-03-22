package product

func MapFunc3[X1, X2, X3, Y any](xs1 []X1, xs2 []X2, xs3 []X3, f func(X1, X2, X3) Y) []Y {
	ys := make([]Y, 0, len(xs1) * len(xs2) * len(xs3))
	for _, x1 := range xs1 {
		for _, x2 := range xs2 {
			for _, x3 := range xs3 {
				ys = append(ys, f(x1, x2, x3))
			}
		}
	}
	return ys
}