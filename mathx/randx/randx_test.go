// Package randx_test は、randx パッケージの検証を行います。
// 乱数の生成そのものが主目的のパッケージである為、自動テスト方針に従い、
// シードを固定しない【統計テスト】を中心に行います。
// エラー系（引数の検証）は【期待値テスト】として行います。
package randx_test

import (
	"math"
	"strings"
	"testing"

	"github.com/sw965/omw/mathx/randx"
	"github.com/sw965/omw/slicesx"
)

const sampleN = 10000

// assertErrMsgSubs は、エラーメッセージに調査に必要な情報が含まれている事を確認する。
func assertErrMsgSubs(t *testing.T, err error, subs []string) {
	t.Helper()
	if err == nil {
		t.Fatal("エラーを期待したが、nilが返された")
	}
	msg := err.Error()
	for _, sub := range subs {
		if !strings.Contains(msg, sub) {
			t.Errorf("エラーメッセージに %q が含まれていない: msg = %s", sub, msg)
		}
	}
}

func TestNewPCGs(t *testing.T) {
	t.Run("正常_個数", func(t *testing.T) {
		rngs := randx.NewPCGs(3)
		if len(rngs) != 3 {
			t.Fatalf("len(rngs)の不一致: got = %d, want = %d", len(rngs), 3)
		}
		for i, rng := range rngs {
			if rng == nil {
				t.Errorf("rngs[%d]がnil", i)
			}
		}
	})

	t.Run("準正常_0個", func(t *testing.T) {
		rngs := randx.NewPCGs(0)
		if len(rngs) != 0 {
			t.Fatalf("len(rngs)の不一致: got = %d, want = %d", len(rngs), 0)
		}
	})

	t.Run("統計_異なる乱数列", func(t *testing.T) {
		// 2つの乱数器が同じ列を生成する確率は無視できるほど小さい
		rngs := randx.NewPCGs(2)
		same := true
		for i := 0; i < 10; i++ {
			if rngs[0].Uint64() != rngs[1].Uint64() {
				same = false
				break
			}
		}
		if same {
			t.Error("2つの乱数器が同一の乱数列を生成した")
		}
	})
}

func TestIntRange(t *testing.T) {
	rng := randx.NewPCG()
	const eps = 0.15

	runCase := func(minVal, maxVal int) {
		t.Run("統計_範囲と平均", func(t *testing.T) {
			t.Helper()
			got := make([]int, sampleN)
			var err error
			for i := range sampleN {
				got[i], err = randx.IntRange(minVal, maxVal, rng)
				if err != nil {
					t.Fatalf("予期せぬエラー: %v", err)
				}
			}

			// 生成値は [min, max) に収まり、min と max-1 は十分な試行で観測されるはず
			wantMin := minVal
			wantMax := maxVal - 1
			wantAvg := float64(wantMin+wantMax) / 2.0

			gotMin, gotMax := got[0], got[0]
			gotSum := 0
			for _, v := range got {
				if v < gotMin {
					gotMin = v
				}
				if v > gotMax {
					gotMax = v
				}
				gotSum += v
			}
			gotAvg := float64(gotSum) / float64(sampleN)

			if gotMin != wantMin {
				t.Errorf("最小値の不一致: got = %d, want = %d", gotMin, wantMin)
			}
			if gotMax != wantMax {
				t.Errorf("最大値の不一致: got = %d, want = %d", gotMax, wantMax)
			}
			if math.Abs(gotAvg-wantAvg) > eps {
				t.Errorf("平均の不一致: got = %v, want = %v(±%v)", gotAvg, wantAvg, eps)
			}
		})
	}

	runCase(0, 10)
	runCase(-5, 5)
	runCase(-10, -5)
}

