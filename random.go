package omw

import (
	"math/rand"
	"github.com/seehuhn/mt19937"
	"time"
)

func NewMt19937() *rand.Rand {
	y := rand.New(mt19937.New())
	y.Seed(time.Now().UnixNano())
	return y
}

func RandomInt(start, end int, r *rand.Rand) int {
	return r.Intn(end-start) + start
}

func RandomIntWithWeight(ws []float64, r *rand.Rand) int {
	sum := Sum(ws...)
	if sum == 0.0 {
		return r.Intn(len(ws))
	}

	threshold := RandomFloat64(0.0, sum, r)
	accum := 0.0
	for i, w := range ws {
		accum += w
		if accum >= threshold {
			return i
		}
	}
	return len(ws) - 1
}

func RandomFloat64(min, max float64, r *rand.Rand) float64 {
	return r.Float64()*(max-min) + min
}

func RandomBool(r *rand.Rand) bool {
	return r.Intn(2) == 0
}

func RandomChoice[X any](xs []T, r *rand.Rand) X {
	idx := r.Intn(len(xs))
	return xs[idx]
}