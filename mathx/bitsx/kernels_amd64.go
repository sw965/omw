//go:build amd64

package bitsx

import "golang.org/x/sys/cpu"

var useAVX512 = cpu.X86.HasAVX512F && cpu.X86.HasAVX512VPOPCNTDQ

// 実装は kernels_amd64.s

func xorPopcntAVX512(a, b *uint64, n int) int
func dotAVX512(leftData, rightData *uint64, leftRows, rightRows, cols, stride int, results *int)
func dotTernaryAVX512(valueData, signData, nonZeroData *uint64, valueRows, signRows, stride int, results *int)
