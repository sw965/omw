package bitsx

import "math/bits"

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
