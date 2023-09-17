package strings_test

import (
	"testing"
	"github.com/sw965/omw"
	omwstrings "github.com/sw965/omw/strings"
)

func TestReplace(t *testing.T) {
	result := omwstrings.Replace(omw.JSON_EXTENSION, "", 1)("golang.json")
	expected := "golang"
	if result != expected {
		t.Errorf("テスト失敗")
	}
}