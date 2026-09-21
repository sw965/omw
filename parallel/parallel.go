package parallel

import (
	"context"
	"errors"
	"fmt"
)

func For(n, p int, f func(workerIdx, itemIdx int) error) error {
	return ForContext(context.Background(), n, p, f)
}

func ForContext(ctx context.Context, n, p int, f func(workerIdx, itemIdx int) error) error {
	if f == nil {
		return errors.New("f が nil")
	}
	if n < 0 {
		return fmt.Errorf("n >= 0 であるべき: n = %d", n)
	}
	if p < 1 {
		return fmt.Errorf("p >= 1 であるべき: p = %d", p)
	}
	if n == 0 {
		return nil
	}
	if p > n {
		p = n
	}

	if ctx == nil {
		ctx = context.Background()
	}

	// qは各workerに均等に配分する量
	// rは均等に配分しきれずに余った量
	q := n / p
	r := n % p

	errCh := make(chan error, p)

	worker := func(workerIdx, start, end int) {
		for itemIdx := start; itemIdx < end; itemIdx++ {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			default:
			}

			if err := f(workerIdx, itemIdx); err != nil {
				errCh <- fmt.Errorf("エラーが発生: workerIdx = %d, itemIdx = %d, err = %w", workerIdx, itemIdx, err)
				return
			}
		}
		errCh <- nil
	}

	start := 0
	for workerIdx := 0; workerIdx < p; workerIdx++ {
		size := q
		// 余った量をworkerIdxが低い順から1つずつ割り当てる
		// 理解がしにくければ、parallel_test.goのTestFor_TableDriven関数の最初のテストケースを参照
		if workerIdx < r {
			size++
		}
		end := start + size
		go worker(workerIdx, start, end)
		start = end
	}

	errs := make([]error, p)
	for i := 0; i < p; i++ {
		errs[i] = <-errCh
	}
	return errors.Join(errs...)
}
