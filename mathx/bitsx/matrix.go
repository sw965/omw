package bitsx

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"math/bits"
)

type Matrix struct {
	// カプセル化する？
	Rows int
	Cols int
	// カプセル化する？
	Stride int // (Cols + 63) / 64 が基本的な値: 例 Cols = 100 の時、uint64 * 2 = 128bitを確保 (Stride = 2)
	Data   []uint64
	// カプセル化する？
	WordMask uint64
}

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
		WordMask: mask,
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
//
//	k = -10 : 0.00049 (約 0.05%)
//	k = -9  : 0.00098 (約 0.10%)
//	k = -8  : 0.00195 (約 0.20%)
//	k = -7  : 0.00391 (約 0.39%)
//	k = -6  : 0.00781 (約 0.78%)
//	k = -5  : 0.01563 (約 1.56%)
//	k = -4  : 0.03125 (約 3.13%)
//	k = -3  : 0.06250 (    6.25%)
//	k = -2  : 0.12500 (   12.50%)
//	k = -1  : 0.25000 (   25.00%)
//	k =  0  : 0.50000 (   50.00%) -> 一様ランダム
//	k =  1  : 0.75000 (   75.00%)
//	k =  2  : 0.87500 (   87.50%)
//	k =  3  : 0.93750 (   93.75%)
//	k =  4  : 0.96875 (   96.88%)
//	k =  5  : 0.98438 (   98.44%)
//	k =  6  : 0.99219 (   99.22%)
//	k =  7  : 0.99609 (   99.61%)
//	k =  8  : 0.99805 (   99.80%)
//	k =  9  : 0.99902 (   99.90%)
//	k = 10  : 0.99951 (   99.95%)
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

func NewSignMatrix(rows, cols int, xs []int) (Matrix, error) {
    sign, err := NewZerosMatrix(rows, cols)
    if err != nil {
        return Matrix{}, err
    }

    stride := sign.Stride
    for r := 0; r < rows; r++ {
        zRowOffset := r * cols
        dataRowOffset := r * stride
        
        for k := 0; k < stride; k++ {
            // このブロックの開始位置を事前に計算
            xBaseIdx := zRowOffset + (k << 6) // k * 64
            
            validBits := 64
            if (k<<6 + 64) > cols {
                validBits = cols - (k << 6)
            }

            var word uint64
            for b := 0; b < validBits; b++ {
                // スライスへの連続アクセス
                if xs[xBaseIdx + b] >= 0 {
                    word |= (uint64(1) << uint(b))
                }
            }
            sign.Data[dataRowOffset + k] = word
        }
    }
    return sign, nil
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
				word &= m.WordMask
			}
			count += bits.OnesCount64(word)
		}
	}
	return count
}

func (m *Matrix) ApplyMask() {
	if m.WordMask == ^uint64(0) {
		return // マスク不要
	}

	for r := 0; r < m.Rows; r++ {
		// 各行の最後のuint64ブロック
		idx := (r * m.Stride) + (m.Stride - 1)
		m.Data[idx] &= m.WordMask
	}
}

func (m Matrix) Dot(other Matrix) ([]int, error) {
	if m.Cols != other.Cols {
		return nil, fmt.Errorf("dimension mismatch: m.Cols %d != other.Cols %d", m.Cols, other.Cols)
	}
	if m.Stride != other.Stride {
		return nil, fmt.Errorf("stride mismatch: m.Stride %d != other.Stride %d", m.Stride, other.Stride)
	}
	if m.WordMask != other.WordMask {
		return nil, fmt.Errorf("mask mismatch: m.WordMask %x != other.WordMask %x", m.WordMask, other.WordMask)
	}

	outRows := m.Rows
	outCols := other.Rows
	counts := make([]int, outRows*outCols)

	mData := m.Data
	oData := other.Data
	stride := m.Stride
	mask := m.WordMask

	for r := range outRows {
		mOffset := r * stride
		yOffset := r * outCols
		for c := range outCols {
			oOffset := c * stride
			count := 0
			for k := range stride {
				mWord := mData[mOffset+k]
				oWord := oData[oOffset+k]

				xnor := ^(mWord ^ oWord)
				if k == stride-1 {
					xnor &= mask
				}
				count += bits.OnesCount64(xnor)
			}
			counts[yOffset+c] = count
		}
	}
	return counts, nil
}

