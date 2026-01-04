package bitsx

import (
	"fmt"
	"math/bits"

	"github.com/sw965/omw/constraints"
)

func IndexErrorMessage(idx, bitSize int) string {
    return fmt.Sprintf("bitsx: index %d out of range [0, %d)", idx, bitSize)
}

func FromIndices[B constraints.Unsigned](idxs []int) (B, error) {
	var b B
	bitSize := bits.Len64(uint64(^B(0)))

	for _, idx := range idxs {
		if idx < 0 || idx >= bitSize {
			return 0, fmt.Errorf(IndexErrorMessage(idx, bitSize))
		}
		b |= 1 << idx
	}
	return b, nil
}

func ToggleBit[B constraints.Unsigned](b B, idx int) (B, error) {
	bitSize := bits.Len64(uint64(^B(0)))
	if idx < 0 || idx >= bitSize {
		return 0, fmt.Errorf(IndexErrorMessage(idx, bitSize))
	}
	return b ^ (1 << idx), nil
}

func SetBit[B constraints.Unsigned](b B, idx int) (B, error) {
	bitSize := bits.Len64(uint64(^B(0)))
	if idx < 0 || idx >= bitSize {
		return 0, fmt.Errorf(IndexErrorMessage(idx, bitSize))
	}
	return b | (1 << idx), nil
}

