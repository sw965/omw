package bitsx

import (
	"fmt"
	"math"
	"math/bits"
	"math/rand/v2"
	"slices"
)

// Matrixはビット列を行列として扱う。
// コンストラクタで生成したMatrixは、端数ビット(Cols % 64 の範囲外)が常に0に保たれる。
type Matrix struct {
	rows int
	cols int
	data []uint64
}

func (m *Matrix) Rows() int {
	return m.rows
}

func (m *Matrix) Cols() int {
	return m.cols
}

func NewZerosMatrix(rows, cols int) (*Matrix, error) {
	if rows <= 0 {
		return nil, fmt.Errorf("rows > 0 であるべき: rows = %d", rows)
	}

	if cols <= 0 {
		return nil, fmt.Errorf("cols > 0 であるべき: cols = %d", cols)
	}

	m := &Matrix{
		rows: rows,
		cols: cols,
	}

	stride := m.Stride()
	if stride <= 0 {
		return nil, fmt.Errorf("colsが大きすぎる: cols = %d", cols)
	}

	n := rows * stride
	m.data = make([]uint64, n)
	return m, nil
}

func NewOnesMatrix(rows, cols int) (*Matrix, error) {
	m, err := NewZerosMatrix(rows, cols)
	if err != nil {
		return nil, err
	}

	for i := range m.data {
		m.data[i] = ^uint64(0)
	}

	m.ApplyTailMask()
	return m, nil
}

func NewRandMatrix(rows, cols int, k int, rng *rand.Rand) (*Matrix, error) {
	m, err := NewZerosMatrix(rows, cols)
	if err != nil {
		return nil, err
	}

	for i := range m.data {
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
		m.data[i] = word
	}

	m.ApplyTailMask()
	return m, nil
}

func NewSignMatrix(rows, cols int, x []int) (*Matrix, error) {
	n := rows * cols
	if n != len(x) {
		return nil, fmt.Errorf("len(x) == rows * cols であるべき: len(x) = %d, rows = %d, cols = %d", len(x), rows, cols)
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
		sign.data[ctx.WordIndex] = signWord
		return nil
	})

	if err != nil {
		return nil, err
	}
	return sign, nil
}

func (m *Matrix) Stride() int {
	return (m.cols + 63) / 64
}

func (m *Matrix) TailMask() uint64 {
	r := m.cols % 64
	if r == 0 {
		return ^uint64(0)
	}
	return (uint64(1) << uint(r)) - 1
}

func (m *Matrix) ApplyTailMask() {
	mask := m.TailMask()
	if mask == ^uint64(0) {
		return // マスク不要
	}

	stride := m.Stride()
	for r := 0; r < m.rows; r++ {
		// 各行の64ビットの余りが出た列にマスクを適用
		idx := (r * stride) + (stride - 1)
		m.data[idx] &= mask
	}
}

func (m *Matrix) Word(idx int) (uint64, error) {
	if idx < 0 || idx >= len(m.data) {
		return 0, fmt.Errorf("0 <= index < %d であるべき: index = %d", len(m.data), idx)
	}
	return m.data[idx], nil
}

func (m *Matrix) SetWord(idx int, word uint64) error {
	if idx < 0 || idx >= len(m.data) {
		return fmt.Errorf("0 <= index < %d であるべき: index = %d", len(m.data), idx)
	}
	if idx%m.Stride() == m.Stride()-1 {
		word &= m.TailMask()
	}
	m.data[idx] = word
	return nil
}

func (m *Matrix) Clone() *Matrix {
	return &Matrix{
		rows: m.rows,
		cols: m.cols,
		data: slices.Clone(m.data),
	}
}

func (m *Matrix) Equal(other *Matrix) bool {
	if m.rows != other.rows || m.cols != other.cols {
		return false
	}
	return slices.Equal(m.data, other.data)
}

func (m *Matrix) ValidateSameShape(other *Matrix) error {
	if m.rows != other.rows || m.cols != other.cols {
		return fmt.Errorf("形状の不一致: (%d x %d) vs (%d x %d)", m.rows, m.cols, other.rows, other.cols)
	}

	if len(m.data) != len(other.data) {
		return fmt.Errorf("データ長の不一致: %d vs %d", len(m.data), len(other.data))
	}
	return nil
}

func (m *Matrix) And(other *Matrix) (*Matrix, error) {
	if err := m.ValidateSameShape(other); err != nil {
		return nil, err
	}

	c := m.Clone()
	for i := range c.data {
		c.data[i] &= other.data[i]
	}
	return c, nil
}

func (m *Matrix) Xor(other *Matrix) (*Matrix, error) {
	if err := m.ValidateSameShape(other); err != nil {
		return nil, err
	}

	c := m.Clone()
	for i := range c.data {
		c.data[i] ^= other.data[i]
	}
	return c, nil
}

