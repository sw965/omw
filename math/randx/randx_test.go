// Package randx_test performs statistical validation for the randx package.
// While random number tests often use fixed seeds for reproducibility,
// this package prioritizes practical verification by conducting intensive
// statistical tests to ensure distribution correctness.
//
// Package randx_test は、randx パッケージの統計的な検証を行います。
// 乱数のテストでは再現性を持たせるためにシード値を固定することも多いですが、
// このパッケージでは実用的な検証を重視し、分布の正確さを確認するための
// 重たい統計テストを中心に行います。
package randx_test

import (
	"errors"
	"fmt"
	"github.com/sw965/omw/constraints"
	"github.com/sw965/omw/math/randx"
	"github.com/sw965/omw/slicesx"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

type number interface {
	constraints.Integer | constraints.Float
}

type invalidRangeCase[T number] struct {
	name       string
	min        T
	max        T
	wantErrMin T
	wantErrMax T
}

// TestIntRange_ErrorとTestFloatRange_Errorを抽象化した関数
func runInvalidRangeErrorTests[T number](
	t *testing.T,
	cases []invalidRangeCase[T],
	f func(min, max T, rng *rand.Rand) (T, error),
) {

	t.Helper()
	rng := randx.NewPCGFromGlobalSeed()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Helper()

			got, gotErr := f(tc.min, tc.max, rng)

			var zero T
			if got != zero {
				t.Errorf("期待される戻り値は%vだが、%vが返された", zero, got)
			}

			if gotErr == nil {
				t.Fatalf("エラーを返す事を期待されるが、nilが返された")
			}

			targetErr := &randx.InvalidRangeError[T]{}
			if !errors.As(gotErr, &targetErr) {
				t.Errorf("期待されるエラー型と異なります: want :%T, got: %T", targetErr, gotErr)
				return
			}

			if targetErr.Min != tc.wantErrMin {
				t.Errorf("gotErrMin :%v, wantErrMin: %v", targetErr.Min, tc.wantErrMin)
			}
			if targetErr.Max != tc.wantErrMax {
				t.Errorf("gotErrMax :%v, wantErrMax: %v", targetErr.Max, tc.wantErrMax)
			}
		})
	}
}

func TestIntRange(t *testing.T) {
	testNum := 10000
	rng := randx.NewPCGFromGlobalSeed()

	runCase := func(minVal, maxVal int) {
		name := fmt.Sprintf("統計テスト_最小値%d_最大値%d", minVal, maxVal)
		t.Run(name, func(t *testing.T) {
			t.Helper()

			got := make([]int, testNum)
			var err error
			for i := range testNum {
				got[i], err = randx.IntRange(minVal, maxVal, rng)
				if err != nil {
					t.Fatalf("予期せぬエラーが発生しました: %v", err)
				}
			}

			wantMin := minVal
			wantMax := maxVal - 1
			wantAvg := float64(wantMin+wantMax) / 2.0

			gotMin := slices.Min(got)
			gotMax := slices.Max(got)
			gotSum := 0
			for _, v := range got {
				gotSum += v
			}
			gotAvg := float64(gotSum) / float64(testNum)

			if gotMin != wantMin {
				t.Errorf("wantMin: %d, gotMin: %d", wantMin, gotMin)
			}

			if gotMax != wantMax {
				t.Errorf("wantMax: %d, gotMax: %d", wantMax, gotMax)
			}

			epsilon := 0.1
			if math.Abs(gotAvg-wantAvg) > epsilon {
				t.Errorf("wantAvg: %v (±%v), gotAvg: %v", wantAvg, epsilon, gotAvg)
			}
		})
	}

	runCase(0, 10)
	runCase(-5, 5)
	runCase(-10, -5)
}

func TestIntRange_Error(t *testing.T) {
	cases := []invalidRangeCase[int]{
		{
			name:       "異常_最小値100_最大値99",
			min:        100,
			max:        99,
			wantErrMin: 100,
			wantErrMax: 99,
		},
		{
			name:       "異常_最小値100_最大値100",
			min:        100,
			max:        100,
			wantErrMin: 100,
			wantErrMax: 100,
		},
	}

	runInvalidRangeErrorTests(t, cases, func(min, max int, rng *rand.Rand) (int, error) {
		return randx.IntRange[int](min, max, rng)
	})
}

