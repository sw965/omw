package rand

import (
	"fmt"
	"math/rand"
	"time"
	"github.com/seehuhn/mt19937"
	omwmath "github.com/sw965/omw/math"
	"golang.org/x/exp/slices"
)

func NewMt19937() *rand.Rand {
	r := rand.New(mt19937.New())
	r.Seed(time.Now().UnixNano())
	return r
}

func Int(min, max int, r *rand.Rand) int {
	return r.Intn(max-min) + min
}

func Ints(n, min, max int, r *rand.Rand) []int {
	s := make([]int, n)
	for i := 0; i < n; i++ {
		s[i] = Int(min, max, r)
	}
	return s
}

func IntByWeight(ws []float64, r *rand.Rand) int {
	sum := omwmath.Sum(ws...)
	if sum == 0.0 {
		return r.Intn(len(ws))
	}

	threshold := Float64(0.0, sum, r)
	total := 0.0
	for i, w := range ws {
		total += w
		if total >= threshold {
			return i
		}
	}
	return len(ws) - 1
}

func Float64(min, max float64, r *rand.Rand) float64 {
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
	y := make(S, n)
	for i := 0; i < n; i++ {
		y[i] = Choice(s, r)
	}
	return y
}

func Shuffle[S ~[]E, E any](s S, r *rand.Rand) {
	r.Shuffle(len(s), func(i, j int) { s[i], s[j] = s[j], s[i] })
}

func IsPercentageMet(p int, r *rand.Rand) (bool, error) {
    if p > 100 {
        return false, fmt.Errorf("IsPercentageMetの第1引数は、100以下でなければならない。")
    }
    if p < 0 {
        return false, fmt.Errorf("IsPercentageMetの第1引数は、0以上でなければならない。")
    }
    return r.Intn(100) < p, nil
}