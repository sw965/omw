package os_test

import (
	osmw "github.com/sw965/omw/os"
	"testing"
	"os"
	"golang.org/x/exp/slices"
)

const (
	TEST_FOLDER = "/test_folder/"
)

func TestDirEntries_Names(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	path := dir + TEST_FOLDER
	dirEntries, err := osmw.NewDirEntries(path)
	if err != nil {
		panic(err)
	}

	result := dirEntries.Names()
	expected := []string{"go.go", "haskell.hs", "python.py"}
	if !slices.Equal(result, expected) {
		t.Errorf("テスト失敗")
	}
}

