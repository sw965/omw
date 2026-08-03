//go:build !amd64

package bitsx

// amd64 以外では常に pure Go 実装へフォールバックする。
var useAVX512 = false

func xorPopcntAVX512(a, b *uint64, n int) int {
	panic("unreachable")
}

func dotAVX512(leftData, rightData *uint64, leftRows, rightRows, cols, stride int, results *int) {
	panic("unreachable")
}

func dotTernaryAVX512(valueData, signData, nonZeroData *uint64, valueRows, signRows, stride int, results *int) {
	panic("unreachable")
}
