package math

import (
	"golang.org/x/exp/constraints"
	"github.com/sw965/omw"
)

func Min[X constraints.Ordered](xs ...X) X {
	y := xs[0]
	for _, x := range xs[1:] {
		if x < y {
			y = x
		}
	}
	return y
}

func Max[X constraints.Ordered](xs ...X) X {
	return omw.Math_Max(xs...)
}

func Sum[X constraints.Ordered](xs ...X) X {
	y := xs[0]
	for _, x := range xs[1:] {
		y += x
	}
	return y
}

func Mean[X constraints.Integer | constraints.Float](xs ...X) X {
	return Sum(xs...) / X(len(xs))
}