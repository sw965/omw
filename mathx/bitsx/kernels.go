package bitsx

import (
	"fmt"
	"math/bits"

	"github.com/sw965/omw/mathx"
)

// validateDotAVX512Familyは、rows/colsとdataの長さの対応が崩れていないかを
// 検査する。Dot/DotTernaryのAVX512パス(dotAVX512/dotTernaryAVX512)は、行数と
// strideの掛け算で計算した値を使ってバッファ内を飛び飛びに読み書きする為、
// この整合性が崩れたまま渡すと範囲外アクセスになる。GobDecodeがここを通すのも、
// デコード結果が後にDot/DotTernaryへ渡り得るからという同じ理由による。
// (xorPopcntAVX512はlen(data)をそのまま渡すだけなので、この検査を必要としない)
func (m *Matrix) validateDotAVX512Family() error {
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

// validateDotAVX512Args は、dotAVX512(kernels_amd64.go)を安全に呼び出す為の
// 前提条件(非負・桁あふれ無し・stride共有先のcolsの一致)を検証し、
// resultsに必要な長さを返す。
func validateDotAVX512Args(left, right *Matrix) (resultLen int, err error) {
	if left.cols != right.cols {
		return 0, fmt.Errorf("列数が不一致: m.Cols = %d, other.Cols = %d", left.cols, right.cols)
	}

	if err := left.validateDotAVX512Family(); err != nil {
		return 0, err
	}

	if err := right.validateDotAVX512Family(); err != nil {
		return 0, err
	}

	resultLen, ok := mathx.Mul(left.rows, right.rows)
	if !ok {
		return 0, fmt.Errorf("結果配列が大きすぎる: leftRows = %d, rightRows = %d", left.rows, right.rows)
	}
	return resultLen, nil
}

// validateDotTernaryAVX512Args は、dotTernaryAVX512(kernels_amd64.go)を安全に
// 呼び出す為の前提条件(非負・桁あふれ無し・stride共有先のcolsの一致)を検証し、
// resultsに必要な長さを返す。
func validateDotTernaryAVX512Args(value, sign, nonZero *Matrix) (resultLen int, err error) {
	if value.cols != sign.cols {
		return 0, fmt.Errorf("列数が不一致: m.Cols = %d, sign.Cols = %d", value.cols, sign.cols)
	}

	if err := sign.ValidateSameShape(nonZero); err != nil {
		return 0, err
	}

	// nonZero は sign と同形状だが、メモリ安全性を他の検査の呼び出し順序に
	// 依存させない為、3つとも明示的に検査する。
	if err := value.validateDotAVX512Family(); err != nil {
		return 0, err
	}

	if err := sign.validateDotAVX512Family(); err != nil {
		return 0, err
	}

	if err := nonZero.validateDotAVX512Family(); err != nil {
		return 0, err
	}

	resultLen, ok := mathx.Mul(value.rows, sign.rows)
	if !ok {
		return 0, fmt.Errorf("結果配列が大きすぎる: valueRows = %d, signRows = %d", value.rows, sign.rows)
	}
	return resultLen, nil
}

func xorPopcntGo(a, b []uint64) int {
	sum := 0
	for i := range a {
		sum += bits.OnesCount64(a[i] ^ b[i])
	}
	return sum
}

func dotGo(leftData, rightData []uint64, leftRows, rightRows, cols, stride int, results []int) {
	// 端数ビットは常に0が前提条件

	for r := range leftRows {
		leftRow := leftData[r*stride : (r+1)*stride]
		resultsRow := results[r*rightRows : (r+1)*rightRows]
		for c := range rightRows {
			rightRow := rightData[c*stride : (c+1)*stride]
			resultsRow[c] = cols - xorPopcntGo(leftRow, rightRow)
		}
	}
}

func dotTernaryGo(valueData, signData, nonZeroData []uint64, valueRows, signRows, stride int, results []int) {
	nonZeroCounts := make([]int, signRows)
	for c := range signRows {
		nonZeroRow := nonZeroData[c*stride : (c+1)*stride]
		count := 0
		for _, w := range nonZeroRow {
			count += bits.OnesCount64(w)
		}
		nonZeroCounts[c] = count
	}

	for r := range valueRows {
		valueRow := valueData[r*stride : (r+1)*stride]
		resultsRow := results[r*signRows : (r+1)*signRows]
		for c := range signRows {
			signRow := signData[c*stride : (c+1)*stride]
			nonZeroRow := nonZeroData[c*stride : (c+1)*stride]
			mismatchCount := 0
			for k := range valueRow {
				mismatchCount += bits.OnesCount64((valueRow[k] ^ signRow[k]) & nonZeroRow[k])
			}
			resultsRow[c] = nonZeroCounts[c] - 2*mismatchCount
		}
	}
}
