package bitsx_test

import (
	"bytes"
	"encoding/gob"
	"math/rand/v2"
	"testing"

	"github.com/sw965/omw/mathx/bitsx"
)

func TestMatrixGobRoundTrip(t *testing.T) {
	rng := rand.New(rand.NewPCG(3, 4))
	want, err := bitsx.NewRandMatrix(3, 130, rng)
	if err != nil {
		t.Fatalf("nilを期待したが、エラーが返された: %v", err)
	}

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(want); err != nil {
		t.Fatalf("エンコード失敗: %v", err)
	}

	var got bitsx.Matrix
	if err := gob.NewDecoder(&buf).Decode(&got); err != nil {
		t.Fatalf("デコード失敗: %v", err)
	}

	if !want.Equal(&got) {
		t.Error("gobの往復で内容が変化した")
	}
}
