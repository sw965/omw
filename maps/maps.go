package maps

import (
	"fmt"
)

func New[M ~map[K]V, KS ~[]K, VS ~[]V, K comparable, V any](ks KS, vs VS) (M, error) {
	if len(ks) != len(vs) {
		return M{}, fmt.Errorf("keysとvaluesのlenが一致しません。len(ks) == len(vs) である必要があります。")
	}
	m := M{}
	for i, k := range ks {
		if _, ok := m[k]; ok {
			msg := fmt.Sprintf("%v が 重複しています。", k)
			return M{}, fmt.Errorf(msg)
		}
		m[k] = vs[i]
	}
	return m, nil
}

func Invert[YM ~map[V]K, XM ~map[K]V, K, V comparable](xm XM) YM {
	ym := YM{}
	for k, v := range xm {
		ym[v] = k
	}
	return ym
}