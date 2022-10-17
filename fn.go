package omw

import (
  "fmt"
)

func MinInt(x1, x2 int) int {
  if x1 < x2 {
    return x1
  }
  return x2
}

func MaxInt(x1, x2 int) int {
  if x1 > x2 {
    return x1
  }
  return x2
}

func MakeIntRange(start, end, step int) ([]int, error) {
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

func IsSliceIntEqual(data1, data2 []int) bool {
  for i, ele := range data1 {
    if ele != data2[i] {
      return false
    }
  }
  return len(data1) == len(data2)
}

func SliceIntContains(data []int, x ...int) bool {
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