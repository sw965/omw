package omw

func InvertMap[YM ~map[V]K, XM ~map[K]V, K, V comparable](xm XM) YM {
	ym := YM{}
	for k, v := range xm {
		ym[v] = k
	}
	return ym
}

func MapValueMapFunc[YM ~map[K]YV, XM ~map[K]XV, K comparable, XV, YV any](xm XM, f func(XV)YV) YM {
	ym := YM{}
	for k, v := range xm {
		ym[k] = f(v)
	}
	return ym
}