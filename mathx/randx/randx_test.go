// Package randx_test は、randx パッケージの検証を行う。
// 乱数の生成そのものが主目的のパッケージである為、統計テストを中心に行う。
package randx_test

import (
	"math"
	"slices"
	"testing"

	"github.com/sw965/omw/mathx"
	"github.com/sw965/omw/mathx/randx"
	"github.com/sw965/omw/slicesx"
)

const sampleN = 10000

func TestNewPCG(t *testing.T) {
	t.Run("異なる乱数列", func(t *testing.T) {
		rng1 := randx.NewPCG()
		rng2 := randx.NewPCG()
		for range 10 {
			if rng1.Uint64() != rng2.Uint64() {
				return
			}
		}
		t.Error("10回中すべて同じ値が生成されました")
	})
}

func TestNewPCGs(t *testing.T) {
	t.Run("生成個数", func(t *testing.T) {
		tests := []struct {
			name string
			n    int
		}{
			{name: "0個", n: 0},
			{name: "1個", n: 1},
			{name: "3個", n: 3},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				rs, err := randx.NewPCGs(tt.n)
				if err != nil {
					t.Fatalf("nilを期待したが、エラーが返された: %v", err)
				}
				if len(rs) != tt.n {
					t.Errorf("要素数の不一致: got = %d want = %d", len(rs), tt.n)
				}
				for i, r := range rs {
					if r == nil {
						t.Errorf("インデックス %d の *rand.Rand が nil です", i)
					}
				}
			})
		}
	})

	t.Run("異なる乱数列", func(t *testing.T) {
		rs, err := randx.NewPCGs(2)
		if err != nil {
			t.Fatalf("nilを期待したが、エラーが返された: %v", err)
		}
		for range 10 {
			if rs[0].Uint64() != rs[1].Uint64() {
				return
			}
		}
		t.Error("10回中すべて同じ値が生成されました")
	})
}

func TestNewPCGs_Error(t *testing.T) {
	tests := []struct {
		name string
		n    int
	}{
		{
			name: "負の数",
			n:    -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs, err := randx.NewPCGs(tt.n)
			if rs != nil {
				t.Errorf("エラー時の戻り値がnilではない: got = %v", rs)
			}
			if err == nil {
				t.Error("エラーを期待したが、nilが返された")
			}
		})
	}
}

func TestIntRange_Statistics(t *testing.T) {
	rng := randx.NewPCG()
	const eps = 0.1

	tests := []struct {
		min int
		max int
	}{
		{min: 0, max: 10},
		{min: -5, max: 5},
		{min: -10, max: -1},
	}

	for _, tt := range tests {
		got := make([]int, sampleN)
		var err error
		for i := range sampleN {
			got[i], err = randx.IntRange(tt.min, tt.max, rng)
			if err != nil {
				t.Fatalf("nilを期待したが、エラーが返された: %v", err)
			}
		}

		// 生成値は [min, max) に収まり、min と max-1 は十分な試行で観測されるはず
		wantMin := tt.min
		wantMax := tt.max - 1
		wantAvg := float64(wantMin+wantMax) / 2.0

		gotMin := slices.Min(got)
		gotMax := slices.Max(got)
		gotSum := 0
		for _, v := range got {
			gotSum += v
		}
		gotAvg := float64(gotSum) / float64(sampleN)

		if gotMin != wantMin {
			t.Errorf("最小値の不一致: got = %d want = %d", gotMin, wantMin)
		}
		if gotMax != wantMax {
			t.Errorf("最大値の不一致: got = %d want = %d", gotMax, wantMax)
		}
		if !mathx.ApproxEqual(gotAvg, wantAvg, eps) {
			t.Errorf("平均の不一致: got = %g want = %g(±%g)", gotAvg, wantAvg, eps)
		}
	}
}

func TestIntRange_Error(t *testing.T) {
	rng := randx.NewPCG()
	tests := []struct {
		name string
		min  int
		max  int
	}{
		{
			name: "minがmaxより大きい",
			min:  5,
			max:  3,
		},
		{
			name: "境界_minとmaxが等しい",
			min:  5,
			max:  5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := randx.IntRange(tt.min, tt.max, rng)
			if got != 0 {
				t.Errorf("エラー時の戻り値がゼロ値ではない: got = %d", got)
			}
			if err == nil {
				t.Error("エラーを期待したが、nilが返された")
			}
		})
	}
}

