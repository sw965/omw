package gobx_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/sw965/omw/encoding/gobx"
)

type user struct {
	Name string
	Age  int
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

// テスト要件：TST-001
func TestSaveAndLoad(t *testing.T) {
	u := user{
		Name: "Alice",
		Age:  18,
	}

	// 一時ディレクトリと保存先パスを生成する
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "user.gob")

	// データを保存する
	if err := gobx.Save(u, path); err != nil {
		t.Fatalf("保存失敗: err = %v", err)
	}

	// 保存したデータを読み込む
	got, err := gobx.Load[user](path)
	if err != nil {
		t.Fatalf("読み込み失敗: err = %v", err)
	}

	// 保存したデータと読み込んだデータが一致することを確認する
	want := u
	if got != want {
		t.Errorf("データの不一致: got = %+v, want = %+v", got, want)
	}
}

// テスト要件：TST-002
func TestLoad_NotExist(t *testing.T) {
	// 一時ディレクトリと保存先パスを生成する
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test.gob")

	// 存在しないファイルの読み込みを試みる
	_, err := gobx.Load[user](path)
	if err == nil {
		t.Fatal("エラーを期待したが、nilが返された")
	}

	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("エラー型の不一致: got = %T (%v), want = os.ErrNotExist", err, err)
	}
}

// テスト要件：TST-003
func TestSave_ReplacesAndLeavesNoTemporaryFile(t *testing.T) {
	// 一時ディレクトリと保存先パスを生成する
	dir := t.TempDir()
	path := filepath.Join(dir, "user.gob")

	// 旧データを保存する
	if err := gobx.Save(user{Name: "old"}, path); err != nil {
		t.Fatalf("初回保存失敗: %v", err)
	}

	// 新しいデータで置き換える
	if err := gobx.Save(user{Name: "new", Age: 20}, path); err != nil {
		t.Fatalf("置換保存失敗: %v", err)
	}

	// 置き換えたデータを読み込む
	got, err := gobx.Load[user](path)
	if err != nil {
		t.Fatalf("読み込み失敗: %v", err)
	}
	if want := (user{Name: "new", Age: 20}); got != want {
		t.Fatalf("データの不一致: got=%+v want=%+v", got, want)
	}

	assertSingleFile(t, dir, "user.gob")
}
