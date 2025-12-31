package randx

import (
	"fmt"
	"github.com/sw965/omw/constraints"
	"github.com/sw965/omw/mathx"
	"math/rand/v2"
)

func NewPCGFromGlobalSeed() *rand.Rand {
	// グローバル乱数を用いてシードを設定
	return rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))
}

func IntRange[I constraints.Integer](minVal, maxVal I, rng *rand.Rand) (I, error) {
	if minVal >= maxVal {
		var zero I
		return zero, fmt.Errorf("min >= max: min=%d max=%d", minVal, maxVal)
	}

	// 差分をuint64で計算することで、int64の全範囲や負の数も安全に扱える
	diff := uint64(maxVal) - uint64(minVal)
	return I(rng.Uint64N(diff)) + minVal, nil
}

func IntByWeight[F constraints.Float](ws []F, rng *rand.Rand) (int, error) {
	n := len(ws)
	if n == 0 {
		return -1, fmt.Errorf("wsが不正: len(ws)=0")
	}

	sum := F(0.0)
	for i, w := range ws {
		if w < 0.0 || mathx.IsNaN(w) || mathx.IsInf(w, 0) {
			return -1, fmt.Errorf("ws[%d]が不正(負/NaN/Inf): ws[%d]=%.6g", i, i, w)
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

	var current F = 0.0
	for i, w := range ws {
		current += w
		if current >= threshold {
			return i, nil
		}
	}

	// 最後の要素のインデックスを返す
	return len(ws) - 1, nil
}

func FloatRange[F constraints.Float](minVal, maxVal F, rng *rand.Rand) (F, error) {
	if minVal >= maxVal {
		var zero F
		return zero, fmt.Errorf("min >= max: min=%.6g max=%.6g", minVal, maxVal)
	}

	if mathx.IsNaN(minVal) || mathx.IsInf(minVal, 0) {
		var zero F
		return zero, fmt.Errorf("minが不正(NaN/Inf): min=%v", minVal)
	}

	if mathx.IsNaN(maxVal) || mathx.IsInf(maxVal, 0) {
		var zero F
		return zero, fmt.Errorf("maxが不正(NaN/Inf): max=%v", maxVal)
	}

	return F(rng.Float64())*(maxVal-minVal) + minVal, nil
}

func Choice[S ~[]E, E any](s S, rng *rand.Rand) (E, error) {
	n := len(s)
	if n == 0 {
		var zero E
		return zero, fmt.Errorf("sが不正: len(s)=0")
	}
	idx := rng.IntN(n)
	return s[idx], nil
}

func Bool(rng *rand.Rand) bool {
	n := rng.IntN(2)
	return n == 0
}
