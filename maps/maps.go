package maps

func Keys[KS ~[]K, M ~map[K]V, K comparable, V any](m M) KS {
	keys := make(KS, 0, len(m))
	for k, _ := range m {
		keys = append(keys, k)
	}
	return keys
}
