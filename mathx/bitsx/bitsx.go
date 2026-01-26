package bitsx

import (
	"fmt"
	"math/bits"
	"math/rand/v2"

	"github.com/sw965/omw/constraints"
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
func NewZerosMatrix(rows, cols int) (Matrix, error) {
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

func NewOnesMatrix(rows, cols int) (Matrix, error) {
	m, err := NewZerosMatrix(rows, cols)
	if err != nil {
		return Matrix{}, err
	}

	for i := range m.Data {
		m.Data[i] = ^uint64(0)
	}

	m.ApplyMask()
	return m, nil
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
	m, err := NewZerosMatrix(rows, cols)
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

func (m Matrix) Xor(other Matrix) (Matrix, error) {
	if m.Rows != other.Rows || m.Cols != other.Cols {
		return Matrix{}, fmt.Errorf("dimension mismatch: (%dx%d) vs (%dx%d)", m.Rows, m.Cols, other.Rows, other.Cols)
	}
	c := m.Clone()
	for i := range c.Data {
		c.Data[i] ^= other.Data[i]
	}
	c.ApplyMask()
	return c, nil
}

func (m Matrix) HammingDistance(other Matrix) (int, error) {
	// 異なるビットであれば1になる
	diff, err := m.Xor(other)
	if err != nil {
		return 0, err
	}
	return diff.PopCount(), nil
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

func (m Matrix) Dot(other Matrix) ([]int, error) {
	if m.Cols != other.Cols {
		return nil, fmt.Errorf("dimension mismatch: m.Cols %d != other.Cols %d", m.Cols, other.Cols)
	}
	if m.Stride != other.Stride {
		return nil, fmt.Errorf("stride mismatch: m.Stride %d != other.Stride %d", m.Stride, other.Stride)
	}
	if m.RowMask != other.RowMask {
		return nil, fmt.Errorf("mask mismatch: m.RowMask %x != other.RowMask %x", m.RowMask, other.RowMask)
	}

	// 出力は (m.Rows x other.Rows) の行列になる
	outRows := m.Rows
	outCols := other.Rows
	results := make([]int, outRows*outCols)

	// 計算の高速化のため、ループ内でスライスアクセスの境界チェックを減らす工夫
	mData := m.Data
	oData := other.Data
	stride := m.Stride
	mask := m.RowMask

	// mの各行について (重みのフィルタ1つ分)
	for r := 0; r < outRows; r++ {
		mRowBase := r * stride
		resRowBase := r * outCols

		// otherの各行について (画像のパッチ1つ分)
		for c := 0; c < outCols; c++ {
			oRowBase := c * stride
			popCount := 0

			for k := 0; k < stride; k++ {
				mWord := mData[mRowBase+k]
				oWord := oData[oRowBase+k]

				xnor := ^(mWord ^ oWord)

				// 最後のブロックのみマスク処理
				if k == stride-1 {
					xnor &= mask
				}
				popCount += bits.OnesCount64(xnor)
			}
			
			// 結果を格納 (i行j列)
			results[resRowBase+c] = popCount
		}
	}

	return results, nil
}

func (m Matrix) DotVec(vec Matrix) ([]int, error) {
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

// Matrix(2値)を1を1, 0を-1と置き換えた時、Ternaryで0の部分を数えないようにする
// 後でもっとわかりやすいコメントを書く。
func (m Matrix) DotTernaryVec(sign, nonZero Matrix) ([]int, error) {
	if sign.Rows != 1 {
		return nil, fmt.Errorf("sign.Rows != 1: sign.Rows = 1 にするべき")
	}

	if m.Cols != sign.Cols {
		return nil, fmt.Errorf("m.Cols != sign.Cols: m.Cols = %d, sign.Cols = %d: m.Cols = sign.Cols にするべき", m.Cols, sign.Cols)
	}

	if m.Stride != sign.Stride {
		return nil, fmt.Errorf("m.Stride != sign.Stride: m.Stride = %d, sign.Stride = %d: m.Stride = sign.Stride にするべき", m.Stride, sign.Stride)
	}

	if m.RowMask != sign.RowMask {
		return nil, fmt.Errorf("m.RowMask != sign.RowMask: m.RowMask = %d, sign.RowMask = %d: m.RowMask = sign.RowMask にするべき", m.RowMask, sign.RowMask)
	}

	if nonZero.Rows != 1 {
		return nil, fmt.Errorf("後でエラーメッセージを書く")
	}

	if nonZero.Cols != m.Cols {
		return nil, fmt.Errorf("後でエラーメッセージを書く")
	}

	if nonZero.RowMask != m.RowMask {
		return nil, fmt.Errorf("後でエラーメッセージを書く")
	}

	counts := make([]int, m.Rows)
	for r := 0; r < m.Rows; r++ {
		start := r * m.Stride
		popCount := 0

		for k := 0; k < m.Stride; k++ {
			matWord := m.Data[start+k]
			vWord := sign.Data[k]
			nonZeroWord := nonZero.Data[k]

			// 最後のブロックのみマスク処理
			if k == m.Stride-1 {
    			// nonZeroWordが綺麗になれば、それを使用するvalidXnorも自動的に綺麗になる
   				nonZeroWord &= m.RowMask 
			}

			xnor := ^(matWord ^ vWord)
			vaildXnor := xnor & nonZeroWord
			popCount += bits.OnesCount64(vaildXnor)
		}
		counts[r] = popCount
	}
	return counts, nil
}

func (m Matrix) Transpose() (Matrix, error) {
    dst, err := NewZerosMatrix(m.Cols, m.Rows)
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

func NewBEFPrototypeMatrices(n, dim int, iters int, rng *rand.Rand) ([]Matrix, error) {
	protos := make([]Matrix, n)
	for i := range n {
		// 一様分布
		m, err := NewRandMatrix(1, dim, 0, rng)
		if err != nil {
			return nil, err
		}
		protos[i] = m
	}

	currentCost, err := CalculateBEFCost(protos)
	if err != nil {
		return nil, err
	}

	for range iters {
		nIdx := rng.IntN(n)
		dIdx := rng.IntN(dim)

		protos[nIdx].Toggle(0, dIdx)
		cost, err := CalculateBEFCost(protos)
		if err != nil {
			return nil, err
		}

		if cost < currentCost {
			currentCost = cost
		} else {
			// コストが増えたら元に戻す
			protos[nIdx].Toggle(0, dIdx)
		}
	}
	return protos, nil
}

func CalculateBEFCost(protos []Matrix) (float64, error) {
	n := len(protos)
    distances := make([]float64, 0, n*n)
    sum := 0.0

    for i := range len(protos) {
        for j := i + 1; j < len(protos); j++ {
            d, err := protos[i].HammingDistance(protos[j])
			if err != nil {
				return 0.0, err
			}

            df := float64(d)
            distances = append(distances, df)
            sum += df
        }
    }

	dn := len(distances)
	// 距離の平均
    mean := sum / float64(dn)

	// 距離の分散
    variance := 0.0
    for _, d := range distances {
        diff := d - mean
        variance += diff * diff
    }
    variance /= float64(dn)

    // コスト = -(合計距離) + (距離の分散)
    // 距離を最大化したいので、合計距離にはマイナスをつけて最小化問題にする
    return -sum + variance, nil
}

type Tensor3 struct {
    Depths int // Z軸
    Rows   int // Y軸
    Cols   int // X軸
    
    Stride int // 行ごとのブロック数 (Cols由来)
    LayerStride int // 層ごとのブロック数 (Rows * Stride)
    
    Data   []uint64
    RowMask uint64
}

// Tensor3 のコンストラクタ（未定義だったため追加）
func NewTensor3(depths, rows, cols int) (Tensor3, error) {
	if depths <= 0 || rows <= 0 || cols <= 0 {
		return Tensor3{}, fmt.Errorf("dimensions must be positive: (%d, %d, %d)", depths, rows, cols)
	}

	stride := (cols + 63) / 64
	layerStride := rows * stride
	totalWords := depths * layerStride

	r := cols % 64
	var mask uint64
	if r == 0 {
		mask = ^uint64(0)
	} else {
		mask = (uint64(1) << r) - 1
	}

	return Tensor3{
		Depths:      depths,
		Rows:        rows,
		Cols:        cols,
		Stride:      stride,
		LayerStride: layerStride,
		Data:        make([]uint64, totalWords),
		RowMask:     mask,
	}, nil
}

// SetBit for Tensor3 (ヘルパー用)
func (t *Tensor3) Set(d, r, c int) {
	idx := d*t.LayerStride + r*t.Stride + (c / 64)
	shift := uint(c % 64)
	t.Data[idx] |= (1 << shift)
}

// copyBits は src のビット列を n ビット分、 dst の指定位置にコピーします。
// srcOffset, dstOffset はそれぞれの uint64 スライス先頭からのビットオフセットです。
// n は 64 以下であることを想定しています。
// 高速化のため、範囲チェック等は呼び出し元で保証してください。
func copyBits(dst []uint64, dstOffset int, src []uint64, srcOffset int, n int) {
	if n <= 0 {
		return
	}

	// 1. ソースから n ビット抽出して uint64 に詰める (LSB寄せ)
	srcWordIdx := srcOffset / 64
	srcBitShift := uint(srcOffset % 64)

	// 必要なビットを取り出すためのマスク
	// n=64のときシフト量が64になるとGoではパニックもしくは0になる挙動への対策が必要だが、
	// カーネルサイズが64を超えるケースは稀なため、ここでは一般的な最適化を優先。
	// 厳密には ^uint64(0) >> (64 - n) ですが、n=64対応を入れるなら分岐が必要。
	// 今回は n < 64 前提、あるいは挙動を理解した上での実装とします。
	valMask := ^uint64(0) >> (64 - n)
	if n == 64 {
		valMask = ^uint64(0)
	}

	var val uint64
	// ビット列がワードを跨ぐかどうかで処理を分ける
	if srcBitShift+uint(n) <= 64 {
		// 1ワード内に収まる場合
		val = (src[srcWordIdx] >> srcBitShift) & valMask
	} else {
		// 2ワードに跨る場合
		// 前半部分
		val = src[srcWordIdx] >> srcBitShift
		// 後半部分を上位ビットに結合
		remaining := 64 - srcBitShift
		// srcWordIdx+1 が範囲外になる可能性は呼び出し元で排除済みとする
		val |= (src[srcWordIdx+1] << remaining)
		val &= valMask
	}

	// 2. 抽出した val を dst の位置に合わせて書き込む
	dstWordIdx := dstOffset / 64
	dstBitShift := uint(dstOffset % 64)

	// 書き込み先がワードを跨ぐか
	if dstBitShift+uint(n) <= 64 {
		dst[dstWordIdx] |= (val << dstBitShift)
	} else {
		dst[dstWordIdx] |= (val << dstBitShift)
		remaining := 64 - dstBitShift
		dst[dstWordIdx+1] |= (val >> remaining)
	}
}

// Im2Col は3次元の2値テンソルを、畳み込み演算用に2次元行列へ展開します。
// カーネルの高さ(kH)と幅(kW)を個別に指定可能です。
// 高速化のため、ビット単位ではなくビットブロック単位でコピーを行います。
func (t *Tensor3) Im2Col(kH, kW, stride, padding int) (Matrix, error) {
	// 出力サイズの計算 (縦はkH, 横はkWを使用)
	outH := (t.Rows + 2*padding - kH) / stride + 1
	outW := (t.Cols + 2*padding - kW) / stride + 1

	if outH <= 0 || outW <= 0 {
		return Matrix{}, fmt.Errorf("output dimensions non-positive: check kernel size, stride, or padding")
	}

	numPatches := outH * outW
	// パッチサイズは 深さ * 高さ * 幅
	patchVectorSize := t.Depths * kH * kW

	res, err := NewZerosMatrix(numPatches, patchVectorSize)
	if err != nil {
		return Matrix{}, err
	}

	// ローカル変数キャッシュ（ポインタ参照を減らす）
	tData := t.Data
	resData := res.Data
	resStrideBits := res.Stride * 64 // 行あたりのビット数
	tStrideBits := t.Stride * 64     // 入力画像 1行あたりのビット数
	tLayerStrideBits := t.LayerStride * 64 // 入力画像 1層(Depth)あたりのビット数

	patchIdx := 0
	for oh := 0; oh < outH; oh++ {
		imgYStart := oh*stride - padding

		for ow := 0; ow < outW; ow++ {
			imgXStart := ow*stride - padding

			// 書き込み先のビットオフセットの初期値 (Matrixの該当行の先頭)
			dstBitBase := patchIdx * resStrideBits
			dstBitOffset := dstBitBase

			for d := 0; d < t.Depths; d++ {
				dOffsetBits := d * tLayerStrideBits

				// カーネルの縦ループ
				for ky := 0; ky < kH; ky++ {
					imgY := imgYStart + ky

					// Y方向のパディング判定
					// 画像範囲外の行であれば、カーネル幅(kW)分だけ書き込み位置を進める(0埋め)
					if imgY < 0 || imgY >= t.Rows {
						dstBitOffset += kW
						continue
					}

					yOffsetBits := dOffsetBits + imgY*tStrideBits

					// --- X方向の最適化ロジック ---
					// ループを使わず、有効範囲を一括コピーする

					// カーネルのX座標 [0, kW) のうち、画像範囲内に重なる部分 [kxStart, kxEnd) を求める
					kxStart := 0
					if imgXStart < 0 {
						kxStart = -imgXStart
					}

					kxEnd := kW
					if imgXStart+kW > t.Cols {
						kxEnd = t.Cols - imgXStart
					}

					// 1. 左側のパディング (範囲外)
					// imgXStart < 0 の場合、その分だけ書き込み位置をスキップ
					if kxStart > 0 {
						dstBitOffset += kxStart
					}

					// 2. 有効データのコピー
					// kxStart から kxEnd までのビットをコピー
					copyLen := kxEnd - kxStart
					if copyLen > 0 {
						// コピー元の絶対ビット位置
						srcBitPos := yOffsetBits + (imgXStart + kxStart)

						// ビット一括転送
						copyBits(resData, dstBitOffset, tData, srcBitPos, copyLen)

						dstBitOffset += copyLen
					}

					// 3. 右側のパディング (範囲外)
					// カーネルの残り幅分をスキップ
					if kxEnd < kW {
						dstBitOffset += (kW - kxEnd)
					}
				}
			}
			patchIdx++
		}
	}

	// 念のためマスク適用
	res.ApplyMask()
	return res, nil
}