package bitsx_test

import (
	"math/rand/v2"
	"testing"

	"github.com/sw965/omw/mathx"
	"github.com/sw965/omw/mathx/bitsx"
)

// 端数ビット(cols % 64 の範囲外)はBitでは観測できない為、ワード単位で0であることを確認する。
func assertTailBitsZero(t *testing.T, m *bitsx.Matrix) {
	t.Helper()

	stride := m.Stride()
	tailMask := m.TailMask()
	for r := range m.Rows() {
		idx := (r * stride) + (stride - 1)
		word, err := m.Word(idx)
		if err != nil {
			t.Fatalf("nilを期待したが、エラーが返された: %v", err)
		}

		if word&^tailMask != 0 {
			t.Errorf("Word(%d)の端数ビットが0でない: got = %#x tailMask = %#x", idx, word, tailMask)
		}
	}
}

func TestNewThermometerMatrices(t *testing.T) {
	tests := []struct {
		name    string
		n       int
		rows    int
		cols    int
		want    bitsx.Matrices
		wantErr bool
	}{
		{
			// totalBits=4なので、1の数は i*4/2 -> 0, 2, 4
			name: "正常_1行",
			n:    3,
			rows: 1,
			cols: 4,
			want: bitsx.Matrices{
				bitsx.NewMatrixForTest(t, 4, [][]int{{}}),
				bitsx.NewMatrixForTest(t, 4, [][]int{{0, 1}}),
				bitsx.NewMatrixForTest(t, 4, [][]int{{0, 1, 2, 3}}),
			},
			wantErr: false,
		},
		{
			// totalBits=6なので、1の数は i*6/3 -> 0, 2, 4, 6
			// 4の時に行をまたぐ為、1が行優先の通し順で並ぶことを確認できる
			name: "正常_複数行_行をまたぐ",
			n:    4,
			rows: 2,
			cols: 3,
			want: bitsx.Matrices{
				bitsx.NewMatrixForTest(t, 3, [][]int{{}, {}}),
				bitsx.NewMatrixForTest(t, 3, [][]int{{0, 1}, {}}),
				bitsx.NewMatrixForTest(t, 3, [][]int{{0, 1, 2}, {0}}),
				bitsx.NewMatrixForTest(t, 3, [][]int{{0, 1, 2}, {0, 1, 2}}),
			},
			wantErr: false,
		},
		{
			// totalBits=5なので、1の数は i*5/3 -> 0, 1, 3, 5 (切り捨て)
			name: "正常_割り切れない",
			n:    4,
			rows: 1,
			cols: 5,
			want: bitsx.Matrices{
				bitsx.NewMatrixForTest(t, 5, [][]int{{}}),
				bitsx.NewMatrixForTest(t, 5, [][]int{{0}}),
				bitsx.NewMatrixForTest(t, 5, [][]int{{0, 1, 2}}),
				bitsx.NewMatrixForTest(t, 5, [][]int{{0, 1, 2, 3, 4}}),
			},
			wantErr: false,
		},
		{
			// totalBits=1なので、1の数は i*1/1 -> 0, 1
			name: "正常_最小",
			n:    2,
			rows: 1,
			cols: 1,
			want: bitsx.Matrices{
				bitsx.NewMatrixForTest(t, 1, [][]int{{}}),
				bitsx.NewMatrixForTest(t, 1, [][]int{{0}}),
			},
			wantErr: false,
		},
		{
			// totalBits=140なので、1の数は i*140/4 -> 0, 35, 70, 105, 140
			// 1行が2ワードにまたがる為、ワード単位の書き込み先を確認できる
			name: "正常_複数ワード_端数あり",
			n:    5,
			rows: 2,
			cols: 70,
			want: bitsx.Matrices{
				bitsx.NewPrefixOnesMatrixForTest(t, 2, 70, 0),
				bitsx.NewPrefixOnesMatrixForTest(t, 2, 70, 35),
				bitsx.NewPrefixOnesMatrixForTest(t, 2, 70, 70),
				bitsx.NewPrefixOnesMatrixForTest(t, 2, 70, 105),
				bitsx.NewPrefixOnesMatrixForTest(t, 2, 70, 140),
			},
			wantErr: false,
		},
		{
			// totalBits=100なので、1の数は i*100/4 -> 0, 25, 50, 75, 100
			// 75は2ワード目(列64〜99)の途中で途切れる
			name: "正常_複数ワード_ワードの途中で途切れる",
			n:    5,
			rows: 1,
			cols: 100,
			want: bitsx.Matrices{
				bitsx.NewPrefixOnesMatrixForTest(t, 1, 100, 0),
				bitsx.NewPrefixOnesMatrixForTest(t, 1, 100, 25),
				bitsx.NewPrefixOnesMatrixForTest(t, 1, 100, 50),
				bitsx.NewPrefixOnesMatrixForTest(t, 1, 100, 75),
				bitsx.NewPrefixOnesMatrixForTest(t, 1, 100, 100),
			},
			wantErr: false,
		},
		{
			name:    "異常_nが2未満",
			n:       1,
			rows:    2,
			cols:    3,
			wantErr: true,
		},
		{
			name:    "異常_nが0",
			n:       0,
			rows:    2,
			cols:    3,
			wantErr: true,
		},
		{
			name:    "異常_rowsが0以下",
			n:       2,
			rows:    0,
			cols:    3,
			wantErr: true,
		},
		{
			name:    "異常_colsが0以下",
			n:       2,
			rows:    2,
			cols:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := bitsx.NewThermometerMatrices(tt.n, tt.rows, tt.cols)
			if tt.wantErr && err == nil {
				t.Fatal("エラーを期待したが、nilが返された")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("nilを期待したが、エラーが返された: %v", err)
			}

			if tt.wantErr {
				if got != nil {
					t.Error("nilを期待したが、結果が返された")
				}
				return
			}

			if len(got) != len(tt.want) {
				t.Fatalf("要素数の不一致: got = %d, want = %d", len(got), len(tt.want))
			}

			for i := range got {
				if !got[i].Equal(tt.want[i]) {
					t.Errorf("ms[%d]の不一致", i)
				}
			}
		})
	}
}

