package omw

func SliceFloat64IndicesAccess(x []float64, indices []int) []float64 {
	result := make([]float64, len(indices))
	for i, index := range indices {
		result[i] = x[index]
	}
	return result
}