func TestIndexByWeights_Statistics(t *testing.T) {
	rng := randx.NewPCG()
	const eps = 0.03

	tests := []struct {
		name       string
		ws         []float64
		wantRatios []float64
	}{
		{
			name:       "合計が1",
			ws:         []float64{0.2, 0.3, 0.1, 0.4},
			wantRatios: []float64{0.2, 0.3, 0.1, 0.4},
		},
		{
			name:       "合計が1ではない",
			ws:         []float64{0.1, 0.1},
			wantRatios: []float64{0.5, 0.5},
		},
		{
			name:       "0を含む",
			ws:         []float64{0.0, 0.3, 0.7},
			wantRatios: []float64{0.0, 0.3, 0.7},
		},
		{
			name:       "全ての重みが0",
			ws:         []float64{0.0, 0.0, 0.0, 0.0},
			wantRatios: []float64{0.25, 0.25, 0.25, 0.25},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := make([]int, sampleN)
			var err error
			for i := range sampleN {
				got[i], err = randx.IndexByWeights(tt.ws, rng)
				if err != nil {
					t.Fatalf("nilを期待したが、エラーが返された: %v", err)
				}
			}

			counts := slicesx.Counts(got)
			for i, wantRatio := range tt.wantRatios {
				gotRatio := float64(counts[i]) / float64(sampleN)
				if !mathx.ApproxEqual(gotRatio, wantRatio, eps) {
					t.Errorf("出現比率の不一致: got = %.3f want = %.3f(±%.3f) idx = %d", gotRatio, wantRatio, eps, i)
				}
			}
		})
	}
}

func TestIndexByWeights_Error(t *testing.T) {
	rng := randx.NewPCG()
	tests := []struct {
		name string
		ws   []float64
	}{
		{
			name: "空の重み",
			ws:   []float64{},
		},
		{
			name: "負の重み",
			ws:   []float64{0.2, -0.1, 0.9},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := randx.IndexByWeights(tt.ws, rng)
			if got != -1 {
				t.Errorf("エラー時の戻り値が-1ではない: got = %d", got)
			}
			if err == nil {
				t.Error("エラーを期待したが、nilが返された")
			}
		})
	}
}

func TestFloatRange_Statistics(t *testing.T) {
	rng := randx.NewPCG()
	const eps = 0.1

	tests := []struct {
		min float64
		max float64
	}{
		{min: -1.0, max: 2.0},
		{min: 0.0, max: 1.0},
		{min: -10.0, max: -5.0},
	}

	for _, tt := range tests {
		got := make([]float64, sampleN)
		var err error
		for i := range sampleN {
			got[i], err = randx.FloatRange(tt.min, tt.max, rng)
			if err != nil {
				t.Fatalf("nilを期待したが、エラーが返された: %v (min=%v, max=%v)", err, tt.min, tt.max)
			}
		}

		// 一様分布の期待平均値は (min + max) / 2
		wantAvg := (tt.min + tt.max) / 2.0

		gotSum := 0.0
		for _, v := range got {
			if v < tt.min || v >= tt.max {
				t.Fatalf("値が範囲外: got = %f want = [%f, %f)", v, tt.min, tt.max)
			}
			gotSum += v
		}
		gotAvg := gotSum / float64(sampleN)

		if !mathx.ApproxEqual(gotAvg, wantAvg, eps) {
			t.Errorf("平均の不一致: got = %.3f want = %.3f(±%.3f)", gotAvg, wantAvg, eps)
		}
	}
}

func TestFloatRange_Error(t *testing.T) {
	rng := randx.NewPCG()
	tests := []struct {
		name string
		min  float64
		max  float64
	}{
		{
			name: "minがmaxより大きい",
			min:  2.0,
			max:  1.0,
		},
		{
			name: "負同士でminがmaxより大きい",
			min:  -1.0,
			max:  -2.0,
		},
		{
			name: "minとmaxが等しい",
			min:  1.0,
			max:  1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := randx.FloatRange(tt.min, tt.max, rng)
			if got != 0.0 {
				t.Errorf("エラー時の戻り値がゼロ値ではない: got = %f", got)
			}
			if err == nil {
				t.Error("エラーを期待したが、nilが返された")
			}
		})
	}
}

func TestChoice_Statistics(t *testing.T) {
	rng := randx.NewPCG()
	const eps = 0.04

	tests := []struct {
		name string
		s    []string
	}{
		{
			name: "重複なし",
			s:    []string{"りんご", "ゴリラ", "ラッパ"},
		},
		{
			name: "重複あり",
			s:    []string{"魚", "魚", "肉"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := make([]string, sampleN)
			for i := range sampleN {
				v, err := randx.Choice(tt.s, rng)
				if err != nil {
					t.Fatalf("nilを期待したが、エラーが返された: %v", err)
				}
				got[i] = v
			}

			// 各要素の期待出現確率 = 要素の重複数 / len(s)
			sCounts := slicesx.Counts(tt.s)
			gotCounts := slicesx.Counts(got)
			n := len(tt.s)

			for k, c := range sCounts {
				wantRatio := float64(c) / float64(n)
				gotRatio := float64(gotCounts[k]) / float64(sampleN)
				if !mathx.ApproxEqual(gotRatio, wantRatio, eps) {
					t.Errorf("%q の出現比率の不一致: got = %.3f want = %.3f(±%.3f)", k, gotRatio, wantRatio, eps)
				}
			}
		})
	}
}

