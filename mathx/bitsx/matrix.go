package bitsx

import (
	"fmt"
	"math"
	"math/bits"
	"math/rand/v2"
	"slices"
)

type Matrix struct {
	rows        int
	cols        int
	stride      int // cols = 100 のとき、uint64 * 2 = 128bitを確保 (stride = 2)
	colTailMask uint64
	Data        []uint64
}

func NewZerosMatrix(rows, cols int) (*Matrix, error) {
	if rows <= 0 {
		return nil, fmt.Errorf("rows <= 0: rows > 0 であるべき")
	}

	if cols <= 0 {
		return nil, fmt.Errorf("cols <= 0: cols > 0 であるべき")
	}

	stride := (cols + 63) / 64

	r := cols % 64
	// colsが64で割り切れない場合、末尾の不要ビットを0にするためのマスク
	var colTailMask uint64
	if r == 0 {
		// colsが64で割り切れる場合、全てのビットが1のマスクを使用
		colTailMask = ^uint64(0)
	} else {
		// r = 2 のとき、uint(1) << r は 100 になる。100を-1で引くと、011になる。
		// r = 5 のとき、uint(1) << r は 100000 になり。100000を-1で引くと、011111になる。
		colTailMask = (uint64(1) << r) - 1
	}

	return &Matrix{
		rows:        rows,
		cols:        cols,
		stride:      stride,
		colTailMask: colTailMask,
		Data:        make([]uint64, rows*stride),
	}, nil
}

