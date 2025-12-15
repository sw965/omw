// Package randx provides utility functions for random number generation
// using Go 1.23's math/rand/v2 package.
//
// Package randx は、Go 1.23 の math/rand/v2 パッケージを使用した
// 乱数生成のユーティリティ関数を提供します。
package randx

import (
	"errors"
	"fmt"
	"github.com/sw965/omw/constraints"
	"math"
	"math/rand/v2"
)

// InvalidRangeError records an error when the minimum value is greater than the maximum value.
//
// InvalidRangeError は、最小値が最大値よりも大きい場合にエラーを記録します。
type InvalidRangeError[T constraints.Integer | constraints.Float] struct {
	Min T
	Max T
}

func (e *InvalidRangeError[T]) Error() string {
	return fmt.Sprintf("無効な範囲です: 最小値(%v) >= 最大値(%v) が 満たさないようにしなければなりません", e.Min, e.Max)
}

var (
	ErrEmptySlice = errors.New("空スライスエラー")
	ErrNaN        = errors.New("Nanエラー")
	ErrNegative   = errors.New("負の値エラー")
)

// NewPCGFromGlobalSeed creates a new PCG random number generator seeded from the global random source.
// This ensures a wide state space (128-bit) initialization without relying on time.Now().
//
// NewPCGFromGlobalSeed は、グローバル乱数ソースからシードを生成し、新しいPCG乱数生成器を作成します。
// これにより、現在時刻(time.Now)に依存することなく、128ビットの状態空間をフルに活用した初期化を保証します。
func NewPCGFromGlobalSeed() *rand.Rand {
	// グローバル乱数（ChaCha8）から64bit整数を2つ取得してシードにする
	return rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))
}

// IntRange returns a random integer in the half-open interval [min, max).
// It supports any integer type and safely handles potential overflows by using uint64 internally.
// If min >= max, it returns min.
//
// IntRange は、半開区間 [min, max) の範囲のランダムな整数を返します。
// 任意の整数型をサポートし、内部で uint64 を使用することでオーバーフローを安全に処理します。
func IntRange[I constraints.Integer](minVal, maxVal I, rng *rand.Rand) (I, error) {
	if minVal >= maxVal {
		return 0, &InvalidRangeError[I]{Min: minVal, Max: maxVal}
	}

	// 差分をuint64で計算することで、int64の全範囲や負の数も安全に扱える
	diff := uint64(maxVal) - uint64(minVal)
	return I(rng.Uint64N(diff)) + minVal, nil
}

// IntByWeight returns a random index selected based on the weights provided.
// If the sum of weights is 0, it selects an index uniformly at random.
// It returns an error if weights are empty, negative, or NaN.
//
// IntByWeight は、提供された重みリストに基づいてランダムにインデックスを選択します。
// 重みの合計が0の場合は、完全ランダム（一様分布）で選択します。
// 重みが空、負の値、またはNaNが含まれる場合はエラーを返します。
func IntByWeight[F constraints.Float](ws []F, rng *rand.Rand) (int, error) {
	if len(ws) == 0 {
		return -1, fmt.Errorf("重みが提供されていません: %w", ErrEmptySlice)
	}

	sum := F(0.0)
	for i, w := range ws {
		if w < 0.0 {
			return -1, fmt.Errorf("インデックス %d の重み (%v) が不正です: %w", i, w, ErrNegative)
		}
		if math.IsNaN(float64(w)) {
			return -1, fmt.Errorf("インデックス %d の重みが不正です: %w", i, ErrNaN)
		}
		sum += w
	}

	// 重みの合計が0なら一様ランダム
	if sum == 0.0 {
		return rng.IntN(len(ws)), nil
	}

	threshold, err := FloatRange(0.0, sum, rng)
	if err != nil {
		return 0, err
	}

	var current F = 0.0
	for i, w := range ws {
		current += w
		if current >= threshold {
			return i, nil
		}
	}

	// ここまで到達したら最後の要素を返す
	return len(ws) - 1, nil
}

// FloatRange returns a random floating-point number in the half-open interval [min, max).
//
// FloatRange は、半開区間 [min, max) の範囲のランダムな浮動小数点数を返します。
func FloatRange[F constraints.Float](minVal, maxVal F, rng *rand.Rand) (F, error) {
	if minVal >= maxVal {
		return 0.0, &InvalidRangeError[F]{Min: minVal, Max: maxVal}
	}
	return F(rng.Float64())*(maxVal-minVal) + minVal, nil
}
