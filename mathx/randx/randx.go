package randx

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"

	"github.com/sw965/omw/constraints"
)

func NewPCG() *rand.Rand {
	return rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))
}

func NewPCGs(n int) ([]*rand.Rand, error) {
	if n < 0 {
		return nil, fmt.Errorf("n >= 0 であるべき: n = %d", n)
	}
	rs := make([]*rand.Rand, n)
	for i := range n {
		rs[i] = NewPCG()
	}
	return rs, nil
}

func IntRange[I constraints.Integer](minVal, maxVal I, r *rand.Rand) (I, error) {
	if minVal >= maxVal {
		var zero I
		return zero, fmt.Errorf("min < max であるべき: min = %d, max = %d", minVal, maxVal)
	}
	diff := uint64(maxVal) - uint64(minVal)
	return I(r.Uint64N(diff)) + minVal, nil
}

func IndexByWeights[F constraints.Float](ws []F, r *rand.Rand) (int, error) {
	n := len(ws)
	if n == 0 {
		return -1, errors.New("len(ws) = 0")
	}

	sum := F(0.0)
	for i, w := range ws {
		if w < 0 {
			return -1, fmt.Errorf("重みは0以上であるべき: ws[%d] = %v", i, w)
		}
		sum += w
	}

	// 一様ランダム
	if sum == 0.0 {
		return r.IntN(n), nil
	}

	threshold, err := FloatRange(0.0, sum, r)
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
	return n - 1, nil
}

func FloatRange[F constraints.Float](minVal, maxVal F, r *rand.Rand) (F, error) {
	if minVal >= maxVal {
		var zero F
		return zero, fmt.Errorf("min < max であるべき: min = %v, max = %v", minVal, maxVal)
	}
	return F(r.Float64())*(maxVal-minVal) + minVal, nil
}

func Choice[S ~[]E, E any](s S, r *rand.Rand) (E, error) {
	n := len(s)
	if n == 0 {
		var zero E
		return zero, errors.New("len(s) = 0")
	}
	idx := r.IntN(n)
	return s[idx], nil
}

func Bool(r *rand.Rand) bool {
	return r.Uint32()&1 == 0
}

func IntNorm[F constraints.Float](minVal, maxVal int, mean, std F, r *rand.Rand) (int, error) {
	if minVal >= maxVal {
		return 0, fmt.Errorf("min < max であるべき: min = %d, max = %d", minVal, maxVal)
	}

	if std < 0 {
		return 0, errors.New("std >= 0 であるべき")
	}

	if mean < F(minVal) || mean >= F(maxVal) {
		return 0, fmt.Errorf("mean = %v: [%d, %d) の間であるべき", mean, minVal, maxVal)
	}

	if std == 0 {
		n := int(math.Round(float64(mean)))
		if n >= minVal && n < maxVal {
			return n, nil
		}
		return 0, fmt.Errorf("std = 0 かつ mean = %v の丸め結果 %d が [%d, %d) の外", mean, n, minVal, maxVal)
	}

	for {
		f := r.NormFloat64()*float64(std) + float64(mean)
		n := int(math.Round(f))

		// 範囲内チェック
		if n >= minVal && n < maxVal {
			return n, nil
		}
		// 範囲外ならやり直し
	}
}
