package bitsx

import (
	"fmt"
	"math"
	"math/bits"
	"math/rand/v2"
	"slices"
)

type Matrix struct {
	Rows int
	Cols int
	Data []uint64
}

func NewZerosMatrix(rows, cols int) (*Matrix, error) {
	if rows <= 0 {
		return nil, fmt.Errorf("rows <= 0: rows > 0 であるべき")
	}

	if cols <= 0 {
		return nil, fmt.Errorf("cols <= 0: cols > 0 であるべき")
	}

	m := &Matrix{
		Rows: rows,
		Cols: cols,
	}
	stride := m.Stride()
	m.Data = make([]uint64, rows*stride)
	return m, nil
}

func NewOnesMatrix(rows, cols int) (*Matrix, error) {
	m, err := NewZerosMatrix(rows, cols)
	if err != nil {
		return nil, err
	}

	for i := range m.Data {
		m.Data[i] = ^uint64(0)
	}

	m.ApplyColTailMask()
	return m, nil
}

func NewRandMatrix(rows, cols int, k int, rng *rand.Rand) (*Matrix, error) {
	m, err := NewZerosMatrix(rows, cols)
	if err != nil {
		return nil, err
	}

	for i := range m.Data {
		word := rng.Uint64()
		if k < 0 {
			// AND演算を繰り返し、確率を1/2ずつ下げる
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

	m.ApplyColTailMask()
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

func (m *Matrix) Stride() int {
	return (m.Cols + 63) / 64
}

func (m *Matrix) ColTailMask() uint64 {
	r := m.Cols % 64
	if r == 0 {
		return ^uint64(0)
	}
	return (uint64(1) << uint(r)) - 1
}

func (m *Matrix) ApplyColTailMask() {
	mask := m.ColTailMask()
	if mask == ^uint64(0) {
		return // マスク不要
	}

	stride := m.Stride()
	for r := 0; r < m.Rows; r++ {
		// 各行の64ビットの余りが出た列にマスクを適用
		idx := (r * stride) + (stride - 1)
		m.Data[idx] &= mask
	}
}

func (m *Matrix) Clone() *Matrix {
	return &Matrix{
		Rows: m.Rows,
		Cols: m.Cols,
		Data: slices.Clone(m.Data),
	}
}

func (m *Matrix) ValidateSameShape(other *Matrix) error {
	if m.Rows != other.Rows || m.Cols != other.Cols {
		return fmt.Errorf("dimension mismatch: (%dx%d) vs (%dx%d)", 
			m.Rows, m.Cols, other.Rows, other.Cols)
	}

	if len(m.Data) != len(other.Data) {
		return fmt.Errorf("internal data length mismatch: %d vs %d (rows:%d, cols:%d)", 
			len(m.Data), len(other.Data), m.Rows, m.Cols)
	}
	return nil
}

func (m *Matrix) And(other *Matrix) (*Matrix, error) {
	if err := m.ValidateSameShape(other); err != nil {
		return nil, err
	}
	c := m.Clone()
	for i := range c.Data {
		c.Data[i] &= other.Data[i]
	}
	c.ApplyColTailMask()
	return c, nil
}

func (m *Matrix) Xor(other *Matrix) (*Matrix, error) {
	if err := m.ValidateSameShape(other); err != nil {
		return nil, err
	}
	c := m.Clone()
	for i := range c.Data {
		c.Data[i] ^= other.Data[i]
	}
	c.ApplyColTailMask()
	return c, nil
}

func (m *Matrix) OnesCount() int {
	count := 0
	m.ScanRowsWord(nil, func(ctx MatrixWordContext) error {
		word := m.Data[ctx.WordIndex]
		if ctx.IsColTail {
			word &= m.ColTailMask()
		}
		count += bits.OnesCount64(word)
		return nil
	})
	return count
}

func (m *Matrix) HammingDistance(other *Matrix) (int, error) {
	// 異なるビットであれば1になる
	diff, err := m.Xor(other)
	if err != nil {
		return 0, err
	}
	return diff.OnesCount(), nil
}

func (m *Matrix) IndexAndShift(r, c int) (int, uint, error) {
	if r < 0 || r >= m.Rows {
		return 0, 0, fmt.Errorf("row が範囲外: row = %d: row < 0 || row >= Rows(=%d) であるべき", r, m.Rows)
	}
	if c < 0 || c >= m.Cols {
		return 0, 0, fmt.Errorf("col が範囲外: col = %d:col >= 0 && col < Cols(=%d) であるべき", c, m.Cols)
	}

	idx := (r * m.Stride()) + (c / 64)
	shift := uint(c % 64)
	return idx, shift, nil
}

func (m *Matrix) Bit(r, c int) (uint64, error) {
	idx, shift, err := m.IndexAndShift(r, c)
	if err != nil {
		return 0, err
	}
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
	m.Data[idx] ^= (1 << shift)
	return nil
}

func (m *Matrix) Dot(other *Matrix) ([]int, error) {
	if m.Cols != other.Cols {
		return nil, fmt.Errorf("dimension mismatch: m.Cols %d != other.Cols %d", m.Cols, other.Cols)
	}

	yRows := m.Rows
	yCols := other.Rows
	counts := make([]int, yRows*yCols)

	mData := m.Data
	oData := other.Data
	stride := m.Stride()
	mask := m.ColTailMask()

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
	if m.Cols != sign.Cols {
		return nil, fmt.Errorf("dimension mismatch: m.Cols %d != otherSign.Cols %d", m.Cols, sign.Cols)
	}

	if err := sign.ValidateSameShape(nonZero); err != nil {
		return nil, err
	}

	zRows := m.Rows
	zCols := sign.Rows
	z := make([]int, zRows*zCols)

	mData := m.Data
	sData := sign.Data
	nzData := nonZero.Data
	stride := m.Stride()
	mask := m.ColTailMask()

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

				// 符号が一致しているか
				signMath := ^(mWord ^ sWord)
				// 符号一致かつ非0
				validMatch := signMath & nzWord

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
	dst, err := NewZerosMatrix(m.Cols, m.Rows)
	if err != nil {
		return nil, err
	}

	var block [64]uint64

	srcStride := m.Stride()
	dstStride := dst.Stride()
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

	dst.ApplyColTailMask()
	return dst, nil
}

func (m *Matrix) ScanRowsWord(rowIdxs []int, f func(ctx MatrixWordContext) error) error {
	rows := m.Rows
	cols := m.Cols
	stride := m.Stride()

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

			var isColTail bool
			if colEnd > cols {
				colEnd = cols
				isColTail = true
			}

			err := f(MatrixWordContext{
				matrixRows:  rows,
				Row:         r,
				WordIndex:   rowWordOffset + s,
				ColStart:    colStart,
				ColEnd:      colEnd,
				GlobalStart: rowBitOffset + colStart,
				GlobalEnd:   rowBitOffset + colEnd,
				IsColTail:   isColTail,
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

func NewRFFMatrices(n, rows, cols int, sigma float32, rng *rand.Rand) (Matrices, error) {
	if n < 2 {
		return nil, fmt.Errorf("n must be at least 2")
	}

	totalBits := rows * cols
	omegas := make([]float32, totalBits)
	phases := make([]float32, totalBits)

	for i := 0; i < totalBits; i++ {
		omegas[i] = float32(rng.NormFloat64()) * sigma
		phases[i] = rng.Float32() * 2 * math.Pi
	}

	ms := make(Matrices, n)
	for i := 0; i < n; i++ {
		m, err := NewZerosMatrix(rows, cols)
		if err != nil {
			return nil, err
		}
		u := float32(i) / float32(n-1)

		err = m.ScanRowsWord(nil, func(ctx MatrixWordContext) error {
			var mWord uint64
			omegaWord := omegas[ctx.GlobalStart:ctx.GlobalEnd]
			phaseWord := phases[ctx.GlobalStart:ctx.GlobalEnd]
			ctx.ScanBits(func(i, col, colT int) error {
				y := float64(omegaWord[i]*u + phaseWord[i])
				z := float32(math.Cos(y))
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
	matrixRows  int
	Row         int
	WordIndex   int
	ColStart    int
	ColEnd      int
	GlobalStart int
	GlobalEnd   int
	IsColTail   bool
}

func (ctx MatrixWordContext) ScanBits(f func(i, col, colT int) error) error {
	colT := (ctx.ColStart * ctx.matrixRows) + ctx.Row
	for i := range ctx.ColEnd - ctx.ColStart {
		col := ctx.ColStart + i
		err := f(i, col, colT)
		if err != nil {
			return err
		}
		colT += ctx.matrixRows
	}
	return nil
}