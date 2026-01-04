package mathx

import (
	"cmp"
	"math"

	"github.com/sw965/omw/constraints"
)

func IsInf[F constraints.Float](f F, sign int) bool {
	return math.IsInf(float64(f), sign)
}

func IsNaN[F constraints.Float](f F) bool {
	return math.IsNaN(float64(f))
}

func Abs[F constraints.Float](x F) F {
	return F(math.Abs(float64(x)))
}

func Log[F constraints.Float](x F) F {
	return F(math.Log(float64(x)))
}

func Sqrt[F constraints.Float](x F) F {
	return F(math.Sqrt(float64(x)))
}

func Sum[T cmp.Ordered](xs ...T) T {
	var sum T
	for _, x := range xs {
		sum += x
	}
	return sum
}