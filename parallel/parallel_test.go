package parallel_test

import (
	"fmt"
	"github.com/sw965/omw/parallel"
	"runtime"
	"testing"
)

func TestFor(t *testing.T) {
	p := runtime.NumCPU()
	dataN := 3 * p
	s := make([]int, dataN)
	for i := range s {
		s[i] = i
	}

	result := make([]int, dataN)
	err := parallel.For(p, dataN, func(workerId, idx int) error {
		result[idx] = 2 * s[idx]
		return nil
	})

	if err != nil {
		panic(err)
	}
	fmt.Println(s)
	fmt.Println(result)
}
