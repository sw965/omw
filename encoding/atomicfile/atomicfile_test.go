package atomicfile_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sw965/omw/encoding/atomicfile"
)

func TestWriteFile(t *testing.T) {
	t.Run("新規保存", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "data.txt")
		if err := atomicfile.WriteFile(path, []byte("new"), 0o640); err != nil {
			t.Fatalf("保存失敗: %v", err)
		}
		got, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("読み込み失敗: %v", err)
		}
		if string(got) != "new" {
			t.Fatalf("内容の不一致: got=%q want=%q", got, "new")
		}
	})

	t.Run("既存ファイルの置換", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "data.txt")
		if err := os.WriteFile(path, []byte("old-long-content"), 0o600); err != nil {
			t.Fatalf("準備失敗: %v", err)
		}
		if err := atomicfile.WriteFile(path, []byte("new"), 0o640); err != nil {
			t.Fatalf("置換失敗: %v", err)
		}
		got, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("読み込み失敗: %v", err)
		}
		if string(got) != "new" {
			t.Fatalf("内容の不一致: got=%q want=%q", got, "new")
		}
		entries, err := os.ReadDir(dir)
		if err != nil {
			t.Fatalf("ディレクトリ読み込み失敗: %v", err)
		}
		if len(entries) != 1 || entries[0].Name() != "data.txt" {
			t.Fatalf("一時ファイルが残っています: %v", entries)
		}
	})

	t.Run("保存先ディレクトリなし", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "missing", "data.txt")
		if err := atomicfile.WriteFile(path, []byte("new"), 0o640); err == nil {
			t.Fatal("エラーを期待したが、nilが返された")
		}
		entries, err := os.ReadDir(dir)
		if err != nil {
			t.Fatalf("ディレクトリ読み込み失敗: %v", err)
		}
		if len(entries) != 0 {
			t.Fatalf("ファイルが残っています: %v", entries)
		}
	})
}

func TestWriteFrom(t *testing.T) {
	t.Run("新規保存", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "data.txt")
		if err := atomicfile.WriteFrom(path, strings.NewReader("streamed"), 0o640); err != nil {
			t.Fatalf("保存失敗: %v", err)
		}
		got, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("読み込み失敗: %v", err)
		}
		if string(got) != "streamed" {
			t.Fatalf("内容の不一致: got=%q want=%q", got, "streamed")
		}
	})

	t.Run("既存ファイルの置換_一時ファイルを残さない", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "data.txt")
		if err := os.WriteFile(path, []byte("old-long-content"), 0o600); err != nil {
			t.Fatalf("準備失敗: %v", err)
		}
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
		entries, err := os.ReadDir(dir)
		if err != nil {
			t.Fatalf("ディレクトリ読み込み失敗: %v", err)
		}
		if len(entries) != 1 || entries[0].Name() != "data.txt" {
			t.Fatalf("一時ファイルが残っています: %v", entries)
		}
	})

	t.Run("保存先ディレクトリなし", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "missing", "data.txt")
		if err := atomicfile.WriteFrom(path, strings.NewReader("streamed"), 0o640); err == nil {
			t.Fatal("エラーを期待したが、nilが返された")
		}
		entries, err := os.ReadDir(dir)
		if err != nil {
			t.Fatalf("ディレクトリ読み込み失敗: %v", err)
		}
		if len(entries) != 0 {
			t.Fatalf("ファイルが残っています: %v", entries)
		}
	})

	t.Run("異常_nilのio.Reader", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "data.txt")
		if err := atomicfile.WriteFrom(path, nil, 0o640); err == nil {
			t.Fatal("エラーを期待したが、nilが返された")
		}
		entries, err := os.ReadDir(dir)
		if err != nil {
			t.Fatalf("ディレクトリ読み込み失敗: %v", err)
		}
		if len(entries) != 0 {
			t.Fatalf("一時ファイルが残っています: %v", entries)
		}
	})
}
