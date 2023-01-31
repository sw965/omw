package omw

import (
	"fmt"
	"math/rand"
)

func RandomInt(start, end int, random *rand.Rand) (int, error) {
	if start >= end {
		return 0, fmt.Errorf("start < end でなければならない")
	}
	return random.Intn(end-start) + start, nil
}

func RandomIntWithWeight(weights []float64, random *rand.Rand) (int, error) {
	weightsSum := Sum(weights...)
	if weightsSum == 0.0 {
		return random.Intn(len(weights)), nil
	}

	rf, err := RandomFloat64(0.0, weightsSum, random)
	if err != nil {
		return 0, err
	}

	threshold := 0.0
	for i, weight := range weights {
		threshold += weight
		if threshold >= rf {
			return i, nil
		}
	}
	return len(weights) - 1, nil
}

func RandomFloat64(min, max float64, random *rand.Rand) (float64, error) {
	if min > max {
		return 0.0, fmt.Errorf("min <= max でなければならない")
	}
	return random.Float64()*(max-min) + min, nil
}

func RandomBool(random *rand.Rand) bool {
	return random.Intn(2) == 0
}

func RandomChoice[T any](x []T, random *rand.Rand) T {
	index := random.Intn(len(x))
	return x[index]
}
