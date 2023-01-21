package omw

import (
	"testing"
)

func TestMapping(t *testing.T) {
	result1 := Mapping(
		[]int{0, 1, 2, 3, 4, 5},
		func(x int) int {return x * x},
	)
	expected1 := []int{0, 1, 4, 9, 16, 25}
	if !Equals(result1, expected1) {
		t.Errorf("テスト失敗")
	}

	result2 := Mapping(
		[]string{"ok", "no", "wtf", "fuck", "omfg"},
		func(x string) string {return StringIndexAccess(x, 0)},
	)
	expected2 := []string{"o", "n", "w", "f", "o"}
	if !Equals(result2, expected2) {
		t.Errorf("テスト失敗")
	}
}

func TestFilter(t *testing.T) {
	result1 := Filter(
		[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		func(x int) bool {return x%2 == 1},
	)
	expected1 := []int{1, 3, 5, 7, 9}
	if !Equals(result1, expected1) {
		t.Errorf("テスト失敗")
	}
}

func TestIsUnique(t *testing.T) {
	result1 := IsUnique([]string{"omfg", "wtf", "no", "ok", "holy"})
	expected1 := true
	if result1 != expected1 {
		t.Errorf("テスト失敗")
	}

	result2 := IsUnique([]string{"omfg", "wtf", "no", "wtf", "holy"})
	expected2 := false
	if result2 != expected2 {
		t.Errorf("テスト失敗")
	}
}