func TestIntByWeight(t *testing.T) {
	testNum := 10000
	rng := randx.NewPCGFromGlobalSeed()

	runCase := func(ws []float64) {
		name := fmt.Sprintf("統計テスト_%v", ws)
		t.Run(name, func(t *testing.T) {
			t.Helper()

			got := make([]int, testNum)
			var err error
			for i := range testNum {
				got[i], err = randx.IntByWeight(ws, rng)
				if err != nil {
					t.Fatalf("予期せぬエラーが発生しました: %v", err)
				}
			}

			// 各整数値毎をキーとして、出現した回数を数える
			counts := slicesx.Counts(got)
			epsilon := 0.03
			wSum := 0.0
			for _, w := range ws {
				wSum += w
			}
			wantRatios := make([]float64, len(ws))
			for i := range wantRatios {
				wantRatios[i] = ws[i] / wSum 
			}

			for i, wantRatio := range wantRatios {
				c := counts[i]
				gotRatio := float64(c) / float64(testNum)

				if math.Abs(gotRatio-wantRatio) > epsilon {
					t.Errorf("index: %d, wantRatio: %.3f (±%.3f), gotRatio: %.3f", i, wantRatio, epsilon, gotRatio)
				}
			}
		})
	}

	// 0が20%, 1が30%, 2が10% 3が40%
	runCase([]float64{0.2, 0.3, 0.1, 0.4})

	// 0が50%, 1が50% ※合計値に対する比率で出現確率が計算される
	runCase([]float64{0.1, 0.1})

	// 0が90%, 1が10%
	runCase([]float64{0.9, 0.1})

	// 0が0%, 1が30%, 2が70%
	runCase([]float64{0.0, 0.3, 0.7})
}

func TestIntByWeight_Error(t *testing.T) {
	rng := randx.NewPCGFromGlobalSeed()

	runCase := func(name string, ws []float64, checkErr func(error) bool) {
		t.Run(name, func(t *testing.T) {
			t.Helper()

			if checkErr == nil {
				t.Fatalf("checkErr は必須です（テストの書き忘れ防止）")
			}

			got, gotErr := randx.IntByWeight(ws, rng)

			if got != -1 {
				t.Errorf("want: %d, got: %d", -1, got)
			}

			if gotErr == nil {
				t.Fatalf("エラーを期待したが、nilが返された")
			}

			if !checkErr(gotErr) {
				t.Errorf("期待するエラー条件を満たしません: %v", gotErr)
			}
		})
	}

	tests := []struct {
		name     string
		ws       []float64
		checkErr func(error) bool
	}{
		{
			name: "異常_空の重み",
			ws:   []float64{},
			checkErr: func(err error) bool {
				return errors.Is(err, randx.ErrEmptySlice)
			},
		},
		{
			name: "異常_負の重み",
			ws:   []float64{0.2, -0.1, 0.9},
			checkErr: func(err error) bool {
				return errors.Is(err, randx.ErrNegative)
			},
		},
		{
			name: "異常_重みがNaN",
			ws:   []float64{0.2, float64(math.NaN()), 0.8},
			checkErr: func(err error) bool {
				return errors.Is(err, randx.ErrNaN)
			},
		},
	}

	for _, tc := range tests {
		runCase(tc.name, tc.ws, tc.checkErr)
	}
}

func TestFloatRange(t *testing.T) {
	testNum := 10000
	rng := randx.NewPCGFromGlobalSeed()

	runCase := func(minVal, maxVal float64) {
		name := fmt.Sprintf("統計テスト_最小値%v_最大値%v", minVal, maxVal)
		t.Run(name, func(t *testing.T) {
			t.Helper()

			// 一様分布の期待平均値は (min + max) / 2
			wantAvg := (minVal + maxVal) / 2.0

			got := make([]float64, testNum)
			var err error
			for i := range testNum {
				got[i], err = randx.FloatRange(minVal, maxVal, rng)
				if err != nil {
					t.Fatalf("予期せぬエラーが発生しました: %v", err)
				}
			}

			gotMin := slices.Min(got)
			gotMax := slices.Max(got)
			gotSum := 0.0
			for _, v := range got {
				gotSum += v
			}
			gotAvg := gotSum / float64(testNum)

			// 境界値テスト(最小)
			if gotMin < minVal {
				t.Errorf("wantMin >= %f, gotMin: %f", minVal, gotMin)
			}

			// 境界値テスト(最大)
			if gotMax >= maxVal {
				t.Errorf("wantMax < %f, gotMax: %f", maxVal, gotMax)
			}

			// 平均値のテスト
			epsilon := 0.1
			if math.Abs(gotAvg-wantAvg) > epsilon {
				t.Errorf("wantAvg: %f (±%f), gotAvg: %f", wantAvg, epsilon, gotAvg)
			}
		})
	}

	runCase(-1.0, 2.0)
	runCase(0.0, 1.0)
	runCase(-10.0, -5.0)
}

