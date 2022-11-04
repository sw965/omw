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

func MakeSliceIntRange(start, end, step int) ([]int, error) {
  if start >= end {
    return []int{}, fmt.Errorf("start < end でなければならない")
  }

  if step <= 0 {
    return []int{}, fmt.Errorf("step > 0 でなければならない")
  }
  
  length := ((end -1 - start) / step) + 1
  result := make([]int, 0, length)
  for i := 0; i < length; i++ {
    result = append(result, start + (step * i))
  }
  return result, nil
}

func SliceIntCopy(data []int) []int {
  result := make([]int, 0, len(data))
  return append(result, data...)
}

func IsSliceIntEqual(data1, data2 []int) bool {
  for i, ele := range data1 {
    if ele != data2[i] {
      return false
    }
  }
  return true
}

func IsSliceIntContains(data []int, x ...int) bool {
	contains := func(xEle int) bool {
		for _, dataEle := range data {
			if dataEle == xEle {
				return true
			}
		}
		return false
	}

	for _, xEle := range x {
	  if !contains(xEle) {
		return false
	  }
	}
	return true
}

func IsSliceIntIndexOutOfRange(data []int, index int) bool {
  return len(data) <= index
}

func SliceIntReverse(data []int) []int {
  length := len(data)
  result := make([]int, 0, length)
  for i := length - 1; i > -1; i-- {
    result = append(result, data[i])
  }
  return result
}

func SliceIntAscendingConsecutiveCount(data []int) int {
  result := 1
  v := data[0]
  for _, ele := range data[1:] {
    if (v + 1) != ele {
      return result
    }
    v = ele
    result += 1
  }
  return result
}

func SliceIntDescendingConsecutiveCount(data []int) int {
  result := 1
  v := data[0]
  for _, ele := range data[1:] {
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
      consecutiveCount := SliceIntDescendingConsecutiveCount(SliceIntReverse(currentNumbers))
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