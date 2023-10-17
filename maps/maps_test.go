
package maps_test

import (
	"testing"
	omwmaps "github.com/sw965/omw/maps"
	"golang.org/x/exp/maps"
)

func TestReverse(t *testing.T) {
	m := map[int]string{0:"4", 1:"3", 2:"2", 3:"1", 4:"0"}
	result := omwmaps.Reverse[map[string]int](m)
	expected := map[string]int{"0":4, "1":3, "2":2, "3":1, "4":0}
	if !maps.Equal(result, expected) {
		t.Errorf("テスト失敗")
	}
}