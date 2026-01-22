package bitsx

import (
	"fmt"
	"math/bits"
	"math/rand/v2"

	"github.com/sw965/omw/constraints"
	"golang.org/x/sys/cpu"
	"slices"
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

type Matrix struct {
	// カプセル化する？
	Rows       int
	Cols       int
	// カプセル化する？
	Stride     int // (Cols + 63) / 64 が基本的な値: 例 Cols = 100 の時、uint64 * 2 = 128bitを確保 (Stride = 2)
	Data       []uint64
	// カプセル化する？
	RowMask    uint64
}

// Zerosにする？
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

// NewRandMatrix は、指定されたバイアス係数 k に基づいて、ビット密度が調整されたランダム行列を生成します。
// k の値によって、各ビットが 1 になる確率 P(1) が以下のように指数関数的に変化します。
//
// k <= 0 の場合: P(1) = (1/2)^(|k|+1)
// k > 0  の場合: P(1) = 1 - (1/2)^(k+1)
//
// k による確率 P(1) の対応表:
//  k = -10 : 0.00049 (約 0.05%)
//  k = -9  : 0.00098 (約 0.10%)
//  k = -8  : 0.00195 (約 0.20%)
//  k = -7  : 0.00391 (約 0.39%)
//  k = -6  : 0.00781 (約 0.78%)
//  k = -5  : 0.01563 (約 1.56%)
//  k = -4  : 0.03125 (約 3.13%)
//  k = -3  : 0.06250 (    6.25%)
//  k = -2  : 0.12500 (   12.50%)
//  k = -1  : 0.25000 (   25.00%)
//  k =  0  : 0.50000 (   50.00%) -> 一様ランダム
//  k =  1  : 0.75000 (   75.00%)
//  k =  2  : 0.87500 (   87.50%)
//  k =  3  : 0.93750 (   93.75%)
//  k =  4  : 0.96875 (   96.88%)
//  k =  5  : 0.98438 (   98.44%)
//  k =  6  : 0.99219 (   99.22%)
//  k =  7  : 0.99609 (   99.61%)
//  k =  8  : 0.99805 (   99.80%)
//  k =  9  : 0.99902 (   99.90%)
//  k = 10  : 0.99951 (   99.95%)
func NewRandMatrix(rows, cols int, k int, rng *rand.Rand) (Matrix, error) {
	m, err := NewMatrix(rows, cols)
	if err != nil {
		return Matrix{}, err
	}

	for i := range m.Data {
		p := rng.Uint64()
		if k < 0 {
			// AND演算を繰り返し、確率1/2ずつ下げる
			iters := -k
			for range iters {
				p &= rng.Uint64()
			}
		} else if k > 0 {
			// OR演算を繰り返し、確率を1/2ずつ上げる
			iters := k
			for range iters {
				p |= rng.Uint64()
			}
		}
		m.Data[i] = p
	}

	m.ApplyMask()
	return m, nil
}

func (m Matrix) Clone() Matrix {
	data := slices.Clone(m.Data)
	m.Data = data
	return m
}

func (m Matrix) And(other Matrix) (Matrix, error) {
	if m.Rows != other.Rows || m.Cols != other.Cols {
		return Matrix{}, fmt.Errorf("dimension mismatch")
	}
	c := m.Clone()
	for i := range c.Data {
		c.Data[i] &= other.Data[i]
	}
	c.ApplyMask()
	return c, nil
}

func (m Matrix) IndexAndShift(r, c int) (int, uint, error) {
	if r < 0 || r >= m.Rows {
		return 0, 0, fmt.Errorf("row が範囲外: row = %d: row < 0 || row >= Rows(=%d) であるべき", r, m.Rows)
	}
	if c < 0 || c >= m.Cols {
		return 0, 0, fmt.Errorf("col が範囲外: col = %d:col >= 0 && col < Cols(=%d) であるべき", c, m.Cols)
	}

	// 2行 * 100列の行列m を例に、idxの計算式を解説する
	// 100列の情報を64ビットで表現するには、2つのuint64が必要
	// よってm.Stride = 2となる
	// m.Dataの中身は次の通り
	// Data[0] は 0行目の0～63列の情報
	// Data[1] は 0行目の64～99列の情報(100～127列はパディング)
	// Data[2] は 1行目の0～63列の情報
	// Data[3] は 1行目の64～99列の情報(100～127列はパディング)
	// ここで、1行目の70列目のビットを取り出す事を考える
	// r = 1, c = 70
	// r は行数を表すが、Dataは行数通りに並んでいないため、r * m.Strideで行数に変換する
	// 次に、cを列のインデックスに変換する方法を考える
	// cが0～63のとき、インデックス0、cが64～127のとき、インデックス1なので、
	// c / 64 で計算出来る。
	idx := (r * m.Stride) + (c / 64)

	// インデックスを特定したうえで、シフト演算などするための値
	// 例えば、70列目というのは、そのインデックスにおいては、先頭から6番目のビットに相当する
	// これは c % 64 で求めることができる
	shift := uint(c % 64)
	return idx, shift, nil
}

func (m Matrix) Bit(r, c int) (uint64, error) {
	idx, shift, err := m.IndexAndShift(r, c)
	if err != nil {
		return 0, err
	}
	// 例:「100100」の3番目のビットが欲しい場合、右に3回ずらして「0001001」にする
	// 1 (000001) と AND演算を行い、右端以外のビットを 0 にして消す
	// これで、n番目のビットの値 (0 or 1) を取得出来る
    return (m.Data[idx] >> shift) & 1, nil
}

func (m *Matrix) Set(r, c int) error {
	idx, shift, err := m.IndexAndShift(r, c)
	if err != nil {
		return err
	}
	m.Data[idx] |= (1 << shift)
	return nil
}

func (m *Matrix) Clear(r, c int) error {
    idx, shift, err := m.IndexAndShift(r, c)
    if err != nil {
        return err
    }
    m.Data[idx] &^= (1 << shift)
    return nil
}

func (m *Matrix) Toggle(r, c int) error {
    idx, shift, err := m.IndexAndShift(r, c)
    if err != nil {
        return err
    }
    
    // XOR (^) 演算を使って、特定位置のビットを反転させる
    // 0 ^ 1 = 1
    // 1 ^ 1 = 0
    m.Data[idx] ^= (1 << shift)
    return nil
}

func (m Matrix) PopCount() int {
    count := 0
    for r := 0; r < m.Rows; r++ {
        start := r * m.Stride
        for k := 0; k < m.Stride; k++ {
            word := m.Data[start+k]
            if k == m.Stride-1 {
                word &= m.RowMask
            }
            count += bits.OnesCount64(word)
        }
    }
    return count
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

func (m Matrix) MulVecAndPopCounts(vec Matrix) ([]int, error) {
	if vec.Rows != 1 {
		return nil, fmt.Errorf("vec.Rows != 1: vec.Rows = 1 にするべき")
	}

	if m.Cols != vec.Cols {
		return nil, fmt.Errorf("m.Cols != vec.Cols: m.Cols = %d, vec.Cols = %d: m.Cols = vec.Cols にするべき", m.Cols, vec.Cols)
	}

	if m.Stride != vec.Stride {
		return nil, fmt.Errorf("m.Stride != vec.Stride: m.Stride = %d, vec.Stride = %d: m.Stride = vec.Stride にするべき", m.Stride, vec.Stride)
	}

	if m.RowMask != vec.RowMask {
		return nil, fmt.Errorf("m.RowMask != vec.RowMask: m.RowMask = %d, vec.RowMask = %d: m.RowMask = vec.RowMask にするべき", m.RowMask, vec.RowMask)
	}

	counts := make([]int, m.Rows)
	for r := 0; r < m.Rows; r++ {
		start := r * m.Stride
		popCount := 0

		for k := 0; k < m.Stride; k++ {
			matWord := m.Data[start+k]
			vWord := vec.Data[k]
			xnor := ^(matWord ^ vWord)
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

func (m Matrix) MulVecAndPopCountsWithMask(vec, mask Matrix) ([]int, error) {
	if vec.Rows != 1 {
		return nil, fmt.Errorf("vec.Rows != 1: vec.Rows = 1 にするべき")
	}

	if m.Cols != vec.Cols {
		return nil, fmt.Errorf("m.Cols != vec.Cols: m.Cols = %d, vec.Cols = %d: m.Cols = vec.Cols にするべき", m.Cols, vec.Cols)
	}

	if m.Stride != vec.Stride {
		return nil, fmt.Errorf("m.Stride != vec.Stride: m.Stride = %d, vec.Stride = %d: m.Stride = vec.Stride にするべき", m.Stride, vec.Stride)
	}

	if m.RowMask != vec.RowMask {
		return nil, fmt.Errorf("m.RowMask != vec.RowMask: m.RowMask = %d, vec.RowMask = %d: m.RowMask = vec.RowMask にするべき", m.RowMask, vec.RowMask)
	}

	if mask.Rows != 1 {
		return nil, fmt.Errorf("後でエラーメッセージを書く")
	}

	if mask.Cols != m.Cols {
		return nil, fmt.Errorf("後でエラーメッセージを書く")
	}

	if mask.RowMask != m.RowMask {
		return nil, fmt.Errorf("後でエラーメッセージを書く")
	}

	counts := make([]int, m.Rows)
	for r := 0; r < m.Rows; r++ {
		start := r * m.Stride
		popCount := 0

		for k := 0; k < m.Stride; k++ {
			matWord := m.Data[start+k]
			vWord := vec.Data[k]
			maskWord := mask.Data[k]

			// 最後のブロックのみマスク処理
			if k == m.Stride-1 {
    			// maskWordが綺麗になれば、それを使用するvalidXnorも自動的に綺麗になる
   				maskWord &= m.RowMask 
			}

			xnor := ^(matWord ^ vWord)
			vaildXnor := xnor & maskWord
			popCount += bits.OnesCount64(vaildXnor)
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

//go:noescape
func mulVecPopCountsAVX512Asm(mat []uint64, vec []uint64, res []int, stride int, mask uint64)

func (m Matrix) MulVecAndPopCountsAVX512(vec Matrix) ([]int, error) {
	if vec.Rows != 1 {
		return nil, fmt.Errorf("vec.Rows != 1: vec.Rows = 1 にするべき")
	}

	if m.Cols != vec.Cols {
		return nil, fmt.Errorf("m.Cols != vec.Cols: m.Cols = %d, vec.Cols = %d: m.Cols = vec.Cols にするべき", m.Cols, vec.Cols)
	}

	if m.Stride != vec.Stride {
		return nil, fmt.Errorf("m.Stride != vec.Stride: m.Stride = %d, vec.Stride = %d: m.Stride = vec.Stride にするべき", m.Stride, vec.Stride)
	}

	if m.RowMask != vec.RowMask {
		return nil, fmt.Errorf("m.RowMask != vec.RowMask: m.RowMask = %d, vec.RowMask = %d: m.RowMask = vec.RowMask にするべき", m.RowMask, vec.RowMask)
	}

	if m.Rows == 0 {
        return nil, nil
    }

	counts := make([]int, m.Rows)

	// AVX-512 F (Foundation) と VPOPCNTDQ (Vector Popcount) の両方が必要
	if cpu.X86.HasAVX512F && cpu.X86.HasAVX512VPOPCNTDQ {
		mulVecPopCountsAVX512Asm(m.Data, vec.Data, counts, m.Stride, m.RowMask)
		return counts, nil
	}

	// フォールバック: AVX-512が使えない場合
	for r := 0; r < m.Rows; r++ {
		start := r * m.Stride
		popCount := 0
		for k := 0; k < m.Stride; k++ {
			matWord := m.Data[start+k]
			vWord := vec.Data[k]
			xnor := ^(matWord ^ vWord)
			if k == m.Stride-1 {
				xnor &= m.RowMask
			}
			popCount += bits.OnesCount64(xnor)
		}
		counts[r] = popCount
	}
	return counts, nil
}