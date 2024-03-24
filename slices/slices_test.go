package slices_test

import (
	omwslices "github.com/sw965/omw/slices"
	"golang.org/x/exp/slices"
	"strings"
	"testing"
)

func TestIsSubset(t *testing.T) {
	a := []string{"a", "b", "c", "d", "e", "f", "g"}
	b := []string{"c", "e", "g"}
	if !omwslices.IsSubset(a, b) {
		t.Errorf("テスト失敗")
	}

	if omwslices.IsSubset(b, a) {
		t.Errorf("テスト失敗")
	}
}

func TestAtIndex(t *testing.T) {
	xs := []int{1, 3, 5, 7, 9, 11}
	result := omwslices.AtIndex(xs)(1)
	expected := 3
	if result != expected {
		t.Errorf("テスト失敗")
	}
}

func TestAtIndices(t *testing.T) {
	xs := []float64{1.0, 0.9, 0.8, 0.7, 0.6, 0.5, 0.4, 0.3, 0.2, 0.1, 0.0}
	result := omwslices.AtIndices(xs)([]int{10, 5, 0})
	expected := []float64{0.0, 0.5, 1.0}
	if !slices.Equal(result, expected) {
		t.Errorf("テスト失敗")
	}
}

func TestCount(t *testing.T) {
	xs := []int{7, 0, 1, 2, 3, 7, 3, 2, 1, 0, 7}
	result := omwslices.Count(xs, 7)
	expected := 3
	if result != expected {
		t.Errorf("テスト失敗")
	}
}

func TestCountFunc(t *testing.T) {
	xs := []string{"abc", "efa", "gge", "ccv", "ukw", "ank"}
	f := func(x string) bool { return strings.Contains(x, "a") }
	result := omwslices.CountFunc(xs, f)
	expected := 3
	if result != expected {
		t.Errorf("テスト失敗")
	}
}

func TestReverse(t *testing.T) {
	xs := []bool{false, true, false, true}
	result := omwslices.Reverse(xs)
	expected := []bool{true, false, true, false}
	if !slices.Equal(result, expected) {
		t.Errorf("テスト失敗")
	}
}

func TestIndices(t *testing.T) {
	xs := []string{"Golang", "Python", "Go", "Java", "GO", "Haskell", "C", "Go", "C++", "Go"}
	result := omwslices.Indices(xs, "Go")
	expected := []int{2, 7, 9}
	if !slices.Equal(result, expected) {
		t.Errorf("テスト失敗")
	}
}

func TestIndicesFunc(t *testing.T) {
	xs := []int{0, 10, 15, 8, 9, 7, 2, 6, 18}
	f := func(x int) bool { return x%2 == 0 && (x >= 10) }
	result := omwslices.IndicesFunc(xs, f)
	expected := []int{1, 8}
	if !slices.Equal(result, expected) {
		t.Errorf("テスト失敗")
	}
}

func TestUnique(t *testing.T) {
	xs := []string{"a", "b", "a", "c", "d", "a", "e", "b", "f", "c", "g"}
	result := omwslices.Unique(xs)
	expected := []string{"a", "b", "c", "d", "e", "f", "g"}
	if !slices.Equal(result, expected) {
		t.Errorf("テスト失敗")
	}
}

func TestIsUnique(t *testing.T) {
	xs := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	if !omwslices.IsUnique(xs) {
		t.Errorf("テスト失敗")
	}

	xs = []int{1, 2, 3, 4, 5, 6, 7, 8, 1, 9}
	if omwslices.IsUnique(xs) {
		t.Errorf("テスト失敗")
	}
}

func TestPermutation(t *testing.T) {
	xs := []string{"太郎", "翔太", "尚子", "博美", "浩二"}
	result := omwslices.Permutation[[][]string, []string](xs, 3)

	if len(result) != 60 {
		t.Errorf("テスト失敗")
	}

breakLabel:
	for _, names1 := range result {
		equalCount := 0
		for _, names2 := range result {
			if slices.Equal(names1, names2) {
				equalCount += 1
			}
			if equalCount > 1 {
				t.Errorf("テスト失敗")
				break breakLabel
			}
		}

		if equalCount != 1 {
			t.Errorf("テスト失敗")
			break breakLabel
		}
	}
}

func TestCombination(t *testing.T) {
	xs := []string{"a", "b", "c", "d", "e"}
	result := omwslices.Combination[[][]string, []string](xs, 3)
	if len(result) != ((5 * 4 * 3) / (3 * 2)) {
		t.Errorf("テスト失敗")
	}

breakLabel:
	for _, as1 := range result {
		equalCount := 0
		for _, as2 := range result {
			sorted1 := omwslices.Sorted(as1)
			sorted2 := omwslices.Sorted(as2)

			if slices.Equal(sorted1, sorted2) {
				equalCount += 1
			}

			if equalCount > 1 {
				t.Errorf("テスト失敗")
				break breakLabel
			}
		}
		if equalCount != 1 {
			t.Errorf("テスト失敗")
			break breakLabel
		}
	}
}

func TestAll(t *testing.T) {
	bs := []bool{true, true, true}
	result := omwslices.All(bs)
	if !result {
		t.Errorf("テスト失敗")
	}

	bs = []bool{true, true, false}
	result = omwslices.All(bs)
	if result {
		t.Errorf("テスト失敗")
	}
}

func TestDeleteAtIndices(t *testing.T) {
	xs := []int{2, 4, 6, 8, 10}
	result1, result2 := omwslices.DeleteAtIndices(xs, 1, 3)
	expected1, expected2 := []int{2, 6, 10}, []int{4, 8}
	if !slices.Equal(result1, expected1) {
		t.Errorf("テスト失敗")
	}
	if !slices.Equal(result2, expected2) {
		t.Errorf("テスト失敗")
	}
}

func TestConcat(t *testing.T) {
	xs1 := []int{0, 1, 2, 3, 4}
	xs2 := []int{5, 6, 7, 8, 9}
	result := omwslices.Concat(xs1, xs2)
	expected := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	if !slices.Equal(result, expected) {
		t.Errorf("テスト失敗")
	}
}