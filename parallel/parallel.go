// Package parallel provides utility functions for parallel execution.
// Package parallel は並列実行のためのユーティリティ関数を提供します。
package parallel

import (
	"errors"
	"fmt"
)

// For executes a loop from 0 to n-1 in parallel using p workers.
// It distributes the iteration range as evenly as possible among the workers.
//
// If n is negative, it returns ErrNegativeN. If p is less than 1, it returns ErrInvalidP.
// If n is 0, it returns nil immediately.
//
// If the callback function f returns an error, the worker stops processing
// subsequent indices assigned to it. Errors from all workers are aggregated
// using errors.Join and returned.
//
// For は p 個のワーカーを使用して、0 から n-1 までのループを並列に実行します。
// 反復範囲は、ワーカー間で可能な限り均等に分配されます。
//
// n が負の場合は ErrNegativeN を返し、p が 1 未満の場合は ErrInvalidP を返します。
// n が 0 の場合は、直ちに nil を返します。
//
// コールバック関数 f がエラーを返した場合、そのワーカーは割り当てられた後続のインデックスの
// 処理を停止します。すべてのワーカーからのエラーは errors.Join を使用して集約され、返されます。
func For(n, p int, f func(workerId, idx int) error) error {
	if n < 0 {
		return fmt.Errorf("nが不正(<0): n=%d", n)
	}
	if p < 1 {
		return fmt.Errorf("pが不正(<1): p=%d", p)
	}
	if n == 0 {
		return nil
	}
	if p > n {
		p = n
	}

	// qは各ワーカーに均等に配分する量
	// rは余った量
	q := n / p
	r := n % p

	errCh := make(chan error, p)

	worker := func(workerId, start, end int) {
		for idx := start; idx < end; idx++ {
			if err := f(workerId, idx); err != nil {
				errCh <- fmt.Errorf("workerId:%d, idx: %d, %w", workerId, idx, err)
				return
			}
		}
		errCh <- nil
	}

	start := 0
	for workerId := 0; workerId < p; workerId++ {
		size := q
		// 余った量をidが低い順から一つずつ割り当てる
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
