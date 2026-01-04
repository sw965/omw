package jsonx_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/sw965/omw/encoding/jsonx"
)

type testUser struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestSaveAndLoad(t *testing.T) {
	user := testUser{
		Name: "Alice",
		Age:  30,
	}

	//一時的なファイルを作って保存
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "user.json")
	if err := jsonx.Save(user, path); err != nil {
		t.Fatalf("failed to save user: %v", err)
	}

	//読み込み
	got, err := jsonx.Load[testUser](path)
	if err != nil {
		t.Fatalf("failed to load user: %v", err)
	}

	want := user
	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestLoad_BOM(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "bom.json")

	bom := []byte{0xEF, 0xBB, 0xBF}
	body := []byte(`{"name":"A","age":1}`)
	// BOMを付ける
	withBOM := append(append([]byte{}, bom...), body...)

	if err := os.WriteFile(path, withBOM, 0644); err != nil {
		t.Fatalf("failed to write file with BOM: %v", err)
	}

	got, err := jsonx.Load[testUser](path)
	if err != nil {
		t.Fatalf("failed to load file with BOM: %v", err)
	}

	want := testUser{Name: "A", Age: 1}
	if got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestLoad_NotExist(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "non_existent.json")

	//t.TempDirはファイルまでは作らないので、「non_existent.json」が存在せず、エラーが起きるはず
	_, err := jsonx.Load[testUser](path)
	if err == nil {
		t.Fatal("expected error for non-existent file, but got nil")
	}

	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected os.ErrNotExist, but got %v", err)
	}
}