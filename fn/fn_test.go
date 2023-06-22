package fn_test

import (
	"testing"
	"github.com/sw965/omw/fn"
	"golang.org/x/exp/slices"
)

func TestMap(t *testing.T) {
	xs := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	add10 := func(x int) int {return x + 10}
	result := fn.Map[[]int, []int](xs, add10)
	expected := []int{10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	if !slices.Equal(result, expected) {
		t.Errorf("テスト失敗")
	}
}

func TestFilter(t *testing.T) {
	xs := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	even := func(x int) bool {return x%2 == 0}
	result := fn.Filter[[]int](xs, even)
	expected := []int{0, 2, 4, 6, 8, 10}
	if !slices.Equal(result, expected) {
		t.Errorf("テスト失敗")
	}
}