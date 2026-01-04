// // Package randx_test performs statistical validation for the randx package.
// // While random number tests often use fixed seeds for reproducibility,
// // this package prioritizes practical verification by conducting intensive
// // statistical tests to ensure distribution correctness.
// //
// // Package randx_test は、randx パッケージの統計的な検証を行います。
// // 乱数のテストでは再現性を持たせるためにシード値を固定することも多いですが、
// // このパッケージでは実用的な検証を重視し、分布の正確さを確認するための
// // 重たい統計テストを中心に行います。
package randx_test

// import (
// 	"fmt"
// 	"github.com/sw965/omw/constraints"
// 	"github.com/sw965/omw/mathx/randx"
// 	"github.com/sw965/omw/slicesx"
// 	"github.com/sw965/omw/mathx"
// 	"math"
// 	"math/rand/v2"
// 	"slices"
// 	"strings"
// 	"testing"
// )

// type invalidRangeCase[T constraints.Number] struct {
// 	name string
// 	min  T
// 	max  T
// }

// // TestIntRange_Error
// // TestFloatRange_Error
// // のヘルパー関数 (抽象化した関数)
// func runInvalidRangeErrorTests[T constraints.Number](
// 	t *testing.T,
// 	cases []invalidRangeCase[T],
// 	f func(min, max T, rng *rand.Rand) (T, error),
// ) {
// 	t.Helper()
// 	rng := randx.NewPCGFromGlobalSeed()

// 	for _, tc := range cases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			t.Helper()

// 			got, gotErr := f(tc.min, tc.max, rng)

// 			if gotErr == nil {
// 				t.Fatalf("エラーを期待したが、nilが返された")
// 			}

// 			// エラーが起きる前提である為、期待する戻り値はゼロ値
// 			var want T
// 			if got != want {
// 				t.Errorf("got: %v, want: %v", got, want)
// 			}

// 			wantErrMsgSubs := []string{
// 				fmt.Sprintf("範囲が不正"),
// 				fmt.Sprintf("min >= max"),
// 				fmt.Sprintf("min = %v", tc.min),
// 				fmt.Sprintf("max = %v", tc.max),
// 				fmt.Sprintf("min < max"),
// 			}

// 			gotErrMsg := gotErr.Error()
// 			for _, sub := range wantErrMsgSubs {
// 				if !strings.Contains(gotErrMsg, sub) {
// 					t.Errorf("gotErrMsg: %s, sub: %s", gotErrMsg, sub)
// 				}
// 			}
// 		})
// 	}
// }

// func collectSamples[T any](t *testing.T, n int, generator func() (T, error)) []T {
// 	t.Helper()
// 	samples := make([]T, n)
// 	for i := range n {
// 		v, err := generator()
// 		if err != nil {
// 			t.Fatalf("予期せぬエラー: %v", err)
// 		}
// 		samples[i] = v
// 	}
// 	return samples
// }

// func testMean[N constraints.Number](t *testing.T, samples []N, want N, eps float64) {
// 	t.Helper()
// 	n := len(samples)
// 	if n == 0 {
// 		t.Errorf("len(samples) = 0: len(samples) > 0 であるべき")
// 	}

// 	sum := mathx.Sum(samples...)
// 	gotMean := float64(sum) / float64(n)

// 	if mathx.Abs(gotMean-want) > eps {
// 		t.Errorf("gotMean: %v, want: %v(±%v)", gotMean, want)
// 	}
// }

// func testDistribution[T comparable](t *testing.T, samples []T, wantRatios map[T]float64, eps float64) {
// 	t.Helper()
// 	n := len(samples)
// 	counts := slicesx.Counts(samples)
// 	gotRatios := map[T]float64{}
// 	for k, c := range counts {
// 		gotRatios[k] = float64(c) / float64(n)
// 	}

// 	for val, wantRatio := range wantRatios {
// 		count := counts[val]
// 		gotRatio := float64(count) / total

// 		if math.Abs(gotRatio-wantRatio) > epsilon {
// 			t.Errorf("Distribution check failed for value '%v': want ratio %.3f(±%.3f), but got %.3f (count: %d/%d)",
// 				val, wantRatio, epsilon, gotRatio, count, int(total))
// 		}
// 	}
// }