func TestNewRFFMatrices(t *testing.T) {
	tests := []struct {
		name    string
		n       int
		rows    int
		cols    int
		sigma   float32
		wantErr bool
	}{
		{
			name:    "正常_端数あり",
			n:       3,
			rows:    2,
			cols:    70,
			sigma:   1.0,
			wantErr: false,
		},
		{
			name:    "正常_最小",
			n:       2,
			rows:    1,
			cols:    1,
			sigma:   1.0,
			wantErr: false,
		},
		{
			name:    "異常_nが2未満",
			n:       1,
			rows:    2,
			cols:    70,
			sigma:   1.0,
			wantErr: true,
		},
		{
			name:    "異常_nが0",
			n:       0,
			rows:    2,
			cols:    70,
			sigma:   1.0,
			wantErr: true,
		},
		{
			name:    "異常_rowsが0以下",
			n:       2,
			rows:    0,
			cols:    70,
			sigma:   1.0,
			wantErr: true,
		},
		{
			name:    "異常_colsが0以下",
			n:       2,
			rows:    2,
			cols:    0,
			sigma:   1.0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rng := rand.New(rand.NewPCG(11, 12))

			got, err := bitsx.NewRFFMatrices(tt.n, tt.rows, tt.cols, tt.sigma, rng)
			if tt.wantErr && err == nil {
				t.Fatal("エラーを期待したが、nilが返された")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("nilを期待したが、エラーが返された: %v", err)
			}

			if tt.wantErr {
				if got != nil {
					t.Error("nilを期待したが、結果が返された")
				}
				return
			}

			if len(got) != tt.n {
				t.Fatalf("要素数の不一致: got = %d, want = %d", len(got), tt.n)
			}

			for i, m := range got {
				if m.Rows() != tt.rows || m.Cols() != tt.cols {
					t.Errorf("ms[%d]の形状の不一致: got = (%d, %d) want = (%d, %d)", i, m.Rows(), m.Cols(), tt.rows, tt.cols)
				}
				assertTailBitsZero(t, m)
			}
		})
	}
}

