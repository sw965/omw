package maps

import (
	"fmt"
)

func Invert[YM ~map[V]K, XM ~map[K]V, K, V comparable](xm XM) YM {
	ym := YM{}
	for k, v := range xm {
		ym[v] = k
	}
	return ym
}