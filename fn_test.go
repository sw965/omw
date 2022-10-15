package omw

import (
	"fmt"
	"testing"
)

func TestMakeIntRange(t *testing.T) {
	result, err := MakeIntRange(10, 29, 3)
	if err != nil {
		t.Errorf("テスト失敗")
		return
	}

	expected := []int{10, 13, 16, 19, 22, 25, 28}
	if !IsSliceIntEqual(result, expected) {
		t.Errorf("テスト失敗")
		return
	}

	result, err = MakeIntRange(10, 28, 3)
	if err != nil {
		t.Errorf("テスト失敗")
		return
	}

	expected = []int{10, 13, 16, 19, 22, 25}
	if !IsSliceIntEqual(result, expected) {
		t.Errorf("テスト失敗")
		return
	}

	_, err = MakeIntRange(100, 101, 1)
	if err != nil {
		t.Errorf("テスト失敗")
		return
	}

	_, err = MakeIntRange(100, 100, 1)
	if err == nil {
		t.Errorf("テスト失敗")
		return
	}

	_, err = MakeIntRange(100, 150, 0)
	if err == nil {
		t.Errorf("テスト失敗")
		return
	}

	fmt.Println("テスト成功")
}