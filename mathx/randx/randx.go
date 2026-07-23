package randx

import (
	"fmt"
	"math"
	"math/rand/v2"

	"github.com/sw965/omw/constraints"
)

func NewPCG() *rand.Rand {
	return rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))
}

func NewPCGs(n int) []*rand.Rand {
	rngs := make([]*rand.Rand, n)
	for i := range n {
		rngs[i] = NewPCG()
	}
	return rngs
}

func IntRange[I constraints.Integer](minVal, maxVal I, rng *rand.Rand) (I, error) {
	if minVal >= maxVal {
		var zero I
		return zero, fmt.Errorf("範囲が不正(min >= max): min = %d, max = %d: min < max であるべき", minVal, maxVal)
	}

	// uint64にキャストする理由:
	// minが負の符号付き整数でも、uint64への変換は2の補数表現のビット列をそのまま保つため、
	// maxVal - minVal の差分はuint64の剰余演算として正しく求まる（差は必ず 0 以上 かつ uint64に収まる）。
	// これにより、符号付き・符号なしのどちらの型でも同じ式で Uint64N を使える。
	// 戻り値の I(...) + minVal も同じ剰余演算の性質で正しい値に戻る。
	diff := uint64(maxVal) - uint64(minVal)
	return I(rng.Uint64N(diff)) + minVal, nil
}

func IntByWeights[F constraints.Float](ws []F, rng *rand.Rand) (int, error) {
	n := len(ws)
	if n == 0 {
		return -1, fmt.Errorf("len(ws) = 0: len(ws) > 0 であるべき")
	}

	sum := F(0.0)
	for i, w := range ws {
		if w < 0.0 || math.IsNaN(float64(w)) || math.IsInf(float64(w), 0) {
			return -1, fmt.Errorf("ws[%d]が不正(負, NaN, Inf): ws[%d] = %v: wsの要素は、非負, 非NaN, 非Inf であるべき", i, i, w)
		}
		sum += w
	}

	// 一様ランダム
	if sum == 0.0 {
		return rng.IntN(n), nil
	}

	threshold, err := FloatRange(0.0, sum, rng)
	if err != nil {
		return -1, err
	}

	var cumulative F = 0.0
	for i, w := range ws {
		cumulative += w
		if cumulative >= threshold {
			return i, nil
		}
	}

	// 最後の要素のインデックスを返す
	return len(ws) - 1, nil
}

func FloatRange[F constraints.Float](minVal, maxVal F, rng *rand.Rand) (F, error) {
	if minVal >= maxVal {
		var zero F
		return zero, fmt.Errorf("範囲が不正(min >= max): min = %v, max = %v: min < max であるべき", minVal, maxVal)
	}

	if math.IsNaN(float64(minVal)) || math.IsInf(float64(minVal), 0) {
		var zero F
		return zero, fmt.Errorf("minが不正(NaN, Inf): min = %v: minは、非NaN, 非Inf であるべき", minVal)
	}

	if math.IsNaN(float64(maxVal)) || math.IsInf(float64(maxVal), 0) {
		var zero F
		return zero, fmt.Errorf("maxが不正(NaN, Inf): max = %v: maxは、非NaN, 非Inf であるべき", maxVal)
	}

	return F(rng.Float64())*(maxVal-minVal) + minVal, nil
}

func Choice[S ~[]E, E any](s S, rng *rand.Rand) (E, error) {
	n := len(s)
	if n == 0 {
		var zero E
		return zero, fmt.Errorf("len(s) = 0: len(s) > 0 であるべき")
	}
	idx := rng.IntN(n)
	return s[idx], nil
}

func Bool(rng *rand.Rand) bool {
	return rng.Uint32()&1 == 0
}

func NormalInt[F constraints.Float](minVal, maxVal int, mean, std F, rng *rand.Rand) (int, error) {
	if minVal > maxVal {
		return 0, fmt.Errorf("範囲が不正(min > max): min = %d, max = %d: min <= max であるべき", minVal, maxVal)
	}

	if std < 0 {
		return 0, fmt.Errorf("std < 0: std >= 0 であるべき")
	}

	if mean < F(minVal) || mean > F(maxVal) {
		return 0, fmt.Errorf("meanが範囲外: mean = %v: [%d, %d] の間であるべき", mean, minVal, maxVal)
	}

	if std == 0 {
		n := int(math.Round(float64(mean)))
		return n, nil
	}

	for {
		f := rng.NormFloat64()*float64(std) + float64(mean)
		n := int(math.Round(f))

		// 範囲内チェック
		if n >= minVal && n <= maxVal {
			return n, nil
		}
		// 範囲外ならやり直し
	}
}
