package omw

import (
	"testing"
)

func TestMakeIntRange(t *testing.T) {
	result, err := MakeIntRange(10, 29, 3)
	if err != nil {
		t.Errorf("テスト失敗")
	}

	expected := []int{10, 13, 16, 19, 22, 25, 28}
	if !IsSliceIntEqual(result, expected) {
		t.Errorf("テスト失敗")
	}

	result, err = MakeIntRange(10, 28, 3)
	if err != nil {
		t.Errorf("テスト失敗")
	}

	expected = []int{10, 13, 16, 19, 22, 25}
	if !IsSliceIntEqual(result, expected) {
		t.Errorf("テスト失敗")
	}

	_, err = MakeIntRange(100, 101, 1)
	if err != nil {
		t.Errorf("テスト失敗")
	}

	_, err = MakeIntRange(100, 100, 1)
	if err == nil {
		t.Errorf("テスト失敗")
	}

	_, err = MakeIntRange(100, 150, 0)
	if err == nil {
		t.Errorf("テスト失敗")
	}
}

func TestSliceIntContains(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}

	result := SliceIntContains(data, 1)
	expected := true

	if result != expected {
		t.Errorf("テスト失敗")
	}

	result = SliceIntContains(data, 5)
	expected = true

	if result != expected {
		t.Errorf("テスト失敗")
	}

	result = SliceIntContains(data, 0)
	expected = false

	if result != expected {
		t.Errorf("テスト失敗")
	}

	result = SliceIntContains(data, 6)
	expected = false

	if result != expected {
		t.Errorf("テスト失敗")
	}

	result = SliceIntContains(data, []int{1, 5}...)
	expected = true

	if result != expected {
		t.Errorf("テスト失敗")
	}

	result = SliceIntContains(data, []int{0, 5}...)
	expected = false

	if result != expected {
		t.Errorf("テスト失敗")
	}

	result = SliceIntContains(data, []int{1, 6}...)
	expected = false

	if result != expected {
		t.Errorf("テスト失敗")
	}
}