package omw

import (
	"testing"
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