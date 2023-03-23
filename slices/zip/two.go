package zip

func NewNestMap2[NM ~map[K1]M, M ~map[K2]V, K1, K2 comparable, V any](ks1 []K1, ks2 []K2, f func(K1, K2) V) NM {
	nm := NM{}
	for _, k1 := range ks1 {
		for _, k2 := range ks2 {
			if _, ok := nm[k1]; !ok {
				nm[k1] = M{}
			}
			nm[k1][k2] = f(k1, k2)
		}
	}
	return nm
}