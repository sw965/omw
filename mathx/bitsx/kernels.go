package bitsx

import (
	"fmt"
	"math/bits"

	"github.com/sw965/omw/mathx"
)

func (m *Matrix) validateDotAVX512Family() error {
	if m.rows <= 0 {
		return fmt.Errorf("行数が不正: Rows = %d: Rows > 0 であるべき", m.rows)
	}

	if m.cols <= 0 {
		return fmt.Errorf("列数が不正: Cols = %d: Cols > 0 であるべき", m.cols)
	}

	// 内部計算が桁あふれして負になるケースをガードする
	stride := m.Stride()
	if stride <= 0 {
		return fmt.Errorf("列数が大きすぎる: Cols = %d", m.cols)
	}

	wantDataLen, ok := mathx.MulOverflowChecked(m.rows, stride)
	if !ok {
		return fmt.Errorf("RowsとColsが大きすぎる: Rows = %d, Cols = %d", m.rows, m.cols)
	}

	if wantDataLen != len(m.data) {
		return fmt.Errorf("内部データ長が不正: len(data) = %d: Rows(=%d) * Stride(=%d) = %d と一致するべき",
			len(m.data), m.rows, stride, wantDataLen)
	}
	return nil
}

func validateDotAVX512Args(left, right *Matrix) (resultsLen int, err error) {
	if left.cols != right.cols {
		return 0, fmt.Errorf("列数が不一致: m.Cols = %d, other.Cols = %d", left.cols, right.cols)
	}

	if err := left.validateDotAVX512Family(); err != nil {
		return 0, err
	}

	if err := right.validateDotAVX512Family(); err != nil {
		return 0, err
	}

	resultsLen, ok := mathx.MulOverflowChecked(left.rows, right.rows)
	if !ok {
		return 0, fmt.Errorf("結果配列が大きすぎる: leftRows = %d, rightRows = %d", left.rows, right.rows)
	}
	return resultsLen, nil
}

func validateDotTernaryAVX512Args(value, sign, nonZero *Matrix) (resultsLen int, err error) {
	if value.cols != sign.cols {
		return 0, fmt.Errorf("列数が不一致: m.Cols = %d, sign.Cols = %d", value.cols, sign.cols)
	}

	if err := sign.ValidateSameShape(nonZero); err != nil {
		return 0, err
	}

	if err := value.validateDotAVX512Family(); err != nil {
		return 0, err
	}

	if err := sign.validateDotAVX512Family(); err != nil {
		return 0, err
	}

	if err := nonZero.validateDotAVX512Family(); err != nil {
		return 0, err
	}

	resultsLen, ok := mathx.MulOverflowChecked(value.rows, sign.rows)
	if !ok {
		return 0, fmt.Errorf("結果配列が大きすぎる: valueRows = %d, signRows = %d", value.rows, sign.rows)
	}
	return resultsLen, nil
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
