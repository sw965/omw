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

	fmt.Println(CombinationNumbers(10, 5))
}
