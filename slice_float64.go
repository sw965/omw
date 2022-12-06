package omw

func SliceFloat64IndicesAccess(x []float64, indices []int) []float64 {
	result := make([]float64, len(indices))
	for i, index := range indices {
		result[i] = x[index]
	}
	return result
}

func SliceFloat64Reverse(x []float64) []float64 {
	length := len(x)
	result := make([]float64, 0, length)
	for i := length - 1; i > -1; i-- {
		result = append(result, x[i])
	}
	return result
}

func SliceFloat64Indices(x []float64, n float64) []int {
	result := make([]int, 0, len(x))
	for i, ele := range x {
		if ele == n {
			result = append(result, i)
		}
	}
	return result
}

func SliceFloat64Equal(x1, x2 []float64) bool {
	for i, ele := range x1 {
		if ele != x2[i] {
			return false
		}
	}
	return true
}