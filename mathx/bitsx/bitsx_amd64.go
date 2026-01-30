//go:build amd64

package bitsx

// dotAVX512 は bitsx_amd64.s で実装される完全アセンブリ版のDot処理です。
//
//go:noescape
func dotAVX512(mData, oData []uint64, out []int, mRows, oRows, stride int, mask uint64)