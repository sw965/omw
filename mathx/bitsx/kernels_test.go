package bitsx

import (
	"math/bits"
	"math/rand/v2"
	"testing"
)

// 差し替え前の実装をリファレンスとして残し、
// 高速版(AVX-512 / pure Go フォールバック)との完全一致を確認する。

func referenceDot(m, other *Matrix) []int {
	yRows := m.Rows
	yCols := other.Rows
	counts := make([]int, yRows*yCols)
	stride := m.Stride()
	mask := m.ColTailMask()

	for r := range yRows {
		for c := range yCols {
			count := 0
			for k := range stride {
				xnor := ^(m.Data[r*stride+k] ^ other.Data[c*stride+k])
				if k == stride-1 {
					xnor &= mask
				}
				count += bits.OnesCount64(xnor)
			}
			counts[r*yCols+c] = count
		}
	}
	return counts
}

func referenceDotTernary(m, sign, nonZero *Matrix) []int {
	zRows := m.Rows
	zCols := sign.Rows
	z := make([]int, zRows*zCols)
	stride := m.Stride()
	mask := m.ColTailMask()

	for r := range zRows {
		for c := range zCols {
			matchCount := 0
			nonZeroCount := 0
			for k := range stride {
				sWord := sign.Data[c*stride+k]
				nzWord := nonZero.Data[c*stride+k]
				validMatch := ^(m.Data[r*stride+k] ^ sWord) & nzWord
				if k == stride-1 {
					validMatch &= mask
					nzWord &= mask
				}
				matchCount += bits.OnesCount64(validMatch)
				nonZeroCount += bits.OnesCount64(nzWord)
			}
			z[r*zCols+c] = 2*matchCount - nonZeroCount
		}
	}
	return z
}

func referenceHammingDistance(a, b *Matrix) int {
	diff, err := a.Xor(b)
	if err != nil {
		panic(err)
	}
	return diff.OnesCount()
}

var kernelTestShapes = []struct{ mRows, cols, oRows int }{
	{1, 1, 1},
	{1, 63, 5},
	{1, 64, 8},
	{1, 65, 3},
	{1, 100, 7},
	{3, 130, 9},
	{1, 784, 512},
	{1, 512, 1024},
	{4, 511, 33},
	{2, 1024, 10},
	{784, 512, 1},
}

func runWithBothImpls(t *testing.T, f func(t *testing.T)) {
	t.Helper()
	saved := useAVX512
	defer func() { useAVX512 = saved }()

	if saved {
		useAVX512 = true
		t.Run("AVX512", f)
	}
	useAVX512 = false
	t.Run("PureGo", f)
}

func TestDotMatchesReference(t *testing.T) {
	runWithBothImpls(t, func(t *testing.T) {
		rng := rand.New(rand.NewPCG(1, 2))
		for _, s := range kernelTestShapes {
			m, err := NewRandMatrix(s.mRows, s.cols, 0, rng)
			if err != nil {
				t.Fatal(err)
			}
			o, err := NewRandMatrix(s.oRows, s.cols, 0, rng)
			if err != nil {
				t.Fatal(err)
			}

			got, err := m.Dot(o)
			if err != nil {
				t.Fatal(err)
			}
			want := referenceDot(m, o)
			for i := range want {
				if got[i] != want[i] {
					t.Fatalf("shape %+v: counts[%d] = %d, want %d", s, i, got[i], want[i])
				}
			}
		}
	})
}

func TestDotTernaryMatchesReference(t *testing.T) {
	runWithBothImpls(t, func(t *testing.T) {
		rng := rand.New(rand.NewPCG(3, 4))
		for _, s := range kernelTestShapes {
			m, err := NewRandMatrix(s.mRows, s.cols, 0, rng)
			if err != nil {
				t.Fatal(err)
			}
			sign, err := NewRandMatrix(s.oRows, s.cols, 0, rng)
			if err != nil {
				t.Fatal(err)
			}
			nonZero, err := NewRandMatrix(s.oRows, s.cols, 0, rng)
			if err != nil {
				t.Fatal(err)
			}

			got, err := m.DotTernary(sign, nonZero)
			if err != nil {
				t.Fatal(err)
			}
			want := referenceDotTernary(m, sign, nonZero)
			for i := range want {
				if got[i] != want[i] {
					t.Fatalf("shape %+v: z[%d] = %d, want %d", s, i, got[i], want[i])
				}
			}
		}
	})
}

// --- ベンチマーク (crow の binary モデルと同じサイズ) ---
// go test -bench Kernel ./mathx/bitsx/

func benchmarkWithBothImpls(b *testing.B, f func(b *testing.B)) {
	b.Helper()
	saved := useAVX512
	defer func() { useAVX512 = saved }()

	if saved {
		useAVX512 = true
		b.Run("AVX512", f)
	}
	useAVX512 = false
	b.Run("PureGo", f)
}

func BenchmarkKernelDot784x512(b *testing.B) {
	rng := rand.New(rand.NewPCG(1, 2))
	m, _ := NewRandMatrix(1, 784, 0, rng)
	o, _ := NewRandMatrix(512, 784, 0, rng)
	benchmarkWithBothImpls(b, func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if _, err := m.Dot(o); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkKernelDotTernary784x512(b *testing.B) {
	rng := rand.New(rand.NewPCG(1, 2))
	m, _ := NewRandMatrix(784, 512, 0, rng)
	sign, _ := NewRandMatrix(1, 512, 0, rng)
	nonZero, _ := NewRandMatrix(1, 512, 0, rng)
	benchmarkWithBothImpls(b, func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if _, err := m.DotTernary(sign, nonZero); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkKernelHamming1024(b *testing.B) {
	rng := rand.New(rand.NewPCG(1, 2))
	x, _ := NewRandMatrix(1, 1024, 0, rng)
	y, _ := NewRandMatrix(1, 1024, 0, rng)
	benchmarkWithBothImpls(b, func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if _, err := x.HammingDistance(y); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func TestHammingDistanceMatchesReference(t *testing.T) {
	runWithBothImpls(t, func(t *testing.T) {
		rng := rand.New(rand.NewPCG(5, 6))
		for _, s := range kernelTestShapes {
			a, err := NewRandMatrix(s.mRows, s.cols, 0, rng)
			if err != nil {
				t.Fatal(err)
			}
			b, err := NewRandMatrix(s.mRows, s.cols, 0, rng)
			if err != nil {
				t.Fatal(err)
			}

			got, err := a.HammingDistance(b)
			if err != nil {
				t.Fatal(err)
			}
			if want := referenceHammingDistance(a, b); got != want {
				t.Fatalf("shape %+v: got %d, want %d", s, got, want)
			}
		}
	})
}
