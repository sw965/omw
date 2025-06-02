package rand

import (
	"fmt"
	"math/rand"
	"time"
	"github.com/seehuhn/mt19937"
	"math"
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

func IntByWeight[F float](ws []F, rng *rand.Rand) (int, error) {
	sum := omath.Sum(ws...)
	if sum == 0.0 {
		return rng.Intn(len(ws)), nil
	}

	threshold := Float[F](0.0, sum, rng)
	var total F = 0.0
	for i, w := range ws {
		if w < 0.0 {
			msg := fmt.Sprintf("%d 番目の重みがマイナスです", i)
			return -1, fmt.Errorf(msg)
		}

		if math.IsNaN(float64(w)) {
			msg := fmt.Sprintf("%d 番目の重みがNanです。", i)
			return -1, fmt.Errorf(msg)
		}

		total += w
		if total >= threshold {
			return i, nil
		}
	}
	return len(ws) - 1, nil
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

func Choice[S ~[]E, E any](s S, rng *rand.Rand) (E, error) {
	if len(s) == 0 {
		var e E
		return e, fmt.Errorf("要素数が0である為、ランダムに取り出す事が出来ません。")
	}
	idx := rng.Intn(len(s))
	return s[idx], nil
}

func Shuffle[S ~[]E, E any](s S, rng *rand.Rand) {
	rng.Shuffle(len(s), func(i, j int) { s[i], s[j] = s[j], s[i] })
}
