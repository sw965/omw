package bitsx_test

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"slices"
	"testing"

	"github.com/sw965/omw/mathx"
	"github.com/sw965/omw/mathx/bitsx"
)

func TestNewZerosMatrix(t *testing.T) {
	tests := []struct {
		name    string
		rows    int
		cols    int
		wantErr bool
	}{
		{
			name:    "正常",
			rows:    2,
			cols:    100,
			wantErr: false,
		},
		{
			name:    "異常_rowsが0以下",
			rows:    0,
			cols:    10,
			wantErr: true,
		},
		{
			name:    "異常_colsが0以下",
			rows:    10,
			cols:    0,
			wantErr: true,
		},
		{
			name:    "異常_colsの桁あふれ",
			rows:    1,
			cols:    math.MaxInt,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := bitsx.NewZerosMatrix(tt.rows, tt.cols)
			if tt.wantErr && err == nil {
				t.Fatal("エラーを期待したが、nilが返された")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("nilを期待したが、エラーが返された: %v", err)
			}

			if tt.wantErr {
				return
			}

			if m.Rows() != tt.rows || m.Cols() != tt.cols {
				t.Errorf("形状の不一致: got = (%d, %d) want = (%d, %d)", m.Rows(), m.Cols(), tt.rows, tt.cols)
			}

			if c := m.OnesCount(); c != 0 {
				t.Errorf("OnesCountの不一致: got = %d, want = 0", c)
			}
		})
	}
}

func TestNewOnesMatrix(t *testing.T) {
	tests := []struct {
		name    string
		rows    int
		cols    int
		wantErr bool
	}{
		{
			name:    "正常_端数ビットあり",
			rows:    3,
			cols:    100,
			wantErr: false,
		},
		{
			name:    "正常_64の倍数",
			rows:    2,
			cols:    64,
			wantErr: false,
		},
		{
			name:    "正常_最小サイズ",
			rows:    1,
			cols:    1,
			wantErr: false,
		},
		{
			name:    "異常_rowsが0以下",
			rows:    0,
			cols:    10,
			wantErr: true,
		},
		{
			name:    "異常_colsが0以下",
			rows:    10,
			cols:    0,
			wantErr: true,
		},
		{
			name:    "異常_colsの桁あふれ",
			rows:    1,
			cols:    math.MaxInt,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := bitsx.NewOnesMatrix(tt.rows, tt.cols)
			if tt.wantErr && err == nil {
				t.Fatal("エラーを期待したが、nilが返された")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("nilを期待したが、エラーが返された: %v", err)
			}

			if tt.wantErr {
				return
			}

			if m.Rows() != tt.rows || m.Cols() != tt.cols {
				t.Errorf("形状の不一致: got = (%d, %d) want = (%d, %d)", m.Rows(), m.Cols(), tt.rows, tt.cols)
			}

			want := tt.rows * tt.cols
			if c := m.OnesCount(); c != want {
				t.Errorf("OnesCountの不一致: got = %d, want = %d", c, want)
			}
		})
	}
}

func TestNewRandMatrixHalfPow(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))

	tests := []struct {
		name    string
		rows    int
		cols    int
		n       int
		wantErr bool
	}{
		{name: "正常_nが1", rows: 3, cols: 100, n: 1},
		{name: "正常_nが3", rows: 3, cols: 100, n: 3},
		{name: "正常_nが0", rows: 3, cols: 100, n: 0},
		{name: "異常_nが負", rows: 3, cols: 100, n: -1, wantErr: true},
		{name: "異常_rowsが0以下", rows: 0, cols: 10, n: 1, wantErr: true},
		{name: "異常_colsが0以下", rows: 10, cols: 0, n: 1, wantErr: true},
		{name: "異常_colsの桁あふれ", rows: 1, cols: math.MaxInt, n: 1, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := bitsx.NewRandMatrixHalfPow(tt.rows, tt.cols, tt.n, rng)
			if tt.wantErr && err == nil {
				t.Fatal("エラーを期待したが、nilが返された")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("nilを期待したが、エラーが返された: %v", err)
			}

			if tt.wantErr {
				return
			}

			if m.Rows() != tt.rows || m.Cols() != tt.cols {
				t.Errorf("形状の不一致: got = (%d, %d) want = (%d, %d)", m.Rows(), m.Cols(), tt.rows, tt.cols)
			}

			totalBits := tt.rows * tt.cols
			c := m.OnesCount()
			if c < 0 || c > totalBits {
				t.Errorf("OnesCountの不一致: got = %d, want = 0 <= c <= %d", c, totalBits)
			}
			// n = 0 は確率 1。端数ワードの余りビットは ApplyTailMask で 0 のため、ちょうど総ビット数になる
			if tt.n == 0 && c != totalBits {
				t.Errorf("n = 0 で全ビット1にならない: got = %d, want = %d", c, totalBits)
			}
		})
	}
}

func TestNewRandMatrix(t *testing.T) {
	t.Run("正常_NewRandMatrixHalfPowのn=1と一致", func(t *testing.T) {
		got, err := bitsx.NewRandMatrix(3, 100, rand.New(rand.NewPCG(5, 6)))
		if err != nil {
			t.Fatalf("nilを期待したが、エラーが返された: %v", err)
		}
		want, err := bitsx.NewRandMatrixHalfPow(3, 100, 1, rand.New(rand.NewPCG(5, 6)))
		if err != nil {
			t.Fatalf("nilを期待したが、エラーが返された: %v", err)
		}
		if !got.Equal(want) {
			t.Error("同じシードで NewRandMatrixHalfPow(n=1) と一致しない")
		}
	})

	t.Run("異常_rowsが0以下", func(t *testing.T) {
		if _, err := bitsx.NewRandMatrix(0, 10, rand.New(rand.NewPCG(1, 2))); err == nil {
			t.Fatal("エラーを期待したが、nilが返された")
		}
	})
}

func TestNewRandMatrixHalfPowStatistics(t *testing.T) {
	rng := rand.New(rand.NewPCG(42, 100))

	const (
		rows = 1000
		cols = 1000
	)
	totalBits := float64(rows * cols)

	tests := []struct {
		name  string
		n     int
		wantP float64
		tol   float64
	}{
		{
			// n = 0 は乱数を使わず全ビット 1 なので、誤差は無い
			name:  "nが0_確率1",
			n:     0,
			wantP: 1.0,
			tol:   0,
		},
		{
			// N=10^6, p=0.5, σ=0.0005（tol=0.015 は 30σ。正しく実装されていれば収まる確率 ≒ 100%）
			name:  "nが1_確率0.5",
			n:     1,
			wantP: 0.5,
			tol:   0.015,
		},
		{
			// N=10^6, p=0.25, σ≈0.000433（tol=0.015 は 34.6σ。正しく実装されていれば収まる確率 ≒ 100%）
			name:  "nが2_確率0.25",
			n:     2,
			wantP: 0.25,
			tol:   0.015,
		},
		{
			// N=10^6, p=0.125, σ≈0.000331（tol=0.015 は 45.3σ。正しく実装されていれば収まる確率 ≒ 100%）
			name:  "nが3_確率0.125",
			n:     3,
			wantP: 0.125,
			tol:   0.015,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := bitsx.NewRandMatrixHalfPow(rows, cols, tt.n, rng)
			if err != nil {
				t.Fatalf("nilを期待したが、エラーが返された: %v", err)
			}

			gotP := float64(m.OnesCount()) / totalBits
			if !mathx.ApproxEqual(gotP, tt.wantP, tt.tol) {
				t.Errorf("確率の不一致: got = %f, want = %f (±%f)", gotP, tt.wantP, tt.tol)
			}
		})
	}
}

