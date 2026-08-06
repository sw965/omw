package bitsx_test

import (
	"bytes"
	"encoding/gob"
	"math"
	"math/rand/v2"
	"testing"

	"github.com/sw965/omw/mathx/bitsx"
)

func TestNewZerosMatrix(t *testing.T) {
	t.Run("正常", func(t *testing.T) {
		m, err := bitsx.NewZerosMatrix(2, 100)
		if err != nil {
			t.Fatalf("予期せぬエラー: %v", err)
		}
		if m.Rows() != 2 || m.Cols() != 100 {
			t.Errorf("形状の不一致: got = (%d, %d), want = (2, 100)", m.Rows(), m.Cols())
		}
		if got := m.OnesCount(); got != 0 {
			t.Errorf("OnesCountの不一致: got = %d, want = 0", got)
		}
	})

	t.Run("異常_rowsが0以下", func(t *testing.T) {
		if _, err := bitsx.NewZerosMatrix(0, 10); err == nil {
			t.Fatal("エラーを期待したが、nilが返された")
		}
	})

	t.Run("異常_colsが0以下", func(t *testing.T) {
		if _, err := bitsx.NewZerosMatrix(10, 0); err == nil {
			t.Fatal("エラーを期待したが、nilが返された")
		}
	})

	t.Run("異常_colsの桁あふれ", func(t *testing.T) {
		if _, err := bitsx.NewZerosMatrix(1, math.MaxInt); err == nil {
			t.Fatal("エラーを期待したが、nilが返された")
		}
	})
}

func TestNewOnesMatrix(t *testing.T) {
	// 列数が64の倍数ではない場合、端数ビットは0に保たれ、
	// OnesCountは論理的なビット数(rows*cols)と一致するはず
	m, err := bitsx.NewOnesMatrix(3, 100)
	if err != nil {
		t.Fatalf("予期せぬエラー: %v", err)
	}
	want := 3 * 100
	if got := m.OnesCount(); got != want {
		t.Errorf("OnesCountの不一致: got = %d, want = %d", got, want)
	}
}

func TestMatrixBitOperations(t *testing.T) {
	m, err := bitsx.NewZerosMatrix(2, 70)
	if err != nil {
		t.Fatalf("予期せぬエラー: %v", err)
	}

	t.Run("正常_SetとBit", func(t *testing.T) {
		if err := m.Set(1, 69); err != nil {
			t.Fatalf("予期せぬエラー: %v", err)
		}
		got, err := m.Bit(1, 69)
		if err != nil {
			t.Fatalf("予期せぬエラー: %v", err)
		}
		if got != 1 {
			t.Errorf("値の不一致: got = %d, want = 1", got)
		}
	})

	t.Run("正常_Clear", func(t *testing.T) {
		if err := m.Clear(1, 69); err != nil {
			t.Fatalf("予期せぬエラー: %v", err)
		}
		got, err := m.Bit(1, 69)
		if err != nil {
			t.Fatalf("予期せぬエラー: %v", err)
		}
		if got != 0 {
			t.Errorf("値の不一致: got = %d, want = 0", got)
		}
	})

	t.Run("正常_Toggle", func(t *testing.T) {
		if err := m.Toggle(0, 0); err != nil {
			t.Fatalf("予期せぬエラー: %v", err)
		}
		got, err := m.Bit(0, 0)
		if err != nil {
			t.Fatalf("予期せぬエラー: %v", err)
		}
		if got != 1 {
			t.Errorf("値の不一致: got = %d, want = 1", got)
		}

		if err := m.Toggle(0, 0); err != nil {
			t.Fatalf("予期せぬエラー: %v", err)
		}
		got, err = m.Bit(0, 0)
		if err != nil {
			t.Fatalf("予期せぬエラー: %v", err)
		}
		if got != 0 {
			t.Errorf("値の不一致: got = %d, want = 0", got)
		}
	})

	t.Run("異常_範囲外", func(t *testing.T) {
		if _, err := m.Bit(-1, 0); err == nil {
			t.Fatal("エラーを期待したが、nilが返された")
		}
		if _, err := m.Bit(0, 70); err == nil {
			t.Fatal("エラーを期待したが、nilが返された")
		}
		if err := m.Set(2, 0); err == nil {
			t.Fatal("エラーを期待したが、nilが返された")
		}
	})
}

