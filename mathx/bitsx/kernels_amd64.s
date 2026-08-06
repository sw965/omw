//go:build amd64

#include "textflag.h"

// 呼び出し側が「列の端数ビットは常に0」という不変条件を保証する為、列マスク処理を省略する。

// Z0 の8レーンを AX へ水平加算する(Z0/Y1/X1 を破壊)。
#define HSUM_Z0_AX \
	VEXTRACTI64X4 $1, Z0, Y1; \
	VPADDQ        Y1, Y0, Y0; \
	VEXTRACTI128  $1, Y0, X1; \
	VPADDQ        X1, X0, X0; \
	VPSRLDQ       $8, X0, X1; \
	VPADDQ        X1, X0, X0; \
	VMOVQ         X0, AX

// func xorPopcntAVX512(aDataFirstElem, bDataFirstElem *uint64, n int) int
TEXT ·xorPopcntAVX512(SB), NOSPLIT, $0-32
	MOVQ aDataFirstElem+0(FP), SI
	MOVQ bDataFirstElem+8(FP), DI
	MOVQ n+16(FP), CX
	VPXORQ Z0, Z0, Z0

loop:
	CMPQ CX, $8
	JLT  tail
	VMOVDQU64 (SI), Z1
	VPXORQ (DI), Z1, Z1
	VPOPCNTQ Z1, Z1
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
	DECL  AX
	KMOVB AX, K1
	VMOVDQU64.Z (SI), K1, Z1
	VMOVDQU64.Z (DI), K1, Z2
	VPXORQ Z2, Z1, Z1
	VPOPCNTQ Z1, Z1
	VPADDQ Z1, Z0, Z0

reduce:
	HSUM_Z0_AX
	MOVQ AX, ret+24(FP)
	VZEROUPPER
	RET

// func dotAVX512(leftDataFirstElem, rightDataFirstElem *uint64, leftRows, rightRows, cols, stride int, resultsFirstElem *int)
//
// results[r*rightRows+c] = cols - popcount(left行r ^ right行c)
TEXT ·dotAVX512(SB), NOSPLIT, $0-56
	MOVQ leftDataFirstElem+0(FP), SI
	MOVQ rightDataFirstElem+8(FP), DI
	MOVQ leftRows+16(FP), R8
	MOVQ rightRows+24(FP), BX
	MOVQ cols+32(FP), R10
	MOVQ stride+40(FP), R9
	MOVQ resultsFirstElem+48(FP), DX

	MOVQ R9, CX
	ANDQ $7, CX
	MOVL $1, AX
	SHLL CX, AX
	DECL AX
	KMOVB AX, K1

leftRowLoop:
	TESTQ R8, R8
	JEQ   dotDone
	MOVQ DI, R11
	MOVQ BX, R12

rightRowLoop:
	TESTQ R12, R12
	JEQ   nextLeftRow
	MOVQ SI, R13
	MOVQ R9, CX
	VPXORQ Z0, Z0, Z0

dotWordLoop:
	CMPQ CX, $8
	JLT  dotWordTail
	VMOVDQU64 (R13), Z1
	VPXORQ (R11), Z1, Z1
	VPOPCNTQ Z1, Z1
	VPADDQ Z1, Z0, Z0
	ADDQ $64, R13
	ADDQ $64, R11
	SUBQ $8, CX
	JMP  dotWordLoop

dotWordTail:
	TESTQ CX, CX
	JEQ   dotReduce
	VMOVDQU64.Z (R13), K1, Z1
	VMOVDQU64.Z (R11), K1, Z2
	VPXORQ Z2, Z1, Z1
	VPOPCNTQ Z1, Z1
	VPADDQ Z1, Z0, Z0
	LEAQ (R11)(CX*8), R11

dotReduce:
	HSUM_Z0_AX
	MOVQ R10, R14
	SUBQ AX, R14
	MOVQ R14, (DX)
	ADDQ $8, DX
	DECQ R12
	JMP  rightRowLoop

nextLeftRow:
	LEAQ (SI)(R9*8), SI
	DECQ R8
	JMP  leftRowLoop

dotDone:
	VZEROUPPER
	RET

// func dotTernaryAVX512(valueDataFirstElem, signDataFirstElem, nonZeroDataFirstElem *uint64, valueRows, signRows, stride int, resultsFirstElem *int)
//
// results[r*signRows+c] = popcount(nonZero行c) - 2*popcount((value行r ^ sign行c) & nonZero行c)
//
// popcount(nonZero行) は value行 に依存しない為、ブロックごとに事前計算して value行 のループで使い回す。
//
// ただし内側を value行 にすると results の書き込みが signRows*8 バイト刻みになり、
// signRows が512の倍数でL1のセット競合とTLBミスを起こす
// (Zen4実測 valueRows=768, cols=768, signRows=512: 2.6 → 5.4 ns/結果)。
// その為 sign行 を8行ずつ処理し、内側で results の連続8個(=1キャッシュライン)を書き切る。
//
// valueRows=1 の場合は、求めた popcount(nonZero行) を1回しか使わない為、
// 事前に求めず1回の走査で両方求める実装より遅い。
//
// 汎用レジスタ13本に空きが無い為、8個の popcount(nonZero行) はフレームに置き、
// 各レジスタの役割には名前を付ける(AXのみ一時用)。
// 不変
#define mainBytes    SI
#define strideBytes  R9
#define nzCountBase  DI

// 走査で進む
#define restSignRows BX
#define blockSign    R12
#define blockNz      R13
#define curValue     R10
#define curRes       R11
#define curSign      R14
#define curNz        DX

// バイトオフセット
#define colOff       R8
#define byteOff      CX

