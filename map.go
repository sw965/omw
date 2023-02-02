package omw

func Keys[K comparable, V any](x map[K]V) []K {
	y := make([]K, 0, len(x))
	for k, _ := range x {
		y = append(y, k)
	}
	return y
}

func Values[K comparable, V any](x map[K]V) []V {
	y := make([]V, 0, len(x))
	for _, v := range x {
		y = append(y, v)
	}
	return y
}

func Items[K comparable, V any](x map[K]V) ([]K, []V) {
	xLen := len(x)
	ks := make([]K, 0, xLen)
	vs := make([]V, 0, xLen)
	for k, v := range x {
		ks = append(ks, k)
		vs = append(vs, v)
	}
	return ks, vs
}
