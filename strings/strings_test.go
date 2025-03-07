package strings_test

import (
	omwstrings "github.com/sw965/omw/strings"
	"testing"
)

func TestPadLeft(t *testing.T) {
	s := "あいうえお"
	ret, err := omwstrings.PadLeft(s, " ", 10)
	if err != nil {
		panic(err)
	}
	expected := "     あいうえお"
	if ret != expected {
		t.Errorf("テスト失敗")
	}
}

func TestPadRight(t *testing.T) {
	s := "かきくけこ"
	ret, err := omwstrings.PadRight(s, "あ", 10)
	if err != nil {
		panic(err)
	}
	expected := "かきくけこあああああ"
	if ret != expected {
		t.Errorf("テスト失敗")
	}
}
