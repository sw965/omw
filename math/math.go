package math

import (
	"golang.org/x/exp/constraints"
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
	y := xs[0]
	for _, x := range xs[1:] {
		if x > y {
			y = x
		}
	}
	return y
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

func CountConsecutiveDecrease[X constraints.Integer](xs ...X) int {
	y := 1
	a := xs[0] - 1
	for _, x := range xs[1:] {
		if x != a {
			break
		}
		y += 1
		a = x - 1
	}
	return y
}

func PermutationCount(n, r int) int {
	c := 1
	for i := 0; i < r; i++ {
		c *= (n - i)
	}
	return c
}

func SequenceCount(n, r int) int {
	c := 1
	for i := 0; i < r; i++ {
		c *= n
	}
	return c
}

func CombinationCount(n, r int) int {
	a := 1
	for i := 0; i < r; i++ {
		a *= (n - i)
	}
	m := 1
	for i := 0; i < r; i++ {
		m *= (r - i)
	}
	return a / m
}