func TestMatrixTranspose(t *testing.T) {
	rng := rand.New(rand.NewPCG(7, 8))

	// 転置の性質: (1) 形状が入れ替わる (2) 全ての(r,c)でm[r][c] == mT[c][r] (3) 二回転置すると元に戻る
	shapes := []struct{ rows, cols int }{
		{1, 1},
		{3, 70},
		{64, 64},
		{100, 130},
	}

	for _, s := range shapes {
		m, err := bitsx.NewRandMatrix(s.rows, s.cols, 0, rng)
		if err != nil {
			t.Fatalf("予期せぬエラー: %v", err)
		}

		mT, err := m.Transpose()
		if err != nil {
			t.Fatalf("予期せぬエラー: %v", err)
		}

		if mT.Rows() != s.cols || mT.Cols() != s.rows {
			t.Fatalf("転置後の形状の不一致: got = (%d, %d), want = (%d, %d)", mT.Rows(), mT.Cols(), s.cols, s.rows)
		}

		for r := 0; r < s.rows; r++ {
			for c := 0; c < s.cols; c++ {
				orig, err := m.Bit(r, c)
				if err != nil {
					t.Fatalf("予期せぬエラー: %v", err)
				}
				transposed, err := mT.Bit(c, r)
				if err != nil {
					t.Fatalf("予期せぬエラー: %v", err)
				}
				if orig != transposed {
					t.Fatalf("shape (%d, %d): (%d, %d)のビットが一致しない", s.rows, s.cols, r, c)
				}
			}
		}

		mTT, err := mT.Transpose()
		if err != nil {
			t.Fatalf("予期せぬエラー: %v", err)
		}
		if !m.Equal(mTT) {
			t.Fatalf("shape (%d, %d): 二回転置しても元に戻らない", s.rows, s.cols)
		}
	}
}

func TestMatrixHammingDistance(t *testing.T) {
	t.Run("正常_自身との距離は0", func(t *testing.T) {
		rng := rand.New(rand.NewPCG(9, 10))
		m, err := bitsx.NewRandMatrix(2, 100, 0, rng)
		if err != nil {
			t.Fatalf("予期せぬエラー: %v", err)
		}
		got, err := m.HammingDistance(m)
		if err != nil {
			t.Fatalf("予期せぬエラー: %v", err)
		}
		if got != 0 {
			t.Errorf("値の不一致: got = %d, want = 0", got)
		}
	})

	t.Run("正常_全ビット反転との距離は総ビット数", func(t *testing.T) {
		zeros, err := bitsx.NewZerosMatrix(2, 100)
		if err != nil {
			t.Fatalf("予期せぬエラー: %v", err)
		}
		ones, err := bitsx.NewOnesMatrix(2, 100)
		if err != nil {
			t.Fatalf("予期せぬエラー: %v", err)
		}
		got, err := zeros.HammingDistance(ones)
		if err != nil {
			t.Fatalf("予期せぬエラー: %v", err)
		}
		if got != 200 {
			t.Errorf("値の不一致: got = %d, want = 200", got)
		}
	})

	t.Run("異常_形状不一致", func(t *testing.T) {
		a, err := bitsx.NewZerosMatrix(1, 10)
		if err != nil {
			t.Fatalf("予期せぬエラー: %v", err)
		}
		b, err := bitsx.NewZerosMatrix(1, 20)
		if err != nil {
			t.Fatalf("予期せぬエラー: %v", err)
		}
		if _, err := a.HammingDistance(b); err == nil {
			t.Fatal("エラーを期待したが、nilが返された")
		}
	})

	t.Run("正常_ゼロ値同士は0", func(t *testing.T) {
		// Dataが空の場合、AVX-512版は&Data[0]を取れない。
		// pure Go版と同じく0を返すこと。
		var a, b bitsx.Matrix
		got, err := a.HammingDistance(&b)
		if err != nil {
			t.Fatalf("予期せぬエラー: %v", err)
		}
		if got != 0 {
			t.Errorf("値の不一致: got = %d, want = 0", got)
		}
	})
}

