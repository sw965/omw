package parallel

func DistributeIndicesEvenly(n, p int) [][]int {
	q := n / p
	r := n % p
	result := make([][]int, p)
	idx := 0
	for i := 0; i < p; i++ {
		size := q
		if i < r {
			size += 1
		}
		group := make([]int, size)
		for j := 0; j < size; j++ {
			group[j] = idx
			idx += 1
		}
		result[i] = group
	}
	return result
}