func FuzzNewRandMatrixHalfPow(f *testing.F) {
	seeds := []struct {
		rows8        uint8
		cols16       uint16
		n8           int8
		seed1, seed2 uint64
	}{
		{3, 100, 1, 1, 2},
		{1, 64, 3, 10, 20},
		{10, 1, 0, 100, 200},
		{2, 5, -1, 1000, 2000},
	}
	for _, s := range seeds {
		f.Add(s.rows8, s.cols16, s.n8, s.seed1, s.seed2)
	}

	f.Fuzz(func(t *testing.T, rows8 uint8, cols16 uint16, n8 int8, seed1, seed2 uint64) {
		rows := int(rows8)
		cols := int(cols16)
		n := int(n8)

		rng := rand.New(rand.NewPCG(seed1, seed2))
		m, err := bitsx.NewRandMatrixHalfPow(rows, cols, n, rng)
		if rows <= 0 || cols <= 0 || n < 0 {
			if err == nil {
				t.Fatal("エラーを期待したが、nilが返された")
			}
			return
		}

		if err != nil {
			t.Fatalf("nilを期待したが、エラーが返された: %v", err)
		}

		if m.Rows() != rows || m.Cols() != cols {
			t.Errorf("形状の不一致: got = (%d, %d) want = (%d, %d)", m.Rows(), m.Cols(), rows, cols)
		}

		totalBits := rows * cols
		if c := m.OnesCount(); c < 0 || c > totalBits {
			t.Errorf("OnesCountの不一致: got = %d, want = 0 <= c <= %d", c, totalBits)
		}
	})
}

func TestNewSignMatrix(t *testing.T) {
	tests := []struct {
		name     string
		rows     int
		cols     int
		x        []int
		wantOnes int
		wantBits [][]int
		wantErr  bool
	}{
		{
			name:     "正常_正負混在",
			rows:     2,
			cols:     3,
			x:        []int{10, -5, 0, -1, 3, -100},
			wantOnes: 3,
			wantBits: [][]int{
				{1, 0, 1},
				{0, 1, 0},
			},
			wantErr: false,
		},
		{
			name:     "正常_すべて正",
			rows:     1,
			cols:     4,
			x:        []int{1, 2, 3, 4},
			wantOnes: 4,
			wantBits: [][]int{
				{1, 1, 1, 1},
			},
			wantErr: false,
		},
		{
			name:     "正常_すべて負",
			rows:     2,
			cols:     2,
			x:        []int{-1, -2, -3, -4},
			wantOnes: 0,
			wantBits: [][]int{
				{0, 0},
				{0, 0},
			},
			wantErr: false,
		},
		{
			name:    "異常_len(x)不一致",
			rows:    2,
			cols:    3,
			x:       []int{1, 2, 3},
			wantErr: true,
		},
		{
			name:    "異常_rowsが0以下",
			rows:    0,
			cols:    3,
			x:       []int{},
			wantErr: true,
		},
		{
			name:    "異常_colsが0以下",
			rows:    2,
			cols:    0,
			x:       []int{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := bitsx.NewSignMatrix(tt.rows, tt.cols, tt.x)
			if tt.wantErr && err == nil {
				t.Fatal("エラーを期待したが、nilが返された")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("nilを期待したが、エラーが返された: %v", err)
			}

			if tt.wantErr {
				return
			}

			if m.Rows() != tt.rows || m.Cols() != tt.cols {
				t.Errorf("形状の不一致: got = (%d, %d) want = (%d, %d)", m.Rows(), m.Cols(), tt.rows, tt.cols)
			}

			if c := m.OnesCount(); c != tt.wantOnes {
				t.Errorf("OnesCountの不一致: got = %d, want = %d", c, tt.wantOnes)
			}

			for r := 0; r < tt.rows; r++ {
				for c := 0; c < tt.cols; c++ {
					gotBit, err := m.Bit(r, c)
					if err != nil {
						t.Fatalf("nilを期待したが、エラーが返された %v", err)
					}
					if gotBit != uint64(tt.wantBits[r][c]) {
						t.Errorf("Bit(%d, %d)の不一致: got = %d want = %d", r, c, gotBit, tt.wantBits[r][c])
					}
				}
			}
		})
	}
}

func TestMatrixTailMask(t *testing.T) {
	tests := []struct {
		name string
		cols int
		want uint64
	}{
		{
			name: "正常_64の倍数_64列",
			cols: 64,
			want: ^uint64(0),
		},
		{
			name: "正常_64の倍数_128列",
			cols: 128,
			want: ^uint64(0),
		},
		{
			name: "正常_端数1ビット",
			cols: 1,
			want: 0b00000001,
		},
		{
			name: "正常_端数6ビット",
			cols: 70,
			want: 0b00111111,
		},
		{
			name: "正常_端数63ビット",
			cols: 63,
			want: 0b01111111_11111111_11111111_11111111_11111111_11111111_11111111_11111111,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := bitsx.NewZerosMatrix(1, tt.cols)
			if err != nil {
				t.Fatalf("nilを期待したが、エラーが返された: %v", err)
			}
			if got := m.TailMask(); got != tt.want {
				t.Errorf("TailMaskの不一致: got = %#x want = %#x", got, tt.want)
			}
		})
	}
}

func TestMatrixApplyTailMask(t *testing.T) {
	tests := []struct {
		name      string
		rows      int
		cols      int
		setup     func(m *bitsx.Matrix) error
		wantWords map[int]uint64
	}{
		{
			name: "正常_端数6ビット",
			rows: 2,
			cols: 70,
			setup: func(m *bitsx.Matrix) error {
				if err := m.SetWord(1, ^uint64(0)); err != nil {
					return err
				}
				return m.SetWord(3, ^uint64(0))
			},
			wantWords: map[int]uint64{
				1: 0b00111111,
				3: 0b00111111,
			},
		},
		{
			name: "正常_64の倍数",
			rows: 2,
			cols: 64,
			setup: func(m *bitsx.Matrix) error {
				if err := m.SetWord(0, ^uint64(0)); err != nil {
					return err
				}
				return m.SetWord(1, ^uint64(0))
			},
			wantWords: map[int]uint64{
				0: ^uint64(0),
				1: ^uint64(0),
			},
		},
		{
			name: "正常_端数1ビット",
			rows: 1,
			cols: 1,
			setup: func(m *bitsx.Matrix) error {
				return m.SetWord(0, ^uint64(0))
			},
			wantWords: map[int]uint64{
				0: 0b00000001,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := bitsx.NewZerosMatrix(tt.rows, tt.cols)
			if err != nil {
				t.Fatalf("nilを期待したが、エラーが返された: %v", err)
			}
			if err := tt.setup(m); err != nil {
				t.Fatalf("nilを期待したが、エラーが返された: %v", err)
			}
			m.ApplyTailMask()

			for idx, want := range tt.wantWords {
				got, err := m.Word(idx)
				if err != nil {
					t.Fatalf("nilを期待したが、エラーが返された: %v", err)
				}
				if got != want {
					t.Errorf("Word(%d)の不一致: got = %#x want = %#x", idx, got, want)
				}
			}
		})
	}
}

func TestMatrixValidateSameShape(t *testing.T) {
	tests := []struct {
		name    string
		a       *bitsx.Matrix
		b       *bitsx.Matrix
		wantErr bool
	}{
		{
			name:    "一致",
			a:       bitsx.NewZerosMatrixForTest(t, 2, 70),
			b:       bitsx.NewZerosMatrixForTest(t, 2, 70),
			wantErr: false,
		},
		{
			name:    "不一致_rows",
			a:       bitsx.NewZerosMatrixForTest(t, 2, 70),
			b:       bitsx.NewZerosMatrixForTest(t, 3, 70),
			wantErr: true,
		},
		{
			name:    "不一致_cols",
			a:       bitsx.NewZerosMatrixForTest(t, 2, 70),
			b:       bitsx.NewZerosMatrixForTest(t, 2, 71),
			wantErr: true,
		},
		{
			name:    "不一致_rowsとcolsが反転",
			a:       bitsx.NewZerosMatrixForTest(t, 2, 70),
			b:       bitsx.NewZerosMatrixForTest(t, 70, 2),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.a.ValidateSameShape(tt.b)
			if tt.wantErr && err == nil {
				t.Fatal("エラーを期待したが、nilが返された")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("nilを期待したが、エラーが返された: %v", err)
			}
		})
	}
}