func (m *Matrix) OnesCount() int {
	count := 0
	err := m.ScanRowsWord(nil, func(ctx MatrixWordContext) error {
		word := m.data[ctx.WordIndex]
		if ctx.IsTail {
			word &= m.TailMask()
		}
		count += bits.OnesCount64(word)
		return nil
	})

	// rowIdxs=nil なので、エラーは発生しないはずだが、念のため
	if err != nil {
		panic(err)
	}
	return count
}

func (m *Matrix) HammingDistance(other *Matrix) (int, error) {
	if err := m.ValidateSameShape(other); err != nil {
		return 0, err
	}

	if len(m.data) == 0 {
		return 0, nil
	}

	if useAVX512 {
		return xorPopcntAVX512(&m.data[0], &other.data[0], len(m.data)), nil
	}
	return xorPopcntGo(m.data, other.data), nil
}

func (m *Matrix) IndexAndShift(r, c int) (int, uint, error) {
	if r < 0 || r >= m.rows {
		return 0, 0, fmt.Errorf("0 <= row < %d であるべき: row = %d", m.rows, r)
	}

	if c < 0 || c >= m.cols {
		return 0, 0, fmt.Errorf("0 <= col < %d であるべき: col = %d", m.cols, c)
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
	return (m.data[idx] >> shift) & 1, nil
}

func (m *Matrix) Set(r, c int) error {
	idx, shift, err := m.IndexAndShift(r, c)
	if err != nil {
		return err
	}
	m.data[idx] |= (1 << shift)
	return nil
}

func (m *Matrix) Clear(r, c int) error {
	idx, shift, err := m.IndexAndShift(r, c)
	if err != nil {
		return err
	}
	m.data[idx] &^= (1 << shift)
	return nil
}

func (m *Matrix) Toggle(r, c int) error {
	idx, shift, err := m.IndexAndShift(r, c)
	if err != nil {
		return err
	}
	m.data[idx] ^= (1 << shift)
	return nil
}

func (m *Matrix) Dot(other *Matrix) ([]int, error) {
	resultsLen, err := validateDotAVX512Args(m, other)
	if err != nil {
		return nil, err
	}

	leftRows := m.rows
	rightRows := other.rows
	stride := m.Stride()
	results := make([]int, resultsLen)

	if useAVX512 {
		dotAVX512(&m.data[0], &other.data[0], leftRows, rightRows, m.cols, stride, &results[0])
	} else {
		dotGo(m.data, other.data, leftRows, rightRows, m.cols, stride, results)
	}
	return results, nil
}

func (m *Matrix) DotTernary(sign, nonZero *Matrix) ([]int, error) {
	resultsLen, err := validateDotTernaryAVX512Args(m, sign, nonZero)
	if err != nil {
		return nil, err
	}

	valueRows := m.rows
	signRows := sign.rows
	stride := m.Stride()
	results := make([]int, resultsLen)

	if useAVX512 {
		dotTernaryAVX512(&m.data[0], &sign.data[0], &nonZero.data[0], valueRows, signRows, stride, &results[0])
	} else {
		dotTernaryGo(m.data, sign.data, nonZero.data, valueRows, signRows, stride, results)
	}
	return results, nil
}

// 64x64ビットブロックの転置で用いる、交換する幅ごとのマスク。
const (
	swapMask32 = 0x00000000FFFFFFFF
	swapMask16 = 0x0000FFFF0000FFFF
	swapMask8  = 0x00FF00FF00FF00FF
	swapMask4  = 0x0F0F0F0F0F0F0F0F
	swapMask2  = 0x3333333333333333
	swapMask1  = 0x5555555555555555
)

// aとbの間で、幅dのビットブロックを交換する。
func swapBitBlock(a, b uint64, d uint, mask uint64) (uint64, uint64) {
	t := (b ^ (a >> d)) & mask
	return a ^ (t << d), b ^ t
}

// 8ワードに対して、幅d1→d2→d3の順に3段のビットブロック交換を適用する。
// 交換する組み合わせは、幅32/16/8の3段でも、幅4/2/1の3段でも同一になる。
func swapBitBlock3Stages(
	w0, w1, w2, w3, w4, w5, w6, w7 uint64,
	d1, d2, d3 uint,
	mask1, mask2, mask3 uint64,
) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	w0, w4 = swapBitBlock(w0, w4, d1, mask1)
	w1, w5 = swapBitBlock(w1, w5, d1, mask1)
	w2, w6 = swapBitBlock(w2, w6, d1, mask1)
	w3, w7 = swapBitBlock(w3, w7, d1, mask1)

	w0, w2 = swapBitBlock(w0, w2, d2, mask2)
	w1, w3 = swapBitBlock(w1, w3, d2, mask2)
	w4, w6 = swapBitBlock(w4, w6, d2, mask2)
	w5, w7 = swapBitBlock(w5, w7, d2, mask2)

	w0, w1 = swapBitBlock(w0, w1, d3, mask3)
	w2, w3 = swapBitBlock(w2, w3, d3, mask3)
	w4, w5 = swapBitBlock(w4, w5, d3, mask3)
	w6, w7 = swapBitBlock(w6, w7, d3, mask3)

	return w0, w1, w2, w3, w4, w5, w6, w7
}