func (m Matrix) DotTernary(otherSign, otherNonZero Matrix) ([]int, error) {
	// 1. 次元の整合性チェック
	if m.Cols != otherSign.Cols {
		return nil, fmt.Errorf("dimension mismatch: m.Cols %d != otherSign.Cols %d", m.Cols, otherSign.Cols)
	}
	if otherSign.Cols != otherNonZero.Cols || otherSign.Rows != otherNonZero.Rows {
		return nil, fmt.Errorf("otherSign and otherNonZero dimension mismatch")
	}
	// ストライドとマスクのチェック (高速化のため一致を前提とする)
	if m.Stride != otherSign.Stride || otherSign.Stride != otherNonZero.Stride {
		return nil, fmt.Errorf("stride mismatch")
	}
	if m.WordMask != otherSign.WordMask || otherSign.WordMask != otherNonZero.WordMask {
		return nil, fmt.Errorf("mask mismatch")
	}

	outRows := m.Rows
	outCols := otherSign.Rows // otherは転置されている前提
	zScores := make([]int, outRows*outCols)

	// ポインタと定数のキャッシュ
	mData := m.Data
	sData := otherSign.Data
	nData := otherNonZero.Data
	stride := m.Stride
	mask := m.WordMask

	// 2. 行列積のループ
	for r := 0; r < outRows; r++ {
		mRowOffset := r * stride
		resRowOffset := r * outCols

		for c := 0; c < outCols; c++ {
			oRowOffset := c * stride

			matchCount := 0
			activeCount := 0

			for k := 0; k < stride; k++ {
				mWord := mData[mRowOffset+k]
				sWord := sData[oRowOffset+k]
				nWord := nData[oRowOffset+k]

				// A. 符号が一致しているか (XNOR)
				sameSign := ^(mWord ^ sWord)

				// B. 有効(NonZero)かつ符号一致 (AND)
				validMatch := sameSign & nWord

				// 最後のブロックのみマスク処理
				if k == stride-1 {
					validMatch &= mask
					nWord &= mask
				}

				matchCount += bits.OnesCount64(validMatch)
				activeCount += bits.OnesCount64(nWord)
			}

			// Zスコアの計算
			// Z = Match - Mismatch
			//   = Match - (Active - Match)
			//   = 2 * Match - Active
			zScores[resRowOffset+c] = 2*matchCount - activeCount
		}
	}

	return zScores, nil
}

// func (m Matrix) Transpose() (Matrix, error) {
// 	dst, err := NewZerosMatrix(m.Cols, m.Rows)
// 	if err != nil {
// 		return Matrix{}, err
// 	}

// 	var (
// 		block [64]uint64
// 		mask  uint64
// 		t     uint64
// 		a, b  uint64
// 		other int
// 	)

// 	for r := 0; r < m.Rows; r += 64 {
// 		for cWord := 0; cWord < m.Stride; cWord++ {
// 			rowsToRead := 64
// 			if r+64 > m.Rows {
// 				rowsToRead = m.Rows - r
// 			}

// 			block = [64]uint64{}
// 			for i := 0; i < rowsToRead; i++ {
// 				srcIdx := (r+i)*m.Stride + cWord
// 				block[i] = m.Data[srcIdx]
// 			}

// 			// 32x32 swap
// 			mask = 0x00000000FFFFFFFF
// 			for i := 0; i < 32; i++ {
// 				other = i + 32
// 				a, b = block[i], block[other]
// 				// B(Top-Right)と C(Bottom-Left) を交換
// 				t = (b ^ (a >> 32)) & mask
// 				block[i] = a ^ (t << 32)
// 				block[other] = b ^ t
// 			}

// 			// 16x16 swap
// 			mask = 0x0000FFFF0000FFFF
// 			for j := 0; j < 64; j += 32 {
// 				for i := j; i < j+16; i++ {
// 					other = i + 16
// 					a, b = block[i], block[other]
// 					t = (b ^ (a >> 16)) & mask
// 					block[i] = a ^ (t << 16)
// 					block[other] = b ^ t
// 				}
// 			}

// 			// 8x8 swap
// 			mask = 0x00FF00FF00FF00FF
// 			for j := 0; j < 64; j += 16 {
// 				for i := j; i < j+8; i++ {
// 					other = i + 8
// 					a, b = block[i], block[other]
// 					t = (b ^ (a >> 8)) & mask
// 					block[i] = a ^ (t << 8)
// 					block[other] = b ^ t
// 				}
// 			}

// 			// 4x4 swap
// 			mask = 0x0F0F0F0F0F0F0F0F
// 			for j := 0; j < 64; j += 8 {
// 				for i := j; i < j+4; i++ {
// 					other = i + 4
// 					a, b = block[i], block[other]
// 					t = (b ^ (a >> 4)) & mask
// 					block[i] = a ^ (t << 4)
// 					block[other] = b ^ t
// 				}
// 			}

