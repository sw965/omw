package omw

func NewMap[K comparable, V any](ks []K, vs []V) map[K]V {
	y := map[K]V{}
	for i, k := range ks {
		y[k] = vs[i]
	}
	return y
}

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
	n := len(x)
	ks := make([]K, 0, n)
	vs := make([]V, 0, n)
	for k, v := range x {
		ks = append(ks, k)
		vs = append(vs, v)
	}
	return ks, vs
}