// 64x64ビットブロックの転置は、幅32→16→8→4→2→1の6段のビットブロック交換で行う。
// 各段の交換相手は、ブロック内の行番号の特定のビットだけを跨ぐ為、次の2組に分けられる。
//
//   - 前半3段(幅32/16/8): 行番号のbit3,4,5のみを跨ぐ → 行 {c, c+8, ..., c+56} の8行で閉じる
//   - 後半3段(幅4/2/1)  : 行番号のbit0,1,2のみを跨ぐ → 行 {8h, 8h+1, ..., 8h+7} の8行で閉じる
//
// 8行をレジスタに載せたまま3段を処理できるので、段ごとに配列へ書き戻す必要がない。
// 更に、前半3段は入力の読み込みと、後半3段は出力の書き込みと融合できる為、
// 64ワードの中間バッファへの往復は1回で済む。

// 前半3段(幅32/16/8)。src側の行を直接読み、blockへ書き出す。
// srcの読み込み位置は base + i*stride (i = 0..63) の64行。
func transpose64BlockHead(src []uint64, base, stride int, block *[64]uint64) {
	step := stride * 8
	s := base
	for c := range 8 {
		w0, w1, w2, w3, w4, w5, w6, w7 := swapBitBlock3Stages(
			src[s], src[s+step], src[s+step*2], src[s+step*3],
			src[s+step*4], src[s+step*5], src[s+step*6], src[s+step*7],
			32, 16, 8, swapMask32, swapMask16, swapMask8,
		)

		block[c] = w0
		block[c+8] = w1
		block[c+16] = w2
		block[c+24] = w3
		block[c+32] = w4
		block[c+40] = w5
		block[c+48] = w6
		block[c+56] = w7
		s += stride
	}
}

// 後半3段(幅4/2/1)。blockから読み、dst側の行へ直接書き出す。
// dstの書き込み位置は base + i*stride (i = 0..63) の64行。
func transpose64BlockTail(block *[64]uint64, dst []uint64, base, stride int) {
	d := base
	for h := 0; h < 64; h += 8 {
		w0, w1, w2, w3, w4, w5, w6, w7 := swapBitBlock3Stages(
			block[h], block[h+1], block[h+2], block[h+3],
			block[h+4], block[h+5], block[h+6], block[h+7],
			4, 2, 1, swapMask4, swapMask2, swapMask1,
		)

		dst[d] = w0
		dst[d+stride] = w1
		dst[d+stride*2] = w2
		dst[d+stride*3] = w3
		dst[d+stride*4] = w4
		dst[d+stride*5] = w5
		dst[d+stride*6] = w6
		dst[d+stride*7] = w7
		d += stride * 8
	}
}

// src側の64x64ビットブロックを転置して、dst側へ書き出す。
// srcRows・dstRowsは、それぞれの側で実際に読み書きできる残り行数。
func transpose64Block(
	src []uint64, srcBase, srcStride, srcRows int,
	dst []uint64, dstBase, dstStride, dstRows int,
	block *[64]uint64,
) {
	// ホットパス: 64行揃っているので、読み込み・書き込みを転置に融合できる
	if srcRows >= 64 && dstRows >= 64 {
		transpose64BlockHead(src, srcBase, srcStride, block)
		transpose64BlockTail(block, dst, dstBase, dstStride)
		return
	}

	// エッジケース: blockを経由する。
	// blockに対して base = 0, stride = 1 を渡すと、前半3段・後半3段のどちらも
	// block内で完結する読み書きになる。
	readRows := min(srcRows, 64)
	s := srcBase
	for i := range readRows {
		block[i] = src[s]
		s += srcStride
	}
	// 足りない部分は0埋め（ゴミデータが混ざらないように）
	for i := readRows; i < 64; i++ {
		block[i] = 0
	}

	transpose64BlockHead(block[:], 0, 1, block)
	transpose64BlockTail(block, block[:], 0, 1)

	writeRows := min(dstRows, 64)
	d := dstBase
	for i := range writeRows {
		dst[d] = block[i]
		d += dstStride
	}
}

