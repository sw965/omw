package omw

import (
	"golang.org/x/exp/constraints"
)

func Identity[X any](x X) X {
	return x
}

func ToNumber[X, Y constraints.Integer | constraints.Float](x X) Y {
	return Y(x)
}

func ToStrTilde[X, Y ~string](x X) Y {
	return Y(x)
}