func TestIntRange_Error(t *testing.T) {
	rng := randx.NewPCG()
	tests := []struct {
		name string
		min  int
		max  int
	}{
		{name: "異常_minがmaxより大きい", min: 100, max: 99},
		{name: "異常_境界_minとmaxが等しい", min: 100, max: 100},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := randx.IntRange(tc.min, tc.max, rng)
			if got != 0 {
				t.Errorf("エラー時の戻り値がゼロ値ではない: got = %d", got)
			}
			assertErrMsgSubs(t, err, []string{"範囲が不正", "min = ", "max = ", "min < max"})
		})
	}
}

func TestIntByWeights(t *testing.T) {
	rng := randx.NewPCG()
	const eps = 0.03

	runCase := func(name string, ws []float64, wantRatios []float64) {
		t.Run(name, func(t *testing.T) {
			t.Helper()
			got := make([]int, sampleN)
			var err error
			for i := range sampleN {
				got[i], err = randx.IntByWeights(ws, rng)
				if err != nil {
					t.Fatalf("予期せぬエラー: %v", err)
				}
			}

			counts := slicesx.Counts(got)
			for i, wantRatio := range wantRatios {
				gotRatio := float64(counts[i]) / float64(sampleN)
				if math.Abs(gotRatio-wantRatio) > eps {
					t.Errorf("インデックス%dの出現比率の不一致: got = %.3f, want = %.3f(±%.3f)", i, gotRatio, wantRatio, eps)
				}
			}
		})
	}

	// 出現確率は、重みの合計値に対する比率で決まる
	runCase("統計_通常", []float64{0.2, 0.3, 0.1, 0.4}, []float64{0.2, 0.3, 0.1, 0.4})
	runCase("統計_合計が1ではない", []float64{0.1, 0.1}, []float64{0.5, 0.5})
	runCase("統計_偏り", []float64{0.9, 0.1}, []float64{0.9, 0.1})
	runCase("統計_重み0を含む", []float64{0.0, 0.3, 0.7}, []float64{0.0, 0.3, 0.7})
	// 重みが全て0の場合は一様ランダム
	runCase("統計_全ての重みが0", []float64{0.0, 0.0, 0.0, 0.0}, []float64{0.25, 0.25, 0.25, 0.25})
}

func TestIntByWeights_Error(t *testing.T) {
	rng := randx.NewPCG()
	tests := []struct {
		name           string
		ws             []float64
		wantErrMsgSubs []string
	}{
		{
			name:           "異常_空の重み",
			ws:             []float64{},
			wantErrMsgSubs: []string{"len(ws) = 0", "len(ws) > 0"},
		},
		{
			name:           "異常_負の重み",
			ws:             []float64{0.2, -0.1, 0.9},
			wantErrMsgSubs: []string{"ws[1] = -0.1", "非負"},
		},
		{
			name:           "異常_重みがNaN",
			ws:             []float64{0.2, 0.8, math.NaN()},
			wantErrMsgSubs: []string{"ws[2] = NaN", "非NaN"},
		},
		{
			name:           "異常_重みが正の無限大",
			ws:             []float64{1.0, 0.5, math.Inf(1), 0.25},
			wantErrMsgSubs: []string{"ws[2] = +Inf", "非Inf"},
		},
		{
			name:           "異常_重みが負の無限大",
			ws:             []float64{math.Inf(-1), 1.0, 0.5, 0.25},
			wantErrMsgSubs: []string{"ws[0] = -Inf", "非Inf"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := randx.IntByWeights(tc.ws, rng)
			if got != -1 {
				t.Errorf("エラー時の戻り値が-1ではない: got = %d", got)
			}
			assertErrMsgSubs(t, err, tc.wantErrMsgSubs)
		})
	}
}