// マスク生成の間だけ CX に端数ワード数が入る(可変シフト量はCLでなければならない為)。
#define tailWords    CX

// フレーム96バイト(空き無し)
//   nzCounts -96..-40 (8スロット) / valueEnd -32 / resRowBytes -24 / resColBase -16 / blockBytes -8
TEXT ·dotTernaryAVX512(SB), NOSPLIT, $96-56
	MOVQ signDataFirstElem+8(FP), blockSign
	MOVQ nonZeroDataFirstElem+16(FP), blockNz
	MOVQ signRows+32(FP), restSignRows
	MOVQ stride+40(FP), AX   // AX = stride (ワード数)

	MOVQ AX, strideBytes
	SHLQ $3, strideBytes

	MOVQ AX, mainBytes
	SHLQ $3, mainBytes
	ANDQ $-64, mainBytes     // 8ワード(=64バイト)単位に切り捨て

	MOVQ AX, tailWords
	ANDQ $7, tailWords
	MOVL $1, AX
	SHLL tailWords, AX
	DECL AX
	KMOVB AX, K1

	MOVQ valueRows+24(FP), AX
	IMULQ strideBytes, AX
	ADDQ valueDataFirstElem+0(FP), AX
	MOVQ AX, valueEnd-32(SP)

	MOVQ restSignRows, AX
	SHLQ $3, AX
	MOVQ AX, resRowBytes-24(SP)

	MOVQ resultsFirstElem+48(FP), AX
	MOVQ AX, resColBase-16(SP)

	LEAQ nzCounts-96(SP), nzCountBase

ternaryBlockLoop:
	TESTQ restSignRows, restSignRows
	JEQ   ternaryDone
	MOVQ restSignRows, AX
	CMPQ AX, $8
	JLE  ternaryBlockSizeOK
	MOVQ $8, AX

ternaryBlockSizeOK:
	SHLQ $3, AX
	MOVQ AX, blockBytes-8(SP)
	MOVQ blockNz, curNz
	XORQ colOff, colOff

ternaryNzRowLoop:
	CMPQ colOff, blockBytes-8(SP)
	JGE  ternaryNzDone
	VPXORQ Z0, Z0, Z0
	XORQ  byteOff, byteOff

ternaryNzWordLoop:
	CMPQ byteOff, mainBytes
	JGE  ternaryNzWordTail
	VPOPCNTQ (curNz)(byteOff*1), Z1
	VPADDQ Z1, Z0, Z0
	ADDQ $64, byteOff
	JMP  ternaryNzWordLoop

ternaryNzWordTail:
	CMPQ byteOff, strideBytes
	JGE  ternaryNzReduce
	VMOVDQU64.Z (curNz)(byteOff*1), K1, Z1
	VPOPCNTQ Z1, Z1
	VPADDQ Z1, Z0, Z0

ternaryNzReduce:
	HSUM_Z0_AX
	MOVQ AX, (nzCountBase)(colOff*1)
	ADDQ strideBytes, curNz
	ADDQ $8, colOff
	JMP  ternaryNzRowLoop

ternaryNzDone:
	MOVQ valueDataFirstElem+0(FP), curValue
	MOVQ resColBase-16(SP), curRes

ternaryValueRowLoop:
	CMPQ curValue, valueEnd-32(SP)
	JCC  ternaryNextBlock
	MOVQ blockSign, curSign
	MOVQ blockNz, curNz
	XORQ colOff, colOff

ternaryColLoop:
	CMPQ colOff, blockBytes-8(SP)
	JGE  ternaryNextValueRow
	VPXORQ Z0, Z0, Z0
	XORQ  byteOff, byteOff

ternaryWordLoop:
	CMPQ byteOff, mainBytes
	JGE  ternaryWordTail
	VMOVDQU64 (curValue)(byteOff*1), Z1
	VPXORQ (curSign)(byteOff*1), Z1, Z1
	VPANDQ (curNz)(byteOff*1), Z1, Z1
	VPOPCNTQ Z1, Z1
	VPADDQ Z1, Z0, Z0
	ADDQ $64, byteOff
	JMP  ternaryWordLoop

ternaryWordTail:
	CMPQ byteOff, strideBytes
	JGE  ternaryReduce
	VMOVDQU64.Z (curValue)(byteOff*1), K1, Z1
	VMOVDQU64.Z (curSign)(byteOff*1), K1, Z2
	VPXORQ Z2, Z1, Z1
	VMOVDQU64.Z (curNz)(byteOff*1), K1, Z2
	VPANDQ Z2, Z1, Z1
	VPOPCNTQ Z1, Z1
	VPADDQ Z1, Z0, Z0

ternaryReduce:
	HSUM_Z0_AX
	SHLQ $1, AX
	NEGQ AX
	ADDQ (nzCountBase)(colOff*1), AX

	// ブロック化の狙い。同一キャッシュラインへ連続して書き込む
	MOVQ AX, (curRes)(colOff*1)

	ADDQ strideBytes, curSign
	ADDQ strideBytes, curNz
	ADDQ $8, colOff
	JMP  ternaryColLoop

ternaryNextValueRow:
	ADDQ resRowBytes-24(SP), curRes
	ADDQ strideBytes, curValue
	JMP  ternaryValueRowLoop

ternaryNextBlock:
	MOVQ blockBytes-8(SP), AX
	ADDQ AX, resColBase-16(SP)
	SHRQ $3, AX              // AX = このブロックで処理した行数
	SUBQ AX, restSignRows
	IMULQ strideBytes, AX
	ADDQ AX, blockSign
	ADDQ AX, blockNz
	JMP  ternaryBlockLoop

ternaryDone:
	VZEROUPPER
	RET
