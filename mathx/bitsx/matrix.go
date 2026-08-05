package bitsx

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math"
	"math/bits"
	"math/rand/v2"
	"slices"

	"github.com/sw965/omw/mathx"
)

type Matrix struct {
	rows int
	cols int
	// 各行の端数ビット(Cols % 64 の範囲外)は常に0に保たれる。
	// Dot / DotTernary / HammingDistance はこの不変条件を前提に
	// 列マスク処理を省略している。
	// 非公開フィールドとし、SetWordのみを書き込み経路とすることでこの不変条件を保証する。
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
		return nil, fmt.Errorf("行数が不正: rows = %d: rows > 0 であるべき", rows)
	}

	if cols <= 0 {
		return nil, fmt.Errorf("列数が不正: cols = %d: cols > 0 であるべき", cols)
	}

	m := &Matrix{
		rows: rows,
		cols: cols,
	}

	stride := m.Stride()
	if stride <= 0 {
		return nil, fmt.Errorf("列数が大きすぎる: cols = %d", cols)
	}

	n, ok := mathx.Mul(rows, stride)
	if !ok {
		return nil, fmt.Errorf("rowsとcolsが大きすぎる: rows = %d, cols = %d", rows, cols)
	}

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
	need, ok := mathx.Mul(rows, cols)
	if !ok {
		return nil, fmt.Errorf("rowsとcolsが大きすぎる: rows = %d, cols = %d", rows, cols)
	}

	if len(x) < need {
		return nil, fmt.Errorf("len(x)が不足: len(x) = %d: %d 以上であるべき", len(x), need)
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
		return 0, fmt.Errorf("idxが範囲外: idx = %d: 0 <= idx < %d であるべき", idx, len(m.data))
	}
	return m.data[idx], nil
}

