package omw

import (
	"testing"
)

func TestStrIndexAccess(t *testing.T) {
	str := "こんにちは世界"
	result := StrIndexAccess(str, 0)
	expected := "こ"

	if result != expected {
		t.Errorf("テスト失敗")
	}

	result = StrIndexAccess(str, 3)
	expected = "ち"
	if result != expected {
		t.Errorf("テスト失敗")
	}

	result = StrIndexAccess(str, 5)
	expected = "世"
	if result != expected {
		t.Errorf("テスト失敗")
	}

	result = StrIndexAccess(str, 6)
	expected = "界"
	if result != expected {
		t.Errorf("テスト失敗")
	}
}
