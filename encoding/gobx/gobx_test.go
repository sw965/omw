package gobx_test

import (
	"errors"
	"github.com/sw965/omw/encoding/gobx"
	"os"
	"path/filepath"
	"testing"
)

type testUser struct {
	Name string
	Age  int
}

func TestSaveAndLoad(t *testing.T) {
	user := testUser{
		Name: "Alice",
		Age:  30,
	}

	//一時的なファイルを作って保存
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "user.gob")
	if err := gobx.Save(user, path); err != nil {
		t.Fatalf("failed to save user: %v", err)
	}

	//読み込み
	got, err := gobx.Load[testUser](path)
	if err != nil {
		t.Fatalf("failed to load user: %v", err)
	}

	if got != user {
		t.Errorf("got %+v, want %+v", got, user)
	}
}

func TestLoad_NotExist(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "non_existent.gob")

	//t.TempDirはファイルまでは作らないので、non_existent.gob が存在せず、エラーが起きるはず
	_, err := gobx.Load[testUser](path)
	if err == nil {
		t.Fatal("expected error for non-existent file, but got nil")
	}

	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected os.ErrNotExist, but got %v", err)
	}
}

func TestSave_InvalidPath(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "no_such_dir", "user.gob")
	user := testUser{Name: "Bob", Age: 25}
	err := gobx.Save(user, path)
	if err == nil {
		t.Fatal("expected error for invalid path, but got nil")
	}

	var pe *os.PathError
	if !errors.As(err, &pe) {
		t.Fatalf("expected *os.PathError, but got %T (%v)", err, err)
	}
}