func TestNewRFFMatricesProperties(t *testing.T) {
	const (
		n    = 4
		rows = 2
		cols = 70
	)
	totalBits := rows * cols

	newRng := func() *rand.Rand {
		return rand.New(rand.NewPCG(21, 22))
	}

	// sigmaが0だと周波数が全て0になり、cos(位相)だけで決まる為、iに依らず全て同じ行列になる
	zeroSigma, err := bitsx.NewRFFMatrices(n, rows, cols, 0, newRng())
	if err != nil {
		t.Fatalf("nilを期待したが、エラーが返された: %v", err)
	}

	for i := 1; i < n; i++ {
		if !zeroSigma[i].Equal(zeroSigma[0]) {
			t.Errorf("sigmaが0なのに、ms[%d]がms[0]と一致しない", i)
		}
	}

	// 全0や全1では上記の一致が自明になってしまう為、0と1が混在していることを確認する
	if c := zeroSigma[0].OnesCount(); c <= 0 || c >= totalBits {
		t.Errorf("OnesCountの不一致: got = %d, want = 0 < c < %d", c, totalBits)
	}

	// ms[0]は u=0 の為、周波数に掛かる係数が0になり、sigmaの値に依存しない
	largeSigma, err := bitsx.NewRFFMatrices(n, rows, cols, 5.0, newRng())
	if err != nil {
		t.Fatalf("nilを期待したが、エラーが返された: %v", err)
	}

	if !largeSigma[0].Equal(zeroSigma[0]) {
		t.Error("ms[0]がsigmaの値によって変化した")
	}

	// sigmaが大きければ、iが進むにつれてms[0]から離れる
	allSame := true
	for i := 1; i < n; i++ {
		if !largeSigma[i].Equal(largeSigma[0]) {
			allSame = false
			break
		}
	}

	if allSame {
		t.Error("sigmaを大きくしても、全ての行列がms[0]と同一になった")
	}

	// 同じ乱数列を与えれば再現する
	again, err := bitsx.NewRFFMatrices(n, rows, cols, 5.0, newRng())
	if err != nil {
		t.Fatalf("nilを期待したが、エラーが返された: %v", err)
	}

	for i := range n {
		if !again[i].Equal(largeSigma[i]) {
			t.Errorf("同じ乱数列なのに、ms[%d]が再現しない", i)
		}
	}
}

func TestMatricesETFCost(t *testing.T) {
	tests := []struct {
		name    string
		ms      bitsx.Matrices
		want    float32
		wantErr bool
	}{
		{
			// 距離は[4]のみ。合計4、分散0なので、コストは -4
			name: "正常_2つ_全0と全1",
			ms: bitsx.Matrices{
				bitsx.NewZerosMatrixForTest(t, 1, 4),
				bitsx.NewOnesMatrixForTest(t, 1, 4),
			},
			want:    -4.0,
			wantErr: false,
		},
		{
			// 距離は[0]のみ。合計0、分散0なので、コストは0
			name: "正常_2つ_同一",
			ms: bitsx.Matrices{
				bitsx.NewOnesMatrixForTest(t, 1, 4),
				bitsx.NewOnesMatrixForTest(t, 1, 4),
			},
			want:    0.0,
			wantErr: false,
		},
		{
			// 距離は[4, 2, 2]。合計8、平均8/3、分散8/9なので、コストは -8 + 8/9
			name: "正常_3つ_分散あり",
			ms: bitsx.Matrices{
				bitsx.NewMatrixForTest(t, 4, [][]int{{}}),
				bitsx.NewMatrixForTest(t, 4, [][]int{{0, 1, 2, 3}}),
				bitsx.NewMatrixForTest(t, 4, [][]int{{0, 1}}),
			},
			want:    -8.0 + 8.0/9.0,
			wantErr: false,
		},
		{
			// 端数ビットは距離に含まれない。距離は[70]のみなので、コストは -70
			name: "正常_端数あり",
			ms: bitsx.Matrices{
				bitsx.NewZerosMatrixForTest(t, 1, 70),
				bitsx.NewOnesMatrixForTest(t, 1, 70),
			},
			want:    -70.0,
			wantErr: false,
		},
		{
			name:    "異常_空",
			ms:      nil,
			wantErr: true,
		},
		{
			name: "異常_1つ",
			ms: bitsx.Matrices{
				bitsx.NewZerosMatrixForTest(t, 1, 4),
			},
			wantErr: true,
		},
		{
			name: "異常_形状不一致",
			ms: bitsx.Matrices{
				bitsx.NewZerosMatrixForTest(t, 1, 4),
				bitsx.NewZerosMatrixForTest(t, 1, 5),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ms.ETFCost()
			if tt.wantErr && err == nil {
				t.Fatal("エラーを期待したが、nilが返された")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("nilを期待したが、エラーが返された: %v", err)
			}

			if tt.wantErr {
				if got != 0.0 {
					t.Errorf("0を期待したが、値が返された: %f", got)
				}
				return
			}

			if !mathx.ApproxEqual(got, tt.want, float32(1e-4)) {
				t.Errorf("コストの不一致: got = %f, want = %f", got, tt.want)
			}
		})
	}
}

