package omw

import (
	"golang.org/x/exp/constraints"
)

func IsRemainderZero[X constraints.Integer](a X) func(X) bool {
	return func(x X) bool { return x%a == 0 }
}

func Identity[X any](x X) X {
	return x
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

func DescendingConsecutiveCount(xs ...int) int {
	y := 1
	expected := xs[0] - 1
	for _, x := range xs[1:] {
		if x != expected {
			return y
		}
		expected = x - 1
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
	end := r - 1
	yn := CombinationTotalNum(n, r)
	y := make([][]int, 0, yn)
	if r == 0 {
		return y
	}
	for i := 0; i < yn; i++ {
		copyNums := MapFunc(nums, Identity[int])
		y = append(y, copyNums)
		max := Max(nums...)
		if max == (n - 1) {
			reverseNums := Reverse(nums)
			count := DescendingConsecutiveCount(reverseNums...)
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