func (m *Matrix) Transpose() (*Matrix, error) {
	dst, err := NewZerosMatrix(m.cols, m.rows)
	if err != nil {
		return nil, err
	}

	var block [64]uint64

	srcStride := m.Stride()
	dstStride := dst.Stride()
	rows := m.rows

	// ブロック単位での処理 (64行ずつ)
	for r := 0; r < rows; r += 64 {
		// 残り行数が64未満なら、エッジケースとして扱われる
		srcRows := rows - r
		// 転置後は、dstの「r列が含まれるワード」に書き込まれる
		// rは常に64の倍数なので単純なシフト
		dstColWord := r / 64
		srcRowOffset := r * srcStride

		// 横方向（Word単位）のループ
		for cWord := range srcStride {
			// 転置後は、dstの「cWord * 64 + (0..63)」行目に書き込まれる
			dstRowBase := cWord * 64

			transpose64Block(
				m.data, srcRowOffset+cWord, srcStride, srcRows,
				dst.data, dstRowBase*dstStride+dstColWord, dstStride, dst.rows-dstRowBase,
				&block,
			)
		}
	}

	dst.ApplyTailMask()
	return dst, nil
}

func (m *Matrix) ScanRowsWord(rowIdxs []int, f func(ctx MatrixWordContext) error) error {
	rows := m.rows
	cols := m.cols
	stride := m.Stride()

	if rowIdxs == nil {
		rowIdxs = make([]int, rows)
		for i := range rows {
			rowIdxs[i] = i
		}
	}

	for _, r := range rowIdxs {
		if r < 0 || r >= rows {
			return fmt.Errorf("row が範囲外: row = %d: 0 <= row < %d であるべき", r, rows)
		}

		rowWordOffset := r * stride
		rowBitOffset := r * cols
		for s := range stride {
			colStart := s << 6
			colEnd := colStart + 64

			var isTail bool
			if colEnd > cols {
				colEnd = cols
				isTail = true
			}

			err := f(MatrixWordContext{
				matrixRows:  rows,
				Row:         r,
				WordIndex:   rowWordOffset + s,
				ColStart:    colStart,
				ColEnd:      colEnd,
				GlobalStart: rowBitOffset + colStart,
				GlobalEnd:   rowBitOffset + colEnd,
				IsTail:      isTail,
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
	if n < 2 {
		return nil, fmt.Errorf("nが不正(n < 2): n = %d", n)
	}

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
		return nil, fmt.Errorf("nが不正(n < 2): n = %d", n)
	}

	if rows <= 0 {
		return nil, fmt.Errorf("行数が不正: rows = %d: rows > 0 であるべき", rows)
	}

	if cols <= 0 {
		return nil, fmt.Errorf("列数が不正: cols = %d: cols > 0 であるべき", cols)
	}

	totalBits := rows * cols

	omegas := make([]float32, totalBits)
	phases := make([]float32, totalBits)

	for i := range totalBits {
		omegas[i] = float32(rng.NormFloat64()) * sigma
		phases[i] = rng.Float32() * 2 * math.Pi
	}

	ms := make(Matrices, n)
	for i := range n {
		m, err := NewZerosMatrix(rows, cols)
		if err != nil {
			return nil, err
		}
		u := float32(i) / float32(n-1)

		err = m.ScanRowsWord(nil, func(ctx MatrixWordContext) error {
			var mWord uint64
			omegaWord := omegas[ctx.GlobalStart:ctx.GlobalEnd]
			phaseWord := phases[ctx.GlobalStart:ctx.GlobalEnd]
			scanErr := ctx.ScanBits(func(i, col, colT int) error {
				y := float64(omegaWord[i]*u + phaseWord[i])
				z := float32(math.Cos(y))
				if z >= 0 {
					mWord |= (1 << uint(i))
				}
				return nil
			})
			if scanErr != nil {
				return scanErr
			}
			m.data[ctx.WordIndex] = mWord
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
		return nil, fmt.Errorf("nが不正(n < 2): n = %d", n)
	}

	ms := make(Matrices, n)
	totalBits := rows * cols

	for i := range n {
		m, err := NewZerosMatrix(rows, cols)
		if err != nil {
			return nil, err
		}

		scaled := i * totalBits
		numOnes := scaled / (n - 1)
		err = m.ScanRowsWord(nil, func(ctx MatrixWordContext) error {
			var word uint64
			scanErr := ctx.ScanBits(func(i, col, colT int) error {
				if (ctx.GlobalStart + i) < numOnes {
					word |= (uint64(1) << uint(i))
				}
				return nil
			})
			if scanErr != nil {
				return scanErr
			}
			m.data[ctx.WordIndex] = word
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
	// n < 2 の場合、距離のペアが1つも存在せず、平均の計算がゼロ除算になる
	if n < 2 {
		return 0.0, fmt.Errorf("nが不正(n < 2): n = %d", n)
	}

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
	// また、距離の分散もコストに含める事で、均等な距離を保ちながら、距離を最大化する事が出来る
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
	IsTail      bool
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
