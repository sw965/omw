package bitsx

import (
	"fmt"
	"math/bits"
	"math/rand/v2"

	"github.com/sw965/omw/constraints"
)

func FromIndices[B constraints.Unsigned](idxs []int) (B, error) {
	var b B
	var err error
	for _, idx := range idxs {
		b, err = Set(b, idx)
		if err != nil {
			return 0, err
		}
	}
	return b, nil
}

// RandHalfPow は、各ビットが確率 (1/2)^n で 1 になる値を返す。
// 各ビットが確率 1/2 で 1 の乱数を n 個 AND する。n = 0 なら全ビット 1。
func RandHalfPow[B constraints.Unsigned](n int, rng *rand.Rand) (B, error) {
	if n < 0 {
		return 0, fmt.Errorf("n >= 0 であるべき: n = %d", n)
	}
	b := ^B(0)
	for range n {
		b &= B(rng.Uint64())
	}
	return b, nil
}

func Size[B constraints.Unsigned]() int {
	return bits.Len64(uint64(^B(0)))
}

func Bit[B constraints.Unsigned](b B, idx int) (B, error) {
	if err := validateIndex[B](idx); err != nil {
		return 0, err
	}
	return (b >> idx) & 1, nil
}

func Set[B constraints.Unsigned](b B, idx int) (B, error) {
	if err := validateIndex[B](idx); err != nil {
		return 0, err
	}
	return b | (1 << idx), nil
}

func Toggle[B constraints.Unsigned](b B, idx int) (B, error) {
	if err := validateIndex[B](idx); err != nil {
		return 0, err
	}
	return b ^ (1 << idx), nil
}

func Clear[B constraints.Unsigned](b B, idx int) (B, error) {
	if err := validateIndex[B](idx); err != nil {
		return 0, err
	}
	return b &^ (1 << idx), nil
}

func ClearLowest[B constraints.Unsigned](b B) B {
	// b - 1 を計算すると、元の値bのビット列は次のように変化する
	// 1. 最も右側の1が0になる
	// 2. その1より右側にあるすべての0が1になる
	// 3. それより左側のビットは一切変わらない
	//
	// 例1: b = 12 = 1100(2進数)の時
	// b - 1 = 11 = 1011(2進数)
	// b & (b - 1) = 1100 & 1011 = 1000
	//
	// 例2: b = 22 = 10110(2進数)の時
	// b - 1 = 21 = 10101(2進数)
	// b & (b - 1) = 10110 & 10101 = 10100
	//
	// どちらの例でも、結果として元のbの「最も右側(下位)にある1」だけが消去されている
	return b & (b - 1)
}

func ExtractLowest[B constraints.Unsigned](b B) B {
	// -bを計算すると、元の値bのビット列は次のように変化する
	// 1. 最も右側の1は、そのまま（1のまま）
	// 2. その1より右側にあるすべての0も、そのまま（0のまま）。
	// 3. それより左側のビットは、すべて反転する。
	//
	// 例1: b = 12 = 1100(2進数)の時
	// -b = -12 = 0100(2進数)
	// b & -b = 1100 & 0100 = 0100
	//
	// 例2: b = 22 = 10110(2進数)の時
	// -b = -22 = 01010(2進数)
	// b & -b = 10110 & 01010 = 00010
	//
	// 結果として、元のbから「最も右側の1」だけを抽出した値が得られる
	return b & -b
}

func Indices[B constraints.Unsigned](b B) []int {
	c := bits.OnesCount64(uint64(b))
	idxs := make([]int, c)
	for i := range c {
		idxs[i] = bits.TrailingZeros64(uint64(b))
		b = ClearLowest(b)
	}
	return idxs
}

func Singles[B constraints.Unsigned](b B) []B {
	c := bits.OnesCount64(uint64(b))
	singles := make([]B, c)
	for i := range c {
		lsb := ExtractLowest(b)
		singles[i] = lsb
		b = ClearLowest(b)
	}
	return singles
}

func IsSubset[B constraints.Unsigned](super, sub B) bool {
	return (super | sub) == super
}

func indexError(idx, size int) error {
	return fmt.Errorf("0 <= index < %d であるべき: index = %d", size, idx)
}

func validateIndex[B constraints.Unsigned](idx int) error {
	size := Size[B]()
	if idx < 0 || idx >= size {
		return indexError(idx, size)
	}
	return nil
}
