package randx

import (
	"fmt"
	"github.com/sw965/omw/constraints"
	"github.com/sw965/omw/mathx"
	"math"
	"math/rand/v2"
)

func NewPCGFromGlobalSeed() *rand.Rand {
	// グローバル乱数を用いてシードを設定
	return rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))
}

func IntRange[I constraints.Integer](minVal, maxVal I, rng *rand.Rand) (I, error) {
	if minVal >= maxVal {
		var zero I
		return zero, fmt.Errorf("範囲が不正(min >= max): min = %d, max = %d: min < max であるべき", minVal, maxVal)
	}

	// uint64にキャストする理由を書く
	diff := uint64(maxVal) - uint64(minVal)
	return I(rng.Uint64N(diff)) + minVal, nil
}

// ws → w に命名変更をするべき？
func IntByWeights[F constraints.Float](ws []F, rng *rand.Rand) (int, error) {
	n := len(ws)
	if n == 0 {
		return -1, fmt.Errorf("len(ws) = 0: len(ws) > 0 であるべき")
	}

	sum := F(0.0)
	for i, w := range ws {
		if w < 0.0 || mathx.IsNaN(w) || mathx.IsInf(w, 0) {
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
		return zero, fmt.Errorf("範囲が不正(min >= max): min = %v, max = %v: min < max であるべき", minVal, maxVal)
	}

	if mathx.IsNaN(minVal) || mathx.IsInf(minVal, 0) {
		var zero F
		return zero, fmt.Errorf("minが不正(NaN, Inf): min = %v: minは、非NaN, 非Inf であるべき", minVal)
	}

	if mathx.IsNaN(maxVal) || mathx.IsInf(maxVal, 0) {
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

func NormalInt(minVal, maxVal int, mean, std float64, rng *rand.Rand) (int, error) {
	if minVal > maxVal {
		return 0, fmt.Errorf("範囲指定が不正(min > max): min=%d, max=%d: min < max であるべき", minVal, maxVal)
	}

	if std < 0 {
		return 0, fmt.Errorf("std < 0: std >= 0 であるべき")
	}

	if mean < float64(minVal) || mean > float64(maxVal) {
        return 0, fmt.Errorf("meanが範囲外: mean=%v: [%d, %d] の間であるべき", mean, minVal, maxVal)
    }

	if std == 0 {
        n := int(math.Round(mean))
        return n, nil
    }

	for {
		f := rng.NormFloat64()*std + mean
		n := int(math.Round(f))

		// 範囲内チェック
		if n >= minVal && n <= maxVal {
			return n, nil
		}
		// 範囲外ならやり直し
	}
}