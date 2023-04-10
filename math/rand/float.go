package rand

import (
	"math/rand"
)

func Float64(min, max float64, r *rand.Rand) float64 {
	return r.Float64()*(max-min) + min
}