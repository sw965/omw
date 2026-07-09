//go:build !amd64

package bitsx

// amd64 以外では常に pure Go 実装へフォールバックする。
var useAVX512 = false

func xorPopcntAVX512(a, b *uint64, n int) int { panic("unreachable") }

func dotRowXorPopcntAVX512(mRow, oData *uint64, stride, oRows int, counts *int) {
	panic("unreachable")
}

func dotTernaryRowAVX512(mRow, sData, nzData *uint64, stride, sRows int, z *int) {
	panic("unreachable")
}
