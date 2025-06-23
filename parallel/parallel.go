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

func For(p, n int, f func(int, int) error) error {
	errCh := make(chan error, p)

	worker := func(workerIdx int, sIdxs []int) {
		for _, sIdx := range sIdxs {
			err := f(workerIdx, sIdx)
			if err != nil {
				errCh <- err
				return
			}
		}
		errCh <- nil
	}

	for workerIdx, sIdxs := range DistributeIndicesEvenly(n, p) {
		go worker(workerIdx, sIdxs)
	}

	for i := 0; i < p; i++ {
		if err := <- errCh; err != nil {
			return err
		}
	}
	return nil
}
