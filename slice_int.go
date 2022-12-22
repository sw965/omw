package omw

import (
	"fmt"
)

func MakeSliceInt(start, end, step int) ([]int, error) {
	if start >= end {
		return []int{}, fmt.Errorf("start < end でなければならない")
	}

	if step <= 0 {
		return []int{}, fmt.Errorf("step > 0 でなければならない")
	}

	length := ((end - 1 - start) / step) + 1
	result := make([]int, 0, length)
	for i := 0; i < length; i++ {
		result = append(result, start+(step*i))
	}
	return result, nil
}

func SliceIntCopy(x []int) []int {
	result := make([]int, 0, len(x))
	return append(result, x...)
}

func SliceIntReverse(x []int) []int {
	length := len(x)
	result := make([]int, 0, length)
	for i := length - 1; i > -1; i-- {
		result = append(result, x[i])
	}
	return result
}

func SliceIntIndicesAccess(x, indices []int) []int {
	result := make([]int, len(indices))
	for i, index := range indices {
		result[i] = x[index]
	}
	return result
}

func SliceIntEqual(x1, x2 []int) bool {
	for i, ele := range x1 {
		if ele != x2[i] {
			return false
		}
	}
	return true
}

func SliceIntContains(x []int, n ...int) bool {
	contains := func(nEle int) bool {
		for _, xEle := range x {
			if xEle == nEle {
				return true
			}
		}
		return false
	}

	for _, nEle := range n {
		if !contains(nEle) {
			return false
		}
	}
	return true
}

func SliceIntIndices(x []int, n int) []int {
	result := make([]int, 0, len(x))
	for i, ele := range x {
		if ele == n {
			result = append(result, i)
		}
	}
	return result
}