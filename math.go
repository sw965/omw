package omw

import (
	"fmt"
)

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

	currentNumbers := make([]int, r)
	for i := 0; i < r; i++ {
		currentNumbers[i] = i
	}

	currentNumbersLength := len(currentNumbers)
	currentNumbersEndIndex := currentNumbersLength - 1

	combinationTotalNum, err := CombinationTotalNum(n, r)
	if err != nil {
		return [][]int{}, err
	}

	result := make([][]int, 0, combinationTotalNum)

	for i := 0; i < combinationTotalNum; i++ {
		copyCurrentNumbers := make([]int, currentNumbersLength)
		for i, v := range currentNumbers {
			copyCurrentNumbers[i] = v
		}

		result = append(result, copyCurrentNumbers)

		max := currentNumbers[0]
		for _, v := range currentNumbers[1:] {
			if v > max {
				max = v
			}
		}

		if max == (n - 1) {
			reverseCurrentNumbers := make([]int, 0, currentNumbersLength)
			for j := currentNumbersLength - 1; j > -1; j-- {
				reverseCurrentNumbers = append(reverseCurrentNumbers, currentNumbers[j])
			}

			consecutiveCount := DescendingConsecutiveCount(reverseCurrentNumbers...)
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