// 			// 2x2 swap
// 			mask = 0x3333333333333333
// 			for j := 0; j < 64; j += 4 {
// 				for i := j; i < j+2; i++ {
// 					other = i + 2
// 					a, b = block[i], block[other]
// 					t = (b ^ (a >> 2)) & mask
// 					block[i] = a ^ (t << 2)
// 					block[other] = b ^ t
// 				}
// 			}

// 			// 1x1 swap
// 			mask = 0x5555555555555555
// 			for j := 0; j < 64; j += 2 {
// 				other = j + 1
// 				a, b = block[j], block[other]
// 				t = (b ^ (a >> 1)) & mask
// 				block[j] = a ^ (t << 1)
// 				block[other] = b ^ t
// 			}

// 			dstRowBase := cWord * 64
// 			dstColWord := r / 64
// 			rowsToWrite := 64
// 			if dstRowBase+64 > dst.Rows {
// 				rowsToWrite = dst.Rows - dstRowBase
// 			}
// 			for i := 0; i < rowsToWrite; i++ {
// 				dstIdx := (dstRowBase+i)*dst.Stride + dstColWord
// 				dst.Data[dstIdx] = block[i]
// 			}
// 		}
// 	}

// 	dst.ApplyMask()
// 	return dst, nil
// }

func transpose64Block(block *[64]uint64) {
	var (
		mask uint64
		t    uint64
		a, b uint64
	)

	// 32x32 swap
	mask = 0x00000000FFFFFFFF
	for j := 0; j < 32; j++ {
		a, b = block[j], block[j+32]
		t = (b ^ (a >> 32)) & mask
		block[j] = a ^ (t << 32)
		block[j+32] = b ^ t
	}

	// 16x16 swap
	mask = 0x0000FFFF0000FFFF
	for j := 0; j < 64; j += 32 {
		for i := j; i < j+16; i++ {
			a, b = block[i], block[i+16]
			t = (b ^ (a >> 16)) & mask
			block[i] = a ^ (t << 16)
			block[i+16] = b ^ t
		}
	}

	// 8x8 swap
	mask = 0x00FF00FF00FF00FF
	for j := 0; j < 64; j += 16 {
		for i := j; i < j+8; i++ {
			a, b = block[i], block[i+8]
			t = (b ^ (a >> 8)) & mask
			block[i] = a ^ (t << 8)
			block[i+8] = b ^ t
		}
	}

	// 4x4 swap
	mask = 0x0F0F0F0F0F0F0F0F
	for j := 0; j < 64; j += 8 {
		for i := j; i < j+4; i++ {
			a, b = block[i], block[i+4]
			t = (b ^ (a >> 4)) & mask
			block[i] = a ^ (t << 4)
			block[i+4] = b ^ t
		}
	}

	// 2x2 swap
	mask = 0x3333333333333333
	for j := 0; j < 64; j += 4 {
		for i := j; i < j+2; i++ {
			a, b = block[i], block[i+2]
			t = (b ^ (a >> 2)) & mask
			block[i] = a ^ (t << 2)
			block[i+2] = b ^ t
		}
	}

	// 1x1 swap
	mask = 0x5555555555555555
	for j := 0; j < 64; j += 2 {
		a, b = block[j], block[j+1]
		t = (b ^ (a >> 1)) & mask
		block[j] = a ^ (t << 1)
		block[j+1] = b ^ t
	}
}

