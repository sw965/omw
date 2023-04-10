package rand

import (
	"math/rand"
)

func Choice[XS ~[]X, X any](xs XS, r *rand.Rand) X {
	idx := r.Intn(len(xs))
	return xs[idx]
}