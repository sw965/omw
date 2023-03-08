package omw

import (
	"golang.org/x/exp/constraints"
)

func ToNumber[X, T constraints.Integer | constraints.Float](x X) T {
	return T(x)
}

func ToStrTilde[X, T ~string](x X) T {
	return T(x)
}