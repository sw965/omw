package rand

import (
	"github.com/sw965/omw"
	"math/rand"
	"github.com/seehuhn/mt19937"
	"time"
	"golang.org/x/exp/slices"
)

func NewMt19937() *rand.Rand {
	y := rand.New(mt19937.New())
	y.Seed(time.Now().UnixNano())
	return y
}

func Bool(r *rand.Rand) bool {
	return r.Intn(2) == 0
}

func Float64(min, max float64, r *rand.Rand) float64 {
	return r.Float64()*(max-min) + min
}

func Int(min, max int, r *rand.Rand) int {
	return r.Intn(max-min) + min
}

func IntWithWeight(ws []float64, r *rand.Rand) int {
	sum := omw.Sum(ws...)
	if sum == 0.0 {
		return r.Intn(len(ws))
	}

	threshold := Float64(0.0, sum, r)
	accum := 0.0
	for i, w := range ws {
		accum += w
		if accum >= threshold {
			return i
		}
	}
	return len(ws) - 1
}

func Choice[XS ~[]X, X any](xs XS, r *rand.Rand) X {
	idx := r.Intn(len(xs))
	return xs[idx]
}

func Shuffled[XS ~[]X, X any](xs XS, r *rand.Rand) XS {
	clone := slices.Clone(xs)
	r.Shuffle(len(clone), func(i, j int) {clone[i], clone[j] = clone[j], clone[i]})
	return clone
}