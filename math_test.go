package omw

import (
	"testing"
	"fmt"
)

func TestDescendingConsecutiveCount(t *testing.T) {
	result := DescendingConsecutiveCount([]int{5, 4, 3, 1, 0}...)
	expected := 3

	if result != expected {
		t.Errorf("テスト失敗")
	}

	result = DescendingConsecutiveCount([]int{7, 5, 3, 2, 1}...)
	expected = 1

	if result != expected {
		t.Errorf("テスト失敗")
	}
}

func TestCombinationTotalNum(t *testing.T) {
	result, err := CombinationTotalNum(5, 3)

	if err != nil {
		t.Errorf("テスト失敗")
	}

	expected := 10

	if result != expected {
		t.Errorf("テスト失敗")
	}

	result, err = CombinationTotalNum(100, 3)

	if err != nil {
		t.Errorf("テスト失敗")
	}

	expected = 161700

	if result != expected {
		t.Errorf("テスト失敗")
	}
}

func TestCombinationNumbers(t *testing.T) {
	result, err :=  CombinationNumbers(5, 3)
	if err != nil {
		panic(err)
	}

	expected := [][]int{
		[]int{0, 1, 2}, []int{0, 1, 3}, []int{0, 1, 4},
		[]int{0, 2, 3}, []int{0, 2, 4},
		[]int{0, 3, 4},

		[]int{1, 2, 3}, []int{1, 2, 4},
		[]int{1, 3, 4},
		[]int{2, 3, 4},
	}

	for i, vs := range result {
		for j, v := range vs {
			if v != expected[i][j] {
				t.Errorf("テスト失敗")
			}
		} 
	}
}


func TestPermutationNumbers(t *testing.T) {
	n, r := 5, 3
	result := PermutationNumbers(n, r)
	for _, v := range result {
		fmt.Println(v)
	}
	if len(result) != PermutationTotalNum(n, r) {
		t.Errorf("テスト失敗")
	}
}