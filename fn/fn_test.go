package fn_test

import (
	"github.com/sw965/omw/fn"
	"golang.org/x/exp/slices"
	"testing"
)

func TestMap(t *testing.T) {
	xs := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	add10 := func(x int) int { return x + 10 }
	result := fn.Map[[]int, []int](xs, add10)
	expected := []int{10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	if !slices.Equal(result, expected) {
		t.Errorf("テスト失敗")
	}
}

func TestMapIndex(t *testing.T) {
	f := func(i int, x bool) int {
		if x {
			return i * i
		} else {
			return i + 10
		}
	}

	xs := []bool{true, false, true, true, false, false, false}
	result := fn.MapIndex[[]int](xs, f, 1)
	expected := []int{1, 12, 9, 16, 15, 16, 17}
	if !slices.Equal(result, expected) {
		t.Errorf("テスト失敗")
	}
}

func TestFilter(t *testing.T) {
	xs := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	even := func(x int) bool { return x%2 == 0 }
	result := fn.Filter[[]int](xs, even)
	expected := []int{0, 2, 4, 6, 8, 10}
	if !slices.Equal(result, expected) {
		t.Errorf("テスト失敗")
	}
}

func TestAny(t *testing.T) {
	xs := []int{0, 2, 4, 6, 8, 9, 10}
	if !fn.Any(xs, func(x int) bool { return x%2 == 1 }) {
		t.Errorf("テスト失敗")
	}

	if fn.Any(xs, func(x int) bool { return x > 10 }) {
		t.Errorf("テスト失敗")
	}
}
