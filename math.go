package omw

import (
	"fmt"
	"golang.org/x/exp/constraints"
)

func Min[T constraints.Ordered](xs ...T) T {
	result := xs[0]
	for _, x := range xs[1:] {
		if x < result {
			result = x
		}
	}
	return result
}

func Max[T constraints.Ordered](xs ...T) T {
	result := xs[0]
	for _, x := range xs[1:] {
		if x > result {
			result = x
		}
	}
	return result
}

func Sum[T constraints.Ordered](xs ...T) T {
	result := xs[0]
	for _, x := range xs[1:] {
		result += x
	}
	return result
}

func Mean[T constraints.Float](xs ...T) T {
	return Sum(xs...) / T(len(xs))
}

func Identity[T any](x T) T {
	return x
}

func DescendingConsecutiveCount(x ...int) int {
	result := 1
	v := x[0]
	for _, ele := range x[1:] {
		expected := v - 1
		if ele != expected {
			return result
		}
		v = ele
		result += 1
	}
	return result
}

func PermutationCombinationShareErrMsg(n, r int) string {
	if n < r {
		return "n >= r でなければならない"
	}

	if n <= 0 {
		return "n は 0より 大きい必要がある"
	}

	if r < 0 {
		return "r は 0以上 である必要がある"
	}
	return ""
}

func PermutationError(n, r int) error {
	errMsg := PermutationCombinationShareErrMsg(n, r)
	if errMsg != "" {
		return fmt.Errorf("Permutationにおいて、" + errMsg)
	} else {
		return nil
	}
}

func CombinationError(n, r int) error {
	errMsg := PermutationCombinationShareErrMsg(n, r)
	if errMsg != "" {
		return fmt.Errorf("Combinationにおいて、" + errMsg)
	} else {
		return nil
	}
}

func PermutationTotalNum(n, r int) (int, error) {
	err := PermutationError(n, r)
	if err != nil {
		return 0, err
	}

	result := 1
	for i := 0; i < r; i++ {
		result *= (n - i)
	}
	return result, nil
}

func PermutationNumberss(n, r int) ([][]int, error) {
	totalNum, err := PermutationTotalNum(n, r)
	if err != nil {
		return [][]int{}, err
	}
	result := make([][]int, 0, totalNum)
	if r == 0 {
		return result, nil
	}
	var f func(int, []int)

	f = func(nest int, numbers []int) {
		if nest == r {
			result = append(result, numbers)
			return
		}

		for i := 0; i < n; i++ {
			isContinue := false

			for _, number := range numbers {
				if number == i {
					isContinue = true
					break
				}
			}

			if isContinue {
				continue
			}

			copyNumbers := make([]int, 0, r)
			for _, number := range numbers {
				copyNumbers = append(copyNumbers, number)
			}

			f(nest+1, append(copyNumbers, i))
		}
	}

	f(0, make([]int, 0, r))
	return result, nil
}

func CombinationTotalNum(n, r int) (int, error) {
	err := CombinationError(n, r)
	if err != nil {
		return 0, err
	}

	numer := 1

	for i := 0; i < r; i++ {
		numer *= (n - i)
	}

	denom := 1

	for i := 0; i < r; i++ {
		denom *= (r - i)
	}

	return numer / denom, nil
}

func CombinationNumberss(n, r int) ([][]int, error) {
	numbers := make([]int, r)
	for i := 0; i < r; i++ {
		numbers[i] = i
	}
	end := r - 1
	totalNum, err := CombinationTotalNum(n, r)
	if err != nil {
		return [][]int{}, err
	}
	result := make([][]int, 0, totalNum)

	if r == 0 {
		return result, nil
	}

	for i := 0; i < totalNum; i++ {
		copyNumbers := MapFunc(numbers, func(x int) int { return x })
		result = append(result, copyNumbers)
		max := Max(numbers...)

		if max == (n - 1) {
			reverseNumbers := Reverse(numbers)
			consecutiveCount := DescendingConsecutiveCount(reverseNumbers...)
			index := end - consecutiveCount

			if index < 0 {
				break
			}

			numbers[index] += 1
			for j := index + 1; j < r; j++ {
				numbers[j] = numbers[index] + j - (index)
			}
		} else {
			numbers[end] += 1
		}
	}
	return result, nil
}
