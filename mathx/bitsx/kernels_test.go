package bitsx

import (
	"math/rand/v2"
	"testing"
)

// カーネルの検証は2段階に分ける。
//  1. pure Go版(dotGo / dotTernaryGo)を、手計算した期待値と性質で検証する。
//  2. AVX-512版が pure Go版と一致することを、多数の形状で検証する。
//
// 1で参照実装を書くと本体と同じ重複になる為、期待値と性質だけで検証する。
//
// なお pure Go版は「端数ビットは常に0」という Matrix の不変条件を前提に
// 列マスク処理を省略している。その不変条件自体は、コンストラクタ側のテスト
// (NewOnesMatrixのOnesCountが rows*cols になること)が担保している。

// newTestMatrix は、行ごとに立てる列を指定して Matrix を作る。
func newTestMatrix(t *testing.T, cols int, setColsPerRow [][]int) *Matrix {
	t.Helper()
	m, err := NewZerosMatrix(len(setColsPerRow), cols)
	if err != nil {
		t.Fatal(err)
	}

	for r, setCols := range setColsPerRow {
		for _, c := range setCols {
			if err := m.Set(r, c); err != nil {
				t.Fatal(err)
			}
		}
	}
	return m
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

func callDotGo(left, right *Matrix) []int {
	results := make([]int, left.Rows*right.Rows)
	dotGo(left.Data, right.Data, left.Rows, right.Rows, left.Cols, left.Stride(), results)
	return results
}

func callDotTernaryGo(value, sign, nonZero *Matrix) []int {
	results := make([]int, value.Rows*sign.Rows)
	dotTernaryGo(value.Data, sign.Data, nonZero.Data, value.Rows, sign.Rows, value.Stride(), results)
	return results
}

func assertResults(t *testing.T, name string, got, want []int) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("%s: 長さの不一致: got = %d, want = %d", name, len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("%s: [%d] = %d, want %d (全体: got = %v, want = %v)", name, i, got[i], want[i], got, want)
		}
	}
}

// dotGo は「left行rとright行cで、値が一致する列の数」を results[r*rightRows+c] に書く。
func TestDotGoExpectedValues(t *testing.T) {
	t.Run("1列", func(t *testing.T) {
		one := newTestMatrix(t, 1, [][]int{{0}})
		zero := newTestMatrix(t, 1, [][]int{{}})
		assertResults(t, "1と1", callDotGo(one, one), []int{1})
		assertResults(t, "1と0", callDotGo(one, zero), []int{0})
		assertResults(t, "0と0", callDotGo(zero, zero), []int{1})
	})

	t.Run("4列_複数行", func(t *testing.T) {
		// left  行0 = 1111, 行1 = 1000
		// right 行0 = 0000, 行1 = 1100, 行2 = 1111
		left := newTestMatrix(t, 4, [][]int{{0, 1, 2, 3}, {0}})
		right := newTestMatrix(t, 4, [][]int{{}, {0, 1}, {0, 1, 2, 3}})

		// 行0: 1111 vs 0000 = 0個, vs 1100 = 2個, vs 1111 = 4個
		// 行1: 1000 vs 0000 = 3個, vs 1100 = 3個, vs 1111 = 1個
		want := []int{0, 2, 4, 3, 3, 1}
		assertResults(t, "4列", callDotGo(left, right), want)
	})

	t.Run("65列_2ワード目にまたがる", func(t *testing.T) {
		// 列64だけを立てる。1ワード目は全て0、2ワード目に1ビットだけ
		left := newTestMatrix(t, 65, [][]int{{64}})
		right := newTestMatrix(t, 65, [][]int{{}, {64}, {0}})

		// vs 全0   : 列64だけ不一致 → 64個
		// vs 列64  : 全て一致       → 65個
		// vs 列0   : 列0と列64が不一致 → 63個
		want := []int{64, 65, 63}
		assertResults(t, "65列", callDotGo(left, right), want)
	})

	t.Run("64列_端数なし", func(t *testing.T) {
		left := newTestMatrix(t, 64, [][]int{{0, 63}})
		right := newTestMatrix(t, 64, [][]int{{63}})
		// 列0だけ不一致 → 63個
		assertResults(t, "64列", callDotGo(left, right), []int{63})
	})
}

