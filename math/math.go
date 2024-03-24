package math

import (
	"github.com/sw965/omw"
	"golang.org/x/exp/constraints"
)

func Add[X constraints.Integer | constraints.Float](a, b X) X {
	return a + b
}

func Sub[X constraints.Integer | constraints.Float](a, b X) X {
	return a - b
}

func Mul[X constraints.Integer | constraints.Float](a, b X) X {
	return a * b
}

func Div[X constraints.Integer | constraints.Float](a, b X) X {
	return a / b
}

func Min[X constraints.Ordered](xs ...X) X {
	return omw.Min(xs...)
}

func Max[X constraints.Ordered](xs ...X) X {
	return omw.Max(xs...)
}

func Sum[X constraints.Ordered](xs ...X) X {
	return omw.Sum(xs...)
}

func Mean[X constraints.Integer | constraints.Float](xs ...X) X {
	return omw.Mean(xs...)
}

type Permutation struct {
	N int
	R int
}

func (p *Permutation) TotalNum() int {
	return omw.PermutationTotalNum(p.N, p.R)
}

func (p *Permutation) Get() [][]int {
	return omw.GetPermutation(p.N, p.R)
}

type Combination struct {
	N int
	R int
}

func (c *Combination) TotalNum() int {
	return omw.CombinationTotalNum(c.N, c.R)
}

func (c *Combination) Get() [][]int {
	return omw.GetCombination(c.N, c.R)
}
