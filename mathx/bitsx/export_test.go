package bitsx

import (
	"math/rand/v2"
	"testing"
)

func NewMatrixForTest(t *testing.T, cols int, oneColIdxsPerRow [][]int) *Matrix {
	t.Helper()
	m := NewZerosMatrixForTest(t, len(oneColIdxsPerRow), cols)
	for r, oneColIdxs := range oneColIdxsPerRow {
		for _, c := range oneColIdxs {
			if err := m.Set(r, c); err != nil {
				t.Fatalf("nilを期待したが、エラーが返された: %v", err)
			}
		}
	}
	return m
}

func NewPrefixOnesMatrixForTest(t *testing.T, rows, cols, ones int) *Matrix {
	t.Helper()
	m := NewZerosMatrixForTest(t, rows, cols)
	for k := range ones {
		if err := m.Set(k/cols, k%cols); err != nil {
			t.Fatalf("nilを期待したが、エラーが返された: %v", err)
		}
	}
	return m
}

func NewZerosMatrixForTest(t *testing.T, rows, cols int) *Matrix {
	t.Helper()
	m, err := NewZerosMatrix(rows, cols)
	if err != nil {
		t.Fatalf("nilを期待したが、エラーが返された: %v", err)
	}
	return m
}

func NewOnesMatrixForTest(t *testing.T, rows, cols int) *Matrix {
	t.Helper()
	m, err := NewOnesMatrix(rows, cols)
	if err != nil {
		t.Fatalf("nilを期待したが、エラーが返された: %v", err)
	}
	return m
}

func NewRandMatrixForTest(t *testing.T, rows, cols int, rng *rand.Rand) *Matrix {
	t.Helper()
	m, err := NewRandMatrix(rows, cols, rng)
	if err != nil {
		t.Fatalf("nilを期待したが、エラーが返された: %v", err)
	}
	return m
}

func NewOnesMatrixWithSetWordForTest(t *testing.T, rows, cols, setIdx int, word uint64) *Matrix {
	t.Helper()
	m := NewOnesMatrixForTest(t, rows, cols)
	if err := m.SetWord(setIdx, word); err != nil {
		t.Fatalf("nilを期待したが、エラーが返された: %v", err)
	}
	return m
}