func TestChoice_Error(t *testing.T) {
	rng := randx.NewPCG()
	tests := []struct {
		name string
		s    []string
	}{
		{
			name: "空スライス",
			s:    []string{},
		},
		{
			name: "nil",
			s:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := randx.Choice(tt.s, rng)
			if got != "" {
				t.Errorf("エラー時の戻り値がゼロ値ではない: got = %s", got)
			}
			if err == nil {
				t.Error("エラーを期待したが、nilが返された")
			}
		})
	}
}

func TestBool(t *testing.T) {
	rng := randx.NewPCG()
	trueCount := 0
	for range sampleN {
		if randx.Bool(rng) {
			trueCount++
		}
	}

	gotRatio := float64(trueCount) / float64(sampleN)
	const want = 0.5
	const eps = 0.025
	if !mathx.ApproxEqual(gotRatio, want, eps) {
		t.Errorf("trueの比率の不一致: got = %.3f want = %.3f(±%.3f)", gotRatio, want, eps)
	}
}

func TestIntNorm_Statistics(t *testing.T) {
	rng := randx.NewPCG()

	t.Run("平均と標準偏差", func(t *testing.T) {
		const mean = 10.0
		const std = 5.0
		got := make([]int, sampleN)
		var err error
		for i := range sampleN {
			// 範囲を十分広く取り、切り捨ての影響をなくす
			got[i], err = randx.IntNorm(-1000, 1000, mean, std, rng)
			if err != nil {
				t.Fatalf("nilを期待したが、エラーが返された: %v", err)
			}
		}

		gotSum := 0.0
		for _, v := range got {
			gotSum += float64(v)
		}
		gotMean := gotSum / float64(sampleN)

		var gotVar float64
		for _, v := range got {
			d := float64(v) - gotMean
			gotVar += d * d
		}
		gotStd := math.Sqrt(gotVar / float64(sampleN))

		const meanEps = 0.3
		if !mathx.ApproxEqual(gotMean, mean, meanEps) {
			t.Errorf("平均の不一致: got = %.3f want = %.3f(±%.3f)", gotMean, mean, meanEps)
		}

		// 整数への丸目で分散はわずかに増える(連続性補正 1/12)ため、許容誤差は広めに取る
		const stdEps = 0.3
		if !mathx.ApproxEqual(gotStd, std, stdEps) {
			t.Errorf("標準偏差の不一致: got = %.3f want = %.3f(±%.3f)", gotStd, std, stdEps)
		}
	})

	t.Run("範囲内に収まる", func(t *testing.T) {
		const minVal, maxVal = 0, 3
		for range sampleN {
			// stdを大きくして、範囲外の再抽選が頻発する状況を作る
			got, err := randx.IntNorm(minVal, maxVal, 1.5, 10.0, rng)
			if err != nil {
				t.Fatalf("nilを期待したが、エラーが返された: %v", err)
			}
			if got < minVal || got >= maxVal {
				t.Fatalf("値が範囲外: got = %d want = [%d, %d)", got, minVal, maxVal)
			}
		}
	})

	t.Run("stdが0", func(t *testing.T) {
		got, err := randx.IntNorm(0, 100, 41.6, 0.0, rng)
		if err != nil {
			t.Fatalf("nilを期待したが、エラーが返された: %v", err)
		}
		// std = 0 の場合、meanを四捨五入した値を返す
		want := 42
		if got != want {
			t.Errorf("値の不一致: got = %d want = %d", got, want)
		}
	})
}

func TestIntNorm_Error(t *testing.T) {
	rng := randx.NewPCG()
	tests := []struct {
		name string
		min  int
		max  int
		mean float64
		std  float64
	}{
		{
			name: "minがmaxより大きい",
			min:  10, max: 0, mean: 5.0, std: 1.0,
		},
		{
			name: "境界_minとmaxが等しい",
			min:  10, max: 10, mean: 5.0, std: 1.0,
		},
		{
			name: "stdが負",
			min:  0, max: 10, mean: 5.0, std: -1.0,
		},
		{
			name: "meanが範囲外",
			min:  0, max: 10, mean: 11.0, std: 1.0,
		},
		{
			name: "stdが0で丸め結果が上限",
			min:  0, max: 1, mean: 0.6, std: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := randx.IntNorm(tt.min, tt.max, tt.mean, tt.std, rng)
			if got != 0 {
				t.Errorf("エラー時の戻り値がゼロ値ではない: got = %d", got)
			}
			if err == nil {
				t.Error("エラーを期待したが、nilが返された")
			}
		})
	}
}
