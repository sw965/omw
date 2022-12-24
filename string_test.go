package omw

import (
	"testing"
)

func TestStringIndexAccess(t *testing.T) {
	str := "こんにちは世界"
	result := StringIndexAccess(str, 0)
	expected := "こ"

	if result != expected {
		t.Errorf("テスト失敗")
	}

	result = StringIndexAccess(str, 3)
	expected = "ち"
	if result != expected {
		t.Errorf("テスト失敗")
	}

	result = StringIndexAccess(str, 5)
	expected = "世"
	if result != expected {
		t.Errorf("テスト失敗")
	}

	result = StringIndexAccess(str, 6)
	expected = "界"
	if result != expected {
		t.Errorf("テスト失敗")
	}
}
