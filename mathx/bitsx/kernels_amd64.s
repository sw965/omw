//go:build amd64

#include "textflag.h"

// AVX-512 (VPOPCNTQ) を使った popcount カーネル群。
// いずれも「XOR(とAND)して立っているビットを数える」だけの単純な処理を、
// 512bit(uint64 x 8ワード)ずつまとめて行う。
// 8ワードに満たない端数は、マスクレジスタ(K1)でゼロ埋めロードして同じ命令列で処理する。
//
// 呼び出し側(kernels.go / matrix.go)は「列の端数ビットは常に0」という
// Matrix の不変条件を前提にしているため、ここでは列マスク処理を行わない。

// func xorPopcntAVX512(a, b *uint64, n int) int
// n ワード分の popcount(a[i] ^ b[i]) の総和を返す。
TEXT ·xorPopcntAVX512(SB), NOSPLIT, $0-32
	MOVQ a+0(FP), SI
	MOVQ b+8(FP), DI
	MOVQ n+16(FP), CX
	VPXORQ Z0, Z0, Z0        // Z0 = 累積カウンタ(64bit x 8レーン)

loop:
	CMPQ CX, $8
	JLT  tail
	VMOVDQU64 (SI), Z1
	VPXORQ (DI), Z1, Z1      // Z1 = a ^ b
	VPOPCNTQ Z1, Z1          // 各64bitレーンのpopcount
	VPADDQ Z1, Z0, Z0
	ADDQ $64, SI
	ADDQ $64, DI
	SUBQ $8, CX
	JMP  loop

tail:
	TESTQ CX, CX
	JEQ   reduce
	MOVL  $1, AX
	SHLL  CX, AX
	DECL  AX                 // AX = (1 << 端数ワード数) - 1
	KMOVB AX, K1
	VMOVDQU64.Z (SI), K1, Z1 // 端数ワードだけロードし残りレーンは0
	VMOVDQU64.Z (DI), K1, Z2
	VPXORQ Z2, Z1, Z1
	VPOPCNTQ Z1, Z1
	VPADDQ Z1, Z0, Z0

reduce:
	// Z0 の8レーンを1つのスカラへ水平加算
	VEXTRACTI64X4 $1, Z0, Y1
	VPADDQ Y1, Y0, Y0
	VEXTRACTI128 $1, Y0, X1
	VPADDQ X1, X0, X0
	VPSRLDQ $8, X0, X1
	VPADDQ X1, X0, X0
	VMOVQ X0, AX
	MOVQ AX, ret+24(FP)
	VZEROUPPER
	RET

// func dotRowXorPopcntAVX512(mRow, oData *uint64, stride, oRows int, counts *int)
// m の1行(strideワード)と、行方向に連続配置された other の全 oRows 行との
// XOR-popcount を counts[0..oRows-1] に書き込む。
TEXT ·dotRowXorPopcntAVX512(SB), NOSPLIT, $0-40
	MOVQ mRow+0(FP), SI
	MOVQ oData+8(FP), DI
	MOVQ stride+16(FP), R9
	MOVQ oRows+24(FP), BX
	MOVQ counts+32(FP), DX

	// 端数ワード用マスク K1 = (1 << (stride % 8)) - 1 を先に作っておく
	MOVQ R9, CX
	ANDQ $7, CX
	MOVL $1, AX
	SHLL CX, AX
	DECL AX
	KMOVB AX, K1

rowLoop:
	TESTQ BX, BX
	JEQ   done
	MOVQ  SI, R10            // m の行の先頭に戻す
	MOVQ  R9, CX             // 残りワード数 = stride
	VPXORQ Z0, Z0, Z0

wordLoop:
	CMPQ CX, $8
	JLT  wordTail
	VMOVDQU64 (R10), Z1
	VPXORQ (DI), Z1, Z1
	VPOPCNTQ Z1, Z1
	VPADDQ Z1, Z0, Z0
	ADDQ $64, R10
	ADDQ $64, DI
	SUBQ $8, CX
	JMP  wordLoop

wordTail:
	TESTQ CX, CX
	JEQ   reduce2
	VMOVDQU64.Z (R10), K1, Z1
	VMOVDQU64.Z (DI), K1, Z2
	VPXORQ Z2, Z1, Z1
	VPOPCNTQ Z1, Z1
	VPADDQ Z1, Z0, Z0
	LEAQ (DI)(CX*8), DI      // other 側を端数分だけ進めて次の行へ