func TestMatrixClone(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))

	tests := []struct {
		name string
		m    *bitsx.Matrix
	}{
		{
			name: "正常_端数なし",
			m:    bitsx.NewRandMatrixForTest(t, 2, 128, rng),
		},
		{
			name: "正常_端数あり",
			m:    bitsx.NewRandMatrixForTest(t, 3, 100, rng),
		},
		{
			name: "正常_最小1x1",
			m:    bitsx.NewRandMatrixForTest(t, 1, 1, rng),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cloned := tt.m.Clone()

			if !tt.m.Equal(cloned) {
				t.Error("クローン前後の行列が一致しない")
			}

			if err := cloned.Toggle(0, 0); err != nil {
				t.Fatalf("nilを期待したが、エラーが返された: %v", err)
			}

			// 参照透過性の確認
			if tt.m.Equal(cloned) {
				t.Error("ディープコピーになっていない")
			}
		})
	}
}

func TestMatrixEqual(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))

	randMat := bitsx.NewRandMatrixForTest(t, 3, 100, rng)
	randMatDiff := randMat.Clone()
	if err := randMatDiff.Toggle(0, 0); err != nil {
		t.Fatalf("nilを期待したが、エラーが返された: %v", err)
	}

	tests := []struct {
		name string
		a    *bitsx.Matrix
		b    *bitsx.Matrix
		want bool
	}{
		{
			name: "一致",
			a:    randMat,
			b:    randMat.Clone(),
			want: true,
		},
		{
			name: "形状の違い",
			a:    bitsx.NewZerosMatrixForTest(t, 3, 100),
			b:    bitsx.NewZerosMatrixForTest(t, 3, 101),
			want: false,
		},
		{
			name: "内容の違い",
			a:    randMat,
			b:    randMatDiff,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Equal(tt.b); got != tt.want {
				t.Errorf("結果の不一致: got = %t want = %t", got, tt.want)
			}
		})
	}
}

func TestMatrixWord(t *testing.T) {
	const alternateBits = 0b01010101_01010101_01010101_01010101_01010101_01010101_01010101_01010101

	tests := []struct {
		name    string
		m       *bitsx.Matrix
		idx     int
		want    uint64
		wantErr bool
	}{
		{
			// cols=70 (stride=2), rows=2 -> len(data)=4
			// 有効範囲: 0 <= idx < 4
			name:    "正常_境界_下限",
			m:       bitsx.NewOnesMatrixForTest(t, 2, 70),
			idx:     0,
			want:    ^uint64(0),
			wantErr: false,
		},
		{
			// cols=128 (stride=2), rows=3 -> len(data)=6
			// 有効範囲: 0 <= idx < 6
			name:    "正常_境界_上限",
			m:       bitsx.NewOnesMatrixWithSetWordForTest(t, 3, 128, 5, alternateBits),
			idx:     5,
			want:    alternateBits,
			wantErr: false,
		},
		{
			// cols=64 (stride=1), rows=1 -> len(data)=1
			// 有効範囲: 0 <= idx < 1
			name:    "異常_境界_下限未満",
			m:       bitsx.NewZerosMatrixForTest(t, 1, 64),
			idx:     -1,
			wantErr: true,
		},
		{
			// cols=70 (stride=2), rows=2 -> len(data)=4
			// 有効範囲: 0 <= idx < 4
			name:    "異常_境界_上限超過",
			m:       bitsx.NewZerosMatrixForTest(t, 2, 70),
			idx:     4,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.Word(tt.idx)
			if tt.wantErr && err == nil {
				t.Fatal("エラーを期待したが、nilが返された")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("nilを期待したが、エラーが返された: %v", err)
			}

			if tt.wantErr {
				return
			}

			if got != tt.want {
				t.Errorf("Word(%d)の不一致: got = %#x want = %#x", tt.idx, got, tt.want)
			}
		})
	}
}

func TestMatrixSetWord(t *testing.T) {
	newZeros := func(t *testing.T, rows, cols int) *bitsx.Matrix {
		t.Helper()
		m, err := bitsx.NewZerosMatrix(rows, cols)
		if err != nil {
			t.Fatalf("nilを期待したが、エラーが返された: %v", err)
		}
		return m
	}

	tests := []struct {
		name    string
		m       *bitsx.Matrix
		idx     int
		word    uint64
		want    uint64
		wantErr bool
	}{
		{
			// cols=70 (stride=2), rows=2 -> len(data)=4
			// 有効範囲: 0 <= idx < 4 (idx 0 は非tailワード)
			name:    "正常_境界_下限",
			m:       newZeros(t, 2, 70),
			idx:     0,
			word:    ^uint64(0),
			want:    ^uint64(0),
			wantErr: false,
		},
		{
			// cols=70 (stride=2), rows=2 -> len(data)=4
			// idx 3 は 2行目のtailワード（有効ビットは 70-64=6 ビット）
			name:    "正常_境界_上限_tailマスク適用",
			m:       newZeros(t, 2, 70),
			idx:     3,
			word:    ^uint64(0),
			want:    0b00111111,
			wantErr: false,
		},
		{
			// cols=64 (stride=1), rows=1 -> len(data)=1
			// 有効範囲: 0 <= idx < 1
			name:    "異常_境界_下限未満",
			m:       newZeros(t, 1, 64),
			idx:     -1,
			word:    0,
			wantErr: true,
		},
		{
			// cols=70 (stride=2), rows=2 -> len(data)=4
			// 有効範囲: 0 <= idx < 4
			name:    "異常_境界_上限超過",
			m:       newZeros(t, 2, 70),
			idx:     4,
			word:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.m.SetWord(tt.idx, tt.word)
			if tt.wantErr && err == nil {
				t.Fatal("エラーを期待したが、nilが返された")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("nilを期待したが、エラーが返された: %v", err)
			}

			if tt.wantErr {
				return
			}

			got, err := tt.m.Word(tt.idx)
			if err != nil {
				t.Fatalf("nilを期待したが、エラーが返された: %v", err)
			}

			if got != tt.want {
				t.Errorf("SetWord(%d)後のWordの不一致: got = %#x want = %#x", tt.idx, got, tt.want)
			}
		})
	}
}

func TestMatrixBit(t *testing.T) {
	tests := []struct {
		name    string
		m       *bitsx.Matrix
		row     int
		col     int
		want    uint64
		wantErr bool
	}{
		{
			name:    "正常_境界_下限",
			m:       bitsx.NewZerosMatrixForTest(t, 2, 70),
			row:     0,
			col:     0,
			want:    0,
			wantErr: false,
		},
		{
			name:    "正常_境界_上限",
			m:       bitsx.NewOnesMatrixForTest(t, 2, 70),
			row:     1,
			col:     69,
			want:    1,
			wantErr: false,
		},
		{
			name:    "異常_境界_r下限未満",
			m:       bitsx.NewZerosMatrixForTest(t, 2, 70),
			row:     -1,
			col:     0,
			wantErr: true,
		},
		{
			name:    "異常_境界_r上限超過",
			m:       bitsx.NewZerosMatrixForTest(t, 2, 70),
			row:     2,
			col:     0,
			wantErr: true,
		},
		{
			name:    "異常_境界_c下限未満",
			m:       bitsx.NewZerosMatrixForTest(t, 2, 70),
			row:     0,
			col:     -1,
			wantErr: true,
		},
		{
			name:    "異常_境界_c上限超過",
			m:       bitsx.NewZerosMatrixForTest(t, 2, 70),
			row:     0,
			col:     70,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.Bit(tt.row, tt.col)
			if tt.wantErr && err == nil {
				t.Fatal("エラーを期待したが、nilが返された")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("nilを期待したが、エラーが返された: %v", err)
			}

			if tt.wantErr {
				return
			}

			if got != tt.want {
				t.Errorf("Bit(%d, %d)の不一致: got = %d want = %d", tt.row, tt.col, got, tt.want)
			}
		})
	}
}

