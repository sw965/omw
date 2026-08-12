package bitsx_test

import (
	"bytes"
	"encoding/gob"
	"math"
	"math/rand/v2"
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

func TestNewRandMatrix(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))

	tests := []struct {
		name    string
		rows    int
		cols    int
		k       int
		rng     *rand.Rand
		wantErr bool
	}{
		{
			name:    "正常_kが0",
			rows:    3,
			cols:    100,
			k:       0,
			rng:     rng,
			wantErr: false,
		},
		{
			name:    "正常_kが負",
			rows:    3,
			cols:    100,
			k:       -2,
			rng:     rng,
			wantErr: false,
		},
		{
			name:    "正常_kが正",
			rows:    3,
			cols:    100,
			k:       2,
			rng:     rng,
			wantErr: false,
		},
		{
			name:    "異常_rowsが0以下",
			rows:    0,
			cols:    10,
			k:       0,
			rng:     rng,
			wantErr: true,
		},
		{
			name:    "異常_colsが0以下",
			rows:    10,
			cols:    0,
			k:       0,
			rng:     rng,
			wantErr: true,
		},
		{
			name:    "異常_colsの桁あふれ",
			rows:    1,
			cols:    math.MaxInt,
			k:       0,
			rng:     rng,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := bitsx.NewRandMatrix(tt.rows, tt.cols, tt.k, tt.rng)
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
			if c := m.OnesCount(); c < 0 || c > totalBits {
				t.Errorf("OnesCountの不一致: got = %d, want = 0 <= c <= %d", c, totalBits)
			}
		})
	}
}

func TestNewRandMatrixStatistics(t *testing.T) {
	rng := rand.New(rand.NewPCG(42, 100))

	const (
		rows = 1000
		cols = 1000
	)
	totalBits := float64(rows * cols)

	tests := []struct {
		name  string
		k     int
		wantP float64
		tol   float64
	}{
		{
			// N=10^6, p=0.5, σ=0.0005（tol=0.015 は 30σ。正しく実装されていれば収まる確率 ≒ 100%）
			name:  "kが0_確率0.5",
			k:     0,
			wantP: 0.5,
			tol:   0.015,
		},
		{
			// N=10^6, p=0.25, σ≈0.000433（tol=0.015 は 34.6σ。正しく実装されていれば収まる確率 ≒ 100%）
			name:  "kが-1_確率0.25",
			k:     -1,
			wantP: 0.25,
			tol:   0.015,
		},
		{
			// N=10^6, p=0.75, σ≈0.000433（tol=0.015 は 34.6σ。正しく実装されていれば収まる確率 ≒ 100%）
			name:  "kが1_確率0.75",
			k:     1,
			wantP: 0.75,
			tol:   0.015,
		},
		{
			// N=10^6, p=0.125, σ≈0.000331（tol=0.015 は 45.3σ。正しく実装されていれば収まる確率 ≒ 100%）
			name:  "kが-2_確率0.125",
			k:     -2,
			wantP: 0.125,
			tol:   0.015,
		},
		{
			// N=10^6, p=0.875, σ≈0.000331（tol=0.015 は 45.3σ。正しく実装されていれば収まる確率 ≒ 100%）
			name:  "kが2_確率0.875",
			k:     2,
			wantP: 0.875,
			tol:   0.015,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := bitsx.NewRandMatrix(rows, cols, tt.k, rng)
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

func FuzzNewRandMatrix(f *testing.F) {
	seeds := []struct {
		rows8        uint8
		cols16       uint16
		k8           int8
		seed1, seed2 uint64
	}{
		{3, 100, 0, 1, 2},
		{1, 64, -2, 10, 20},
		{10, 1, 2, 100, 200},
	}
	for _, s := range seeds {
		f.Add(s.rows8, s.cols16, s.k8, s.seed1, s.seed2)
	}

	f.Fuzz(func(t *testing.T, rows8 uint8, cols16 uint16, k8 int8, seed1, seed2 uint64) {
		rows := int(rows8)
		cols := int(cols16)
		k := int(k8)

		rng := rand.New(rand.NewPCG(seed1, seed2))
		m, err := bitsx.NewRandMatrix(rows, cols, k, rng)
		if rows <= 0 || cols <= 0 {
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
				t.Errorf("OnesCount()の不一致: got = %d, want = %d", got, tt.want)
			}
		})
	}
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
		b.Fatalf("%v", err)
	}

	for b.Loop() {
		if _, err := m.Transpose(); err != nil {
			b.Fatalf("%v", err)
		}
	}
}
