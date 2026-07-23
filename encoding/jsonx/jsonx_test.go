package jsonx_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/sw965/omw/encoding/jsonx"
)

type user struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestSave_ReplacesWithoutTemporaryFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "user.json")
	if err := jsonx.Save(user{Name: "old"}, path); err != nil {
		t.Fatalf("初回保存失敗: %v", err)
	}
	if err := jsonx.Save(user{Name: "new", Age: 20}, path); err != nil {
		t.Fatalf("置換保存失敗: %v", err)
	}

	got, err := jsonx.Load[user](path)
	if err != nil {
		t.Fatalf("読み込み失敗: %v", err)
	}
	if want := (user{Name: "new", Age: 20}); got != want {
		t.Fatalf("データの不一致: got=%+v want=%+v", got, want)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ディレクトリ読み込み失敗: %v", err)
	}
	if len(entries) != 1 || entries[0].Name() != "user.json" {
		t.Fatalf("一時ファイルが残っています: %v", entries)
	}
}

func TestSaveAndLoad(t *testing.T) {
	u := user{
		Name: "Alice",
		Age:  18,
	}

	// 一時的なファイルを作って保存
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "user.json")
	if err := jsonx.Save(u, path); err != nil {
		t.Fatalf("保存失敗: err = %v", err)
	}

	// 読み込み
	got, err := jsonx.Load[user](path)
	if err != nil {
		t.Fatalf("読み込み失敗: err = %v", err)
	}

	// 保存したデータと読み込んだデータが一致しているかをチェック
	want := u
	if got != want {
		t.Errorf("データの不一致: got = %+v, want = %+v", got, want)
	}
}

func TestLoad_BOM(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "bom.json")

	bom := []byte{0xEF, 0xBB, 0xBF}
	body := []byte(`{"name":"Bob","age":16}`)
	// BOMを付ける
	withBOM := append(append([]byte{}, bom...), body...)

	// BOMを付けたデータを保存
	if err := os.WriteFile(path, withBOM, 0644); err != nil {
		t.Fatalf("書き込み失敗: err = %v", err)
	}

	got, err := jsonx.Load[user](path)
	if err != nil {
		t.Fatalf("読み込み失敗: err = %v", err)
	}

	want := user{Name: "Bob", Age: 16}
	if got != want {
		t.Errorf("データの不一致: got = %+v, want = %+v", got, want)
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "broken.json")

	// JSONとして不正な内容を保存
	if err := os.WriteFile(path, []byte(`{"name": "Alice", `), 0644); err != nil {
		t.Fatalf("書き込み失敗: err = %v", err)
	}

	_, err := jsonx.Load[user](path)
	if err == nil {
		t.Fatal("読み込み時の想定外の非エラー")
	}
}

func TestLoad_NotExist(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test.json")

	// 「test.json」が存在せず、エラーが起きるはず
	_, err := jsonx.Load[user](path)
	if err == nil {
		t.Fatal("読み込み時の想定外の非エラー")
	}

	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("エラー型の不一致: got = %T (%v), want = os.ErrNotExist", err, err)
	}
}