func TestMatrixSet(t *testing.T) {
	tests := []struct {
		name    string
		m       *bitsx.Matrix
		row     int
		col     int
		wantErr bool
	}{
		{
			name:    "正常_境界_下限",
			m:       bitsx.NewZerosMatrixForTest(t, 2, 70),
			row:     0,
			col:     0,
			wantErr: false,
		},
		{
			name:    "正常_境界_上限",
			m:       bitsx.NewZerosMatrixForTest(t, 2, 70),
			row:     1,
			col:     69,
			wantErr: false,
		},
		{
			name:    "正常_1の場所にSet",
			m:       bitsx.NewOnesMatrixForTest(t, 2, 70),
			row:     0,
			col:     0,
			wantErr: false,
		},
		{
			name:    "異常_境界_r下限未満",
			m:       bitsx.NewZerosMatrixForTest(t, 2, 70),
			row:     -1,
			col:     0,
			wantErr: true,
		},
		{
			name:    "異常_境界_r上限超過",
			m:       bitsx.NewZerosMatrixForTest(t, 2, 70),
			row:     2,
			col:     0,
			wantErr: true,
		},
		{
			name:    "異常_境界_c下限未満",
			m:       bitsx.NewZerosMatrixForTest(t, 2, 70),
			row:     0,
			col:     -1,
			wantErr: true,
		},
		{
			name:    "異常_境界_c上限超過",
			m:       bitsx.NewZerosMatrixForTest(t, 2, 70),
			row:     0,
			col:     70,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.m.Set(tt.row, tt.col)
			if tt.wantErr && err == nil {
				t.Fatal("エラーを期待したが、nilが返された")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("nilを期待したが、エラーが返された: %v", err)
			}

			if tt.wantErr {
				return
			}

			got, err := tt.m.Bit(tt.row, tt.col)
			if err != nil {
				t.Fatalf("nilを期待したが、エラーが返された: %v", err)
			}

			if got != 1 {
				t.Errorf("Set(%d, %d)後のBitの不一致: got = %d want = 1", tt.row, tt.col, got)
			}
		})
	}
}

func TestMatrixClear(t *testing.T) {
	tests := []struct {
		name    string
		m       *bitsx.Matrix
		row     int
		col     int
		wantErr bool
	}{
		{
			name:    "正常_境界_下限",
			m:       bitsx.NewOnesMatrixForTest(t, 2, 70),
			row:     0,
			col:     0,
			wantErr: false,
		},
		{
			name:    "正常_境界_上限",
			m:       bitsx.NewOnesMatrixForTest(t, 2, 70),
			row:     1,
			col:     69,
			wantErr: false,
		},
		{
			name:    "正常_0の場所をClear",
			m:       bitsx.NewZerosMatrixForTest(t, 2, 70),
			row:     0,
			col:     0,
			wantErr: false,
		},
		{
			name:    "異常_境界_r下限未満",
			m:       bitsx.NewOnesMatrixForTest(t, 2, 70),
			row:     -1,
			col:     0,
			wantErr: true,
		},
		{
			name:    "異常_境界_r上限超過",
			m:       bitsx.NewOnesMatrixForTest(t, 2, 70),
			row:     2,
			col:     0,
			wantErr: true,
		},
		{
			name:    "異常_境界_c下限未満",
			m:       bitsx.NewOnesMatrixForTest(t, 2, 70),
			row:     0,
			col:     -1,
			wantErr: true,
		},
		{
			name:    "異常_境界_c上限超過",
			m:       bitsx.NewOnesMatrixForTest(t, 2, 70),
			row:     0,
			col:     70,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.m.Clear(tt.row, tt.col)
			if tt.wantErr && err == nil {
				t.Fatal("エラーを期待したが、nilが返された")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("nilを期待したが、エラーが返された: %v", err)
			}

			if tt.wantErr {
				return
			}

			got, err := tt.m.Bit(tt.row, tt.col)
			if err != nil {
				t.Fatalf("nilを期待したが、エラーが返された: %v", err)
			}

			if got != 0 {
				t.Errorf("Clear(%d, %d)後のBitの不一致: got = %d want = 0", tt.row, tt.col, got)
			}
		})
	}
}

func TestMatrixToggle(t *testing.T) {
	tests := []struct {
		name    string
		m       *bitsx.Matrix
		row     int
		col     int
		want    uint64
		wantErr bool
	}{
		{
			name:    "正常_境界_下限_0から1へ反転",
			m:       bitsx.NewZerosMatrixForTest(t, 2, 70),
			row:     0,
			col:     0,
			want:    1,
			wantErr: false,
		},
		{
			name:    "正常_境界_上限_1から0へ反転",
			m:       bitsx.NewOnesMatrixForTest(t, 2, 70),
			row:     1,
			col:     69,
			want:    0,
			wantErr: false,
		},
		{
			name:    "異常_境界_r下限未満",
			m:       bitsx.NewZerosMatrixForTest(t, 2, 70),
			row:     -1,
			col:     0,
			wantErr: true,
		},
		{
			name:    "異常_境界_r上限超過",
			m:       bitsx.NewZerosMatrixForTest(t, 2, 70),
			row:     2,
			col:     0,
			wantErr: true,
		},
		{
			name:    "異常_境界_c下限未満",
			m:       bitsx.NewZerosMatrixForTest(t, 2, 70),
			row:     0,
			col:     -1,
			wantErr: true,
		},
		{
			name:    "異常_境界_c上限超過",
			m:       bitsx.NewZerosMatrixForTest(t, 2, 70),
			row:     0,
			col:     70,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.m.Toggle(tt.row, tt.col)
			if tt.wantErr && err == nil {
				t.Fatal("エラーを期待したが、nilが返された")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("nilを期待したが、エラーが返された: %v", err)
			}

			if tt.wantErr {
				return
			}

			got, err := tt.m.Bit(tt.row, tt.col)
			if err != nil {
				t.Fatalf("nilを期待したが、エラーが返された: %v", err)
			}

			if got != tt.want {
				t.Errorf("Toggle(%d, %d)後のBitの不一致: got = %d want = %d", tt.row, tt.col, got, tt.want)
			}
		})
	}
}

func TestMatrixAnd(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	randMatA := bitsx.NewRandMatrixForTest(t, 3, 70, rng)
	randMatB := bitsx.NewRandMatrixForTest(t, 3, 70, rng)

	tests := []struct {
		name    string
		a       *bitsx.Matrix
		b       *bitsx.Matrix
		want    *bitsx.Matrix
		wantErr bool
	}{
		{
			name:    "正常_全1と全1",
			a:       bitsx.NewOnesMatrixForTest(t, 3, 70),
			b:       bitsx.NewOnesMatrixForTest(t, 3, 70),
			want:    bitsx.NewOnesMatrixForTest(t, 3, 70),
			wantErr: false,
		},
		{
			name:    "正常_全1と全0",
			a:       bitsx.NewOnesMatrixForTest(t, 3, 70),
			b:       bitsx.NewZerosMatrixForTest(t, 3, 70),
			want:    bitsx.NewZerosMatrixForTest(t, 3, 70),
			wantErr: false,
		},
		{
			name:    "正常_全0と全1",
			a:       bitsx.NewZerosMatrixForTest(t, 3, 70),
			b:       bitsx.NewOnesMatrixForTest(t, 3, 70),
			want:    bitsx.NewZerosMatrixForTest(t, 3, 70),
			wantErr: false,
		},
		{
			name:    "正常_ランダム行列",
			a:       randMatA,
			b:       randMatB,
			want:    nil,
			wantErr: false,
		},
		{
			name:    "異常_形状不一致_rows",
			a:       bitsx.NewZerosMatrixForTest(t, 2, 70),
			b:       bitsx.NewZerosMatrixForTest(t, 3, 70),
			wantErr: true,
		},
		{
			name:    "異常_形状不一致_cols",
			a:       bitsx.NewZerosMatrixForTest(t, 2, 70),
			b:       bitsx.NewZerosMatrixForTest(t, 2, 71),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aClone := tt.a.Clone()
			bClone := tt.b.Clone()

			got, err := tt.a.And(tt.b)
			if tt.wantErr && err == nil {
				t.Fatal("エラーを期待したが、nilが返された")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("nilを期待したが、エラーが返された: %v", err)
			}

			if tt.wantErr {
				return
			}

			if !tt.a.Equal(aClone) || !tt.b.Equal(bClone) {
				t.Error("元の行列が変更された")
			}

			if tt.want != nil {
				if !got.Equal(tt.want) {
					t.Error("結果の不一致")
				}
			} else {
				for r := range 3 {
					for c := range 70 {
						bitA, aErr := tt.a.Bit(r, c)
						if aErr != nil {
							t.Fatalf("Bit(%d, %d)で期待しないエラーが発生した: %v", r, c, aErr)
						}

						bitB, bErr := tt.b.Bit(r, c)
						if bErr != nil {
							t.Fatalf("Bit(%d, %d)で期待しないエラーが発生した: %v", r, c, bErr)
						}

						gotBit, gotErr := got.Bit(r, c)
						if gotErr != nil {
							t.Fatalf("Bit(%d, %d)で期待しないエラーが発生した: %v", r, c, gotErr)
						}

						wantBit := bitA & bitB
						if gotBit != wantBit {
							t.Errorf("Bit(%d, %d)の不一致: got = %d, want = %d", r, c, gotBit, wantBit)
						}
					}
				}
			}
		})
	}
}

