package omw

import (
	"testing"
	"fmt"
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

func TestMakeIntRange(t *testing.T) {
	result, err := MakeSliceIntRange(10, 29, 3)
	if err != nil {
		t.Errorf("テスト失敗")
	}

	expected := []int{10, 13, 16, 19, 22, 25, 28}
	if !IsSliceIntEqual(result, expected) {
		t.Errorf("テスト失敗")
	}

	result, err = MakeSliceIntRange(10, 28, 3)
	if err != nil {
		t.Errorf("テスト失敗")
	}

	expected = []int{10, 13, 16, 19, 22, 25}
	if !IsSliceIntEqual(result, expected) {
		t.Errorf("テスト失敗")
	}

	_, err = MakeSliceIntRange(100, 101, 1)
	if err != nil {
		t.Errorf("テスト失敗")
	}

	_, err = MakeSliceIntRange(100, 100, 1)
	if err == nil {
		t.Errorf("テスト失敗")
	}

	_, err = MakeSliceIntRange(100, 150, 0)
	if err == nil {
		t.Errorf("テスト失敗")
	}
}

func TestSliceIntContains(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}

	result := IsSliceIntContains(data, 1)
	expected := true

	if result != expected {
		t.Errorf("テスト失敗")
	}

	result = IsSliceIntContains(data, 5)
	expected = true

	if result != expected {
		t.Errorf("テスト失敗")
	}

	result = IsSliceIntContains(data, 0)
	expected = false

	if result != expected {
		t.Errorf("テスト失敗")
	}

	result = IsSliceIntContains(data, 6)
	expected = false

	if result != expected {
		t.Errorf("テスト失敗")
	}

	result = IsSliceIntContains(data, []int{1, 5}...)
	expected = true

	if result != expected {
		t.Errorf("テスト失敗")
	}

	result = IsSliceIntContains(data, []int{0, 5}...)
	expected = false

	if result != expected {
		t.Errorf("テスト失敗")
	}

	result = IsSliceIntContains(data, []int{1, 6}...)
	expected = false

	if result != expected {
		t.Errorf("テスト失敗")
	}
}

func TestSliceIntReverse(t *testing.T) {
	result := SliceIntReverse([]int{0, 2, 4, 6, 8, 10})
	expected := []int{10, 8, 6, 4, 2, 0}
	
	if !IsSliceIntEqual(result, expected) {
		t.Errorf("テスト失敗")
	}
}

func TestSliceIntAscendingConsecutiveCount(t *testing.T) {
	result := SliceIntAscendingConsecutiveCount([]int{0, 2, 3, 4, 5})
	expected := 1

	if result != expected {
		t.Errorf("テスト失敗")
	}

	result = SliceIntAscendingConsecutiveCount([]int{1, 2, 3, 5, 6, 7})
	expected = 3

	if result != expected {
		t.Errorf("テスト失敗")
	}
}

func TestSliceIntDescendingConsecutiveCount(t *testing.T) {
	result := SliceIntDescendingConsecutiveCount([]int{5, 4, 3, 1, 0})
	expected := 3

	if result != expected {
		t.Errorf("テスト失敗")
	}

	result = SliceIntDescendingConsecutiveCount([]int{7, 5, 3, 2, 1})
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

	fmt.Println(CombinationNumbers(10, 5))
}