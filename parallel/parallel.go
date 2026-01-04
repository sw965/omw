package parallel

import (
	"fmt"
	"errors"
)

func For(n, p int, f func(workerId, idx int) error) error {
	if n < 0 {
		return fmt.Errorf("nが不正(n < 0): n = %d: n >= 0 であるべき", n)
	}
	if p < 1 {
		return fmt.Errorf("pが不正(p < 1): p = %d: p >= 1 であるべき", p)
	}
	if n == 0 {
		return nil
	}
	if p > n {
		p = n
	}

	// qは各workerに均等に配分する量
	// rは均等に配分しきれずに余った量
	q := n / p
	r := n % p

	errCh := make(chan error, p)

	worker := func(workerId, start, end int) {
		for idx := start; idx < end; idx++ {
			if err := f(workerId, idx); err != nil {
				errCh <- fmt.Errorf("worker %d failed at index %d: %w", workerId, idx, err)
				return
			}
		}
		errCh <- nil
	}

	start := 0
	for workerId := 0; workerId < p; workerId++ {
		size := q
		// 余った量をworkerIdが低い順から1つずつ割り当てる
		// 理解がしにくければ、parallel_test.goのTestFor関数の最初のテストケースを見るとわかりやすいかも
		if workerId < r {
			size++
		}
		end := start + size
		go worker(workerId, start, end)
		start = end
	}

	errs := make([]error, p)
	for i := 0; i < p; i++ {
		errs[i] = <-errCh
	}
	return errors.Join(errs...)
}
