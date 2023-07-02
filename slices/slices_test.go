package slices_test

import (
	"testing"
	omws "github.com/sw965/omw/slices"
	"golang.org/x/exp/slices"
	"strings"
	"fmt"
)

func TestMakeFunc(t *testing.T) {
	f := func(i int) int {return i*2}
	result := omws.MakeFunc[[]int](5, f)
	expected := []int{0, 2, 4, 6, 8}
	if !slices.Equal(result, expected) {
		t.Errorf("テスト失敗")
	}
}

func TestIntegerRangeStep1(t *testing.T) {
	rng := omws.IntegerRange[[]int, int]{Start:0, End:5, Step:1}
	result := rng.Make()
	expected := []int{0, 1, 2, 3, 4}
	if !slices.Equal(result, expected) {
		t.Errorf("テスト失敗")
	}
}

func TestIntegerRangeStep2(t *testing.T) {
	rng := omws.IntegerRange[[]int, int]{Start:0, End:10, Step:2}
	result := rng.Make()
	expected := []int{0, 2, 4, 6, 8}
	if !slices.Equal(result, expected) {
		t.Errorf("テスト失敗")
	}
}

func TestIntegerRangeStep3(t *testing.T) {
	rng := omws.IntegerRange[[]int, int]{Start:3, End:29, Step:3}
	result := rng.Make()
	expected := []int{3, 6, 9, 12, 15, 18, 21, 24, 27}
	if !slices.Equal(result, expected) {
		t.Errorf("テスト失敗")
	}
}

func TestIsSubset(t *testing.T) {
	a := []string{"a", "b", "c", "d", "e", "f", "g"}
	b := []string{"c", "e", "g"}
	if !omws.IsSubset(a, b) {
		t.Errorf("テスト失敗")
	}

	if omws.IsSubset(b, a) {
		t.Errorf("テスト失敗")
	}
}

func TestAccess(t *testing.T) {
	xs := []float64{1.0, 0.9, 0.8, 0.7, 0.6, 0.5, 0.4, 0.3, 0.2, 0.1, 0.0}
	result := omws.Access(xs, 10, 5, 0)
	expected := []float64{0.0, 0.5, 1.0}
	if !slices.Equal(result, expected) {
		t.Errorf("テスト失敗")
	}
}

func TestCount(t *testing.T) {
	xs := []int{7, 0, 1, 2, 3, 7, 3, 2, 1, 0, 7}
	result := omws.Count(xs, 7)
	expected := 3
	if result != expected {
		t.Errorf("テスト失敗")
	}
}

func TestCountFunc(t *testing.T) {
	xs := []string{"abc", "efa", "gge", "ccv", "ukw", "ank"}
	f := func(x string) bool { return strings.Contains(x, "a")}
	result := omws.CountFunc(xs, f)
	expected := 3
	if result != expected {
		t.Errorf("テスト失敗")
	}
}

func TestReverse(t *testing.T) {
	xs := []bool{false, true, false, true}
	result := omws.Reverse(xs)
	expected := []bool{true, false, true, false}
	if !slices.Equal(result, expected) {
		t.Errorf("テスト失敗")
	}
}

func TestIndices(t *testing.T) {
	xs := []string{"Golang", "Python", "Go", "Java", "GO", "Haskell", "C", "Go", "C++", "Go"}
	result := omws.Indices(xs, "Go")
	expected := []int{2, 7, 9}
	if !slices.Equal(result, expected) {
		t.Errorf("テスト失敗")
	}
}

func TestIndicesFunc(t *testing.T) {
	xs := []int{0, 10, 15, 8, 9, 7, 2, 6, 18}
	f := func(x int) bool { return x%2==0 && (x>=10) }
	result := omws.IndicesFunc(xs, f)
	expected := []int{1, 8}
	if !slices.Equal(result, expected) {
		t.Errorf("テスト失敗")
	}
}

func TestToUnique(t *testing.T) {
	xs := []string{"a", "b", "a", "c", "d", "a", "e", "b", "f", "c", "g"}
	result := omws.ToUnique(xs)
	expected := []string{"a", "b", "c", "d", "e", "f", "g"}
	if !slices.Equal(result, expected) {
		t.Errorf("テスト失敗")
	}
}

func TestIsUnique(t *testing.T) {
	xs := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	if !omws.IsUnique(xs) {
		t.Errorf("テスト失敗")
	}

	xs = []int{1, 2, 3, 4, 5, 6, 7, 8, 1, 9}
	if omws.IsUnique(xs) {
		t.Errorf("テスト失敗")
	}
}

func TestPermutation(t *testing.T) {
	xs := []string{"太郎", "翔太", "尚子", "博美", "浩二"}
	result := omws.Permutation[[][]string, []string](xs, 3)

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
	result := omws.Combination[[][]string, []string](xs, 3)

	if len(result) != ( (5 * 4 * 3) / (3 * 2) ) {
		t.Errorf("テスト失敗")
	}

	breakLabel:
	for _, as1 := range result {
		equalCount := 0
		for _, as2 := range result {
			sorted1 := omws.Sorted(as1)
			sorted2 := omws.Sorted(as2)

			if slices.Equal(sorted1, sorted2){
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

func TestPop(t *testing.T) {
	xs := []float64{0.0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7}
	result1, result2 := omws.Pop(xs, 5)
	expected1, expected2 := []float64{0.0, 0.1, 0.2, 0.3, 0.4, 0.6, 0.7}, 0.5

	if !slices.Equal(result1, expected1) {
		t.Errorf("テスト失敗")
	}

	if result2 != expected2 {
		t.Errorf("テスト失敗")
	}
}

func TestPops(t *testing.T) {
	xs := []string{"a", "b", "c", "d", "e", "f"}
	result, ys := omws.Pops(xs, []int{2, 5})

	if !slices.Equal(result, []string{"a", "b", "d", "e"}) {
		t.Errorf("テスト失敗")
	}

	if !slices.Equal(ys, []string{"c", "f"}) {
		t.Errorf("テスト失敗")
	}
}

func TestSorted(t *testing.T) {
	xs := []string{"c", "a", "b", "d", "f", "e"}
	result := omws.Sorted(xs)
	expected := []string{"a", "b", "c", "d", "e", "f"}
	if !slices.Equal(result, expected) {
		t.Errorf("テスト失敗")
	}
}

func TestProduct(t *testing.T) {
	xs1 := []string{"a", "b", "c"}
	xs2 := []string{"d", "e"}
	xs3 := []string{"f", "g"}
	result := omws.Product[[][]string](xs1, xs2, xs3)
	fmt.Println(result)
}

func TestAll(t *testing.T) {
	bs := []bool{true, true, true}
	result := omws.All(bs)
	if !result {
		t.Errorf("テスト失敗")
	}

	bs = []bool{true, true, false}
	result = omws.All(bs)
	if result {
		t.Errorf("テスト失敗")
	}
}