package omw

import (
	"fmt"
	"testing"
)

func TestMinInt(t *testing.T) {
	result := MinInt(1, 2, 0, 3, 5, 4)
	expected := 0
	if result != expected {
		t.Errorf("テスト失敗")
	}
}

func TestMaxInt(t *testing.T) {
	result := MaxInt([]int{4, 0, 2, 3, 3, 5}...)
	expected := 5
	if result != expected {
		t.Errorf("テスト失敗")
	}
}

func TestMinIntIndices(t *testing.T) {
	result := MinIntIndices([]int{0, 1, 2, 3, 4, 5, 5, 4, 3, 2, 1, 0}...)
	expected := []int{0, 11}
	if !SliceIntEqual(result, expected) {
		t.Errorf("テスト失敗")
	}
}

func TestMaxIntIndices(t *testing.T) {
	result := MaxIntIndices([]int{0, 1, 2, 3, 4, 5, 5, 3, 4, 3, 2, 1, 0}...)
	expected := []int{5, 6}
	if !SliceIntEqual(result, expected) {
		t.Errorf("テスト失敗")
	}
}

func TestSumInt(t *testing.T) {
	x, err := MakeSliceIntRange(1, 11, 1)
	if err != nil {
		panic(err)
	}

	result := SumInt(x...)
	expected := 55
	if result != expected {
		t.Errorf("テスト失敗")
	}
}

func TestMinFloat64(t *testing.T) {
	result := MinFloat64([]float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.0, 0.1, 0.2, 0.3, -1.0, 3.5}...)
	expected := -1.0
	if result != expected {
		t.Errorf("テスト失敗")
	}
}

func TestMaxFloat64(t *testing.T) {
	result := MaxFloat64([]float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.0, 0.1, 0.2, 0.3, -1.0, 3.5}...)
	expected := 3.5
	if result != expected {
		t.Errorf("テスト失敗")
	}
}

func TestMinFloat64Indices(t *testing.T) {
	result := MinFloat64Indices([]float64{0.1, 0.11, 0.12, 0.2, 0.3, 0.3, 0.2, 0.12, 0.11, 0.1}...)
	expected := []int{0, 9}
	if !SliceIntEqual(result, expected) {
		t.Errorf("テスト失敗")
	}
}

func TestMaxFloat64Indices(t *testing.T) {
	result := MaxFloat64Indices([]float64{0.1, 0.11, 0.12, 0.2, 0.3, 0.3, 0.2, 0.12, 0.11, 0.1}...)
	expected := []int{4, 5}
	if !SliceIntEqual(result, expected) {
		t.Errorf("テスト失敗")
	}
}

func TestSumFloat64(t *testing.T) {
	result := SumFloat64([]float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0}...)
	fmt.Println(result, "≒ 5.5")
}

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
	n, r := 100, 3
	result, err := CombinationNumbers(n, r)

	if err != nil {
		t.Errorf("テスト失敗")
	}

	combinationTotalNum, err := CombinationTotalNum(n, r)

	if err != nil {
		panic(err)
	}

	if len(result) != combinationTotalNum {
		t.Errorf("テスト失敗")
	}

	fmt.Println(CombinationNumbers(7, 5))
}