func (m Matrix) Transpose() (Matrix, error) {
	dst, err := NewZerosMatrix(m.Cols, m.Rows)
	if err != nil {
		return Matrix{}, err
	}

	var block [64]uint64
	
	srcStride := m.Stride
	dstStride := dst.Stride
	srcData := m.Data
	dstData := dst.Data
	rows := m.Rows
	
	// ブロック単位での処理 (64行ずつ)
	for r := 0; r < rows; r += 64 {
		// 残り行数が64未満かどうか
		remainingRows := rows - r
		isFullBlock := remainingRows >= 64
		rowsToProcess := 64
		if !isFullBlock {
			rowsToProcess = remainingRows
		}

		// 横方向（Word単位）のループ
		for cWord := 0; cWord < srcStride; cWord++ {
			// 1. 読み込み (Read)
			// Optimize: インデックス計算の乗算を避けるため、ベースオフセットを計算
			srcBaseIdx := r*srcStride + cWord
			
			if isFullBlock {
				// ホットパス: 分岐なしで64回読み込む
				// コンパイラによるBounds Check Eliminationが効きやすくなる
				for i := 0; i < 64; i++ {
					block[i] = srcData[srcBaseIdx]
					srcBaseIdx += srcStride
				}
			} else {
				// エッジケース: 慎重に読み込む
				for i := 0; i < rowsToProcess; i++ {
					block[i] = srcData[srcBaseIdx]
					srcBaseIdx += srcStride
				}
				// 足りない部分は0埋め（ゴミデータが混ざらないように）
				for i := rowsToProcess; i < 64; i++ {
					block[i] = 0
				}
			}

			// 2. CPU内転置 (Process)
			transpose64Block(&block)

			// 3. 書き込み (Write)
			// 転置後は、dstの「cWord行目」の「r列が含まれるブロック」に書き込まれる
			// dstの行インデックス: cWord * 64 + (0..63)
			// dstの列ワードインデックス: r / 64
			
			dstRowBase := cWord * 64
			dstColWord := r / 64 // rは常に64の倍数なので単純なシフト
			
			// 書き込み先の行数チェック
			dstRowsToWrite := 64
			if dstRowBase+64 > dst.Rows {
				dstRowsToWrite = dst.Rows - dstRowBase
			}

			dstBaseIdx := dstRowBase*dstStride + dstColWord

			if dstRowsToWrite == 64 {
				// ホットパス
				for i := 0; i < 64; i++ {
					dstData[dstBaseIdx] = block[i]
					dstBaseIdx += dstStride
				}
			} else {
				// エッジケース
				for i := 0; i < dstRowsToWrite; i++ {
					dstData[dstBaseIdx] = block[i]
					dstBaseIdx += dstStride
				}
			}
		}
	}

	dst.ApplyMask()
	return dst, nil
}

func (m *Matrix) ScanRowsWord(rowIdxs []int, f func(ctx MatrixWordContext) error) error {
	rows, cols, stride := m.Rows, m.Cols, m.Stride
	wordMask := m.WordMask

	if rowIdxs == nil {
		rowIdxs = make([]int, rows)
		for i := range rows {
			rowIdxs[i] = i
		}
	}

	for _, r := range rowIdxs {
		if r < 0 || r >= rows {
			return fmt.Errorf("後でエラーメッセージを書く")
		}

		rowWordOffset := r * stride
		rowBitOffset  := r * cols
		for s := 0; s < stride; s++ {
			colStart := s << 6
			colEnd   := colStart + 64
			mask       := ^uint64(0)

			if colEnd > cols {
				colEnd = cols
				mask = wordMask
			}

			err := f(MatrixWordContext{
				rows:rows,
				Row:         r,
				WordIndex:   rowWordOffset + s,
				ColStart:    colStart,
				ColEnd:      colEnd,
				GlobalStart: rowBitOffset + colStart,
				GlobalEnd:   rowBitOffset + colEnd,
				Mask:        mask,
			})

			if err != nil {
				return err
			}
		}
	}
	return nil
}

type Matrices []Matrix

func NewBEFPrototypeMatrices(n, rows, cols int, iters int, rng *rand.Rand) (Matrices, error) {
	protos := make(Matrices, n)
	for i := range n {
		m, err := NewRandMatrix(rows, cols, 0, rng)
		if err != nil {
			return nil, err
		}
		protos[i] = m
	}

	currentCost, err := protos.CalculateBEFCost()
	if err != nil {
		return nil, err
	}

	for range iters {
		nIdx := rng.IntN(n)
		rIdx := rng.IntN(rows)
		cIdx := rng.IntN(cols)

		protos[nIdx].Toggle(rIdx, cIdx)
		cost, err := protos.CalculateBEFCost()
		if err != nil {
			return nil, err
		}

		if cost < currentCost {
			currentCost = cost
		} else {
			protos[nIdx].Toggle(rIdx, cIdx)
		}
	}
	return protos, nil
}

func (ms Matrices) CalculateBEFCost() (float64, error) {
	n := len(ms)
	distances := make([]float64, 0, n*n)
	sum := 0.0

	for i := range len(ms) {
		for j := i + 1; j < len(ms); j++ {
			d, err := ms[i].HammingDistance(ms[j])
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

type MatrixWordContext struct {
	rows        int
	Row         int
	WordIndex   int
	ColStart  int
	ColEnd    int
	GlobalStart int
	GlobalEnd   int
	Mask        uint64
}

func (ctx MatrixWordContext) ScanBits(f func(i, col, colT int)) {
	colT := (ctx.ColStart * ctx.rows) + ctx.Row
	for i := range ctx.ColEnd-ctx.ColStart {
		col := ctx.ColStart + i
		f(i, col, colT)
		colT += ctx.rows
	}
}