package rand

import (
	"math/rand"
	"time"
	"github.com/seehuhn/mt19937"
	omath "github.com/sw965/omw/math"
)

func NewMt19937() *rand.Rand {
	rng := rand.New(mt19937.New())
	rng.Seed(time.Now().UnixNano())
	return rng
}

type integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func Integer[I integer](min, max I, rng *rand.Rand) I {
	ri := rng.Intn(int(max - min))
	return I(ri) + min
}

func IntByWeight[F float](ws []F, rng *rand.Rand) int {
	sum := omath.Sum(ws...)
	if sum == 0.0 {
		return rng.Intn(len(ws))
	}

	threshold := Float[F](0.0, sum, rng)
	var total F = 0.0
	for i, w := range ws {
		total += w
		if total >= threshold {
			return i
		}
	}
	return len(ws) - 1
}

type float interface {
	~float32 | ~float64
}

func Float[F float](min, max F, rng *rand.Rand) F {
	rf := F(rng.Float64())
	return rf*(max-min) + min
}

func Bool(rng *rand.Rand) bool {
	return rng.Intn(2) == 0
}

func Choice[S ~[]E, E any](s S, rng *rand.Rand) E {
	idx := rng.Intn(len(s))
	return s[idx]
}

func Shuffle[S ~[]E, E any](s S, rng *rand.Rand) {
	rng.Shuffle(len(s), func(i, j int) { s[i], s[j] = s[j], s[i] })
}