func TestFloatRange(t *testing.T) {
	rng := randx.NewPCG()
	const eps = 0.1

	runCase := func(minVal, maxVal float64) {
		t.Run("統計_範囲と平均", func(t *testing.T) {
			t.Helper()
			got := make([]float64, sampleN)
			var err error
			for i := range sampleN {
				got[i], err = randx.FloatRange(minVal, maxVal, rng)
				if err != nil {
					t.Fatalf("予期せぬエラー: %v", err)
				}
			}

			// 一様分布の期待平均値は (min + max) / 2
			wantAvg := (minVal + maxVal) / 2.0

			gotSum := 0.0
			for _, v := range got {
				// 境界: [min, max) に収まる事
				if v < minVal || v >= maxVal {
					t.Fatalf("値が範囲外: got = %f, want = [%f, %f)", v, minVal, maxVal)
				}
				gotSum += v
			}
			gotAvg := gotSum / float64(sampleN)

			if math.Abs(gotAvg-wantAvg) > eps {
				t.Errorf("平均の不一致: got = %.3f, want = %.3f(±%.3f)", gotAvg, wantAvg, eps)
			}
		})
	}

	runCase(-1.0, 2.0)
	runCase(0.0, 1.0)
	runCase(-10.0, -5.0)
}

func TestFloatRange_Error(t *testing.T) {
	rng := randx.NewPCG()

	t.Run("異常_範囲が不正", func(t *testing.T) {
		tests := []struct {
			name string
			min  float64
			max  float64
		}{
			{name: "異常_minがmaxより大きい", min: 2.0, max: 1.0},
			{name: "異常_負同士でminがmaxより大きい", min: -1.0, max: -2.0},
			{name: "異常_境界_minとmaxが等しい", min: 1.0, max: 1.0},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				got, err := randx.FloatRange(tc.min, tc.max, rng)
				if got != 0.0 {
					t.Errorf("エラー時の戻り値がゼロ値ではない: got = %f", got)
				}
				assertErrMsgSubs(t, err, []string{"範囲が不正", "min < max"})
			})
		}
	})

	t.Run("異常_NaNとInf", func(t *testing.T) {
		tests := []struct {
			name           string
			min            float64
			max            float64
			wantErrMsgSubs []string
		}{
			{name: "異常_minがNaN", min: math.NaN(), max: 1.0, wantErrMsgSubs: []string{"minが不正"}},
			{name: "異常_maxがNaN", min: 0.0, max: math.NaN(), wantErrMsgSubs: []string{"maxが不正"}},
			{name: "異常_minが負の無限大", min: math.Inf(-1), max: 1.0, wantErrMsgSubs: []string{"minが不正"}},
			{name: "異常_maxが正の無限大", min: 0.0, max: math.Inf(1), wantErrMsgSubs: []string{"maxが不正"}},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				_, err := randx.FloatRange(tc.min, tc.max, rng)
				assertErrMsgSubs(t, err, tc.wantErrMsgSubs)
			})
		}
	})
}

func TestChoice(t *testing.T) {
	rng := randx.NewPCG()
	const eps = 0.04

	tests := []struct {
		name string
		s    []string
	}{
		{name: "統計_重複なし", s: []string{"りんご", "ゴリラ", "ラッパ"}},
		{name: "統計_重複あり", s: []string{"魚", "魚", "肉"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := make([]string, sampleN)
			for i := range sampleN {
				v, err := randx.Choice(tc.s, rng)
				if err != nil {
					t.Fatalf("予期せぬエラー: %v", err)
				}
				got[i] = v
			}

			// 各要素の期待出現確率 = 要素の重複数 / len(s)
			sCounts := slicesx.Counts(tc.s)
			gotCounts := slicesx.Counts(got)
			n := len(tc.s)

			for k, c := range sCounts {
				wantRatio := float64(c) / float64(n)
				gotRatio := float64(gotCounts[k]) / float64(sampleN)
				if math.Abs(gotRatio-wantRatio) > eps {
					t.Errorf("%q の出現比率の不一致: got = %.3f, want = %.3f(±%.3f)", k, gotRatio, wantRatio, eps)
				}
			}
		})
	}
}

