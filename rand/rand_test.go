package rand_test

import (
	"fmt"
	"github.com/sw965/omw"
	"github.com/sw965/omw/rand"
	"testing"
)

func TestBool(t *testing.T) {
	r := rand.NewMt19937()
	testNum := 1280
	count := 0
	for i := 0; i < testNum; i++ {
		if rand.Bool(r) {
			count += 1
		}
	}
	fmt.Println(float64(testNum)/2.0, "≒", count)
}

func TestFloat64(t *testing.T) {
	r := rand.NewMt19937()
	min, max := 1.0, 5.0
	testNum := 1280
	results := make([]float64, testNum)
	for i := 0; i < testNum; i++ {
		v := rand.Float64(min, max, r)
		results[i] = v
	}

	fmt.Println(omw.Min(results...), "≒", min)
	fmt.Println(omw.Max(results...), "≒", max)
	fmt.Println(omw.Mean(results...), "≒", (min+max)/2.0)
}

func TestInt(t *testing.T) {
	r := rand.NewMt19937()
	min, max := 100, 251
	testNum := 2560
	results := make([]int, testNum)
	for i := 0; i < testNum; i++ {
		v := rand.Int(min, max, r)
		results[i] = v
	}

	fmt.Println(omw.Min(results...), "≒", min)
	fmt.Println(omw.Max(results...), "≒", max-1)
	fmt.Println(omw.Mean(results...), "≒", (min+max-1)/2.0)
}

func TestIntWithWeight(t *testing.T) {
	r := rand.NewMt19937()
	ws := []float64{0.1, 0.2, 0.5, 0.15, 0.05}
	testNum := 10000
	counts := make([]int, len(ws))
	for i := 0; i < testNum; i++ {
		idx := rand.IntWithWeight(ws, r)
		counts[idx] += 1
	}
	expected := []int{1000, 2000, 5000, 1500, 500}
	fmt.Println(counts, "≒", expected)
}

func TestChoice(t *testing.T) {
	r := rand.NewMt19937()
	xs := []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5, 5}
	testNum := 10000
	counts := make([]int, 6)
	for i := 0; i < testNum; i++ {
		num := rand.Choice(xs, r)
		counts[num] += 1
	}

	n := float64(len(xs))
	ps := make([]float64, len(counts))
	for i, c := range counts {
		ps[i] = float64(c) / float64(testNum)
	}
	expected := []float64{0.0, 1.0 / n, 2.0 / n, 3.0 / n, 4.0 / n, 5.0 / n}
	fmt.Println(ps, "≒", expected)
}