reduce2:
	VEXTRACTI64X4 $1, Z0, Y1
	VPADDQ Y1, Y0, Y0
	VEXTRACTI128 $1, Y0, X1
	VPADDQ X1, X0, X0
	VPSRLDQ $8, X0, X1
	VPADDQ X1, X0, X0
	VMOVQ X0, AX
	MOVQ AX, (DX)
	ADDQ $8, DX
	DECQ BX
	JMP  rowLoop

done:
	VZEROUPPER
	RET

// func dotTernaryRowAVX512(mRow, sData, nzData *uint64, stride, sRows int, z *int)
// m の1行と、sign/nonZero の全 sRows 行について
//   z[c] = popcount(nz行) - 2 * popcount((m行 ^ s行) & nz行)
// を書き込む。DotTernary の 2*matchCount - nonZeroCount と同値。
TEXT ·dotTernaryRowAVX512(SB), NOSPLIT, $0-48
	MOVQ mRow+0(FP), SI
	MOVQ sData+8(FP), DI
	MOVQ nzData+16(FP), R11
	MOVQ stride+24(FP), R9
	MOVQ sRows+32(FP), BX
	MOVQ z+40(FP), DX

	// 端数ワード用マスク K1 = (1 << (stride % 8)) - 1
	MOVQ R9, CX
	ANDQ $7, CX
	MOVL $1, AX
	SHLL CX, AX
	DECL AX
	KMOVB AX, K1

rowLoop:
	TESTQ BX, BX
	JEQ   done
	MOVQ  SI, R10            // m の行の先頭に戻す
	MOVQ  R9, CX             // 残りワード数 = stride
	VPXORQ Z0, Z0, Z0        // Z0 = popcount((m^s)&nz) の累積
	VPXORQ Z5, Z5, Z5        // Z5 = popcount(nz) の累積

wordLoop:
	CMPQ CX, $8
	JLT  wordTail
	VMOVDQU64 (R10), Z1
	VPXORQ (DI), Z1, Z1      // m ^ s
	VMOVDQU64 (R11), Z2      // nz
	VPANDQ Z2, Z1, Z1        // (m ^ s) & nz
	VPOPCNTQ Z1, Z1
	VPADDQ Z1, Z0, Z0
	VPOPCNTQ Z2, Z2
	VPADDQ Z2, Z5, Z5
	ADDQ $64, R10
	ADDQ $64, DI
	ADDQ $64, R11
	SUBQ $8, CX
	JMP  wordLoop

wordTail:
	TESTQ CX, CX
	JEQ   reduce3
	VMOVDQU64.Z (R10), K1, Z1
	VMOVDQU64.Z (DI), K1, Z3
	VMOVDQU64.Z (R11), K1, Z2
	VPXORQ Z3, Z1, Z1
	VPANDQ Z2, Z1, Z1
	VPOPCNTQ Z1, Z1
	VPADDQ Z1, Z0, Z0
	VPOPCNTQ Z2, Z2
	VPADDQ Z2, Z5, Z5
	LEAQ (DI)(CX*8), DI      // sign/nonZero 側を端数分だけ進めて次の行へ
	LEAQ (R11)(CX*8), R11

reduce3:
	// Z0 -> AX (不一致数), Z5 -> R8 (非ゼロ数)
	VEXTRACTI64X4 $1, Z0, Y1
	VPADDQ Y1, Y0, Y0
	VEXTRACTI128 $1, Y0, X1
	VPADDQ X1, X0, X0
	VPSRLDQ $8, X0, X1
	VPADDQ X1, X0, X0
	VMOVQ X0, AX

	VEXTRACTI64X4 $1, Z5, Y1
	VPADDQ Y1, Y5, Y5
	VEXTRACTI128 $1, Y5, X1
	VPADDQ X1, X5, X5
	VPSRLDQ $8, X5, X1
	VPADDQ X1, X5, X5
	VMOVQ X5, R8

	// z = nzCount - 2*mismatch
	NEGQ AX
	LEAQ (R8)(AX*2), AX
	MOVQ AX, (DX)
	ADDQ $8, DX
	DECQ BX
	JMP  rowLoop

done:
	VZEROUPPER
	RET