func ClearBit[B constraints.Unsigned](b B, idx int) (B, error) {
	bitSize := bits.Len64(uint64(^B(0)))
	if idx < 0 || idx >= bitSize {
		return 0, fmt.Errorf(IndexErrorMessage(idx, bitSize))
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

func XNOR[B constraints.Unsigned](a, b B) B {
	return ^(a ^ b)
}

type Matrix struct {
	Rows       int
	Cols       int
	Stride     int // (Cols + 63) / 64 が基本的な値: 例 Cols = 100 の時、uint64 * 2 = 128bitを確保 (Stride = 2)
	Data       []uint64
	RowMask    uint64
}

func NewMatrix(rows, cols int) (Matrix, error) {
	if rows <= 0 {
		return Matrix{}, fmt.Errorf("rows <= 0: rows > 0 であるべき")
	}

	if cols <= 0 {
		return Matrix{}, fmt.Errorf("cols <= 0: cols > 0 であるべき")
	}

	stride := (cols + 63) / 64
	
	r := cols % 64
	// 列数が64で割り切れない場合、末尾の不要ビットを0にするためのマスク
	var mask uint64
	if r == 0 {
		// 64で割り切れる場合、全てのビットが1のマスクを使用
		mask = ^uint64(0)
	} else {
	    // r = 2 のとき、uint(1) << r は 100 になり、-1をすると、011になる。
        // r = 5 のとき、uint(1) << r は 100000 になり、-1をすると、011111になる。
		mask = (uint64(1) << r) - 1
	}

	return Matrix{
		Rows:    rows,
		Cols:    cols,
		Stride:  stride,
		Data:    make([]uint64, rows*stride),
		RowMask: mask,
	}, nil
}

func (m Matrix) OnesCount() int {
	c := 0
	for _, v := range m.Data {
		c += bits.OnesCount64(v)
	}
	return c
}

func (m *Matrix) SetBit(r, c int) error {
	if r < 0 || r >= m.Rows {
		return fmt.Errorf("r が範囲外: r >= 0 && r < Rows であるべき")
	}

	if c < 0 || c >= m.Cols {
		return fmt.Errorf("c が範囲外: c >= 0 && c < Cols であるべき")
	}

	idx := (r * m.Stride) + (c / 64)
	shift := uint(c % 64)
	m.Data[idx] |= (1 << shift)
	return nil
}

func (m *Matrix) ApplyMask() {
    if m.RowMask == ^uint64(0) {
        return // マスク不要
    }
    
    for r := 0; r < m.Rows; r++ {
        // 各行の最後のuint64ブロック
        idx := (r * m.Stride) + (m.Stride - 1)
        m.Data[idx] &= m.RowMask
    }
}

func (m Matrix) MulVecAndPopCount(vec Matrix) ([]int, error) {
	if vec.Rows != 1 {
		return nil, fmt.Errorf("vec.Rows != 1: vec.Rows = 1 にするべき")
	}

	if m.Cols != vec.Cols {
		return nil, fmt.Errorf("m.Cols != vec.Cols: m.Cols = vec.Cols にするべき")
	}

	if m.Stride != vec.Stride {
		return nil, fmt.Errorf("m.Stride != vec.Stride: m.Stride = vec.Stride にするべき")
	}

	if m.RowMask != vec.RowMask {
		return nil, fmt.Errorf("m.RowMask != vec.RowMask: m.RowMask = vec.RowMask にするべき")
	}

	counts := make([]int, m.Rows)
	for r := 0; r < m.Rows; r++ {
		start := r * m.Stride
		popCount := 0

		for k := 0; k < m.Stride; k++ {
			a := vec.Data[k]
			b := m.Data[start+k]
			xnor := ^(a ^ b)

			// 最後のブロックのみマスク処理
			if k == m.Stride-1 {
				xnor &= m.RowMask
			}
			popCount += bits.OnesCount64(xnor)
		}
		counts[r] = popCount
	}
	return counts, nil
}

func (m Matrix) Transpose() (Matrix, error) {
    dst, err := NewMatrix(m.Cols, m.Rows)
    if err != nil {
        return Matrix{}, err
    }

    var (
        block [64]uint64
        mask  uint64
        t     uint64
        a, b  uint64
        other int
    )

    for r := 0; r < m.Rows; r += 64 {
        for cWord := 0; cWord < m.Stride; cWord++ {
            rowsToRead := 64
            if r+64 > m.Rows {
                rowsToRead = m.Rows - r
            }

            block = [64]uint64{}
            for i := 0; i < rowsToRead; i++ {
                srcIdx := (r + i) * m.Stride + cWord
                block[i] = m.Data[srcIdx]
            }

            // 32x32 swap
            mask = 0x00000000FFFFFFFF
            for i := 0; i < 32; i++ {
                other = i + 32
                a, b = block[i], block[other]
                // B(Top-Right)と C(Bottom-Left) を交換
                t = (b ^ (a >> 32)) & mask
                block[i] = a ^ (t << 32)
                block[other] = b ^ t
            }

            // 16x16 swap
            mask = 0x0000FFFF0000FFFF
            for j := 0; j < 64; j += 32 {
                for i := j; i < j+16; i++ {
                    other = i + 16
                    a, b = block[i], block[other]
                    t = (b ^ (a >> 16)) & mask
                    block[i] = a ^ (t << 16)
                    block[other] = b ^ t
                }
            }

            // 8x8 swap
            mask = 0x00FF00FF00FF00FF
            for j := 0; j < 64; j += 16 {
                for i := j; i < j+8; i++ {
                    other = i + 8
                    a, b = block[i], block[other]
                    t = (b ^ (a >> 8)) & mask
                    block[i] = a ^ (t << 8)
                    block[other] = b ^ t
                }
            }

            // 4x4 swap
            mask = 0x0F0F0F0F0F0F0F0F
            for j := 0; j < 64; j += 8 {
                for i := j; i < j+4; i++ {
                    other = i + 4
                    a, b = block[i], block[other]
                    t = (b ^ (a >> 4)) & mask
                    block[i] = a ^ (t << 4)
                    block[other] = b ^ t
                }
            }

            // 2x2 swap
            mask = 0x3333333333333333
            for j := 0; j < 64; j += 4 {
                for i := j; i < j+2; i++ {
                    other = i + 2
                    a, b = block[i], block[other]
                    t = (b ^ (a >> 2)) & mask
                    block[i] = a ^ (t << 2)
                    block[other] = b ^ t
                }
            }

            // 1x1 swap
            mask = 0x5555555555555555
            for j := 0; j < 64; j += 2 {
                other = j + 1
                a, b = block[j], block[other]
                t = (b ^ (a >> 1)) & mask
                block[j] = a ^ (t << 1)
                block[other] = b ^ t
            }

            dstRowBase := cWord * 64
            dstColWord := r / 64
            rowsToWrite := 64
            if dstRowBase+64 > dst.Rows {
                rowsToWrite = dst.Rows - dstRowBase
            }
            for i := 0; i < rowsToWrite; i++ {
                dstIdx := (dstRowBase + i) * dst.Stride + dstColWord
                dst.Data[dstIdx] = block[i]
            }
        }
    }

    dst.ApplyMask()
    return dst, nil
}