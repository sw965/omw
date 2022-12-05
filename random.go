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

func MakeRandomSliceInt(size, start, end int, random *rand.Rand) ([]int, error) {
	result := make([]int, size)
	for i := 0; i < size; i++ {
		v, err := RandomInt(start, end, random)
		if err != nil {
			return []int{}, err
		}
		result[i] = v
	}
	return result, nil
}

func RandomChoiceInt(random *rand.Rand, x ...int) int {
	index := random.Intn(len(x))
	return x[index]
}

func RandomIntWithWeight(weight []float64, random *rand.Rand) int {
	weightSum := SumFloat64(weight...)
	if weightSum == 0.0 {
		return random.Intn(len(weight))
	}

	r, err := RandomFloat64(0.0, weightSum, random)
	if err != nil {
		panic(err)
	}

	currentSum := 0.0
	for i, w := range weight {
		currentSum += w
		if currentSum >= r {
			return i
		}
	}
	return len(weight)
}

func RandomFloat64(min, max float64, random *rand.Rand) (float64, error) {
	if min > max {
		return 0.0, fmt.Errorf("第一引数(min) > 第二引数(max) になっている")
	}
	return random.Float64()*(max-min) + min, nil
}

func MakeRandomSliceFloat64(size int, min, max float64, random *rand.Rand) ([]float64, error) {
	result := make([]float64, size)
	for i := 0; i < size; i++ {
		v, err := RandomFloat64(min, max, random)
		if err != nil {
			return []float64{}, err
		}
		result[i] = v
	}
	return result, nil
}

func RandomBool(random *rand.Rand) bool {
	return random.Intn(2) == 0
}