func TestMatrixXor(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))

	randMatA := bitsx.NewRandMatrixForTest(t, 3, 70, rng)
	randMatB := bitsx.NewRandMatrixForTest(t, 3, 70, rng)

	tests := []struct {
		name    string
		a       *bitsx.Matrix
		b       *bitsx.Matrix
		want    *bitsx.Matrix
		wantErr bool
	}{
		{
			name:    "正常_全1と全1",
			a:       bitsx.NewOnesMatrixForTest(t, 3, 70),
			b:       bitsx.NewOnesMatrixForTest(t, 3, 70),
			want:    bitsx.NewZerosMatrixForTest(t, 3, 70),
			wantErr: false,
		},
		{
			name:    "正常_全1と全0",
			a:       bitsx.NewOnesMatrixForTest(t, 3, 70),
			b:       bitsx.NewZerosMatrixForTest(t, 3, 70),
			want:    bitsx.NewOnesMatrixForTest(t, 3, 70),
			wantErr: false,
		},
		{
			name:    "正常_全0と全1",
			a:       bitsx.NewZerosMatrixForTest(t, 3, 70),
			b:       bitsx.NewOnesMatrixForTest(t, 3, 70),
			want:    bitsx.NewOnesMatrixForTest(t, 3, 70),
			wantErr: false,
		},
		{
			name:    "正常_全0と全0",
			a:       bitsx.NewZerosMatrixForTest(t, 3, 70),
			b:       bitsx.NewZerosMatrixForTest(t, 3, 70),
			want:    bitsx.NewZerosMatrixForTest(t, 3, 70),
			wantErr: false,
		},
		{
			name:    "正常_ランダム行列",
			a:       randMatA,
			b:       randMatB,
			want:    nil,
			wantErr: false,
		},
		{
			name:    "異常_形状不一致_rows",
			a:       bitsx.NewZerosMatrixForTest(t, 2, 70),
			b:       bitsx.NewZerosMatrixForTest(t, 3, 70),
			wantErr: true,
		},
		{
			name:    "異常_形状不一致_cols",
			a:       bitsx.NewZerosMatrixForTest(t, 2, 70),
			b:       bitsx.NewZerosMatrixForTest(t, 2, 71),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aClone := tt.a.Clone()
			bClone := tt.b.Clone()

			got, err := tt.a.Xor(tt.b)
			if tt.wantErr && err == nil {
				t.Fatal("エラーを期待したが、nilが返された")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("nilを期待したが、エラーが返された: %v", err)
			}

			if tt.wantErr {
				return
			}

			if !tt.a.Equal(aClone) || !tt.b.Equal(bClone) {
				t.Error("元の行列が変更された")
			}

			if tt.want != nil {
				if !got.Equal(tt.want) {
					t.Error("結果の不一致")
				}
			} else {
				for r := range 3 {
					for c := range 70 {
						bitA, aErr := tt.a.Bit(r, c)
						if aErr != nil {
							t.Fatalf("Bit(%d, %d)で期待しないエラーが発生した: %v", r, c, aErr)
						}

						bitB, bErr := tt.b.Bit(r, c)
						if bErr != nil {
							t.Fatalf("Bit(%d, %d)で期待しないエラーが発生した: %v", r, c, bErr)
						}

						gotBit, gotErr := got.Bit(r, c)
						if gotErr != nil {
							t.Fatalf("Bit(%d, %d)で期待しないエラーが発生した: %v", r, c, gotErr)
						}

						wantBit := bitA ^ bitB
						if gotBit != wantBit {
							t.Errorf("Bit(%d, %d)の不一致: got = %d, want = %d", r, c, gotBit, wantBit)
						}
					}
				}
			}
		})
	}
}

func TestMatrixOnesCount(t *testing.T) {
	tests := []struct {
		name string
		m    *bitsx.Matrix
		want int
	}{
		{
			name: "正常_全0行列",
			m:    bitsx.NewZerosMatrixForTest(t, 3, 70),
			want: 0,
		},
		{
			name: "正常_全1行列_端数なし",
			m:    bitsx.NewOnesMatrixForTest(t, 2, 128),
			want: 256,
		},
		{
			name: "正常_全1行列_端数あり",
			m:    bitsx.NewOnesMatrixForTest(t, 3, 70),
			want: 210,
		},
		{
			name: "正常_部分指定ビット",
			m:    bitsx.NewMatrixForTest(t, 70, [][]int{{0}, {69}}),
			want: 2,
		},
		{
			name: "正常_tailマスク適用後",
			m:    bitsx.NewOnesMatrixWithSetWordForTest(t, 1, 70, 1, ^uint64(0)),
			want: 70,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.OnesCount(); got != tt.want {
				t.Errorf("値の不一致: got = %d, want = %d", got, tt.want)
			}
		})
	}
}

func TestMatrixHammingDistance(t *testing.T) {
	rng := rand.New(rand.NewPCG(9, 10))

	randMat := bitsx.NewRandMatrixForTest(t, 2, 100, rng)
	randMatToggled := randMat.Clone()
	if err := randMatToggled.Toggle(0, 0); err != nil {
		t.Fatalf("nilを期待したが、エラーが返された: %v", err)
	}

	var zeroMatrixA, zeroMatrixB bitsx.Matrix

	tests := []struct {
		name    string
		a       *bitsx.Matrix
		b       *bitsx.Matrix
		want    int
		wantErr bool
	}{
		{
			name:    "正常_自身との距離",
			a:       randMat,
			b:       randMat,
			want:    0,
			wantErr: false,
		},
		{
			name:    "正常_全0と全1_端数あり",
			a:       bitsx.NewZerosMatrixForTest(t, 2, 70),
			b:       bitsx.NewOnesMatrixForTest(t, 2, 70),
			want:    140,
			wantErr: false,
		},
		{
			name:    "正常_1ビット反転",
			a:       randMat,
			b:       randMatToggled,
			want:    1,
			wantErr: false,
		},
		{
			name:    "正常_ゼロ値同士",
			a:       &zeroMatrixA,
			b:       &zeroMatrixB,
			want:    0,
			wantErr: false,
		},
		{
			name:    "異常_形状不一致_rows",
			a:       bitsx.NewZerosMatrixForTest(t, 1, 10),
			b:       bitsx.NewZerosMatrixForTest(t, 2, 10),
			wantErr: true,
		},
		{
			name:    "異常_形状不一致_cols",
			a:       bitsx.NewZerosMatrixForTest(t, 1, 10),
			b:       bitsx.NewZerosMatrixForTest(t, 1, 20),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.a.HammingDistance(tt.b)
			if tt.wantErr && err == nil {
				t.Fatal("エラーを期待したが、nilが返された")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("nilを期待したが、エラーが返された: %v", err)
			}

			if tt.wantErr {
				return
			}

			if got != tt.want {
				t.Errorf("値の不一致: got = %d, want = %d", got, tt.want)
			}
		})
	}
}

