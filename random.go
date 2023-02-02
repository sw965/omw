package omw

import (
	"github.com/seehuhn/mt19937"
	"math/rand"
	"time"
)

func NewMt19937(random *rand.Rand) *rand.Rand {
	y := rand.New(mt19937.New())
	y.Seed(time.Now().UnixNano())
	return y
}

func RandomInt(start, end int, random *rand.Rand) int {
	return random.Intn(end-start) + start
}

func RandomIntWithWeight(ws []float64, random *rand.Rand) int {
	wsSum := 0.0
	for _, w := range ws {
		wsSum += w
	}
	if wsSum == 0.0 {
		return random.Intn(len(ws))
	}
	threshold := RandomFloat64(0.0, wsSum, random)
	accumWs := 0.0
	for i, w := range ws {
		accumWs += w
		if accumWs >= threshold {
			return i
		}
	}
	return len(ws) - 1
}

func RandomFloat64(min, max float64, random *rand.Rand) float64 {
	return random.Float64()*(max-min) + min
}

func RandomBool(random *rand.Rand) bool {
	return random.Intn(2) == 0
}
