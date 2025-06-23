package bit64

import (
	"fmt"
	"math/bits"
)

func ToggleBit64[B ~uint64](b B, idx int) (B, error) {
	if idx < 0 || idx > 63 {
		return 0, fmt.Errorf("idxは0_63")
	}
	return b ^ (1 << idx), nil
}

//最下位の1ビットをクリア
func ClearLowestBit64[B ~uint64](b B) B {
	b &= b - 1
	return b
}

//最下位の1ビットを抽出
func ExtractLowestBit64[B ~uint64](b B) B {
	return b & -b
}

func OneIndices64[B ~uint64](b B) []int {
	c := bits.OnesCount64(uint64(b))
	idxs := make([]int, c)
	for i := 0; i < c; i++ {
		// 最下位の1ビットの位置を求める
		idxs[i] = bits.TrailingZeros64(uint64(b))
		b = ClearLowestBit64(b)
	}
	return idxs
}

func ToSingles64[B ~uint64](b B) []B {
	c := bits.OnesCount64(uint64(b))
	singles := make([]B, c)
	for i := 0; i < c; i++ {
		lsb := ExtractLowestBit64(b)
		singles[i] = lsb
		//抽出したビットをクリア
		b ^= lsb
	}
	return singles
}
