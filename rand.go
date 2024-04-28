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

func RandBool(r *rand.Rand) bool {
	return r.Intn(2) == 0
}

func RandFloat64(min, max float64, r *rand.Rand) float64 {
	return r.Float64()*(max-min) + min
}

func RandInt(min, max int, r *rand.Rand) int {
	return r.Intn(max-min) + min
}

func RandIntns(n, max int, r *rand.Rand) []int {
	ret := make([]int, n)
	for i := 0; i < n; i++ {
		ret[i] = r.Intn(max)
	}
	return ret
}

func RandIntByWeight(ws []float64, r *rand.Rand) int {
	wSum := Sum(ws...)
	if wSum == 0.0 {
		return r.Intn(len(ws))
	}
	threshold := RandFloat64(0.0, wSum, r)
	total := 0.0
	for i, w := range ws {
		total += w
		if total >= threshold {
			return i
		}
	}
	return len(ws) - 1
}

func RandChoice[XS ~[]X, X any](xs XS, r *rand.Rand) X {
	idx := r.Intn(len(xs))
	return xs[idx]
}

func RandSample[XS ~[]X, X any](n int, xs XS, r *rand.Rand) XS {
	ret := make(XS, n)
	for i := range ret {
		ret[i] = RandChoice(xs, r)
	}
	return ret
}

func ShuffleSlice[XS ~[]X, X any](xs XS, r *rand.Rand) {
	r.Shuffle(len(xs), func(i, j int) { xs[i], xs[j] = xs[j], xs[i] })
}