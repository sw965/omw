package jsonx_test

import (
	"github.com/sw965/omw/encoding/jsonx"
	"path/filepath"
	"testing"
	"os"
	"errors"
)

type testUser struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestSaveAndLoad(t *testing.T) {
	user := testUser{
		Name: "Bob",
		Age:  24,
	}

	//一時的なファイルを作って保存
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test.json")
	if err := jsonx.Save(user, path); err != nil {
		t.Fatalf("予期せぬエラー: %v", err)
	}

	//読み込み
	got, err := jsonx.Load[testUser](path)
	if err != nil {
		t.Fatalf("予期せぬエラー: %v", err)
	}

	want := user
	if got != want {
		t.Errorf("want: %v, got: %v", want, got)
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
		t.Fatal(err)
	}

	got, err := jsonx.Load[testUser](path)
	if err != nil {
		t.Fatalf("予期せぬエラーが発生した: %v", err)
	}

	want := testUser{Name: "A", Age: 1}
	if got != want {
		t.Fatalf("want: %+v, got: %+v", want, got)
	}
}

func TestLoad_NotExist(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "nope.json")

	//t.TempDirはファイルまでは作らないので、nope.gobが存在せず、エラーが起きるはず
	_, err := jsonx.Load[testUser](path)
	if err == nil {
		t.Fatal("エラーを期待したが、nilが返された")
	}

	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("期待されるエラー型と異なります。want: %T, got: %T", os.ErrNotExist, err)
	}
}