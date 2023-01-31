package omw

func Keys[K comparable, V any](kv map[K]V) []K {
	result := make([]K, len(kv))
	for k, _ := range kv {
		result = append(result, k)
	}
	return result
}

func Values[K comparable, V any](kv map[K]V) []V {
	result := make([]V, len(kv))
	for _, v := range kv {
		result = append(result, v)
	}
	return result
}

func Items[K comparable, V any](kv map[K]V) ([]K, []V) {
	length := len(kv)
	ks := make([]K, length)
	vs := make([]V, length)
	for k, v := range kv {
		ks = append(ks, k)
		vs = append(vs, v)
	}
	return ks, vs
}