func TestMatricesValidation(t *testing.T) {
	rng := rand.New(rand.NewPCG(11, 12))

	t.Run("異常_NewETFMatricesのnが2未満", func(t *testing.T) {
		if _, err := bitsx.NewETFMatrices(1, 4, 8, 10, rng); err == nil {
			t.Fatal("エラーを期待したが、nilが返された")
		}
	})

	t.Run("異常_NewRFFMatricesのnが2未満", func(t *testing.T) {
		if _, err := bitsx.NewRFFMatrices(1, 4, 8, 1.0, rng); err == nil {
			t.Fatal("エラーを期待したが、nilが返された")
		}
	})

	t.Run("異常_NewRFFMatricesのrowsが0以下", func(t *testing.T) {
		if _, err := bitsx.NewRFFMatrices(3, 0, 8, 1.0, rng); err == nil {
			t.Fatal("エラーを期待したが、nilが返された")
		}
		if _, err := bitsx.NewRFFMatrices(3, -4, 8, 1.0, rng); err == nil {
			t.Fatal("エラーを期待したが、nilが返された")
		}
	})

	t.Run("異常_NewRFFMatricesのcolsが0以下", func(t *testing.T) {
		if _, err := bitsx.NewRFFMatrices(3, 4, 0, 1.0, rng); err == nil {
			t.Fatal("エラーを期待したが、nilが返された")
		}
		if _, err := bitsx.NewRFFMatrices(3, 4, -8, 1.0, rng); err == nil {
			t.Fatal("エラーを期待したが、nilが返された")
		}
	})

	t.Run("異常_NewThermometerMatricesのnが2未満", func(t *testing.T) {
		if _, err := bitsx.NewThermometerMatrices(1, 4, 8); err == nil {
			t.Fatal("エラーを期待したが、nilが返された")
		}
	})

	t.Run("異常_ETFCostの行列数が2未満", func(t *testing.T) {
		m, err := bitsx.NewZerosMatrix(4, 8)
		if err != nil {
			t.Fatalf("予期せぬエラー: %v", err)
		}
		ms := bitsx.Matrices{m}
		if _, err := ms.ETFCost(); err == nil {
			t.Fatal("エラーを期待したが、nilが返された")
		}
	})

	t.Run("正常_生成された行列の個数と形状", func(t *testing.T) {
		ms, err := bitsx.NewETFMatrices(3, 2, 70, 50, rng)
		if err != nil {
			t.Fatalf("予期せぬエラー: %v", err)
		}
		if len(ms) != 3 {
			t.Fatalf("個数の不一致: got = %d, want = 3", len(ms))
		}
		for i, m := range ms {
			if m.Rows() != 2 || m.Cols() != 70 {
				t.Errorf("ms[%d]の形状の不一致: got = (%d, %d), want = (2, 70)", i, m.Rows(), m.Cols())
			}
		}
	})
}

func TestMatrixWord(t *testing.T) {
	// cols=70 なら Stride()=2 なので、rows=2の内部データ長は4
	m, err := bitsx.NewOnesMatrix(2, 70)
	if err != nil {
		t.Fatalf("予期せぬエラー: %v", err)
	}

	t.Run("正常_有効なidx", func(t *testing.T) {
		got, err := m.Word(0)
		if err != nil {
			t.Fatalf("予期せぬエラー: %v", err)
		}
		if got != ^uint64(0) {
			t.Errorf("値の不一致: got = %#x, want = %#x", got, ^uint64(0))
		}
	})

	t.Run("異常_負のidx", func(t *testing.T) {
		if _, err := m.Word(-1); err == nil {
			t.Fatal("エラーを期待したが、nilが返された")
		}
	})

	t.Run("異常_範囲外のidx", func(t *testing.T) {
		if _, err := m.Word(4); err == nil {
			t.Fatal("エラーを期待したが、nilが返された")
		}
	})
}

