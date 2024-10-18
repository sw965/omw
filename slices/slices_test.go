package slices_test

import (
	"fmt"
	omwslices "github.com/sw965/omw/slices"
	"golang.org/x/exp/slices"
	"testing"
)

func TestCartesianProduct(t *testing.T) {
	ss := [][]string{
		[]string{"a", "b", "c"},
		[]string{"d", "e"},
		[]string{"f", "g", "h"},
	}
	result := omwslices.CartesianProduct[[][]string, []string](ss...)
	expected := [][]string{
		[]string{"a", "d", "f"},
		[]string{"a", "d", "g"},
		[]string{"a", "d", "h"},

		[]string{"a", "e", "f"},
		[]string{"a", "e", "g"},
		[]string{"a", "e", "h"},

		[]string{"b", "d", "f"},
		[]string{"b", "d", "g"},
		[]string{"b", "d", "h"},

		[]string{"b", "e", "f"},
		[]string{"b", "e", "g"},
		[]string{"b", "e", "h"},

		[]string{"c", "d", "f"},
		[]string{"c", "d", "g"},
		[]string{"c", "d", "h"},

		[]string{"c", "e", "f"},
		[]string{"c", "e", "g"},
		[]string{"c", "e", "h"},
	}
	fmt.Println(result)

	for i, s := range result {
		if !slices.Equal(s, expected[i]) {
			t.Errorf("テスト失敗")
			break
		}
	}
}