func NewOnesMatrix(rows, cols int) (*Matrix, error) {
	m, err := NewZerosMatrix(rows, cols)
	if err != nil {
		return nil, err
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
func NewRandMatrix(rows, cols int, k int, rng *rand.Rand) (*Matrix, error) {
	m, err := NewZerosMatrix(rows, cols)
	if err != nil {
		return nil, err
	}

	for i := range m.Data {
		word := rng.Uint64()
		if k < 0 {
			// AND演算を繰り返し、確率1/2ずつ下げる
			iters := -k
			for range iters {
				word &= rng.Uint64()
			}
		} else if k > 0 {
			// OR演算を繰り返し、確率を1/2ずつ上げる
			iters := k
			for range iters {
				word |= rng.Uint64()
			}
		}
		m.Data[i] = word
	}

	m.ApplyMask()
	return m, nil
}

func NewSignMatrix(rows, cols int, x []int) (*Matrix, error) {
	if len(x) < rows*cols {
		return nil, fmt.Errorf("len(x) が不足しています: %d < %d", len(x), rows*cols)
	}

	sign, err := NewZerosMatrix(rows, cols)
	if err != nil {
		return nil, err
	}

	err = sign.ScanRowsWord(nil, func(ctx MatrixWordContext) error {
		var signWord uint64
		xWord := x[ctx.GlobalStart:ctx.GlobalEnd]
		err := ctx.ScanBits(func(i, col, colT int) error {
			if xWord[i] >= 0 {
				signWord |= (uint64(1) << uint(i))
			}
			return nil
		})
		if err != nil {
			return err
		}
		sign.Data[ctx.WordIndex] = signWord
		return nil
	})

	if err != nil {
		return nil, err
	}
	return sign, nil
}

func (m *Matrix) Rows() int {
	return m.rows
}

func (m *Matrix) Cols() int {
	return m.cols
}

func (m *Matrix) Stride() int {
	return m.stride
}

func (m *Matrix) ColTailMask() uint64 {
	return m.colTailMask
}

func (m *Matrix) Clone() *Matrix {
	return &Matrix{
		rows:        m.rows,
		cols:        m.cols,
		stride:      m.stride,
		colTailMask: m.colTailMask,
		Data:        slices.Clone(m.Data),
	}
}

func (m *Matrix) And(other *Matrix) (*Matrix, error) {
	if m.rows != other.rows || m.cols != other.cols {
		return nil, fmt.Errorf("dimension mismatch")
	}
	c := m.Clone()
	for i := range c.Data {
		c.Data[i] &= other.Data[i]
	}
	c.ApplyMask()
	return c, nil
}

func (m *Matrix) Xor(other *Matrix) (*Matrix, error) {
	if m.rows != other.rows || m.cols != other.cols {
		return nil, fmt.Errorf("dimension mismatch: (%dx%d) vs (%dx%d)", m.rows, m.cols, other.rows, other.cols)
	}
	c := m.Clone()
	for i := range c.Data {
		c.Data[i] ^= other.Data[i]
	}
	c.ApplyMask()
	return c, nil
}

func (m *Matrix) HammingDistance(other *Matrix) (int, error) {
	// 異なるビットであれば1になる
	diff, err := m.Xor(other)
	if err != nil {
		return 0, err
	}
	return diff.OnesCount64(), nil
}

func (m *Matrix) IndexAndShift(r, c int) (int, uint, error) {
	if r < 0 || r >= m.rows {
		return 0, 0, fmt.Errorf("row が範囲外: row = %d: row < 0 || row >= Rows(=%d) であるべき", r, m.rows)
	}
	if c < 0 || c >= m.cols {
		return 0, 0, fmt.Errorf("col が範囲外: col = %d:col >= 0 && col < Cols(=%d) であるべき", c, m.cols)
	}

	// 2行 * 100列の行列m を例に、idxの計算式を解説する
	// 100列の情報を64ビットで表現するには、2つのuint64が必要
	// よってm.stride = 2となる
	// m.Dataの中身は次の通り
	// Data[0] は 0行目の0～63列の情報
	// Data[1] は 0行目の64～99列の情報(100～127列はパディング)
	// Data[2] は 1行目の0～63列の情報
	// Data[3] は 1行目の64～99列の情報(100～127列はパディング)
	// ここで、1行目の70列目のビットを取り出す事を考える
	// r = 1, c = 70
	// r は行数を表すが、Dataは行数通りに並んでいないため、r * m.strideで行数に変換する
	// 次に、cを列のインデックスに変換する方法を考える
	// cが0～63のとき、インデックス0、cが64～127のとき、インデックス1なので、
	// c / 64 で計算出来る。
	idx := (r * m.stride) + (c / 64)

	// インデックスを特定したうえで、シフト演算などするための値
	// 例えば、70列目というのは、そのインデックスにおいては、先頭から6番目のビットに相当する
	// これは c % 64 で求めることができる
	shift := uint(c % 64)
	return idx, shift, nil
}

func (m *Matrix) Bit(r, c int) (uint64, error) {
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

func (m *Matrix) OnesCount64() int {
	count := 0
	for r := 0; r < m.rows; r++ {
		start := r * m.stride
		for k := 0; k < m.stride; k++ {
			word := m.Data[start+k]
			if k == m.stride-1 {
				word &= m.colTailMask
			}
			count += bits.OnesCount64(word)
		}
	}
	return count
}

func (m *Matrix) ApplyMask() {
	if m.colTailMask == ^uint64(0) {
		return // マスク不要
	}

	for r := 0; r < m.rows; r++ {
		// 各行の最後のuint64ブロック
		idx := (r * m.stride) + (m.stride - 1)
		m.Data[idx] &= m.colTailMask
	}
}

func (m *Matrix) Dot(other *Matrix) ([]int, error) {
	if m.cols != other.cols {
		return nil, fmt.Errorf("dimension mismatch: m.cols %d != other.cols %d", m.cols, other.cols)
	}
	if m.stride != other.stride {
		return nil, fmt.Errorf("stride mismatch: m.stride %d != other.stride %d", m.stride, other.stride)
	}
	if m.colTailMask != other.colTailMask {
		return nil, fmt.Errorf("mask mismatch: m.colTailMask %x != other.colTailMask %x", m.colTailMask, other.colTailMask)
	}

	yRows := m.rows
	yCols := other.rows
	counts := make([]int, yRows*yCols)

	mData := m.Data
	oData := other.Data
	stride := m.stride
	mask := m.colTailMask

	for r := range yRows {
		mOffset := r * stride
		yOffset := r * yCols
		for c := range yCols {
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

func (m *Matrix) DotTernary(sign, nonZero *Matrix) ([]int, error) {
	if m.cols != sign.cols {
		return nil, fmt.Errorf("dimension mismatch: m.cols %d != otherSign.cols %d", m.cols, sign.cols)
	}

	if sign.rows != nonZero.rows || sign.cols != nonZero.cols {
		return nil, fmt.Errorf("otherSign and otherNonZero dimension mismatch")
	}

	if m.stride != sign.stride || sign.stride != nonZero.stride {
		return nil, fmt.Errorf("stride mismatch")
	}

	if m.colTailMask != sign.colTailMask || sign.colTailMask != nonZero.colTailMask {
		return nil, fmt.Errorf("mask mismatch")
	}

	zRows := m.rows
	zCols := sign.rows
	z := make([]int, zRows*zCols)

	mData := m.Data
	sData := sign.Data
	nzData := nonZero.Data
	stride := m.stride
	mask := m.colTailMask

	for r := 0; r < zRows; r++ {
		mOffset := r * stride
		zOffset := r * zCols
		for c := 0; c < zCols; c++ {
			ternaryOffset := c * stride
			matchCount := 0
			nonZeroCount := 0
			for k := 0; k < stride; k++ {
				mWord := mData[mOffset+k]
				sWord := sData[ternaryOffset+k]
				nzWord := nzData[ternaryOffset+k]

				// 符号が一致しているか (XNOR)
				sameSign := ^(mWord ^ sWord)

				// 有効(NonZero)かつ符号一致 (AND)
				validMatch := sameSign & nzWord

				// 最後のブロックのみマスク処理
				if k == stride-1 {
					validMatch &= mask
					nzWord &= mask
				}

				matchCount += bits.OnesCount64(validMatch)
				nonZeroCount += bits.OnesCount64(nzWord)
			}
			z[zOffset+c] = 2*matchCount - nonZeroCount
		}
	}
	return z, nil
}

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

func (m *Matrix) Transpose() (*Matrix, error) {
	dst, err := NewZerosMatrix(m.cols, m.rows)
	if err != nil {
		return nil, err
	}

	var block [64]uint64

	srcStride := m.stride
	dstStride := dst.stride
	srcData := m.Data
	dstData := dst.Data
	rows := m.rows

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
			if dstRowBase+64 > dst.rows {
				dstRowsToWrite = dst.rows - dstRowBase
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
	rows, cols, stride := m.rows, m.cols, m.stride
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
		rowBitOffset := r * cols
		for s := 0; s < stride; s++ {
			colStart := s << 6
			colEnd := colStart + 64
			mask := ^uint64(0)

			if colEnd > cols {
				colEnd = cols
				mask = m.colTailMask
			}

			err := f(MatrixWordContext{
				rows:        rows,
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

type Matrices []*Matrix

func NewETFMatrices(n, rows, cols int, iters int, rng *rand.Rand) (Matrices, error) {
	ms := make(Matrices, n)
	for i := range n {
		m, err := NewRandMatrix(rows, cols, 0, rng)
		if err != nil {
			return nil, err
		}
		ms[i] = m
	}

	currentCost, err := ms.ETFCost()
	if err != nil {
		return nil, err
	}

	for range iters {
		nIdx := rng.IntN(n)
		rIdx := rng.IntN(rows)
		cIdx := rng.IntN(cols)

		err := ms[nIdx].Toggle(rIdx, cIdx)
		if err != nil {
			return nil, err
		}

		cost, err := ms.ETFCost()
		if err != nil {
			return nil, err
		}

		if cost < currentCost {
			currentCost = cost
		} else {
			err := ms[nIdx].Toggle(rIdx, cIdx)
			if err != nil {
				return nil, err
			}
		}
	}
	return ms, nil
}

func NewRFFMatrices(n, rows, cols int, sigma float64, rng *rand.Rand) (Matrices, error) {
	if n < 2 {
		return nil, fmt.Errorf("n must be at least 2")
	}

	totalBits := rows * cols
	omegas := make([]float64, totalBits)
	phases := make([]float64, totalBits)

	for i := 0; i < totalBits; i++ {
		omegas[i] = rng.NormFloat64() * sigma
		phases[i] = rng.Float64() * 2 * math.Pi
	}

	ms := make(Matrices, n)
	for i := 0; i < n; i++ {
		y := float64(i) / float64(n-1)
		m, err := NewZerosMatrix(rows, cols)
		if err != nil {
			return nil, err
		}

		err = m.ScanRowsWord(nil, func(ctx MatrixWordContext) error {
			var mWord uint64
			omegaWord := omegas[ctx.GlobalStart:ctx.GlobalEnd]
			phaseWord := phases[ctx.GlobalStart:ctx.GlobalEnd]
			ctx.ScanBits(func(i, col, colT int) error {
				z := math.Cos(omegaWord[i]*y + phaseWord[i])
				if z >= 0 {
					mWord |= (1 << uint(i))
				}
				return nil
			})
			m.Data[ctx.WordIndex] = mWord
			return nil
		})

		if err != nil {
			return nil, err
		}
		ms[i] = m
	}
	return ms, nil
}

func NewThermometerMatrices(n, rows, cols int) (Matrices, error) {
	if n < 2 {
		return nil, fmt.Errorf("n must be at least 2")
	}

	ms := make(Matrices, n)
	totalBits := rows * cols

	for i := 0; i < n; i++ {
		m, err := NewZerosMatrix(rows, cols)
		if err != nil {
			return nil, err
		}

		numOnes := (i * totalBits) / (n - 1)
		err = m.ScanRowsWord(nil, func(ctx MatrixWordContext) error {
			var word uint64
			ctx.ScanBits(func(i, col, colT int) error {
				if (ctx.GlobalStart + i) < numOnes {
					word |= (uint64(1) << uint(i))
				}
				return nil
			})
			m.Data[ctx.WordIndex] = word
			return nil
		})

		if err != nil {
			return nil, err
		}
		ms[i] = m
	}
	return ms, nil
}

func (ms Matrices) ETFCost() (float32, error) {
	n := len(ms)
	distances := make([]float32, 0, n*n)
	sum := float32(0.0)
	for i := range len(ms) {
		for j := i + 1; j < len(ms); j++ {
			distance, err := ms[i].HammingDistance(ms[j])
			if err != nil {
				return 0.0, err
			}
			d := float32(distance)
			distances = append(distances, d)
			sum += d
		}
	}

	dn := len(distances)
	dnf := float32(dn)
	// 距離の平均
	mean := sum / dnf

	// 距離の分散
	variance := float32(0.0)
	for _, d := range distances {
		deviation := d - mean
		variance += deviation * deviation
	}
	variance /= dnf

	// コスト = -(合計距離) + (距離の分散)
	// 距離を最大化したいので、合計距離にはマイナスをつけて最小化問題にする
	cost := -sum + variance
	return cost, nil
}

type MatrixWordContext struct {
	rows        int
	Row         int
	WordIndex   int
	ColStart    int
	ColEnd      int
	GlobalStart int
	GlobalEnd   int
	Mask        uint64
}

func (ctx MatrixWordContext) ScanBits(f func(i, col, colT int) error) error {
	colT := (ctx.ColStart * ctx.rows) + ctx.Row
	for i := range ctx.ColEnd - ctx.ColStart {
		col := ctx.ColStart + i
		err := f(i, col, colT)
		if err != nil {
			return err
		}
		colT += ctx.rows
	}
	return nil
}