// func TestIntRange(t *testing.T) {
// 	testNum := 10000
// 	rng := randx.NewPCGFromGlobalSeed()
// 	const epsilon = 0.1

// 	runCase := func(minVal, maxVal int) {
// 		name := fmt.Sprintf("統計テスト 最小値%d 最大値%d", minVal, maxVal)
// 		t.Run(name, func(t *testing.T) {
// 			t.Helper()

// 			got := make([]int, testNum)
// 			var err error
// 			for i := range testNum {
// 				got[i], err = randx.IntRange(minVal, maxVal, rng)
// 				if err != nil {
// 					t.Fatalf("予期せぬエラーが発生: %v", err)
// 				}
// 			}

// 			wantMin := minVal
// 			wantMax := maxVal - 1
// 			wantAvg := float64(wantMin+wantMax) / 2.0

// 			gotMin := slices.Min(got)
// 			gotMax := slices.Max(got)
// 			gotSum := 0
// 			for _, v := range got {
// 				gotSum += v
// 			}
// 			gotAvg := float64(gotSum) / float64(testNum)

// 			if gotMin != wantMin {
// 				t.Errorf("wantMin %d gotMin %d", wantMin, gotMin)
// 			}

// 			if gotMax != wantMax {
// 				t.Errorf("wantMax %d gotMax %d", wantMax, gotMax)
// 			}

// 			if math.Abs(gotAvg-wantAvg) > epsilon {
// 				t.Errorf("wantAvg %v(±%v) gotAvg %v", wantAvg, epsilon, gotAvg)
// 			}
// 		})
// 	}

// 	runCase(0, 10)
// 	runCase(-5, 5)
// 	runCase(-10, -5)
// }

// func TestIntRange_Error(t *testing.T) {
// 	cases := []invalidRangeCase[int]{
// 		{
// 			name: "異常 最小値100 最大値99",
// 			min:  100,
// 			max:  99,
// 		},
// 		{
// 			name: "異常 最小値100 最大値100",
// 			min:  100,
// 			max:  100,
// 		},
// 	}

// 	runInvalidRangeErrorTests(t, cases, func(min, max int, rng *rand.Rand) (int, error) {
// 		return randx.IntRange[int](min, max, rng)
// 	})
// }

// func TestIntByWeight(t *testing.T) {
// 	testNum := 10000
// 	rng := randx.NewPCGFromGlobalSeed()
// 	const epsilon = 0.03

// 	runCase := func(ws []float64) {
// 		name := fmt.Sprintf("統計テスト %v", ws)
// 		t.Run(name, func(t *testing.T) {
// 			t.Helper()

// 			got := make([]int, testNum)
// 			var err error
// 			for i := range testNum {
// 				got[i], err = randx.IntByWeight(ws, rng)
// 				if err != nil {
// 					t.Fatalf("予期せぬエラー: %v", err)
// 				}
// 			}

// 			counts := slicesx.Counts(got)
// 			wSum := 0.0
// 			for _, w := range ws {
// 				wSum += w
// 			}
// 			wantRatios := make([]float64, len(ws))
// 			for i := range wantRatios {
// 				wantRatios[i] = ws[i] / wSum
// 			}

// 			for i, wantRatio := range wantRatios {
// 				c := counts[i]
// 				gotRatio := float64(c) / float64(testNum)

// 				if math.Abs(gotRatio-wantRatio) > epsilon {
// 					t.Errorf("index %d wantRatio %.3f(±%.3f), gotRatio %.3f", i, wantRatio, epsilon, gotRatio)
// 				}
// 			}
// 		})
// 	}

// 	// 0が20%, 1が30%, 2が10% 3が40%
// 	runCase([]float64{0.2, 0.3, 0.1, 0.4})

// 	// 0が50%, 1が50% ※合計値に対する比率で出現確率が計算される
// 	runCase([]float64{0.1, 0.1})

// 	// 0が90%, 1が10%
// 	runCase([]float64{0.9, 0.1})

// 	// 0が0%, 1が30%, 2が70%
// 	runCase([]float64{0.0, 0.3, 0.7})
// }

// func TestIntByWeight_Error(t *testing.T) {
// 	rng := randx.NewPCGFromGlobalSeed()

