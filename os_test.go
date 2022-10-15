package omw

import (
	"testing"
)

const (
	TEST_FOLDER = "./test_folder/"
)

func TestListDIr(t *testing.T) {
	expected := []string{
		"a.foo", "b.baz", "baz.b", "foo.a",
		"test_write.txt", "test_write_lines.txt", "text.txt",
	}

	result, err := ListDir(TEST_FOLDER)

	if err != nil {
		panic(err)
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("テスト失敗")
			break
		}
	}
}

func TestWriteText(t *testing.T) {
	data := "Go Python Java Haskell Rust JavaScript"
	WriteText(TEST_FOLDER+"test_write.txt", data)
}

func TestWriteTextLines(t *testing.T) {
	data := []string{"Go Python", "Haskell Rust", "Java C++"}
	WriteTextLines(TEST_FOLDER+"test_write_lines.txt", data)
}
