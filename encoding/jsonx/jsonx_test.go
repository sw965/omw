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

func TestSaveAndLoad(t *testing.T) {
	u := user{
		Name: "Alice",
		Age:  18,
	}

	// 一時ディレクトリと保存先パスを生成する
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "user.json")

	// データを保存する
	if err := jsonx.Save(u, path); err != nil {
		t.Fatalf("保存失敗: err = %v", err)
	}

	// 保存したデータを読み込む
	got, err := jsonx.Load[user](path)
	if err != nil {
		t.Fatalf("読み込み失敗: err = %v", err)
	}

	// 保存したデータと読み込んだデータが一致することを確認する
	want := u
	if got != want {
		t.Errorf("データの不一致: got = %+v, want = %+v", got, want)
	}
}

func TestLoad_BOM(t *testing.T) {
	// 一時ディレクトリと保存先パスを生成する
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "bom.json")

	bom := []byte{0xEF, 0xBB, 0xBF}
	body := []byte(`{"name":"Bob","age":16}`)
	// JSONにBOMを付加する
	withBOM := append(append([]byte{}, bom...), body...)

	// BOM付きのJSONを保存する
	if err := os.WriteFile(path, withBOM, 0o644); err != nil {
		t.Fatalf("書き込み失敗: err = %v", err)
	}

	// BOM付きのJSONを読み込む
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
	// 一時ディレクトリと保存先パスを生成する
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "broken.json")

	// JSONとして不正な内容を保存する
	if err := os.WriteFile(path, []byte(`{"name": "Alice", `), 0o644); err != nil {
		t.Fatalf("書き込み失敗: err = %v", err)
	}

	// 不正なJSONの読み込みを試みる
	_, err := jsonx.Load[user](path)
	if err == nil {
		t.Fatalf("エラーを期待したが、nilが返された")
	}
}

func TestLoad_NotExist(t *testing.T) {
	// 一時ディレクトリと保存先パスを生成する
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test.json")

	// 存在しないファイルの読み込みを試みる
	_, err := jsonx.Load[user](path)
	if err == nil {
		t.Fatalf("エラーを期待したが、nilが返された")
	}

	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("エラー型の不一致: got = %T (%v), want = os.ErrNotExist", err, err)
	}
}

func TestSave_ReplacesAndLeavesNoTemporaryFile(t *testing.T) {
	// 一時ディレクトリと保存先パスを生成する
	dir := t.TempDir()
	path := filepath.Join(dir, "user.json")

	// 旧データを保存する
	if err := jsonx.Save(user{Name: "old"}, path); err != nil {
		t.Fatalf("初回保存失敗: %v", err)
	}

	// 新しいデータで置き換える
	if err := jsonx.Save(user{Name: "new", Age: 20}, path); err != nil {
		t.Fatalf("置換保存失敗: %v", err)
	}

	// 置き換えたデータを読み込む
	got, err := jsonx.Load[user](path)
	if err != nil {
		t.Fatalf("読み込み失敗: %v", err)
	}
	if want := (user{Name: "new", Age: 20}); got != want {
		t.Fatalf("データの不一致: got = %+v want = %+v", got, want)
	}

	assertSingleFile(t, dir, "user.json")
}