func TestNewETFMatrices(t *testing.T) {
	tests := []struct {
		name    string
		n       int
		rows    int
		cols    int
		iters   int
		wantErr bool
	}{
		{
			name:    "正常_反復なし",
			n:       3,
			rows:    2,
			cols:    70,
			iters:   0,
			wantErr: false,
		},
		{
			name:    "正常_反復あり",
			n:       3,
			rows:    2,
			cols:    70,
			iters:   50,
			wantErr: false,
		},
		{
			name:    "異常_nが2未満",
			n:       1,
			rows:    2,
			cols:    70,
			iters:   10,
			wantErr: true,
		},
		{
			name:    "異常_nが0",
			n:       0,
			rows:    2,
			cols:    70,
			iters:   10,
			wantErr: true,
		},
		{
			name:    "異常_rowsが0以下",
			n:       2,
			rows:    0,
			cols:    70,
			iters:   10,
			wantErr: true,
		},
		{
			name:    "異常_colsが0以下",
			n:       2,
			rows:    2,
			cols:    0,
			iters:   10,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rng := rand.New(rand.NewPCG(31, 32))

			got, err := bitsx.NewETFMatrices(tt.n, tt.rows, tt.cols, tt.iters, rng)
			if tt.wantErr && err == nil {
				t.Fatal("エラーを期待したが、nilが返された")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("nilを期待したが、エラーが返された: %v", err)
			}

			if tt.wantErr {
				if got != nil {
					t.Error("nilを期待したが、結果が返された")
				}
				return
			}

			if len(got) != tt.n {
				t.Fatalf("要素数の不一致: got = %d, want = %d", len(got), tt.n)
			}

			for i, m := range got {
				if m.Rows() != tt.rows || m.Cols() != tt.cols {
					t.Errorf("ms[%d]の形状の不一致: got = (%d, %d) want = (%d, %d)", i, m.Rows(), m.Cols(), tt.rows, tt.cols)
				}
				assertTailBitsZero(t, m)
			}
		})
	}
}

func TestNewETFMatricesOptimization(t *testing.T) {
	const (
		n     = 4
		rows  = 4
		cols  = 64
		iters = 1000
	)

	newRng := func() *rand.Rand {
		return rand.New(rand.NewPCG(33, 34))
	}

	// itersが0なら初期状態のまま返る為、同じ乱数列でNewRandMatrixをn回呼んだ結果と一致する
	initial, err := bitsx.NewETFMatrices(n, rows, cols, 0, newRng())
	if err != nil {
		t.Fatalf("nilを期待したが、エラーが返された: %v", err)
	}

	rng := newRng()
	for i := range n {
		want := bitsx.NewRandMatrixForTest(t, rows, cols, rng)
		if !initial[i].Equal(want) {
			t.Errorf("ms[%d]の不一致: 初期状態が乱数行列と一致しない", i)
		}
	}

	initialCost, err := initial.ETFCost()
	if err != nil {
		t.Fatalf("nilを期待したが、エラーが返された: %v", err)
	}

	optimized, err := bitsx.NewETFMatrices(n, rows, cols, iters, newRng())
	if err != nil {
		t.Fatalf("nilを期待したが、エラーが返された: %v", err)
	}

	optimizedCost, err := optimized.ETFCost()
	if err != nil {
		t.Fatalf("nilを期待したが、エラーが返された: %v", err)
	}

	// コストが下がる場合だけ反転を採用する為、初期状態より悪化することはない。
	// 反復回数が十分にあれば、必ず改善する
	if optimizedCost >= initialCost {
		t.Errorf("コストが改善していない: got = %f, initial = %f", optimizedCost, initialCost)
	}

	// 同じ乱数列を与えれば再現する
	again, err := bitsx.NewETFMatrices(n, rows, cols, iters, newRng())
	if err != nil {
		t.Fatalf("nilを期待したが、エラーが返された: %v", err)
	}

	for i := range n {
		if !again[i].Equal(optimized[i]) {
			t.Errorf("同じ乱数列なのに、ms[%d]が再現しない", i)
		}
	}
}