func TestMatrixDot(t *testing.T) {
	rng := rand.New(rand.NewPCG(13, 14))
	randMatA := bitsx.NewRandMatrixForTest(t, 4, 130, rng)
	randMatB := bitsx.NewRandMatrixForTest(t, 3, 130, rng)
	randWant := make([]int, randMatA.Rows()*randMatB.Rows())
	for aRow := range randMatA.Rows() {
		for bRow := range randMatB.Rows() {
			for col := range randMatA.Cols() {
				aBit, err := randMatA.Bit(aRow, col)
				if err != nil {
					t.Fatalf("Bit(%d, %d)で期待しないエラーが発生した: %v", aRow, col, err)
				}

				bBit, err := randMatB.Bit(bRow, col)
				if err != nil {
					t.Fatalf("Bit(%d, %d)で期待しないエラーが発生した: %v", bRow, col, err)
				}

				if aBit == bBit {
					randWant[aRow*randMatB.Rows()+bRow]++
				}
			}
		}
	}

	var zeroMatrixA, zeroMatrixB bitsx.Matrix

	tests := []struct {
		name    string
		a       *bitsx.Matrix
		b       *bitsx.Matrix
		want    []int
		wantErr bool
	}{
		{
			name: "正常_複数行_結果の並び順",
			a: bitsx.NewMatrixForTest(t, 5, [][]int{
				{0, 2, 4},
				{1, 3},
			}),
			b: bitsx.NewMatrixForTest(t, 5, [][]int{
				{0, 2, 4},
				{0, 1, 2, 3, 4},
				{2, 3},
			}),
			want:    []int{5, 3, 2, 0, 2, 3},
			wantErr: false,
		},
		{
			name:    "正常_全0と全1_端数あり",
			a:       bitsx.NewZerosMatrixForTest(t, 1, 70),
			b:       bitsx.NewOnesMatrixForTest(t, 1, 70),
			want:    []int{0},
			wantErr: false,
		},
		{
			name: "正常_64ビット境界直後",
			a:    bitsx.NewMatrixForTest(t, 65, [][]int{{64}}),
			b: bitsx.NewMatrixForTest(t, 65, [][]int{
				{64},
				{},
			}),
			want:    []int{65, 64},
			wantErr: false,
		},
		{
			name:    "正常_固定乱数_複数ワード",
			a:       randMatA,
			b:       randMatB,
			want:    randWant,
			wantErr: false,
		},
		{
			name:    "異常_列数不一致_同一stride",
			a:       bitsx.NewZerosMatrixForTest(t, 2, 65),
			b:       bitsx.NewZerosMatrixForTest(t, 3, 127),
			wantErr: true,
		},
		{
			name:    "異常_ゼロ値同士",
			a:       &zeroMatrixA,
			b:       &zeroMatrixB,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aClone := tt.a.Clone()
			bClone := tt.b.Clone()

			got, err := tt.a.Dot(tt.b)
			if tt.wantErr && err == nil {
				t.Fatal("エラーを期待したが、nilが返された")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("nilを期待したが、エラーが返された: %v", err)
			}

			if !tt.a.Equal(aClone) || !tt.b.Equal(bClone) {
				t.Error("元の行列が変更された")
			}

			if tt.wantErr {
				if got != nil {
					t.Errorf("nilを期待したが、結果が返された: %v", got)
				}
				return
			}

			if !slices.Equal(got, tt.want) {
				t.Errorf("結果の不一致: got = %v, want = %v", got, tt.want)
			}
		})
	}
}

func TestMatrixDotTernary(t *testing.T) {
	rng := rand.New(rand.NewPCG(15, 16))
	randValue := bitsx.NewRandMatrixForTest(t, 4, 130, rng)
	randSign := bitsx.NewRandMatrixForTest(t, 3, 130, rng)
	randNonZero := bitsx.NewRandMatrixForTest(t, 3, 130, rng)
	randWant := make([]int, randValue.Rows()*randSign.Rows())
	for valueRow := range randValue.Rows() {
		for signRow := range randSign.Rows() {
			for col := range randValue.Cols() {
				nonZeroBit, err := randNonZero.Bit(signRow, col)
				if err != nil {
					t.Fatalf("nonZero.Bit(%d, %d)で期待しないエラーが発生した: %v", signRow, col, err)
				}
				if nonZeroBit == 0 {
					continue
				}

				valueBit, err := randValue.Bit(valueRow, col)
				if err != nil {
					t.Fatalf("value.Bit(%d, %d)で期待しないエラーが発生した: %v", valueRow, col, err)
				}

				signBit, err := randSign.Bit(signRow, col)
				if err != nil {
					t.Fatalf("sign.Bit(%d, %d)で期待しないエラーが発生した: %v", signRow, col, err)
				}

				wantIdx := valueRow*randSign.Rows() + signRow
				if valueBit == signBit {
					randWant[wantIdx]++
				} else {
					randWant[wantIdx]--
				}
			}
		}
	}

	var zeroValue, zeroSign, zeroNonZero bitsx.Matrix

	tests := []struct {
		name    string
		value   *bitsx.Matrix
		sign    *bitsx.Matrix
		nonZero *bitsx.Matrix
		want    []int
		wantErr bool
	}{
		{
			name: "正常_複数行_正負と結果の並び順",
			value: bitsx.NewMatrixForTest(t, 5, [][]int{
				{0, 2, 4},
				{1, 3},
			}),
			sign: bitsx.NewMatrixForTest(t, 5, [][]int{
				{0, 2, 4},
				{0, 1, 2, 3, 4},
				{2, 3},
			}),
			nonZero: bitsx.NewMatrixForTest(t, 5, [][]int{
				{0, 1, 2, 3, 4},
				{0, 1, 4},
				{1, 2, 3},
			}),
			want:    []int{5, 1, 1, -5, -1, -1},
			wantErr: false,
		},
		{
			name:    "正常_nonZeroが全0_端数あり",
			value:   bitsx.NewOnesMatrixForTest(t, 1, 70),
			sign:    bitsx.NewZerosMatrixForTest(t, 1, 70),
			nonZero: bitsx.NewZerosMatrixForTest(t, 1, 70),
			want:    []int{0},
			wantErr: false,
		},
		{
			name:  "正常_64ビット境界直後とマスク外のsign",
			value: bitsx.NewMatrixForTest(t, 65, [][]int{{64}}),
			sign: bitsx.NewMatrixForTest(t, 65, [][]int{
				{0, 64},
				{0},
			}),
			nonZero: bitsx.NewMatrixForTest(t, 65, [][]int{
				{64},
				{64},
			}),
			want:    []int{1, -1},
			wantErr: false,
		},
		{
			name:    "正常_固定乱数_複数ワード",
			value:   randValue,
			sign:    randSign,
			nonZero: randNonZero,
			want:    randWant,
			wantErr: false,
		},
		{
			name:    "異常_valueとsignの列数不一致_同一stride",
			value:   bitsx.NewZerosMatrixForTest(t, 2, 65),
			sign:    bitsx.NewZerosMatrixForTest(t, 3, 127),
			nonZero: bitsx.NewZerosMatrixForTest(t, 3, 127),
			wantErr: true,
		},
		{
			name:    "異常_signとnonZeroの行数不一致",
			value:   bitsx.NewZerosMatrixForTest(t, 2, 65),
			sign:    bitsx.NewZerosMatrixForTest(t, 3, 65),
			nonZero: bitsx.NewZerosMatrixForTest(t, 2, 65),
			wantErr: true,
		},
		{
			name:    "異常_signとnonZeroの列数不一致",
			value:   bitsx.NewZerosMatrixForTest(t, 2, 65),
			sign:    bitsx.NewZerosMatrixForTest(t, 3, 65),
			nonZero: bitsx.NewZerosMatrixForTest(t, 3, 127),
			wantErr: true,
		},
		{
			name:    "異常_ゼロ値",
			value:   &zeroValue,
			sign:    &zeroSign,
			nonZero: &zeroNonZero,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valueClone := tt.value.Clone()
			signClone := tt.sign.Clone()
			nonZeroClone := tt.nonZero.Clone()

			got, err := tt.value.DotTernary(tt.sign, tt.nonZero)
			if tt.wantErr && err == nil {
				t.Fatal("エラーを期待したが、nilが返された")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("nilを期待したが、エラーが返された: %v", err)
			}

			if !tt.value.Equal(valueClone) || !tt.sign.Equal(signClone) || !tt.nonZero.Equal(nonZeroClone) {
				t.Error("元の行列が変更された")
			}

			if tt.wantErr {
				if got != nil {
					t.Errorf("nilを期待したが、結果が返された: %v", got)
				}
				return
			}

			if !slices.Equal(got, tt.want) {
				t.Errorf("結果の不一致: got = %v, want = %v", got, tt.want)
			}
		})
	}
}

