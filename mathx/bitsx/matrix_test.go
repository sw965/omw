package bitsx_test

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"

	"github.com/sw965/omw/mathx/bitsx"
)

func TestNewZerosMatrix(t *testing.T) {
	t.Run("正常", func(t *testing.T) {
		m, err := bitsx.NewZerosMatrix(2, 100)
		if err != nil {
			t.Fatalf("予期せぬエラー: %v", err)
		}
		if m.Rows != 2 || m.Cols != 100 {
			t.Errorf("形状の不一致: got = (%d, %d), want = (2, 100)", m.Rows, m.Cols)
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

	// rows*stride が桁あふれすると、Rowsに対してDataが短い行列が出来てしまう
	t.Run("異常_rowsとcolsの積が桁あふれ", func(t *testing.T) {
		// cols=200 なら Stride()=4 なので、rows*stride は 2^64+4 に折り返して 4 になる
		if _, err := bitsx.NewZerosMatrix(1<<62+1, 200); err == nil {
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

		if mT.Rows != s.cols || mT.Cols != s.rows {
			t.Fatalf("転置後の形状の不一致: got = (%d, %d), want = (%d, %d)", mT.Rows, mT.Cols, s.cols, s.rows)
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
		if !slices.Equal(m.Data, mTT.Data) {
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

// Dot / DotTernary は境界チェックを持たないアセンブリ実装へ生ポインタを渡す為、
// 不整合な Matrix を弾けないと範囲外アクセスになる。
// 検査が壊れた場合、このテストは失敗ではなくクラッシュする可能性がある。
func TestMatrixValidate(t *testing.T) {
	// cols = 100 なら Stride() = 2 なので、rows行には rows*2 ワード必要
	newMatrix := func(rows, cols, dataLen int) *bitsx.Matrix {
		return &bitsx.Matrix{Rows: rows, Cols: cols, Data: make([]uint64, dataLen)}
	}

	t.Run("正常_必要な長さちょうど", func(t *testing.T) {
		if err := newMatrix(3, 100, 6).Validate(); err != nil {
			t.Fatalf("予期せぬエラー: %v", err)
		}
	})

	t.Run("異常_1ワード不足", func(t *testing.T) {
		if err := newMatrix(3, 100, 5).Validate(); err == nil {
			t.Fatal("エラーを期待したが、nilが返された")
		}
	})

	// 商だけを比べると、strideに満たない余りが切り捨てられてすり抜ける
	t.Run("異常_1ワード過剰", func(t *testing.T) {
		if err := newMatrix(3, 100, 7).Validate(); err == nil {
			t.Fatal("エラーを期待したが、nilが返された")
		}
	})

	t.Run("異常_1行分過剰", func(t *testing.T) {
		if err := newMatrix(3, 100, 8).Validate(); err == nil {
			t.Fatal("エラーを期待したが、nilが返された")
		}
	})

	t.Run("異常_Dataが空", func(t *testing.T) {
		if err := newMatrix(1, 64, 0).Validate(); err == nil {
			t.Fatal("エラーを期待したが、nilが返された")
		}
	})

	t.Run("異常_Rowsが0以下", func(t *testing.T) {
		if err := newMatrix(0, 100, 6).Validate(); err == nil {
			t.Fatal("エラーを期待したが、nilが返された")
		}
		if err := newMatrix(-1, 100, 6).Validate(); err == nil {
			t.Fatal("エラーを期待したが、nilが返された")
		}
	})

	t.Run("異常_Colsが0以下", func(t *testing.T) {
		if err := newMatrix(3, 0, 6).Validate(); err == nil {
			t.Fatal("エラーを期待したが、nilが返された")
		}
		if err := newMatrix(3, -1, 6).Validate(); err == nil {
			t.Fatal("エラーを期待したが、nilが返された")
		}
	})

	t.Run("異常_Colsの桁あふれ", func(t *testing.T) {
		// Cols + 63 が桁あふれし、Stride() が負になる
		if err := newMatrix(1, math.MaxInt, 8).Validate(); err == nil {
			t.Fatal("エラーを期待したが、nilが返された")
		}
	})

	// Goの符号付き整数は、あふれても未定義動作にはならず2の補数で折り返す。
	// Rows * Stride() を乗算で判定すると、折り返した値が正常に見えてすり抜ける。
	t.Run("異常_Rowsの桁あふれが負になる場合", func(t *testing.T) {
		// MaxInt * 2 = -2
		if err := newMatrix(math.MaxInt, 100, 6).Validate(); err == nil {
			t.Fatal("エラーを期待したが、nilが返された")
		}
	})

	t.Run("異常_Rowsの桁あふれが小さい正の値になる場合", func(t *testing.T) {
		// (2^62+1) * 4 = 4 となり、len(Data)=8 に収まって見える
		if err := newMatrix(1<<62+1, 200, 8).Validate(); err == nil {
			t.Fatal("エラーを期待したが、nilが返された")
		}
	})

	t.Run("異常_Dotが不整合な行列を弾く", func(t *testing.T) {
		valid, err := bitsx.NewZerosMatrix(3, 100)
		if err != nil {
			t.Fatalf("予期せぬエラー: %v", err)
		}
		short := newMatrix(3, 100, 5)

		if _, err := valid.Dot(short); err == nil {
			t.Fatal("エラーを期待したが、nilが返された")
		}
		if _, err := short.Dot(valid); err == nil {
			t.Fatal("エラーを期待したが、nilが返された")
		}
	})

	t.Run("正常_コンストラクタが作る行列は検査を通る", func(t *testing.T) {
		for _, shape := range []struct{ rows, cols int }{{1, 1}, {3, 63}, {3, 64}, {3, 65}, {7, 130}} {
			m, err := bitsx.NewZerosMatrix(shape.rows, shape.cols)
			if err != nil {
				t.Fatalf("予期せぬエラー: %v", err)
			}
			if err := m.Validate(); err != nil {
				t.Errorf("(%d, %d): %v", shape.rows, shape.cols, err)
			}

			tr, err := m.Transpose()
			if err != nil {
				t.Fatalf("予期せぬエラー: %v", err)
			}
			if err := tr.Validate(); err != nil {
				t.Errorf("(%d, %d)の転置: %v", shape.rows, shape.cols, err)
			}
		}
	})

	t.Run("異常_DotTernaryが不整合な行列を弾く", func(t *testing.T) {
		valid, err := bitsx.NewZerosMatrix(3, 100)
		if err != nil {
			t.Fatalf("予期せぬエラー: %v", err)
		}
		short := newMatrix(3, 100, 5)

		if _, err := short.DotTernary(valid, valid); err == nil {
			t.Fatal("エラーを期待したが、nilが返された")
		}
		if _, err := valid.DotTernary(short, short); err == nil {
			t.Fatal("エラーを期待したが、nilが返された")
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
			if m.Rows != 2 || m.Cols != 70 {
				t.Errorf("ms[%d]の形状の不一致: got = (%d, %d), want = (2, 70)", i, m.Rows, m.Cols)
			}
		}
	})
}