func TestChoice_Error(t *testing.T) {
	rng := randx.NewPCG()

	t.Run("異常_空スライス", func(t *testing.T) {
		got, err := randx.Choice([]string{}, rng)
		if got != "" {
			t.Errorf("エラー時の戻り値がゼロ値ではない: got = %s", got)
		}
		assertErrMsgSubs(t, err, []string{"len(s) = 0", "len(s) > 0"})
	})

	t.Run("異常_nil", func(t *testing.T) {
		_, err := randx.Choice[[]string](nil, rng)
		assertErrMsgSubs(t, err, []string{"len(s) = 0", "len(s) > 0"})
	})
}

func TestBool(t *testing.T) {
	t.Run("統計_比率", func(t *testing.T) {
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
		if math.Abs(gotRatio-want) > eps {
			t.Errorf("trueの比率の不一致: got = %.3f, want = %.3f(±%.3f)", gotRatio, want, eps)
		}
	})
}

func TestNormalInt(t *testing.T) {
	rng := randx.NewPCG()

	t.Run("統計_平均と標準偏差", func(t *testing.T) {
		const mean = 10.0
		const std = 5.0
		got := make([]int, sampleN)
		var err error
		for i := range sampleN {
			// 範囲を十分広く取り、切り捨ての影響をなくす
			got[i], err = randx.NormalInt(-1000, 1000, mean, std, rng)
			if err != nil {
				t.Fatalf("予期せぬエラー: %v", err)
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
		if math.Abs(gotMean-mean) > meanEps {
			t.Errorf("平均の不一致: got = %.3f, want = %.3f(±%.3f)", gotMean, mean, meanEps)
		}

		// 整数への丸めで分散はわずかに増える(連続性補正 1/12)ため、許容誤差は広めに取る
		const stdEps = 0.3
		if math.Abs(gotStd-std) > stdEps {
			t.Errorf("標準偏差の不一致: got = %.3f, want = %.3f(±%.3f)", gotStd, std, stdEps)
		}
	})

	t.Run("性質_範囲内に収まる", func(t *testing.T) {
		const minVal, maxVal = 0, 3
		for range sampleN {
			// stdを大きくして、範囲外の再抽選が頻発する状況を作る
			got, err := randx.NormalInt(minVal, maxVal, 1.5, 10.0, rng)
			if err != nil {
				t.Fatalf("予期せぬエラー: %v", err)
			}
			if got < minVal || got > maxVal {
				t.Fatalf("値が範囲外: got = %d, want = [%d, %d]", got, minVal, maxVal)
			}
		}
	})

	t.Run("正常_stdが0", func(t *testing.T) {
		got, err := randx.NormalInt(0, 100, 41.6, 0.0, rng)
		if err != nil {
			t.Fatalf("予期せぬエラー: %v", err)
		}
		// std = 0 の場合、meanを四捨五入した値を返す
		want := 42
		if got != want {
			t.Errorf("値の不一致: got = %d, want = %d", got, want)
		}
	})
}

func TestNormalInt_Error(t *testing.T) {
	rng := randx.NewPCG()
	tests := []struct {
		name           string
		min            int
		max            int
		mean           float64
		std            float64
		wantErrMsgSubs []string
	}{
		{
			name: "異常_minがmaxより大きい",
			min:  10, max: 0, mean: 5.0, std: 1.0,
			wantErrMsgSubs: []string{"範囲が不正", "min <= max"},
		},
		{
			name: "異常_stdが負",
			min:  0, max: 10, mean: 5.0, std: -1.0,
			wantErrMsgSubs: []string{"std < 0", "std >= 0"},
		},
		{
			name: "異常_meanが範囲外",
			min:  0, max: 10, mean: 11.0, std: 1.0,
			wantErrMsgSubs: []string{"meanが範囲外", "mean = 11"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := randx.NormalInt(tc.min, tc.max, tc.mean, tc.std, rng)
			if got != 0 {
				t.Errorf("エラー時の戻り値がゼロ値ではない: got = %d", got)
			}
			assertErrMsgSubs(t, err, tc.wantErrMsgSubs)
		})
	}
}
