package slices_test

import (
	"fmt"
	omwslices "github.com/sw965/omw/slices"
	"golang.org/x/exp/slices"
	"testing"
)

func TestMinIndex(t *testing.T) {
	s := []int{0, 1, 2, 2, 3, 3, 3, 2, 2, 1, 0}
	result := omwslices.MinIndex(s)
	expected := 0
	if expected != result {
		t.Errorf("テスト失敗")
	}
}

func TestMinIndices(t *testing.T) {
	s := []int{0, 1, 2, 3, 3, 2, 1, 0}
	result := omwslices.MinIndices(s)
	expected := []int{0, 7}
	if !slices.Equal(expected, result) {
		t.Errorf("テスト失敗")
	}
}

func TestMaxIndex(t *testing.T) {
	s := []int{0, 1, 2, 2, 3, 3, 3, 2, 2, 1, 0}
	result := omwslices.MaxIndex(s)
	expected := 4
	if expected != result {
		t.Errorf("テスト失敗")
	}
}

func TestMaxIndices(t *testing.T) {
	s := []int{0, 1, 2, 3, 3, 2, 1, 0}
	result := omwslices.MaxIndices(s)
	expected := []int{3, 4}
	if !slices.Equal(expected, result) {
		t.Errorf("テスト失敗")
	}
}

func TestCartesianProduct(t *testing.T) {
	ss := [][]string{
		[]string{"a", "b", "c"},
		[]string{"d", "e"},
		[]string{"f", "g", "h"},
	}
	result := omwslices.CartesianProduct[[][]string, []string](ss...)
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
	result := omwslices.Argsort(x)
	expected := []int{0, 2, 4, 3, 1}
	if !slices.Equal(result, expected) {
		t.Errorf("テスト失敗")
	}
}
