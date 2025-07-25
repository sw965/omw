package slices_test

import (
	"fmt"
	oslices "github.com/sw965/omw/slices"
	"golang.org/x/exp/slices"
	"testing"
)

func TestMinIndices(t *testing.T) {
	s := []int{0, 1, 2, 3, 3, 2, 1, 0}
	result := oslices.MinIndices(s)
	expected := []int{0, 7}
	if !slices.Equal(expected, result) {
		t.Errorf("テスト失敗")
	}
}

func TestMaxIndices(t *testing.T) {
	s := []int{0, 1, 2, 3, 3, 2, 1, 0}
	result := oslices.MaxIndices(s)
	expected := []int{3, 4}
	if !slices.Equal(expected, result) {
		t.Errorf("テスト失敗")
	}
}

func TestCartesianProducts(t *testing.T) {
	ss := [][]string{
		[]string{"a", "b", "c"},
		[]string{"d", "e"},
		[]string{"f", "g", "h"},
	}
	result := oslices.CartesianProducts[[]string](ss...)
	expected := [][]string{
		[]string{"a", "d", "f"},
		[]string{"a", "d", "g"},
		[]string{"a", "d", "h"},

		[]string{"a", "e", "f"},
		[]string{"a", "e", "g"},
		[]string{"a", "e", "h"},

		[]string{"b", "d", "f"},
		[]string{"b", "d", "g"},
		[]string{"b", "d", "h"},

		[]string{"b", "e", "f"},
		[]string{"b", "e", "g"},
		[]string{"b", "e", "h"},

		[]string{"c", "d", "f"},
		[]string{"c", "d", "g"},
		[]string{"c", "d", "h"},

		[]string{"c", "e", "f"},
		[]string{"c", "e", "g"},
		[]string{"c", "e", "h"},
	}
	fmt.Println(result)

	for i, s := range result {
		if !slices.Equal(s, expected[i]) {
			t.Errorf("テスト失敗")
			break
		}
	}
}

func TestArgSort(t *testing.T) {
	x := []float64{0.1, 1.0, 0.3, 0.7, 0.5}
	result := oslices.Argsort(x)
	expected := []int{0, 2, 4, 3, 1}
	if !slices.Equal(result, expected) {
		t.Errorf("テスト失敗")
	}
}

func TestToUnique(t *testing.T) {
	s := []int{0, 1, 2, 2, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5, 5, 4, 3, 2, 1, 0}
	result := oslices.ToUnique(s)
	expected := []int{0, 1, 2, 3, 4, 5}
	if !slices.Equal(result, expected) {
		t.Errorf("テスト失敗")
	}
}

func TestUniqueFirstIndices(t *testing.T) {
	s := []int{0, 1, 2, 2, 3, 3, 3, 5, 4, 3, 2, 1, 0}
	result := oslices.UniqueFirstIndices(s)
	expected := []int{0, 1, 2, 4, 7, 8}

	if !slices.Equal(result, expected) {
		fmt.Println(result)
		t.Errorf("テスト失敗")
	}
}

func TestTranspose(t *testing.T) {
	result := oslices.Transpose([][]int{
		[]int{0, 1, 2},
		[]int{3, 4},
		[]int{5, 6, 7, 8},
	})

	expected := [][]int{
		[]int{0, 3, 5},
		[]int{1, 4, 6},
		[]int{2, 7},
		[]int{8},
	}

	eq := oslices.AllFuncI(expected, func(s []int, i int) bool {
		return slices.Equal(s, result[i])
	})

	if !eq {
		t.Errorf("テスト失敗")
	}

	result = oslices.Transpose([][]int{
		[]int{10, 11, 12},
		[]int{13, 14, 15},
	})

	expected = [][]int{
		[]int{10, 13},
		[]int{11, 14},
		[]int{12, 15},
	}

	eq = oslices.AllFuncI(expected, func(s []int, i int) bool {
		return slices.Equal(s, result[i])
	})

	if !eq {
		t.Errorf("テスト失敗")
	}
}