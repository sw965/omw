package omw

import (
	"testing"
)

func TestAll(t *testing.T) {
	result := All([]bool{true, true, true, true, true}...)
	expected := true

	if result != expected {
		t.Errorf("テスト失敗")
	}

	result = All([]bool{true, true, false, true, true}...)
	expected = false

	if result != expected {
		t.Errorf("テスト失敗")
	}
}

func TestAny(t *testing.T) {
	result := Any([]bool{true, true, true, true, true}...)
	expected := true

	if result != expected {
		t.Errorf("テスト失敗")
	}

	result = Any([]bool{true, true, false, true, true}...)
	expected = true

	if result != expected {
		t.Errorf("テスト失敗")
	}

	result = Any([]bool{false, false, false, false, false}...)
	expected = false

	if result != expected {
		t.Errorf("テスト失敗")
	}
}