func TestMatrixSetWord(t *testing.T) {
	// cols=70 なら Stride()=2。idx=0は非tail語、idx=1はtail語(有効ビットは70-64=6)
	newZeros := func(t *testing.T) *bitsx.Matrix {
		t.Helper()
		m, err := bitsx.NewZerosMatrix(1, 70)
		if err != nil {
			t.Fatalf("予期せぬエラー: %v", err)
		}
		return m
	}

	t.Run("正常_非tail語はそのまま書き込む", func(t *testing.T) {
		m := newZeros(t)
		if err := m.SetWord(0, ^uint64(0)); err != nil {
			t.Fatalf("予期せぬエラー: %v", err)
		}
		got, err := m.Word(0)
		if err != nil {
			t.Fatalf("予期せぬエラー: %v", err)
		}
		if got != ^uint64(0) {
			t.Errorf("値の不一致: got = %#x, want = %#x", got, ^uint64(0))
		}
	})

	t.Run("正常_tail語は端数ビットが0にマスクされる", func(t *testing.T) {
		m := newZeros(t)
		if err := m.SetWord(1, ^uint64(0)); err != nil {
			t.Fatalf("予期せぬエラー: %v", err)
		}
		got, err := m.Word(1)
		if err != nil {
			t.Fatalf("予期せぬエラー: %v", err)
		}
		if want := m.TailMask(); got != want {
			t.Errorf("値の不一致: got = %#x, want = %#x", got, want)
		}
	})

	t.Run("異常_負のidx", func(t *testing.T) {
		if err := newZeros(t).SetWord(-1, 0); err == nil {
			t.Fatal("エラーを期待したが、nilが返された")
		}
	})

	t.Run("異常_範囲外のidx", func(t *testing.T) {
		if err := newZeros(t).SetWord(2, 0); err == nil {
			t.Fatal("エラーを期待したが、nilが返された")
		}
	})
}

func TestMatrixEqual(t *testing.T) {
	t.Run("正常_同じ内容ならtrue", func(t *testing.T) {
		rng := rand.New(rand.NewPCG(1, 2))
		a, err := bitsx.NewRandMatrix(3, 100, 0, rng)
		if err != nil {
			t.Fatalf("予期せぬエラー: %v", err)
		}
		if !a.Equal(a.Clone()) {
			t.Error("同じ内容なのにfalseが返された")
		}
	})

	t.Run("異常_形状が違えばfalse", func(t *testing.T) {
		a, err := bitsx.NewZerosMatrix(3, 100)
		if err != nil {
			t.Fatalf("予期せぬエラー: %v", err)
		}
		b, err := bitsx.NewZerosMatrix(3, 101)
		if err != nil {
			t.Fatalf("予期せぬエラー: %v", err)
		}
		if a.Equal(b) {
			t.Error("形状が違うのにtrueが返された")
		}
	})

	t.Run("異常_内容が違えばfalse", func(t *testing.T) {
		a, err := bitsx.NewZerosMatrix(3, 100)
		if err != nil {
			t.Fatalf("予期せぬエラー: %v", err)
		}
		b := a.Clone()
		if err := b.Set(0, 0); err != nil {
			t.Fatalf("予期せぬエラー: %v", err)
		}
		if a.Equal(b) {
			t.Error("内容が違うのにtrueが返された")
		}
	})
}

func TestMatrixGobRoundTrip(t *testing.T) {
	rng := rand.New(rand.NewPCG(3, 4))
	want, err := bitsx.NewRandMatrix(3, 130, 0, rng)
	if err != nil {
		t.Fatalf("予期せぬエラー: %v", err)
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

const (
	benchTransposeRows = 768
	benchTransposeCols = 768
)

// 環境: AMD Ryzen 5 8600G
// 形状: 768行 × 768列
// 期待: おおよそ 30〜36 µs/op
// ※ 上記のコメントは勝手に変更・削除してはならない
func BenchmarkMatrixTranspose(b *testing.B) {
	rng := rand.New(rand.NewPCG(11, 12))
	m, err := bitsx.NewRandMatrix(benchTransposeRows, benchTransposeCols, 0, rng)
	if err != nil {
		b.Fatal(err)
	}

	for b.Loop() {
		if _, err := m.Transpose(); err != nil {
			b.Fatal(err)
		}
	}
}
