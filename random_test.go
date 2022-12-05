package omw

import (
	"fmt"
	"github.com/seehuhn/mt19937"
	"math/rand"
	"testing"
	"time"
)

func TestRandomInt(t *testing.T) {
	mtRandom := rand.New(mt19937.New())
	mtRandom.Seed(time.Now().UnixNano())

	start := 0
	end := 100

	for i := 0; i < 1280; i++ {
		result, err := RandomInt(start, end, mtRandom)
		if err != nil {
			panic(err)
		}

		if !(result >= start && result < end) {
			t.Errorf("テスト失敗")
			break
		}
	}

	_, err := RandomInt(100, 99, mtRandom)
	if err == nil {
		t.Errorf("テスト失敗")
	}
}

func TestMakeRandomSliceInt(t *testing.T) {
	mtRandom := rand.New(mt19937.New())
	mtRandom.Seed(time.Now().UnixNano())

	testNum := 10
	for i := 0; i < testNum; i++ {
		result, err := MakeRandomSliceInt(10, 0, 10, mtRandom)
		if err != nil {
			panic(err)
		}
		fmt.Println(result)
	}
}

func TestRandomIntWithWeight(t *testing.T) {
	mtRandom := rand.New(mt19937.New())
	mtRandom.Seed(time.Now().UnixNano())
	weight := []float64{0.1, 0.2, 0.3, 0.4}
	results := make([]int, len(weight))
	testNum := 1000000

	for i := 0; i < testNum; i++ {
		result := RandomIntWithWeight(weight, mtRandom)
		results[result] += 1
	}
	fmt.Println("results = ", results, "weight = ", weight)
}

func TestMakeRandomSliceFloat64(t *testing.T) {
	mtRandom := rand.New(mt19937.New())
	mtRandom.Seed(time.Now().UnixNano())

	testNum := 10
	for i := 0; i < testNum; i++ {
		result, err := MakeRandomSliceFloat64(10, 0, 5.0, mtRandom)
		if err != nil {
			panic(err)
		}
		fmt.Println(result)
	}
}

func TestRandomBool(t *testing.T) {
	mtRandom := rand.New(mt19937.New())
	mtRandom.Seed(time.Now().UnixNano())

	simuNum := 12800
	trueCount := 0

	for i := 0; i < simuNum; i++ {
		if RandomBool(mtRandom) {
			trueCount += 1
		}
	}

	falseCount := simuNum - trueCount
	testResultMsg := fmt.Sprintf("%v ≒ %v", trueCount, falseCount)
	fmt.Println(testResultMsg)
}