func TestFloatRange_Error(t *testing.T) {
	cases := []invalidRangeCase[float64]{
		{
			name:       "異常_最小値2.0_最大値1.0",
			min:        2.0,
			max:        1.0,
			wantErrMin: 2.0,
			wantErrMax: 1.0,
		},
		{
			name:       "異常_最小値-1.0_最大値-2.0",
			min:        -1.0,
			max:        -2.0,
			wantErrMin: -1.0,
			wantErrMax: -2.0,
		},
		{
			name:       "異常_最小値1.0_最大値1.0",
			min:        1.0,
			max:        1.0,
			wantErrMin: 1.0,
			wantErrMax: 1.0,
		},
	}

	runInvalidRangeErrorTests(t, cases, func(min, max float64, rng *rand.Rand) (float64, error) {
		return randx.FloatRange[float64](min, max, rng)
	})
}

func TestChoice(t *testing.T) {
	tests := []struct {
		name string
		s    []string
	}{
		{name: "正常_重複なし", s: []string{"りんご", "ゴリラ", "ラッパ"}},
		{name: "正常_重複あり", s: []string{"魚", "魚", "肉"}},
	}

	testNum := 10000
	rng := randx.NewPCGFromGlobalSeed()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Helper()
	
			got := make([]string, testNum)
			for i := range testNum {
				v, err := randx.Choice(tc.s, rng)
				if err != nil {
					t.Fatalf("予期せぬエラー: %v", err)
				}
				got[i] = v
			}

			sCounts := slicesx.Counts(tc.s)
			gotCounts := slicesx.Counts(got)
			epsilon := 0.04
			sn := len(tc.s)

			for k, v := range sCounts {
				// tc.s が {"りんご", "みかん", "ぶどう"} のように要素の重複がない場合
				// それぞれが出現する期待値は 0.33... である
				// これは 1.0 / len(tc.s) で 計算出来る	・・・この式をAとする
				// もし ts.s が {"緑", "緑", "白"} のように要素の重複がある場合
				// 緑の出現確率は 0.66... 白の出現確率 0.33... である
				// これは緑も白も 要素の重複数 / len(tc.s) で 計算できる ・・・この式をBとする
				// 要素が重複していない時、要素の重複数 = 1 であるから
				// A と B より、要素の重複数 / len(tc.s) で 期待される出現確率を計算できる
				wantRatio := float64(v) / float64(sn)
				gotRatio := float64(gotCounts[k]) / float64(testNum)

				if math.Abs(gotRatio-wantRatio) > epsilon {
					t.Errorf("wantRatio=%.3f (±%.3f) gotRatio=%.3f", wantRatio, epsilon, gotRatio)
				}
			}
		})
	}
}

func TestChoice_Error(t *testing.T) {
	rng := randx.NewPCGFromGlobalSeed()

	got, err := randx.Choice([]string{}, rng)
	if got != "" {
		t.Errorf("空文字を期待したが、%s が返された", got)
	}

	if err == nil {
		t.Fatalf("エラーを期待したが、nilが返された")
	}

	if !errors.Is(err, randx.ErrEmptySlice) {
		t.Errorf("期待されるエラー型が埋め込まれていなかった。want: %T, got: %v", randx.ErrEmptySlice, err)
	}
}

func TestBool(t *testing.T) {
	testNum := 10000
	rng := randx.NewPCGFromGlobalSeed()

	got := make([]bool, testNum)
	for i := range testNum {
		got[i] = randx.Bool(rng)
	}

	counts := slicesx.Counts(got)

	trueRatio := float64(counts[true]) / float64(testNum)
	falseRatio := float64(counts[false]) / float64(testNum)

	epsilon := 0.02
	want := 0.5

	if math.Abs(trueRatio-want) > epsilon {
		t.Errorf("trueRatio=%.3f, want=%.3f (±%.3f)", trueRatio, want, epsilon)
	}

	if math.Abs(falseRatio-want) > epsilon {
		t.Errorf("falseRatio=%.3f, want=%.3f (±%.3f)", falseRatio, want, epsilon)
	}
}