package omw

import (
	"math/rand"
	"time"
	"github.com/seehuhn/mt19937"
)

func NewMt19937() *rand.Rand {
	r := rand.New(mt19937.New())
	r.Seed(time.Now().UnixNano())
	return r
}

func RandIntUniform(min, max int, r *rand.Rand) int {
	return r.Intn(max-min) + min
}

func RandIntsUniform(n, min, max int, r *rand.Rand) []int {
	ret := make([]int, n)
	for i := 0; i < n; i++ {
		ret[i] = RandIntUniform(min, max, r)
	}
	return ret
}

func RandIntByWeight(ws []float64, r *rand.Rand) int {
	wSum := Sum(ws...)
	if wSum == 0.0 {
		return r.Intn(len(ws))
	}
	threshold := RandFloat64Uniform(0.0, wSum, r)
	total := 0.0
	for i, w := range ws {
		total += w
		if total >= threshold {
			return i
		}
	}
	return len(ws) - 1
}

func RandFloat64Uniform(min, max float64, r *rand.Rand) float64 {
	return r.Float64()*(max-min) + min
}

func RandBool(r *rand.Rand) bool {
	return r.Intn(2) == 0
}

func RandChoice[S ~[]E, E any](s S, r *rand.Rand) E {
	idx := r.Intn(len(s))
	return s[idx]
}

func ShuffleSlice[S ~[]E, E any](s S, r *rand.Rand) {
	r.Shuffle(len(s), func(i, j int) { s[i], s[j] = s[j], s[i] })
}