func TestMatrixTranspose(t *testing.T) {
	var zeroMatrix bitsx.Matrix

	tests := []struct {
		name    string
		m       *bitsx.Matrix
		want    *bitsx.Matrix
		wantErr bool
	}{
		{
			name:    "正常_最小1x1",
			m:       bitsx.NewMatrixForTest(t, 1, [][]int{{0}}),
			want:    bitsx.NewMatrixForTest(t, 1, [][]int{{0}}),
			wantErr: false,
		},
		{
			// 1 0 1        1 0
			// 0 1 0   ->   0 1
			//              1 0
			name: "正常_非正方_2x3",
			m: bitsx.NewMatrixForTest(t, 3, [][]int{
				{0, 2},
				{1},
			}),
			want: bitsx.NewMatrixForTest(t, 2, [][]int{
				{0},
				{1},
				{0},
			}),
			wantErr: false,
		},
		{
			// 転置前後の両方に端数ビットがある（cols=70 -> 6ビット、cols=3 -> 3ビット）
			name:    "正常_全1_端数あり_3x70",
			m:       bitsx.NewOnesMatrixForTest(t, 3, 70),
			want:    bitsx.NewOnesMatrixForTest(t, 70, 3),
			wantErr: false,
		},
		{
			name:    "異常_ゼロ値",
			m:       &zeroMatrix,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.Transpose()
			if tt.wantErr && err == nil {
				t.Fatal("エラーを期待したが、nilが返された")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("nilを期待したが、エラーが返された: %v", err)
			}

			if tt.wantErr {
				if got != nil {
					t.Error("nilを期待したが、行列が返された")
				}
				return
			}

			if !got.Equal(tt.want) {
				t.Error("結果の不一致")
			}
		})
	}
}

// 転置が満たすべき性質を検証する。
// (1) 元の行列が変更されない (2) 形状が入れ替わる (3) 全ての(r, c)で m[r][c] == mT[c][r]
// (4) 端数ビットが0に保たれる (5) 二回転置すると元に戻る
func assertTransposed(t *testing.T, m *bitsx.Matrix) {
	t.Helper()

	clone := m.Clone()

	mT, err := m.Transpose()
	if err != nil {
		t.Fatalf("nilを期待したが、エラーが返された: %v", err)
	}

	if !m.Equal(clone) {
		t.Fatal("元の行列が変更された")
	}

	if mT.Rows() != m.Cols() || mT.Cols() != m.Rows() {
		t.Fatalf("形状の不一致: got = (%d, %d) want = (%d, %d)", mT.Rows(), mT.Cols(), m.Cols(), m.Rows())
	}

	for r := range m.Rows() {
		for c := range m.Cols() {
			want, err := m.Bit(r, c)
			if err != nil {
				t.Fatalf("nilを期待したが、エラーが返された: %v", err)
			}

			got, err := mT.Bit(c, r)
			if err != nil {
				t.Fatalf("nilを期待したが、エラーが返された: %v", err)
			}

			if got != want {
				t.Fatalf("Bit(%d, %d)の不一致: got = %d want = %d", c, r, got, want)
			}
		}
	}

	// 端数ビットはBitでは観測できない為、ワード単位で0であることを確認する
	stride := mT.Stride()
	tailMask := mT.TailMask()
	for r := range mT.Rows() {
		idx := (r * stride) + (stride - 1)
		word, err := mT.Word(idx)
		if err != nil {
			t.Fatalf("nilを期待したが、エラーが返された: %v", err)
		}

		if word&^tailMask != 0 {
			t.Errorf("Word(%d)の端数ビットが0でない: got = %#x tailMask = %#x", idx, word, tailMask)
		}
	}

	mTT, err := mT.Transpose()
	if err != nil {
		t.Fatalf("nilを期待したが、エラーが返された: %v", err)
	}

	if !mTT.Equal(m) {
		t.Error("二回転置しても元に戻らない")
	}
}

func TestMatrixTransposeShapes(t *testing.T) {
	rng := rand.New(rand.NewPCG(7, 8))

	// 64行ブロック単位で転置する為、64の境界前後で処理が変わる。
	// 行数と列数はそれぞれ独立にブロックの分割へ効く為、組み合わせを網羅する。
	// 1: 最小 / 63: 1ブロック未満 / 64: 1ブロック丁度 / 65: 1ブロック+1行
	// 127: 1ブロック+63行 / 128: 2ブロック丁度 / 129: 2ブロック+1行
	dims := []int{1, 63, 64, 65, 127, 128, 129}

	for _, rows := range dims {
		for _, cols := range dims {
			t.Run(fmt.Sprintf("%dx%d", rows, cols), func(t *testing.T) {
				assertTransposed(t, bitsx.NewRandMatrixForTest(t, rows, cols, rng))
			})
		}
	}
}

// MatrixWordContextは非公開フィールドを持ち、そのままでは比較できない為、期待値は専用の型で保持する。
type wantMatrixWordContext struct {
	row         int
	wordIndex   int
	colStart    int
	colEnd      int
	globalStart int
	globalEnd   int
	isTail      bool
}

func assertMatrixWordContext(t *testing.T, idx int, got bitsx.MatrixWordContext, want wantMatrixWordContext) {
	t.Helper()

	if got.Row != want.row {
		t.Errorf("contexts[%d].Rowの不一致: got = %d, want = %d", idx, got.Row, want.row)
	}

	if got.WordIndex != want.wordIndex {
		t.Errorf("contexts[%d].WordIndexの不一致: got = %d, want = %d", idx, got.WordIndex, want.wordIndex)
	}

	if got.ColStart != want.colStart {
		t.Errorf("contexts[%d].ColStartの不一致: got = %d, want = %d", idx, got.ColStart, want.colStart)
	}

	if got.ColEnd != want.colEnd {
		t.Errorf("contexts[%d].ColEndの不一致: got = %d, want = %d", idx, got.ColEnd, want.colEnd)
	}

	if got.GlobalStart != want.globalStart {
		t.Errorf("contexts[%d].GlobalStartの不一致: got = %d, want = %d", idx, got.GlobalStart, want.globalStart)
	}

	if got.GlobalEnd != want.globalEnd {
		t.Errorf("contexts[%d].GlobalEndの不一致: got = %d, want = %d", idx, got.GlobalEnd, want.globalEnd)
	}

	if got.IsTail != want.isTail {
		t.Errorf("contexts[%d].IsTailの不一致: got = %t, want = %t", idx, got.IsTail, want.isTail)
	}
}

