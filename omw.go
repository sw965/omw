package omw

import (
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

func MapFunc[XS ~[]X, YS ~[]Y, X, Y any](xs XS, f func(X) Y) YS {
	ys := make(YS, len(xs))
	for i, x := range xs {
		ys[i] = f(x)
	}
	return ys
}

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

func PermutationTotalNum(n, r int) int {
	y := 1
	for i := 0; i < r; i++ {
		y *= (n - i)
	}
	return y
}

func GetPermutation(n, r int) [][]int {
	yn := PermutationTotalNum(n, r)
	y := make([][]int, 0, yn)
	if r == 0 {
		return y
	}
	var f func(int, []int)
	f = func(nest int, nums []int) {
		if nest == r {
			y = append(y, nums)
			return
		}
		for i := 0; i < n; i++ {
			isContinue := false
			for _, num := range nums {
				if num == i {
					isContinue = true
					break
				}
			}
			if isContinue {
				continue
			}
			clone := slices.Clone(nums)
			f(nest+1, append(clone, i))
		}
	}
	f(0, make([]int, 0, r))
	return y
}

func CombinationTotalNum(n, r int) int {
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

func GetCombination(n, r int) [][]int {
	nums := make([]int, r)
	for i := 0; i < r; i++ {
		nums[i] = i
	}

	yn := CombinationTotalNum(n, r)
	y := make([][]int, 0, yn)
	if r == 0 {
		return y
	}

	end := r - 1
	for i := 0; i < yn; i++ {
		clone := slices.Clone(nums)
		y = append(y, clone)
		max := Max(nums...)
		if max == (n - 1) {
			reverse := Reverse(nums)
			count := descendingConsecutiveCount(reverse...)
			idx := end - count
			if idx < 0 {
				break
			}
			nums[idx] += 1
			for j := idx + 1; j < r; j++ {
				nums[j] = nums[idx] + j - (idx)
			}
		} else {
			nums[end] += 1
		}
	}
	return y
}

func descendingConsecutiveCount[X constraints.Integer](xs ...X) int {
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

func Reverse[XS ~[]X, X any](xs XS) XS {
	n := len(xs)
	y := make(XS, 0, n)
	for i := n - 1; i > -1; i-- {
		y = append(y, xs[i])
	}
	return y
}