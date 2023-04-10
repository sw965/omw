package rand

import (
	"math/rand"
	"github.com/sw965/omw/math"
)

func Int(start, end int, r *rand.Rand) int {
	return r.Intn(end-start) + start
}

func IntWithWeight(ws []float64, r *rand.Rand) int {
	sum := math.Sum(ws...)
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