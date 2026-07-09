//go:build amd64

package bitsx

import "golang.org/x/sys/cpu"

// useAVX512 は実行時にCPUの対応状況を見て決まる。
// VPOPCNTQ には AVX512F に加えて AVX512VPOPCNTDQ 拡張が必要。
var useAVX512 = cpu.X86.HasAVX512F && cpu.X86.HasAVX512VPOPCNTDQ

// 実装は kernels_amd64.s を参照。

// xorPopcntAVX512 は n ワード分の popcount(a[i]^b[i]) の総和を返す。
func xorPopcntAVX512(a, b *uint64, n int) int

// dotRowXorPopcntAVX512 は mRow(strideワード)と oData の全 oRows 行との
// XOR-popcount を counts に書き込む。
func dotRowXorPopcntAVX512(mRow, oData *uint64, stride, oRows int, counts *int)

// dotTernaryRowAVX512 は mRow と sign/nonZero の全 sRows 行について
// z[c] = popcount(nz行) - 2*popcount((mRow ^ s行) & nz行) を書き込む。
func dotTernaryRowAVX512(mRow, sData, nzData *uint64, stride, sRows int, z *int)
