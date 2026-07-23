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

func TestSave_ReplacesWithoutTemporaryFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "user.gob")
	if err := gobx.Save(user{Name: "old"}, path); err != nil {
		t.Fatalf("初回保存失敗: %v", err)
	}
	if err := gobx.Save(user{Name: "new", Age: 20}, path); err != nil {
		t.Fatalf("置換保存失敗: %v", err)
	}

	got, err := gobx.Load[user](path)
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
	if len(entries) != 1 || entries[0].Name() != "user.gob" {
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
	path := filepath.Join(tmpDir, "user.gob")
	if err := gobx.Save(u, path); err != nil {
		t.Fatalf("保存失敗: err = %v", err)
	}

	// 読み込み
	got, err := gobx.Load[user](path)
	if err != nil {
		t.Fatalf("読み込み失敗: err = %v", err)
	}

	// 保存したデータと読み込んだデータが一致しているかをチェック
	want := u
	if got != want {
		t.Errorf("データの不一致: got = %+v, want = %+v", got, want)
	}
}

func TestLoad_NotExist(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test.gob")

	// 「test.gob」が存在せず、エラーが起きるはず
	_, err := gobx.Load[user](path)
	if err == nil {
		t.Fatal("読み込み時の想定外の非エラー")
	}

	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("エラー型の不一致: got = %T (%v), want = os.ErrNotExist", err, err)
	}
}
