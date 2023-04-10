package rand

import (
	"math/rand"
	"github.com/seehuhn/mt19937"
	"time"
)

func NewMt19937() *rand.Rand {
	y := rand.New(mt19937.New())
	y.Seed(time.Now().UnixNano())
	return y
}