func TestDotGoProperties(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	colsList := []int{1, 63, 64, 65, 100, 130, 512}

	for _, cols := range colsList {
		zeros, err := NewZerosMatrix(3, cols)
		if err != nil {
			t.Fatal(err)
		}
		ones, err := NewOnesMatrix(3, cols)
		if err != nil {
			t.Fatal(err)
		}

		for i, v := range callDotGo(zeros, zeros) {
			if v != cols {
				t.Fatalf("cols=%d: 全0同士[%d] = %d, want %d", cols, i, v, cols)
			}
		}
		for i, v := range callDotGo(ones, ones) {
			if v != cols {
				t.Fatalf("cols=%d: 全1同士[%d] = %d, want %d", cols, i, v, cols)
			}
		}
		for i, v := range callDotGo(zeros, ones) {
			if v != 0 {
				t.Fatalf("cols=%d: 全0と全1[%d] = %d, want 0", cols, i, v)
			}
		}

		a, err := NewRandMatrix(4, cols, 0, rng)
		if err != nil {
			t.Fatal(err)
		}
		b, err := NewRandMatrix(6, cols, 0, rng)
		if err != nil {
			t.Fatal(err)
		}

		// 自分自身との対角は、必ず全列一致
		self := callDotGo(a, a)
		for r := range a.Rows {
			if got := self[r*a.Rows+r]; got != cols {
				t.Fatalf("cols=%d: 対角[%d] = %d, want %d", cols, r, got, cols)
			}
		}

		// 一致数は左右を入れ替えても変わらない
		ab := callDotGo(a, b)
		ba := callDotGo(b, a)
		for r := range a.Rows {
			for c := range b.Rows {
				if ab[r*b.Rows+c] != ba[c*a.Rows+r] {
					t.Fatalf("cols=%d: 対称性が崩れた (r=%d, c=%d)", cols, r, c)
				}
			}
		}
	}
}

// dotTernaryGo は「nonZero行cで有効な列のうち、一致した数 - 不一致だった数」を
// results[r*signRows+c] に書く。
func TestDotTernaryGoExpectedValues(t *testing.T) {
	t.Run("4列_単一行", func(t *testing.T) {
		// value = 1100, sign = 1010, nonZero = 1110
		// 列0: 1と1で一致(+1) / 列1: 1と0で不一致(-1) / 列2: 0と1で不一致(-1) / 列3: 無効
		value := newTestMatrix(t, 4, [][]int{{0, 1}})
		sign := newTestMatrix(t, 4, [][]int{{0, 2}})
		nonZero := newTestMatrix(t, 4, [][]int{{0, 1, 2}})
		assertResults(t, "4列", callDotTernaryGo(value, sign, nonZero), []int{-1})
	})

	t.Run("nonZeroが全0なら常に0", func(t *testing.T) {
		value := newTestMatrix(t, 4, [][]int{{0, 1, 2, 3}})
		sign := newTestMatrix(t, 4, [][]int{{}})
		nonZero := newTestMatrix(t, 4, [][]int{{}})
		assertResults(t, "全て無効", callDotTernaryGo(value, sign, nonZero), []int{0})
	})

	t.Run("4列_複数行", func(t *testing.T) {
		// value   行0 = 1100, 行1 = 0000
		// sign    行0 = 1100, 行1 = 0011, 行2 = 0000
		// nonZero 行0 = 1111, 行1 = 1111, 行2 = 1100
		value := newTestMatrix(t, 4, [][]int{{0, 1}, {}})
		sign := newTestMatrix(t, 4, [][]int{{0, 1}, {2, 3}, {}})
		nonZero := newTestMatrix(t, 4, [][]int{{0, 1, 2, 3}, {0, 1, 2, 3}, {0, 1}})

		// value行0: 全一致=+4 / 全不一致=-4 / 有効2列とも不一致=-2
		// value行1: 一致2・不一致2=0 / 一致2・不一致2=0 / 有効2列とも一致=+2
		want := []int{4, -4, -2, 0, 0, 2}
		assertResults(t, "4列複数行", callDotTernaryGo(value, sign, nonZero), want)
	})

	t.Run("65列_2ワード目にまたがる", func(t *testing.T) {
		value := newTestMatrix(t, 65, [][]int{{64}})
		sign := newTestMatrix(t, 65, [][]int{{64}, {}})
		nonZero := newTestMatrix(t, 65, [][]int{{64}, {64}})
		// 列64のみ有効。一致で+1、不一致で-1
		assertResults(t, "65列", callDotTernaryGo(value, sign, nonZero), []int{1, -1})
	})
}

