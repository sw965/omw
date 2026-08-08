package bitsx

import (
	"math"
	"testing"
)

// Dot / DotTernary は境界チェックを持たないアセンブリ実装へ生ポインタを渡す為、
// 不整合な Matrix を弾けないと範囲外アクセスになる。
// 検査が壊れた場合、このテストは失敗ではなくクラッシュする可能性がある。
//
// data は非公開フィールドの為、意図的に不整合な内部長を持つ Matrix を組み立てる
// このテストは package bitsx 内部に置く必要がある。
func TestMatrixValidate(t *testing.T) {
	// cols = 100 なら Stride() = 2 なので、rows行には rows*2 ワード必要
	newMatrix := func(rows, cols, dataLen int) *Matrix {
		return &Matrix{rows: rows, cols: cols, data: make([]uint64, dataLen)}
	}

	t.Run("正常_必要な長さちょうど", func(t *testing.T) {
		if err := newMatrix(3, 100, 6).validateDotAVX512Family(); err != nil {
			t.Fatalf("予期せぬエラー: %v", err)
		}
	})

	t.Run("異常_1ワード不足", func(t *testing.T) {
		if err := newMatrix(3, 100, 5).validateDotAVX512Family(); err == nil {
			t.Fatalf("エラーを期待したが、nilが返された")
		}
	})

	// 商だけを比べると、strideに満たない余りが切り捨てられてすり抜ける
	t.Run("異常_1ワード過剰", func(t *testing.T) {
		if err := newMatrix(3, 100, 7).validateDotAVX512Family(); err == nil {
			t.Fatalf("エラーを期待したが、nilが返された")
		}
	})

	t.Run("異常_1行分過剰", func(t *testing.T) {
		if err := newMatrix(3, 100, 8).validateDotAVX512Family(); err == nil {
			t.Fatalf("エラーを期待したが、nilが返された")
		}
	})

	t.Run("異常_dataが空", func(t *testing.T) {
		if err := newMatrix(1, 64, 0).validateDotAVX512Family(); err == nil {
			t.Fatalf("エラーを期待したが、nilが返された")
		}
	})

	t.Run("異常_Rowsが0以下", func(t *testing.T) {
		if err := newMatrix(0, 100, 6).validateDotAVX512Family(); err == nil {
			t.Fatalf("エラーを期待したが、nilが返された")
		}
		if err := newMatrix(-1, 100, 6).validateDotAVX512Family(); err == nil {
			t.Fatalf("エラーを期待したが、nilが返された")
		}
	})

	t.Run("異常_Colsが0以下", func(t *testing.T) {
		if err := newMatrix(3, 0, 6).validateDotAVX512Family(); err == nil {
			t.Fatalf("エラーを期待したが、nilが返された")
		}
		if err := newMatrix(3, -1, 6).validateDotAVX512Family(); err == nil {
			t.Fatalf("エラーを期待したが、nilが返された")
		}
	})

	t.Run("異常_Colsの桁あふれ", func(t *testing.T) {
		// Cols + 63 が桁あふれし、Stride() が負になる
		if err := newMatrix(1, math.MaxInt, 8).validateDotAVX512Family(); err == nil {
			t.Fatalf("エラーを期待したが、nilが返された")
		}
	})

	// Goの符号付き整数は、あふれても未定義動作にはならず2の補数で折り返す。
	// Rows * Stride() を乗算で判定すると、折り返した値が正常に見えてすり抜ける。
	t.Run("異常_Rowsの桁あふれが負になる場合", func(t *testing.T) {
		// MaxInt * 2 = -2
		if err := newMatrix(math.MaxInt, 100, 6).validateDotAVX512Family(); err == nil {
			t.Fatalf("エラーを期待したが、nilが返された")
		}
	})

	t.Run("異常_Rowsの桁あふれが小さい正の値になる場合", func(t *testing.T) {
		// (2^62+1) * 4 = 4 となり、len(data)=8 に収まって見える
		if err := newMatrix(1<<62+1, 200, 8).validateDotAVX512Family(); err == nil {
			t.Fatalf("エラーを期待したが、nilが返された")
		}
	})

	t.Run("異常_Dotが不整合な行列を弾く", func(t *testing.T) {
		valid, err := NewZerosMatrix(3, 100)
		if err != nil {
			t.Fatalf("予期せぬエラー: %v", err)
		}
		short := newMatrix(3, 100, 5)

		if _, err := valid.Dot(short); err == nil {
			t.Fatalf("エラーを期待したが、nilが返された")
		}
		if _, err := short.Dot(valid); err == nil {
			t.Fatalf("エラーを期待したが、nilが返された")
		}
	})

	t.Run("正常_コンストラクタが作る行列は検査を通る", func(t *testing.T) {
		for _, shape := range []struct{ rows, cols int }{{1, 1}, {3, 63}, {3, 64}, {3, 65}, {7, 130}} {
			m, err := NewZerosMatrix(shape.rows, shape.cols)
			if err != nil {
				t.Fatalf("予期せぬエラー: %v", err)
			}
			if err := m.validateDotAVX512Family(); err != nil {
				t.Errorf("(%d, %d): %v", shape.rows, shape.cols, err)
			}

			tr, err := m.Transpose()
			if err != nil {
				t.Fatalf("予期せぬエラー: %v", err)
			}
			if err := tr.validateDotAVX512Family(); err != nil {
				t.Errorf("(%d, %d)の転置: %v", shape.rows, shape.cols, err)
			}
		}
	})

	t.Run("異常_DotTernaryが不整合な行列を弾く", func(t *testing.T) {
		valid, err := NewZerosMatrix(3, 100)
		if err != nil {
			t.Fatalf("予期せぬエラー: %v", err)
		}
		short := newMatrix(3, 100, 5)

		if _, err := short.DotTernary(valid, valid); err == nil {
			t.Fatalf("エラーを期待したが、nilが返された")
		}
		if _, err := valid.DotTernary(short, short); err == nil {
			t.Fatalf("エラーを期待したが、nilが返された")
		}
	})
}
