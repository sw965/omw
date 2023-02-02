package omw

import (
	"golang.org/x/exp/constraints"
)

func Identity[T any](x T) T {
	return x
}

func Min[T constraints.Ordered](xs ...T) T {
	y := xs[0]
	for _, x := range xs[1:] {
		if x < y {
			y = x
		}
	}
	return y
}

func Max[T constraints.Ordered](xs ...T) T {
	y := xs[0]
	for _, x := range xs[1:] {
		if x > y {
			y = x
		}
	}
	return y
}

func DescendingConsecutiveCount(xs ...int) int {
	y := 1
	xExpected := xs[0] - 1
	for _, x := range xs[1:] {
		if x != xExpected {
			return y
		}
		xExpected = x - 1
		y += 1
	}
	return y
}

func PermutationTotalNum(n, r int) int {
	y := 1
	for i := 0; i < r; i++ {
		y *= (n - i)
	}
	return y
}

func PermutationNumberss(n, r int) [][]int {
	yLen := PermutationTotalNum(n, r)
	y := make([][]int, 0, yLen)
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
			copyNums := MapFunc(nums, Identity[int])
			f(nest+1, append(copyNums, i))
		}
	}
	f(0, make([]int, 0, r))
	return y
}

func CombinationTotalNum(n, r int) int {
	numer := 1
	for i := 0; i < r; i++ {
		numer *= (n - i)
	}
	denom := 1
	for i := 0; i < r; i++ {
		denom *= (r - i)
	}
	return numer / denom
}

func CombinationNumberss(n, r int) [][]int {
	nums := make([]int, r)
	for i := 0; i < r; i++ {
		nums[i] = i
	}
	endIdx := r - 1
	yLen := CombinationTotalNum(n, r)
	y := make([][]int, 0, yLen)
	if r == 0 {
		return y
	}
	for i := 0; i < yLen; i++ {
		copyNums := MapFunc(nums, Identity[int])
		y = append(y, copyNums)
		max := Max(nums...)
		if max == (n - 1) {
			reverseNums := Reverse(nums)
			consecutiveCount := DescendingConsecutiveCount(reverseNums...)
			idx := endIdx - consecutiveCount
			if idx < 0 {
				break
			}
			nums[idx] += 1
			for j := idx + 1; j < r; j++ {
				nums[j] = nums[idx] + j - (idx)
			}
		} else {
			nums[endIdx] += 1
		}
	}
	return y
}