func TestDotTernaryGoProperties(t *testing.T) {
	rng := rand.New(rand.NewPCG(3, 4))
	colsList := []int{1, 63, 64, 65, 100, 130, 512}

	for _, cols := range colsList {
		value, err := NewRandMatrix(4, cols, 0, rng)
		if err != nil {
			t.Fatal(err)
		}
		sign, err := NewRandMatrix(5, cols, 0, rng)
		if err != nil {
			t.Fatal(err)
		}

		// nonZeroが全0なら、有効な列が無いので常に0
		allZero, err := NewZerosMatrix(5, cols)
		if err != nil {
			t.Fatal(err)
		}
		for i, v := range callDotTernaryGo(value, sign, allZero) {
			if v != 0 {
				t.Fatalf("cols=%d: nonZeroが全0なのに[%d] = %d", cols, i, v)
			}
		}

		// nonZeroが全1なら、全列が有効。
		// 一致数をa個とすると 結果 = a - (cols - a) = 2a - cols となり、
		// aは dotGo が返す値そのものになる。
		allOne, err := NewOnesMatrix(5, cols)
		if err != nil {
			t.Fatal(err)
		}
		got := callDotTernaryGo(value, sign, allOne)
		matched := callDotGo(value, sign)
		for i := range got {
			if want := 2*matched[i] - cols; got[i] != want {
				t.Fatalf("cols=%d: [%d] = %d, want %d (一致数 = %d)", cols, i, got[i], want, matched[i])
			}
		}
	}
}

// AVX-512版が pure Go版と一致することを、多数の形状で確認する。
func TestDotAVX512MatchesGo(t *testing.T) {
	if !useAVX512 {
		t.Skip("AVX-512(VPOPCNTQ)非対応のCPU")
	}

	rng := rand.New(rand.NewPCG(1, 2))
	forEachKernelShape(t, func(t *testing.T, leftRows, cols, rightRows int) {
		left, err := NewRandMatrix(leftRows, cols, 0, rng)
		if err != nil {
			t.Fatal(err)
		}
		right, err := NewRandMatrix(rightRows, cols, 0, rng)
		if err != nil {
			t.Fatal(err)
		}

		got := make([]int, leftRows*rightRows)
		dotAVX512(&left.Data[0], &right.Data[0], leftRows, rightRows, cols, left.Stride(), &got[0])
		assertResults(t, "AVX512とpure Go", got, callDotGo(left, right))
	})
}

func TestDotTernaryAVX512MatchesGo(t *testing.T) {
	if !useAVX512 {
		t.Skip("AVX-512(VPOPCNTQ)非対応のCPU")
	}

	rng := rand.New(rand.NewPCG(3, 4))
	forEachKernelShape(t, func(t *testing.T, valueRows, cols, signRows int) {
		value, err := NewRandMatrix(valueRows, cols, 0, rng)
		if err != nil {
			t.Fatal(err)
		}
		sign, err := NewRandMatrix(signRows, cols, 0, rng)
		if err != nil {
			t.Fatal(err)
		}
		nonZero, err := NewRandMatrix(signRows, cols, 0, rng)
		if err != nil {
			t.Fatal(err)
		}

		got := make([]int, valueRows*signRows)
		dotTernaryAVX512(&value.Data[0], &sign.Data[0], &nonZero.Data[0],
			valueRows, signRows, value.Stride(), &got[0])
		assertResults(t, "AVX512とpure Go", got, callDotTernaryGo(value, sign, nonZero))
	})
}

// forEachKernelShape は、実サイズの形状に加えて、
// ワード境界(64の倍数)と8行ブロックの端数を網羅する小さい形状を総当たりする。
func forEachKernelShape(t *testing.T, f func(t *testing.T, aRows, cols, bRows int)) {
	t.Helper()
	for _, s := range kernelTestShapes {
		f(t, s.mRows, s.cols, s.oRows)
	}

	for _, cols := range []int{1, 63, 64, 65, 127, 128, 129, 511, 512, 513, 576} {
		for aRows := 1; aRows <= 4; aRows++ {
			for bRows := 1; bRows <= 17; bRows++ {
				f(t, aRows, cols, bRows)
			}
		}
	}
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
	m, err := NewRandMatrix(1, 784, 0, rng)
	if err != nil {
		b.Fatal(err)
	}
	o, err := NewRandMatrix(512, 784, 0, rng)
	if err != nil {
		b.Fatal(err)
	}
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
	m, err := NewRandMatrix(784, 512, 0, rng)
	if err != nil {
		b.Fatal(err)
	}
	sign, err := NewRandMatrix(1, 512, 0, rng)
	if err != nil {
		b.Fatal(err)
	}
	nonZero, err := NewRandMatrix(1, 512, 0, rng)
	if err != nil {
		b.Fatal(err)
	}
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
	x, err := NewRandMatrix(1, 1024, 0, rng)
	if err != nil {
		b.Fatal(err)
	}
	y, err := NewRandMatrix(1, 1024, 0, rng)
	if err != nil {
		b.Fatal(err)
	}
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