func TestMatrixScanRowsWord(t *testing.T) {
	dummyErr := errors.New("dummy error")

	tests := []struct {
		name         string
		m            *bitsx.Matrix
		rowIdxs      []int
		callback     func(ctx bitsx.MatrixWordContext) error
		wantContexts []wantMatrixWordContext
		wantErr      bool
		// コールバックのエラーがそのまま返る場合に、その値を指定する
		wantErrIs error
	}{
		{
			// cols=70 (stride=2) の為、各行の2ワード目が端数ワードになる
			name:    "正常_全行スキャン_端数あり",
			m:       bitsx.NewZerosMatrixForTest(t, 3, 70),
			rowIdxs: nil,
			wantContexts: []wantMatrixWordContext{
				{row: 0, wordIndex: 0, colStart: 0, colEnd: 64, globalStart: 0, globalEnd: 64, isTail: false},
				{row: 0, wordIndex: 1, colStart: 64, colEnd: 70, globalStart: 64, globalEnd: 70, isTail: true},
				{row: 1, wordIndex: 2, colStart: 0, colEnd: 64, globalStart: 70, globalEnd: 134, isTail: false},
				{row: 1, wordIndex: 3, colStart: 64, colEnd: 70, globalStart: 134, globalEnd: 140, isTail: true},
				{row: 2, wordIndex: 4, colStart: 0, colEnd: 64, globalStart: 140, globalEnd: 204, isTail: false},
				{row: 2, wordIndex: 5, colStart: 64, colEnd: 70, globalStart: 204, globalEnd: 210, isTail: true},
			},
			wantErr: false,
		},
		{
			// cols=128 は64の倍数の為、最終ワードも端数ワードにはならない
			name:    "正常_全行スキャン_端数なし",
			m:       bitsx.NewZerosMatrixForTest(t, 2, 128),
			rowIdxs: nil,
			wantContexts: []wantMatrixWordContext{
				{row: 0, wordIndex: 0, colStart: 0, colEnd: 64, globalStart: 0, globalEnd: 64, isTail: false},
				{row: 0, wordIndex: 1, colStart: 64, colEnd: 128, globalStart: 64, globalEnd: 128, isTail: false},
				{row: 1, wordIndex: 2, colStart: 0, colEnd: 64, globalStart: 128, globalEnd: 192, isTail: false},
				{row: 1, wordIndex: 3, colStart: 64, colEnd: 128, globalStart: 192, globalEnd: 256, isTail: false},
			},
			wantErr: false,
		},
		{
			// 指定した行を、指定した順序で走査する
			name:    "正常_行インデックス指定",
			m:       bitsx.NewZerosMatrixForTest(t, 3, 70),
			rowIdxs: []int{2, 0},
			wantContexts: []wantMatrixWordContext{
				{row: 2, wordIndex: 4, colStart: 0, colEnd: 64, globalStart: 140, globalEnd: 204, isTail: false},
				{row: 2, wordIndex: 5, colStart: 64, colEnd: 70, globalStart: 204, globalEnd: 210, isTail: true},
				{row: 0, wordIndex: 0, colStart: 0, colEnd: 64, globalStart: 0, globalEnd: 64, isTail: false},
				{row: 0, wordIndex: 1, colStart: 64, colEnd: 70, globalStart: 64, globalEnd: 70, isTail: true},
			},
			wantErr: false,
		},
		{
			name:    "正常_行インデックス重複",
			m:       bitsx.NewZerosMatrixForTest(t, 3, 70),
			rowIdxs: []int{1, 1},
			wantContexts: []wantMatrixWordContext{
				{row: 1, wordIndex: 2, colStart: 0, colEnd: 64, globalStart: 70, globalEnd: 134, isTail: false},
				{row: 1, wordIndex: 3, colStart: 64, colEnd: 70, globalStart: 134, globalEnd: 140, isTail: true},
				{row: 1, wordIndex: 2, colStart: 0, colEnd: 64, globalStart: 70, globalEnd: 134, isTail: false},
				{row: 1, wordIndex: 3, colStart: 64, colEnd: 70, globalStart: 134, globalEnd: 140, isTail: true},
			},
			wantErr: false,
		},
		{
			// nilは全行走査を意味するが、空スライスは1行も走査しない
			name:         "正常_行インデックスが空",
			m:            bitsx.NewZerosMatrixForTest(t, 3, 70),
			rowIdxs:      []int{},
			wantContexts: nil,
			wantErr:      false,
		},
		{
			name:         "異常_行インデックス下限未満",
			m:            bitsx.NewZerosMatrixForTest(t, 3, 70),
			rowIdxs:      []int{-1},
			wantContexts: nil,
			wantErr:      true,
		},
		{
			name:         "異常_行インデックス上限超過",
			m:            bitsx.NewZerosMatrixForTest(t, 3, 70),
			rowIdxs:      []int{3},
			wantContexts: nil,
			wantErr:      true,
		},
		{
			// エラーを返した時点で中断する為、渡されるコンテキストは最初の1つだけ
			name:     "異常_コールバックエラーで早期中断",
			m:        bitsx.NewZerosMatrixForTest(t, 3, 70),
			rowIdxs:  nil,
			callback: func(ctx bitsx.MatrixWordContext) error { return dummyErr },
			wantContexts: []wantMatrixWordContext{
				{row: 0, wordIndex: 0, colStart: 0, colEnd: 64, globalStart: 0, globalEnd: 64, isTail: false},
			},
			wantErr:   true,
			wantErrIs: dummyErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotContexts []bitsx.MatrixWordContext

			err := tt.m.ScanRowsWord(tt.rowIdxs, func(ctx bitsx.MatrixWordContext) error {
				gotContexts = append(gotContexts, ctx)
				if tt.callback != nil {
					return tt.callback(ctx)
				}
				return nil
			})

			if tt.wantErr && err == nil {
				t.Fatal("エラーを期待したが、nilが返された")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("nilを期待したが、エラーが返された: %v", err)
			}

			if tt.wantErrIs != nil && !errors.Is(err, tt.wantErrIs) {
				t.Errorf("エラーの不一致: got = %v, want = %v", err, tt.wantErrIs)
			}

			if len(gotContexts) != len(tt.wantContexts) {
				t.Fatalf("コールバック実行回数の不一致: got = %d, want = %d", len(gotContexts), len(tt.wantContexts))
			}

			for idx, got := range gotContexts {
				assertMatrixWordContext(t, idx, got, tt.wantContexts[idx])
			}
		})
	}
}

func TestMatrixWordContextScanBits(t *testing.T) {
	dummyErr := errors.New("dummy error")

	// ScanRowsWord 経由でコンテキストを取得
	m := bitsx.NewZerosMatrixForTest(t, 3, 70)
	var nonTailCtx, tailCtx bitsx.MatrixWordContext

	err := m.ScanRowsWord([]int{1}, func(ctx bitsx.MatrixWordContext) error {
		if ctx.IsTail {
			tailCtx = ctx
		} else {
			nonTailCtx = ctx
		}
		return nil
	})
	if err != nil {
		t.Fatalf("nilを期待したが、エラーが返された: %v", err)
	}

	tests := []struct {
		name      string
		ctx       bitsx.MatrixWordContext
		callback  func(i, col, colT int) error
		wantCount int
		wantErr   bool
		// コールバックのエラーがそのまま返る場合に、その値を指定する
		wantErrIs error
	}{
		{
			name:      "正常_非tailワード_64ビット",
			ctx:       nonTailCtx, // Row=1, ColStart=0, ColEnd=64
			wantCount: 64,
			wantErr:   false,
		},
		{
			name:      "正常_tailワード_6ビット",
			ctx:       tailCtx, // Row=1, ColStart=64, ColEnd=70
			wantCount: 6,
			wantErr:   false,
		},
		{
			// ゼロ値はColStartとColEndが共に0の為、一度も呼ばれない
			name:      "正常_ゼロ値",
			ctx:       bitsx.MatrixWordContext{},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:      "異常_コールバックエラーで早期中断",
			ctx:       nonTailCtx,
			callback:  func(i, col, colT int) error { return dummyErr },
			wantCount: 1,
			wantErr:   true,
			wantErrIs: dummyErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := 0
			var bitsInfo [][3]int

			cb := tt.callback
			if cb == nil {
				cb = func(i, col, colT int) error {
					bitsInfo = append(bitsInfo, [3]int{i, col, colT})
					return nil
				}
			}

			err := tt.ctx.ScanBits(func(i, col, colT int) error {
				count++
				return cb(i, col, colT)
			})

			if tt.wantErr && err == nil {
				t.Fatal("エラーを期待したが、nilが返された")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("nilを期待したが、エラーが返された: %v", err)
			}

			if tt.wantErrIs != nil && !errors.Is(err, tt.wantErrIs) {
				t.Errorf("エラーの不一致: got = %v, want = %v", err, tt.wantErrIs)
			}

			if count != tt.wantCount {
				t.Errorf("ScanBits実行回数の不一致: got = %d, want = %d", count, tt.wantCount)
			}

			if tt.wantErr {
				return
			}

			// i, col, colT の値の整合性検証。
			// colTは転置後の通し位置(列 * 行数 + 行)を表す。
			// 走査対象のctxはmから取得している為、行数はm.Rows()と等しい。
			for idx, info := range bitsInfo {
				i, col, colT := info[0], info[1], info[2]
				wantCol := tt.ctx.ColStart + idx
				wantColT := (wantCol * m.Rows()) + tt.ctx.Row

				if i != idx {
					t.Errorf("bitsInfo[%d].iの不一致: got = %d, want = %d", idx, i, idx)
				}
				if col != wantCol {
					t.Errorf("bitsInfo[%d].colの不一致: got = %d, want = %d", idx, col, wantCol)
				}
				if colT != wantColT {
					t.Errorf("bitsInfo[%d].colTの不一致: got = %d, want = %d", idx, colT, wantColT)
				}
			}
		})
	}
}
