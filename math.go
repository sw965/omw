package omw

import (
	"golang.org/x/exp/slices"
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
	y := 1
	for i := 0; i < r; i++ {
		y *= (n - i)
	}
	return y
}

func IntPermutations(n, r int) [][]int {
	c := PermutationCount(n, r)
	ret := make([][]int, 0, c)
	if r == 0 {
		return ret
	}
	var f func(int, []int)
	f = func(nest int, nums []int) {
		if nest == r {
			ret = append(ret, nums)
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
	return ret
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

func IntCombinations(n, r int) [][]int {
	nums := make([]int, r)
	for i := 0; i < r; i++ {
		nums[i] = i
	}

	c := CombinationCount(n, r)
	ret := make([][]int, 0, c)
	if r == 0 {
		return ret
	}

	end := r - 1
	for i := 0; i < c; i++ {
		clone := slices.Clone(nums)
		ret = append(ret, clone)
		max := Max(nums...)
		if max == (n - 1) {
			reversed := ReverseElement(nums)
			count := CountConsecutiveDecrease(reversed...)
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
	return ret
}