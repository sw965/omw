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

	end := r - 1

	combinationTotalNum, err := CombinationTotalNum(n, r)
	if err != nil {
		return [][]int{}, err
	}

	result := make([][]int, 0, combinationTotalNum)

	for i := 0; i < combinationTotalNum; i++ {
		copyCurrentNumbers := make([]int, r)
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
			reverseCurrentNumbers := make([]int, 0, r)
			for j := r - 1; j > -1; j-- {
				reverseCurrentNumbers = append(reverseCurrentNumbers, currentNumbers[j])
			}

			consecutiveCount := DescendingConsecutiveCount(reverseCurrentNumbers...)
			rightMoveIndex := end - consecutiveCount

			if rightMoveIndex < 0 {
				break
			}

			currentNumbers[rightMoveIndex] += 1
			for j := rightMoveIndex + 1; j < r; j++ {
				currentNumbers[j] = currentNumbers[rightMoveIndex] + j - (rightMoveIndex)
			}
		} else {
			currentNumbers[end] += 1
		}
	}
	return result, nil
}

func PermutationTotalNum(n, r int) int {
	result := 1
	for i := 0; i < r; i++ {
		result *= (n - i)
	}
	return result
}

func PermutationNumbers(n, r int) [][]int {
	permutationTotalNum := PermutationTotalNum(n, r)
	result := make([][]int, 0, permutationTotalNum)
	var f func(int, []int)

	f = func(currentNest int, currentNumbers []int) {
		if currentNest == r {
			result = append(result, currentNumbers)
			return
		}
	
		for i := 0; i < n; i++ {
			isContinue := false
	
			for _, v := range currentNumbers {
				if v == i {
					isContinue = true
					break
				}
			}
	
			if isContinue {
				continue
			}

			copyCurrentNumbers := make([]int, 0, r)
			for _, v := range currentNumbers {
				copyCurrentNumbers = append(copyCurrentNumbers, v)
			}

			f(currentNest + 1, append(copyCurrentNumbers, i))
		}
	}

	f(0, make([]int, 0, r))
	return result
}