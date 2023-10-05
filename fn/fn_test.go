package fn_test

import (
	"github.com/sw965/omw/fn"
	"golang.org/x/exp/slices"
	"golang.org/x/exp/maps"
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

func TestProduct2(t *testing.T) {
	xs1 := []int{1, 2, 3}
	xs2 := []int{5, 6, 7}
	f := func(x1, x2 int) int {
		return x1 * x2
	}

	result := fn.Product2[[]int](xs1, xs2, f)
	expected := []int{5, 6, 7, 10, 12, 14, 15, 18, 21}
	if !slices.Equal(result, expected) {
		t.Errorf("テスト失敗")
	}
}

func TestProduct3(t *testing.T) {
	xs1 := []int{3, 2}
	xs2 := []int{7, 6}
	xs3 := []int{10, 11}
	f := func(x1, x2, x3 int) int {
		return (x1 * x2) + x3
	}
	result := fn.Product3[[]int](xs1, xs2, xs3, f)
	expected := []int{31, 32, 28, 29, 24, 25, 22, 23}

	if !slices.Equal(result, expected) {
		t.Errorf("テスト失敗")
	}
}

func TestProduct4(t *testing.T) {
	xs1 := []int{1}
	xs2 := []int{2}
	xs3 := []int{3}
	xs4 := []int{4, 5}
	f := func(x1, x2, x3, x4 int) int {
		return x1 + x2 + x3 + x4
	}
	result := fn.Product4[[]int](xs1, xs2, xs3, xs4, f)
	expected := []int{10, 11}
	if !slices.Equal(result, expected) {
		t.Errorf("テスト失敗")
	}
}

func TestMemo(t *testing.T) {
	xs := []int{0, 1, 2, 3, 4, 5}
	f := func(x int) int {
		return x * x
	}
	result := fn.Memo[map[int]int, []int](xs, f)
	expected := map[int]int{0:0, 1:1, 2:4, 3:9, 4:16, 5:25}
	if !maps.Equal(result, expected) {
		t.Errorf("テスト失敗")
	}
}