package omw

import (
	"fmt"
)

func MinInt(x ...int) int {
	result := x[0]
	for _, ele := range x[1:] {
		if ele < result {
			result = ele
		}
	}
	return result
}

func MaxInt(x ...int) int {
	result := x[0]
	for _, ele := range x[1:] {
		if ele > result {
			result = ele
		}
	}
	return result
}

func MinIntIndices(x ...int) []int {
	min := MinInt(x...)
	result := make([]int, 0, len(x))
	for i, ele := range x {
		if ele == min {
			result = append(result, i)
		}
	}
	return result
}

func MaxIntIndices(x ...int) []int {
	max := MaxInt(x...)
	result := make([]int, 0, len(x))
	for i, ele := range x {
		if ele == max {
			result = append(result, i)
		}
	}
	return result
}

func SumInt(x ...int) int {
	result := 0
	for _, ele := range x {
		result += ele
	}
	return result
}

func MinFloat64(x ...float64) float64 {
	result := x[0]
	for _, ele := range x[1:] {
		if result > ele {
			result = ele
		}
	}
	return result
}

func MaxFloat64(x ...float64) float64 {
	result := x[0]
	for _, ele := range x[1:] {
		if result < ele {
			result = ele
		}
	}
	return result
}

func MinFloat64Indices(x ...float64) []int {
	min := MinFloat64(x...)
	result := make([]int, 0, len(x))
	for i, ele := range x {
		if ele == min {
			result = append(result, i)
		}
	}
	return result
}

func MaxFloat64Indices(x ...float64) []int {
	max := MaxFloat64(x...)
	result := make([]int, 0, len(x))
	for i, ele := range x {
		if ele == max {
			result = append(result, i)
		}
	}
	return result
}

func SumFloat64(x ...float64) float64 {
	result := 0.0
	for _, ele := range x {
		result += ele
	}
	return result
}

func DescendingConsecutiveCount(x ...int) int {
	result := 1
	v := x[0]
	for _, ele := range x[1:] {
		if (v - 1) != ele {
			return result
		}
		v = ele
		result += 1
	}
	return result
}

func CombinationTotalNum(n, r int) (int, error) {
	if n < r {
		return 0, fmt.Errorf("CombinationTotal の 引数は、n >= r でなければならない")
	}

	if r <= 0 {
		return 0, fmt.Errorf("CombinationTotal の 引数は、r > 0 でなければならない")
	}

	if r == 0 {
		return 1, nil
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

func CombinationNumbers(n, r int) ([][]int, error) {
	if n < r {
		return [][]int{}, fmt.Errorf("CombinationNumbers の 引数は、n >= r でなければならない")
	}

	if r <= 0 {
		return [][]int{}, fmt.Errorf("CombinationNumbers の 引数 r は 0より大きい必要がある")
	}

	currentNumbers, err := MakeSliceIntRange(0, r, 1)

	if err != nil {
		return [][]int{}, err
	}

	currentNumbersLength := len(currentNumbers)
	currentNumbersEndIndex := currentNumbersLength - 1

	combinationTotalNum, err := CombinationTotalNum(n, r)
	if err != nil {
		return [][]int{}, err
	}

	result := make([][]int, 0, combinationTotalNum)

	for i := 0; i < combinationTotalNum; i++ {
		result = append(result, SliceIntCopy(currentNumbers))
		max := MaxInt(currentNumbers...)

		if max == (n - 1) {
			consecutiveCount := DescendingConsecutiveCount(SliceIntReverse(currentNumbers)...)
			rightMoveIndex := currentNumbersEndIndex - consecutiveCount

			if rightMoveIndex < 0 {
				break
			}

			currentNumbers[rightMoveIndex] += 1
			for j := rightMoveIndex + 1; j < currentNumbersLength; j++ {
				currentNumbers[j] = currentNumbers[rightMoveIndex] + j - (rightMoveIndex)
			}
		} else {
			currentNumbers[currentNumbersEndIndex] += 1
		}
	}
	return result, nil
}
