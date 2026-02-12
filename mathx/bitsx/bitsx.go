package bitsx

import (
	"fmt"
	"github.com/sw965/omw/constraints"
	"math/bits"
)

func IndexErrorMessage(idx, bitSize int) string {
	return fmt.Sprintf("bitsx: index %d out of range [0, %d)", idx, bitSize)
}

func FromIndices[B constraints.Unsigned](idxs []int) (B, error) {
	var b B
	bitSize := bits.Len64(uint64(^B(0)))

	for _, idx := range idxs {
		if idx < 0 || idx >= bitSize {
			return 0, fmt.Errorf("%s", IndexErrorMessage(idx, bitSize))
		}
		b |= 1 << idx
	}
	return b, nil
}

// 丸々Bitという関数名のBitを消す？

func ToggleBit[B constraints.Unsigned](b B, idx int) (B, error) {
	bitSize := bits.Len64(uint64(^B(0)))
	if idx < 0 || idx >= bitSize {
		return 0, fmt.Errorf("%s", IndexErrorMessage(idx, bitSize))
	}
	return b ^ (1 << idx), nil
}

func SetBit[B constraints.Unsigned](b B, idx int) (B, error) {
	bitSize := bits.Len64(uint64(^B(0)))
	if idx < 0 || idx >= bitSize {
		return 0, fmt.Errorf("%s", IndexErrorMessage(idx, bitSize))
	}
	return b | (1 << idx), nil
}

func ClearBit[B constraints.Unsigned](b B, idx int) (B, error) {
	bitSize := bits.Len64(uint64(^B(0)))
	if idx < 0 || idx >= bitSize {
		return 0, fmt.Errorf("%s", IndexErrorMessage(idx, bitSize))
	}
	return b &^ (1 << idx), nil
}

func ClearLowestBit[B constraints.Unsigned](b B) B {
	return b & (b - 1)
}

func ExtractLowestBit[B constraints.Unsigned](b B) B {
	return b & -b
}

func Indices[B constraints.Unsigned](b B) []int {
	val := uint64(b)
	c := bits.OnesCount64(val)
	idxs := make([]int, c)

	for i := 0; i < c; i++ {
		zeros := bits.TrailingZeros64(val)
		idxs[i] = zeros
		// 最下位ビットを消して次へ
		val &= (val - 1)
	}
	return idxs
}

func Singles[B constraints.Unsigned](b B) []B {
	val := uint64(b)
	c := bits.OnesCount64(val)
	singles := make([]B, c)

	for i := 0; i < c; i++ {
		lsb := val & -val // 最下位ビット抽出
		singles[i] = B(lsb)
		val ^= lsb // 抽出したビットを消去
	}
	return singles
}

func IsSubset[B constraints.Unsigned](super, sub B) bool {
	return (super | sub) == super
}
