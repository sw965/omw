package bitsx

import (
	"fmt"
	"math"
	"math/rand/v2"
)

type Matrices []*Matrix

func NewETFMatrices(n, rows, cols int, iters int, rng *rand.Rand) (Matrices, error) {
	if n < 2 {
		return nil, fmt.Errorf("n >= 2 であるべき: n = %d", n)
	}

	ms := make(Matrices, n)
	for i := range n {
		m, err := NewRandMatrix(rows, cols, 0, rng)
		if err != nil {
			return nil, err
		}
		ms[i] = m
	}

	currentCost, err := ms.ETFCost()
	if err != nil {
		return nil, err
	}

	for range iters {
		nIdx := rng.IntN(n)
		rIdx := rng.IntN(rows)
		cIdx := rng.IntN(cols)

		err := ms[nIdx].Toggle(rIdx, cIdx)
		if err != nil {
			return nil, err
		}

		cost, err := ms.ETFCost()
		if err != nil {
			return nil, err
		}

		if cost < currentCost {
			currentCost = cost
		} else {
			err := ms[nIdx].Toggle(rIdx, cIdx)
			if err != nil {
				return nil, err
			}
		}
	}
	return ms, nil
}

func NewRFFMatrices(n, rows, cols int, sigma float32, rng *rand.Rand) (Matrices, error) {
	if n < 2 {
		return nil, fmt.Errorf("n >= 2 であるべき: n = %d", n)
	}

	if rows <= 0 {
		return nil, fmt.Errorf("rows > 0 であるべき: rows = %d", rows)
	}

	if cols <= 0 {
		return nil, fmt.Errorf("cols > 0 であるべき: cols = %d", cols)
	}

	totalBits := rows * cols

	omegas := make([]float32, totalBits)
	phases := make([]float32, totalBits)

	for i := range totalBits {
		omegas[i] = float32(rng.NormFloat64()) * sigma
		phases[i] = rng.Float32() * 2 * math.Pi
	}

	ms := make(Matrices, n)
	for i := range n {
		m, err := NewZerosMatrix(rows, cols)
		if err != nil {
			return nil, err
		}
		u := float32(i) / float32(n-1)

		err = m.ScanRowsWord(nil, func(ctx MatrixWordContext) error {
			var mWord uint64
			omegaWord := omegas[ctx.GlobalStart:ctx.GlobalEnd]
			phaseWord := phases[ctx.GlobalStart:ctx.GlobalEnd]
			scanErr := ctx.ScanBits(func(i, col, colT int) error {
				y := float64(omegaWord[i]*u + phaseWord[i])
				z := float32(math.Cos(y))
				if z >= 0 {
					mWord |= (1 << uint(i))
				}
				return nil
			})
			if scanErr != nil {
				return scanErr
			}
			m.data[ctx.WordIndex] = mWord
			return nil
		})

		if err != nil {
			return nil, err
		}
		ms[i] = m
	}
	return ms, nil
}

func NewThermometerMatrices(n, rows, cols int) (Matrices, error) {
	if n < 2 {
		return nil, fmt.Errorf("n >= 2 であるべき: n = %d", n)
	}

	ms := make(Matrices, n)
	totalBits := rows * cols

	for i := range n {
		m, err := NewZerosMatrix(rows, cols)
		if err != nil {
			return nil, err
		}

		scaled := i * totalBits
		numOnes := scaled / (n - 1)
		err = m.ScanRowsWord(nil, func(ctx MatrixWordContext) error {
			var word uint64
			scanErr := ctx.ScanBits(func(i, col, colT int) error {
				if (ctx.GlobalStart + i) < numOnes {
					word |= (uint64(1) << uint(i))
				}
				return nil
			})
			if scanErr != nil {
				return scanErr
			}
			m.data[ctx.WordIndex] = word
			return nil
		})

		if err != nil {
			return nil, err
		}
		ms[i] = m
	}
	return ms, nil
}

func (ms Matrices) ETFCost() (float32, error) {
	n := len(ms)
	if n < 2 {
		return 0.0, fmt.Errorf("n >= 2 であるべき: n = %d", n)
	}

	distances := make([]float32, 0, n*n)
	sum := float32(0.0)
	for i := range len(ms) {
		for j := i + 1; j < len(ms); j++ {
			distance, err := ms[i].HammingDistance(ms[j])
			if err != nil {
				return 0.0, err
			}
			d := float32(distance)
			distances = append(distances, d)
			sum += d
		}
	}

	dn := len(distances)
	dnf := float32(dn)
	// 距離の平均
	mean := sum / dnf

	// 距離の分散
	variance := float32(0.0)
	for _, d := range distances {
		deviation := d - mean
		variance += deviation * deviation
	}
	variance /= dnf

	// コスト = -(合計距離) + (距離の分散)
	// 距離を最大化したいので、合計距離にはマイナスをつけて最小化問題にする
	// また、距離の分散もコストに含める事で、均等な距離を保ちながら、距離を最大化する事が出来る
	cost := -sum + variance
	return cost, nil
}
