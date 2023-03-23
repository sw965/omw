package product

func MapFunc4[X1, X2, X3, X4](xs1 []X1, xs2 []X2, xs3 []X3, xs4 []X4, f func(X1, X2, X3, X4) Y) []Y {
	n := len(xs1) * len(xs2) * len(xs3) * len(xs4)
	ys := make([]Y, 0, n)
	for _, x1 := range xs1 {
		for _, x2 := range xs2 {
			for _, x3 := range xs3 {
				for _, x4 := range xs4 {
					ys = append(ys, f(x1, x2, x3, x4))
				}
			}
		}
	}
	return ys
}