func (m *Matrix) SetWord(idx int, word uint64) error {
	if idx < 0 || idx >= len(m.data) {
		return fmt.Errorf("idxが範囲外: idx = %d: 0 <= idx < %d であるべき", idx, len(m.data))
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

// gobEncodedMatrixは、非公開のdataフィールドをencoding/gobでやり取りする為の中継用の型。
type gobEncodedMatrix struct {
	Rows int
	Cols int
	Data []uint64
}

func (m *Matrix) GobEncode() ([]byte, error) {
	buf := &bytes.Buffer{}
	payload := gobEncodedMatrix{Rows: m.rows, Cols: m.cols, Data: m.data}
	if err := gob.NewEncoder(buf).Encode(payload); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (m *Matrix) GobDecode(b []byte) error {
	var payload gobEncodedMatrix
	if err := gob.NewDecoder(bytes.NewReader(b)).Decode(&payload); err != nil {
		return err
	}

	decoded := &Matrix{rows: payload.Rows, cols: payload.Cols, data: payload.Data}
	if err := decoded.Validate(); err != nil {
		return fmt.Errorf("デコードされたMatrixが不正: %w", err)
	}

	// 直列化元が不変条件(端数ビットは常に0)を満たさないデータであっても、
	// ここで強制する。
	decoded.ApplyTailMask()

	*m = *decoded
	return nil
}

func (m *Matrix) Validate() error {
	if m.rows <= 0 {
		return fmt.Errorf("行数が不正: Rows = %d: Rows > 0 であるべき", m.rows)
	}

	if m.cols <= 0 {
		return fmt.Errorf("列数が不正: Cols = %d: Cols > 0 であるべき", m.cols)
	}

	// Colsが極端に大きいと Stride() の内部計算が桁あふれして負になる
	stride := m.Stride()
	if stride <= 0 {
		return fmt.Errorf("列数が大きすぎる: Cols = %d", m.cols)
	}

	need, ok := mathx.Mul(m.rows, stride)
	if !ok {
		return fmt.Errorf("RowsとColsが大きすぎる: Rows = %d, Cols = %d", m.rows, m.cols)
	}

	if need != len(m.data) {
		return fmt.Errorf("内部データ長が不正: len(data) = %d: Rows(=%d) * Stride(=%d) = %d と一致するべき",
			len(m.data), m.rows, stride, need)
	}
	return nil
}

func (m *Matrix) ValidateSameShape(other *Matrix) error {
	if m.rows != other.rows || m.cols != other.cols {
		return fmt.Errorf("形状が不一致: (%dx%d) vs (%dx%d)",
			m.rows, m.cols, other.rows, other.cols)
	}

	if len(m.data) != len(other.data) {
		return fmt.Errorf("内部データ長が不一致: len(data) = %d vs %d (Rows = %d, Cols = %d)",
			len(m.data), len(other.data), m.rows, m.cols)
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

	c.ApplyTailMask()
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

	c.ApplyTailMask()
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
	if err != nil {
		// rowIdxsをnilで渡している為、行範囲エラーは発生し得ない
		panic(err)
	}
	return count
}

func (m *Matrix) HammingDistance(other *Matrix) (int, error) {
	if err := m.ValidateSameShape(other); err != nil {
		return 0, err
	}

	// 読むのは len(m.data) ワードだけで、その長さが等しいことは ValidateSameShape が
	// 保証している為、これ以上の検査は要らない。空の場合のみ &data[0] が取れない。
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
		return 0, 0, fmt.Errorf("row が範囲外: row = %d: row < 0 || row >= Rows(=%d) であるべき", r, m.rows)
	}
	if c < 0 || c >= m.cols {
		return 0, 0, fmt.Errorf("col が範囲外: col = %d: col >= 0 && col < Cols(=%d) であるべき", c, m.cols)
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
	if m.cols != other.cols {
		return nil, fmt.Errorf("列数が不一致: m.Cols = %d, other.Cols = %d", m.cols, other.cols)
	}

	if err := m.Validate(); err != nil {
		return nil, err
	}

	if err := other.Validate(); err != nil {
		return nil, err
	}

	leftRows := m.rows
	rightRows := other.rows
	stride := m.Stride()
	resultLen, ok := mathx.Mul(leftRows, rightRows)
	if !ok {
		return nil, fmt.Errorf("結果配列が大きすぎる: leftRows = %d, rightRows = %d", leftRows, rightRows)
	}
	results := make([]int, resultLen)

	if useAVX512 {
		dotAVX512(&m.data[0], &other.data[0], leftRows, rightRows, m.cols, stride, &results[0])
	} else {
		dotGo(m.data, other.data, leftRows, rightRows, m.cols, stride, results)
	}
	return results, nil
}

func (m *Matrix) DotTernary(sign, nonZero *Matrix) ([]int, error) {
	if m.cols != sign.cols {
		return nil, fmt.Errorf("列数が不一致: m.Cols = %d, sign.Cols = %d", m.cols, sign.cols)
	}

	if err := sign.ValidateSameShape(nonZero); err != nil {
		return nil, err
	}

	// nonZero は sign と同形状だが、メモリ安全性を他の検査の呼び出し順序に
	// 依存させない為、3つとも明示的に検査する。
	if err := m.Validate(); err != nil {
		return nil, err
	}

	if err := sign.Validate(); err != nil {
		return nil, err
	}

	if err := nonZero.Validate(); err != nil {
		return nil, err
	}

	valueRows := m.rows
	signRows := sign.rows
	stride := m.Stride()
	resultLen, ok := mathx.Mul(valueRows, signRows)
	if !ok {
		return nil, fmt.Errorf("結果配列が大きすぎる: valueRows = %d, signRows = %d", valueRows, signRows)
	}
	results := make([]int, resultLen)

	if useAVX512 {
		dotTernaryAVX512(&m.data[0], &sign.data[0], &nonZero.data[0], valueRows, signRows, stride, &results[0])
	} else {
		dotTernaryGo(m.data, sign.data, nonZero.data, valueRows, signRows, stride, results)
	}
	return results, nil
}

func transpose64Block(block *[64]uint64) {
	var (
		mask uint64
		t    uint64
		a, b uint64
	)

	// 32x32 swap
	mask = 0x00000000FFFFFFFF
	for j := range 32 {
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

	srcStride := m.Stride()
	dstStride := dst.Stride()
	srcData := m.data
	dstData := dst.data
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
		for cWord := range srcStride {
			// 1. 読み込み (Read)
			// Optimize: インデックス計算の乗算を避けるため、ベースオフセットを計算
			srcBaseIdx := r*srcStride + cWord

			if isFullBlock {
				// ホットパス: 分岐なしで64回読み込む
				// コンパイラによるBounds Check Eliminationが効きやすくなる
				for i := range 64 {
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
				for i := range 64 {
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

	totalBits, ok := mathx.Mul(rows, cols)
	if !ok {
		return nil, fmt.Errorf("rowsとcolsが大きすぎる: rows = %d, cols = %d", rows, cols)
	}

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
	totalBits, ok := mathx.Mul(rows, cols)
	if !ok {
		return nil, fmt.Errorf("rowsとcolsが大きすぎる: rows = %d, cols = %d", rows, cols)
	}

	for i := range n {
		m, err := NewZerosMatrix(rows, cols)
		if err != nil {
			return nil, err
		}

		scaled, ok := mathx.Mul(i, totalBits)
		if !ok {
			return nil, fmt.Errorf("n * rows * cols が大きすぎる: n = %d, rows = %d, cols = %d", n, rows, cols)
		}

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

	dsn, ok := mathx.Mul(n, n)
	if !ok {
		return 0.0, fmt.Errorf("nが大きすぎる")
	}

	distances := make([]float32, 0, dsn)
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
