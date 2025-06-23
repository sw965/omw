package parallel_test

import (
	"testing"
	"fmt"
	"github.com/sw965/omw/parallel"
	"runtime"
)

func TestFor(t *testing.T) {
	p := runtime.NumCPU()
	dataN := 3 * p
	s := make([]int, dataN)
	for i := range s {
		s[i] = i
	}

	result := make([]int, dataN)
	err := parallel.For(p, dataN, func(workerIdx, sliceIdx int) error {
		result[sliceIdx] = 2 * s[sliceIdx]
		return nil
	})

	if err != nil {
		panic(err)
	}
	fmt.Println(s)
	fmt.Println(result)
}