package rand

import (
	"fmt"
	"math/rand"
	"time"
	"github.com/seehuhn/mt19937"
	"golang.org/x/exp/slices"
	omwmath "github.com/sw965/omw/math"
)

func NewMt19937() *rand.Rand {
	r := rand.New(mt19937.New())
	r.Seed(time.Now().UnixNano())
	return r
}

func IntUniform(min, max int, r *rand.Rand) int {
	return r.Intn(max-min) + min
}

func IntsUniform(n, min, max int, r *rand.Rand) []int {
	ret := make([]int, n)
	for i := 0; i < n; i++ {
		ret[i] = IntUniform(min, max, r)
	}
	return ret
}

func IntByWeight(ws []float64, r *rand.Rand) int {
	wSum := omwmath.Sum(ws...)
	if wSum == 0.0 {
		return r.Intn(len(ws))
	}
	threshold := Float64Uniform(0.0, wSum, r)
	total := 0.0
	for i, w := range ws {
		total += w
		if total >= threshold {
			return i
		}
	}
	return len(ws) - 1
}

func Float64Uniform(min, max float64, r *rand.Rand) float64 {
	return r.Float64()*(max-min) + min
}

func Bool(r *rand.Rand) bool {
	return r.Intn(2) == 0
}

func Choice[S ~[]E, E any](s S, r *rand.Rand) E {
	idx := r.Intn(len(s))
	return s[idx]
}

func Sample[S ~[]E, E any](s S, n int, r *rand.Rand) S {
	ret := make(S, n)
	for i := 0; i < n; i++ {
		ret[i] = Choice(s, r)
	}
	return ret
}

func Shuffle[S ~[]E, E any](s S, r *rand.Rand) {
	r.Shuffle(len(s), func(i, j int) { s[i], s[j] = s[j], s[i] })
}

func Shuffled[S ~[]E, E any](s S, r *rand.Rand) S {
	ret := slices.Clone(s)
	Shuffle(ret, r)
	return ret
}

func IsPercentageMet(percent int, r *rand.Rand) (bool, error) {
    if percent > 100 {
        return false, fmt.Errorf("IsPercentageMetの第1引数は、100以下でなければならない。")
    }
    if percent < 0 {
        return false, fmt.Errorf("IsPercentageMetの第1引数は、0以上でなければならない。")
    }
    return r.Intn(100) < percent, nil
}