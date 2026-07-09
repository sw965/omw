package bitsx

import "math/bits"

// xorPopcntGo は xorPopcntAVX512 と同じ計算の pure Go 版。
// AVX-512 非対応のCPU向けフォールバック。
func xorPopcntGo(a, b []uint64) int {
	sum := 0
	for i := range a {
		sum += bits.OnesCount64(a[i] ^ b[i])
	}
	return sum
}
