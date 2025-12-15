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
		Age:  18,
	}

	//一時的なファイルを作って保存
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test.gob")
	if err := gobx.Save(user, path); err != nil {
		t.Fatalf("予期せぬエラー: %v", err)
	}

	//読み込み
	got, err := gobx.Load[testUser](path)
	if err != nil {
		t.Fatalf("予期せぬエラー: %v", err)
	}

	if got != user {
		t.Errorf("want: %v, got: %v", user, got)
	}
}

func TestLoad_NotExist(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "nope.gob")

	//t.TempDirはファイルまでは作らないので、nope.gobが存在せず、エラーが起きるはず
	_, err := gobx.Load[testUser](path)
	if err == nil {
		t.Fatal("エラーを期待したが、nilが返された")
	}

	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("期待されるエラー型と異なります。want: %T, got: %T", os.ErrNotExist, err)
	}
}

func TestSave_InvalidPath(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "no_such_dir", "x.gob")
	user := testUser{Name: "Bob", Age: 15}
	err := gobx.Save(user, path)
	if err == nil {
		t.Fatal("エラーを期待したが、nilが返された")
	}

	var pe *os.PathError
	if !errors.As(err, &pe) {
		t.Fatalf("期待されるエラー型と異なります。want *os.PathError, got: %T (%v)", err, err)
	}
}
