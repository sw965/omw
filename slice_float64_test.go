package omw

import (
	"testing"
	"fmt"
)

func TestSliceFloat64Reverse(t *testing.T) {
	result := SliceFloat64Reverse([]float64{0.7, 0.1, 0.7, 0.2, 0.2, 0.7, 0.3, 0.3, 0.3, 0.7})
	expected := []float64{0.7, 0.3, 0.3, 0.3, 0.7, 0.2, 0.2, 0.7, 0.1, 0.7}
	if !SliceFloat64Equal(result, expected) {
		t.Errorf("テスト失敗")
	}
}

func TestSliceFloat64Indices(t *testing.T) {
	result := SliceFloat64Indices([]float64{0.7, 0.1, 0.7, 0.2, 0.2, 0.7, 0.3, 0.3, 0.3, 0.7}, 0.7)
	expected := []int{0, 2, 5, 9}
	if !SliceIntEqual(result, expected) {
		t.Errorf("テスト失敗")
	}
}

func TestSliceFloat64Mean(t *testing.T) {
	result := SliceFloat64Mean([]float64{0.1, 0.2, 0.3, 0.4, 0.5})
	expected := 0.3
	testMsg := fmt.Sprintf("%v ≒ %v", result, expected)
	fmt.Println(testMsg)

	result = SliceFloat64Mean([]float64{0.0, 1.0})
	expected = 0.5
	testMsg = fmt.Sprintf("%v ≒ %v", result, expected)
	fmt.Println(testMsg)
}