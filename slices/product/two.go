package product

func MapFunc2[X1, X2, Y any](xs1 []X1, xs2 []X2, f func(X1, X2) Y) []Y {
	ys := make([]Y, 0, len(xs1) * len(xs2))
	for _, x1 := range xs1 {
		for _, x2 := range xs2 {
			ys = append(ys, f(x1, x2))
		}
	}
	return ys
}

func NewMap2[NM ~map[K1]M, M ~map[K2]V, K1, K2 comparable, V any](ks1 []K1, ks2 []K2, f func(K1, K2) V) NM {
	nm := NM{}
	for _, k1 := range ks1 {
		for _, k2 := range ks2 {
			if _, ok := nm[k1]; !ok {
				nm[k1] = map[K2]V{}
			}
			nm[k1][k2] = f(k1, k2)
		}
	}
	return nm
}