// 	runCase := func(name string, ws []float64, wantErrMsgSubs []string) {
// 		t.Run(name, func(t *testing.T) {
// 			t.Helper()

// 			got, gotErr := randx.IntByWeight(ws, rng)

// 			if gotErr == nil {
// 				t.Fatalf("エラーを期待したが、nilが返された")
// 			}

// 			want := -1
// 			if got != want {
// 				t.Errorf("got: %d, want: %d", got, want)
// 			}

// 			gotErrMsg := gotErr.Error()
// 			for _, sub := range wantErrMsgSubs {
// 				if !strings.Contains(gotErrMsg, sub) {
// 					t.Errorf("gotErrMsg: %s, sub: %s", gotErrMsg, sub)
// 				}
// 			}
// 		})
// 	}

// 	tests := []struct {
// 		name           string
// 		ws             []float64
// 		wantErrMsgSubs []string
// 	}{
// 		{
// 			name: "異常 空の重み",
// 			ws:   []float64{},
// 			wantErrMsgSubs: []string{
// 				"len(ws) = 0",
// 				"len(ws) > 0",
// 			},
// 		},
// 		{
// 			name: "異常 負の重み",
// 			ws:   []float64{0.2, -0.1, 0.9},
// 			wantErrMsgSubs: []string{
// 				"ws[1] = -0.1",
// 				"負",
// 				"非負",
// 			},
// 		},
// 		{
// 			name: "異常 重みがNaN",
// 			ws:   []float64{0.2, 0.8, math.NaN()},
// 			wantErrMsgSubs: []string{
// 				"ws[2] = NaN",
// 				"NaN",
// 				"非NaN",
// 			},
// 		},
// 		{
// 			name: "異常 重みが+Inf",
// 			ws: []float64{1.0, 0.5, math.Inf(0), 0.25},
// 			wantErrMsgSubs: []string{
// 				"ws[2] = +Inf",
// 				"Inf",
// 				"非Inf",
// 			},
// 		},
// 		{
// 			name: "異常 重みが-Inf",
// 			ws: []float64{math.Inf(-1), 1.0, 0.5, 0.25},
// 			wantErrMsgSubs: []string{
// 				"ws[0] = -Inf",
// 				"-Inf",
// 				"非Inf",
// 			},
// 		},
// 	}

// 	for _, tc := range tests {
// 		runCase(tc.name, tc.ws, tc.wantErrMsgSubs)
// 	}
// }

// func TestFloatRange(t *testing.T) {
// 	testNum := 10000
// 	rng := randx.NewPCGFromGlobalSeed()
// 	const epsilon = 0.1

// 	runCase := func(minVal, maxVal float64) {
// 		name := fmt.Sprintf("統計テスト 最小値%v 最大値%v", minVal, maxVal)
// 		t.Run(name, func(t *testing.T) {
// 			t.Helper()

// 			// 一様分布の期待平均値は (min + max) / 2
// 			wantAvg := (minVal + maxVal) / 2.0

// 			got := make([]float64, testNum)
// 			var err error
// 			for i := range testNum {
// 				got[i], err = randx.FloatRange(minVal, maxVal, rng)
// 				if err != nil {
// 					t.Fatalf("予期せぬエラーが発生: %v", err)
// 				}
// 			}

// 			gotMin := slices.Min(got)
// 			gotMax := slices.Max(got)
// 			gotSum := 0.0
// 			for _, v := range got {
// 				gotSum += v
// 			}
// 			gotAvg := gotSum / float64(testNum)

// 			// 境界値テスト(最小)
// 			if gotMin < minVal {
// 				t.Errorf("wantMin >= %f gotMin %f", minVal, gotMin)
// 			}

// 			// 境界値テスト(最大)
// 			if gotMax >= maxVal {
// 				t.Errorf("wantMax < %f gotMax %f", maxVal, gotMax)
// 			}

// 			// 平均値のテスト
// 			if math.Abs(gotAvg-wantAvg) > epsilon {
// 				t.Errorf("wantAvg: %.3f (±%.3f), gotAvg: %.3f", wantAvg, epsilon, gotAvg)
// 			}
// 		})
// 	}

// 	runCase(-1.0, 2.0)
// 	runCase(0.0, 1.0)
// 	runCase(-10.0, -5.0)
// }

