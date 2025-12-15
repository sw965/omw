// Package bitsx provides generic bitwise operations for unsigned integers.
//
// Package bitsx は、符号なし整数に対するジェネリックなビット演算機能を提供します。
package bitsx

import (
	"fmt"
	"math/bits"

	"github.com/sw965/omw/constraints"
)

// IndexError records an error concerning an invalid index.
//
// IndexError は、無効なインデックスに関するエラー情報を記録します。
type IndexError struct {
	Index   int
	BitSize int
}

// Error returns the string representation of the error.
//
// Error はエラーの文字列表現を返します。
func (e *IndexError) Error() string {
	return fmt.Sprintf("インデックス %d はこの型の範囲外です (0-%d)", e.Index, e.BitSize-1)
}

// FromIndices creates a bitset with bits set at the specified indices.
// It returns an error if any index is out of range for the type B.
//
// FromIndices は、指定されたインデックスのビットを立てたビットセットを作成します。
// インデックスが型 B の範囲外である場合はエラーを返します。
func FromIndices[B constraints.Unsigned](idxs []int) (B, error) {
	var b B
	bitSize := bits.Len64(uint64(^B(0)))

	for _, idx := range idxs {
		if idx < 0 || idx >= bitSize {
			return 0, &IndexError{Index: idx, BitSize: bitSize}
		}
		b |= 1 << idx
	}
	return b, nil
}

// ToggleBit toggles the bit at the specified index.
//
// ToggleBit は、指定されたインデックスのビットを反転させます。
func ToggleBit[B constraints.Unsigned](b B, idx int) (B, error) {
	bitSize := bits.Len64(uint64(^B(0)))
	if idx < 0 || idx >= bitSize {
		return 0, &IndexError{Index: idx, BitSize: bitSize}
	}
	return b ^ (1 << idx), nil
}

// SetBit sets the bit at the specified index to 1.
//
// SetBit は、指定されたインデックスのビットを1に設定します。
func SetBit[B constraints.Unsigned](b B, idx int) (B, error) {
	bitSize := bits.Len64(uint64(^B(0)))
	if idx < 0 || idx >= bitSize {
		return 0, &IndexError{Index: idx, BitSize: bitSize}
	}
	return b | (1 << idx), nil
}

// ClearBit sets the bit at the specified index to 0.
//
// ClearBit は、指定されたインデックスのビットを0に設定します。
func ClearBit[B constraints.Unsigned](b B, idx int) (B, error) {
	bitSize := bits.Len64(uint64(^B(0)))
	if idx < 0 || idx >= bitSize {
		return 0, &IndexError{Index: idx, BitSize: bitSize}
	}
	return b &^ (1 << idx), nil
}

// ClearLowestBit clears the least significant bit that is set to 1.
// (e.g., 1010 -> 1000)
//
// ClearLowestBit は、1にセットされている最下位ビットをクリアします。
func ClearLowestBit[B constraints.Unsigned](b B) B {
	return b & (b - 1)
}

// ExtractLowestBit isolates the least significant bit that is set to 1.
// (e.g., 1010 -> 0010)
//
// ExtractLowestBit は、1にセットされている最下位ビットのみを抽出します。
func ExtractLowestBit[B constraints.Unsigned](b B) B {
	return b & -b
}

// Indices returns a slice of indices where bits are set to 1.
//
// Indices は、ビットが1にセットされている位置（インデックス）のスライスを返します。
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

// Singles decomposes the bitset into a slice of bitsets, each containing a single bit.
// (e.g., 1010 -> [0010, 1000])
//
// Singles は、ビットセットを「1ビットだけが立った値」のスライスに分解します。
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

// IsSubset reports whether u is a subset of a (meaning all bits set in a are also set in u).
// Note: The original logic `(u|a) == u` checks if `a` is a subset of `u`.
// Let's align with the standard set theory: Is `sub` a subset of `super`?
//
// IsSubset は、sub の全てのビットが super に含まれているか（部分集合か）を判定します。
func IsSubset[B constraints.Unsigned](super, sub B) bool {
	return (super | sub) == super
}
