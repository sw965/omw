//相互import対策

package omw

import (
	"golang.org/x/exp/slices"
	"golang.org/x/exp/constraints"
)

func Fn_Map[XS ~[]X, YS ~[]Y, X, Y any](xs []X, f func(X) Y) YS {
	ys := make(YS, len(xs))
	for i, x := range xs {
		ys[i] = f(x)
	}
	return ys
}


func descendingConsecutiveCount[X constraints.Integer](xs ...X) int {
	y := 1
	expected := xs[0] - 1
	for _, x := range xs[1:] {
		if x != expected {
			break
		}
		y += 1
		expected = x - 1
	}
	return y
}

func Math_Max[X constraints.Ordered](xs ...X) X {
	y := xs[0]
	for _, x := range xs[1:] {
		if x > y {
			y = x
		}
	}
	return y
}

func Math_CombinationTotalNum(n, r int) int {
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

func Math_Combination(n, r int) [][]int {
	nums := make([]int, r)
	for i := 0; i < r; i++ {
		nums[i] = i
	}
	end := r - 1

	yn := Math_CombinationTotalNum(n, r)
	y := make([][]int, 0, yn)
	if r == 0 {
		return y
	}
	for i := 0; i < yn; i++ {
		clone := slices.Clone(nums)
		y = append(y, clone)
		max := Math_Max(nums...)
		if max == (n - 1) {
			reverse := Slices_Reverse(nums)
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

func Math_PermutationTotalNum(n, r int) int {
	y := 1
	for i := 0; i < r; i++ {
		y *= (n - i)
	}
	return y
}

func Math_Permutation(n, r int) [][]int {
	yn := Math_PermutationTotalNum(n, r)
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

func Slices_MakeIntegerRange[XS ~[]X, X constraints.Integer](start, end, step X) XS {
	n := int((end - start) / step)
	y := make(XS, n)
	for i := 0; i < n; i++ {
		y[i] = start + (step * X(i))
	}
	return y
}

func Slices_Reverse[XS ~[]X, X any](xs XS) XS {
	n := len(xs)
	y := make(XS, 0, n)
	for i := n - 1; i > -1; i-- {
		y = append(y, xs[i])
	}
	return y
}