// func TestFloatRange_Error(t *testing.T) {
// 	cases := []invalidRangeCase[float64]{
// 		{
// 			name:       "異常 最小値2.0 最大値1.0",
// 			min:        2.0,
// 			max:        1.0,
// 		},
// 		{
// 			name:       "異常 最小値-1.0 最大値-2.0",
// 			min:        -1.0,
// 			max:        -2.0,
// 		},
// 		{
// 			name:       "異常 最小値1.0 最大値1.0",
// 			min:        1.0,
// 			max:        1.0,
// 		},
// 	}

// 	runInvalidRangeErrorTests(t, cases, func(min, max float64, rng *rand.Rand) (float64, error) {
// 		return randx.FloatRange[float64](min, max, rng)
// 	})
// }

// func TestChoice(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		s    []string
// 	}{
// 		{
// 			name: "正常 重複なし",
// 			s: []string{"りんご", "ゴリラ", "ラッパ"},
// 		},
// 		{
// 			name: "正常 重複あり",
// 			s: []string{"魚", "魚", "肉"},
// 		},
// 	}

// 	testNum := 10000
// 	rng := randx.NewPCGFromGlobalSeed()

// 	for _, tc := range tests {
// 		t.Run(tc.name, func(t *testing.T) {
// 			t.Helper()

// 			got := make([]string, testNum)
// 			for i := range testNum {
// 				v, err := randx.Choice(tc.s, rng)
// 				if err != nil {
// 					t.Fatalf("予期せぬエラー: %v", err)
// 				}
// 				got[i] = v
// 			}

// 			sCounts := slicesx.Counts(tc.s)
// 			gotCounts := slicesx.Counts(got)
// 			epsilon := 0.04
// 			sn := len(tc.s)

// 			for k, v := range sCounts {
// 				// tc.s が {"りんご", "みかん", "ぶどう"} のように要素の重複がない場合
// 				// それぞれが出現する期待値は 0.33... である
// 				// これは 1.0 / len(tc.s) で 計算出来る	・・・この式をAとする
// 				// もし ts.s が {"緑", "緑", "白"} のように要素の重複がある場合
// 				// 緑の出現確率は 0.66... 白の出現確率 0.33... である
// 				// これは緑も白も 要素の重複数 / len(tc.s) で 計算できる ・・・この式をBとする
// 				// 要素が重複していない時、要素の重複数 = 1 であるから
// 				// A と B より、要素の重複数 / len(tc.s) で 期待される出現確率を計算できる
// 				wantRatio := float64(v) / float64(sn)
// 				gotRatio := float64(gotCounts[k]) / float64(testNum)

// 				if mathx.Abs(gotRatio-wantRatio) > epsilon {
// 					t.Errorf("wantRatio=%.3f (±%.3f) gotRatio=%.3f", wantRatio, epsilon, gotRatio)
// 				}
// 			}
// 		})
// 	}
// }

// func TestChoice_Error(t *testing.T) {
// 	rng := randx.NewPCGFromGlobalSeed()

// 	got, err := randx.Choice([]string{}, rng)
// 	if got != "" {
// 		t.Errorf("空スライスを期待したが、%s が返された", got)
// 	}

// 	if err == nil {
// 		t.Fatalf("エラーを期待したが、nilが返された")
// 	}

// 	wantErrMsgSubs := []string{
// 		"len(s) = 0",
// 		"len(s) > 0",
// 	}

// 	errMsg := err.Error()
// 	for _, sub := range wantErrMsgSubs {
// 		if !strings.Contains(errMsg, sub) {
// 			t.Errorf("errMsg %s, sub %s", errMsg, sub)
// 		}
// 	}
// }

// func TestBool(t *testing.T) {
// 	testNum := 10000
// 	rng := randx.NewPCGFromGlobalSeed()

// 	got := make([]bool, testNum)
// 	for i := range testNum {
// 		got[i] = randx.Bool(rng)
// 	}

// 	counts := slicesx.Counts(got)

// 	trueRatio := float64(counts[true]) / float64(testNum)

// 	epsilon := 0.02
// 	want := 0.5

// 	if math.Abs(trueRatio-want) > epsilon {
// 		t.Errorf("trueRatio=%.3f, want=%.3f (±%.3f)", trueRatio, want, epsilon)
// 	}
// }
