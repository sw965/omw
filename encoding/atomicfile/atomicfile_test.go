package atomicfile_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sw965/omw/encoding/atomicfile"
)

func assertSingleFile(t *testing.T, dir string, expectedName string) {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ディレクトリ読み込み失敗: %v", err)
	}
	if len(entries) != 1 || entries[0].Name() != expectedName {
		t.Fatalf("一時ファイルが残っている、またはファイル構成が不正: %v", entries)
	}
}

func assertEmptyDir(t *testing.T, dir string) {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ディレクトリ読み込み失敗: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("ファイルが残っている: %v", entries)
	}
}

func TestWriteFile(t *testing.T) {
	t.Run("新規保存", func(t *testing.T) {
		// 一時ディレクトリと保存先パスを生成する
		dir := t.TempDir()
		path := filepath.Join(dir, "data.txt")

		// ファイルを新規作成
		data := []byte("new")
		if err := atomicfile.WriteFile(path, data, 0o640); err != nil {
			t.Fatalf("保存失敗: %v", err)
		}

		// 新規作成したファイルを読み込む
		got, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("読み込み失敗: %v", err)
		}

		if string(got) != "new" {
			t.Fatalf("内容の不一致: got=%q want=%q", got, "new")
		}

		assertSingleFile(t, dir, "data.txt")
	})

	t.Run("既存ファイルの置換", func(t *testing.T) {
		// 一時ディレクトリと保存先パスを生成する
		dir := t.TempDir()
		path := filepath.Join(dir, "data.txt")

		// 置換前の旧ファイルを作成しておく
		oldData := []byte("old")
		if err := os.WriteFile(path, oldData, 0o600); err != nil {
			t.Fatalf("準備失敗: %v", err)
		}

		// 旧ファイルを置き換える
		newData := []byte("new")
		if err := atomicfile.WriteFile(path, newData, 0o640); err != nil {
			t.Fatalf("置換失敗: %v", err)
		}

		// 置き換えたファイルを読み込む
		got, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("読み込み失敗: %v", err)
		}

		if string(got) != "new" {
			t.Fatalf("内容の不一致: got=%q want=%q", got, "new")
		}

		assertSingleFile(t, dir, "data.txt")
	})

	t.Run("保存先ディレクトリなし", func(t *testing.T) {
		dir := t.TempDir()
		// 存在しないディレクトリを含む保存先パスを生成する
		path := filepath.Join(dir, "missing", "data.txt")

		// ファイルの保存を試みる
		if err := atomicfile.WriteFile(path, []byte("new"), 0o640); err == nil {
			t.Fatal("エラーを期待したが、nilが返された")
		}

		assertEmptyDir(t, dir)
	})
}

func TestWriteFrom(t *testing.T) {
	t.Run("新規保存", func(t *testing.T) {
		// 一時ディレクトリと保存先パスを生成する
		dir := t.TempDir()
		path := filepath.Join(dir, "data.txt")

		// ファイルを新規作成
		if err := atomicfile.WriteFrom(path, strings.NewReader("streamed"), 0o640); err != nil {
			t.Fatalf("保存失敗: %v", err)
		}

		// ファイルを読み込む
		got, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("読み込み失敗: %v", err)
		}

		if string(got) != "streamed" {
			t.Fatalf("内容の不一致: got=%q want=%q", got, "streamed")
		}

		assertSingleFile(t, dir, "data.txt")
	})

	t.Run("既存ファイルの置換", func(t *testing.T) {
		// 一時ディレクトリと保存先パスを生成する
		dir := t.TempDir()
		path := filepath.Join(dir, "data.txt")

		// 置換前の旧ファイルを作成しておく
		if err := os.WriteFile(path, []byte("old-long-content"), 0o600); err != nil {
			t.Fatalf("準備失敗: %v", err)
		}

		// 旧ファイルを置き換える
		if err := atomicfile.WriteFrom(path, strings.NewReader("streamed"), 0o640); err != nil {
			t.Fatalf("置換失敗: %v", err)
		}

		got, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("読み込み失敗: %v", err)
		}

		if string(got) != "streamed" {
			t.Fatalf("内容の不一致: got=%q want=%q", got, "streamed")
		}

		assertSingleFile(t, dir, "data.txt")
	})

	t.Run("保存先ディレクトリなし", func(t *testing.T) {
		dir := t.TempDir()
		// 存在しないディレクトリを含む保存先パスを生成する
		path := filepath.Join(dir, "missing", "data.txt")

		// ファイルの保存を試みる
		if err := atomicfile.WriteFrom(path, strings.NewReader("streamed"), 0o640); err == nil {
			t.Fatal("エラーを期待したが、nilが返された")
		}

		assertEmptyDir(t, dir)
	})

	t.Run("nilのio.Readerを渡す", func(t *testing.T) {
		// 一時ディレクトリと保存先パスを生成する
		dir := t.TempDir()
		path := filepath.Join(dir, "data.txt")

		// nilのio.Readerを渡す
		if err := atomicfile.WriteFrom(path, nil, 0o640); err == nil {
			t.Fatal("エラーを期待したが、nilが返された")
		}
		assertEmptyDir(t, dir)
	})
}
