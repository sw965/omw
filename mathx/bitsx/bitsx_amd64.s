//go:build amd64
#include "textflag.h"

// func dotAVX512(mData, oData []uint64, out []int, mRows, oRows, stride int, mask uint64)
TEXT ·dotAVX512(SB), NOSPLIT, $0-104
	// 引数のロード
	MOVQ	mData_base+0(FP), AX   // mData.Data
	MOVQ	oData_base+24(FP), BX  // oData.Data
	MOVQ	out_base+48(FP), CX    // out.Data
	
	MOVQ	mRows+72(FP), SI
	MOVQ	oRows+80(FP), DI
	MOVQ	stride+88(FP), R12     // R12: stride(int) -> 後でアキュムレータとして再利用
	MOVQ	mask+96(FP), R8

	// ストライドをバイト単位に変換 (stride * 8)
	MOVQ	R12, DX
	SHLQ	$3, DX

	// strideが1未満(0以下)なら何もしない
	TESTQ	R12, R12
	JLE	ret_end

	// loop_limit = (stride - 1) * 8
	MOVQ	DX, R13
	SUBQ	$8, R13 

	// 外側のループ 1: Matrix A の各行 (i = 0 to mRows)
	XORQ	R14, R14 // i counter
loop_rows_m:
	CMPQ	R14, SI
	JGE	ret_end

	// Matrix A の現在の行の先頭アドレスを計算
	MOVQ	AX, R9 

	// 外側のループ 2: Matrix B の各行 (j = 0 to oRows)
	MOVQ	BX, R10      // Reset B pointer to base
	XORQ	R15, R15     // j counter

loop_rows_o:
	CMPQ	R15, DI
	JGE	next_row_m

	// --- 内積計算 (Dot Product) ---
	VPXORD	Z0, Z0, Z0   // Accumulator = 0
	XORQ	R11, R11     // k (byte offset) = 0

	// 1. AVX-512 Vector Loop
vector_loop:
	// ベクトルループ内での制限チェック
	MOVQ	R13, BX      // temp limit
	SUBQ	R11, BX      // remaining bytes
	CMPQ	BX, $64
	JL	scalar_loop_setup

	VMOVDQU64 (R9)(R11*1), Z1
	VMOVDQU64 (R10)(R11*1), Z2
	VPTERNLOGQ $0x99, Z2, Z1, Z1 // XNOR
	VPOPCNTQ Z1, Z1
	VPADDQ	Z1, Z0, Z0

	ADDQ	$64, R11
	JMP	vector_loop

scalar_loop_setup:
	// 水平加算 (ZMM -> Register)
	// Z0 (512bit) の上位256bitを Y1 に抽出
	VEXTRACTI64X4 $1, Z0, Y1
	// 下位256bit(Y0) と Y1 を足す -> 結果は Y0 に入る
	VPADDQ	Y1, Y0, Y0
	
	// Y0 (256bit) の上位128bitを X1 に抽出
	VEXTRACTI128 $1, Y0, X1
	// 下位128bit(X0) と X1 を足す -> 結果は X0 に入る
	VPADDQ	X1, X0, X0
	
	// X0 (128bit) の中の上位64bitと下位64bitを足す
	VPSHUFD	$0xEE, X0, X1 // 上位を X1 の下位にコピー
	VPADDQ	X1, X0, X0
	
	// 結果(下位64bit)を R12 に移動
	MOVQ	X0, R12 

	// 2. Scalar Loop (stride-1 まで)
scalar_loop:
	CMPQ	R11, R13
	JGE	final_block

	MOVQ	(R9)(R11*1), BX
	MOVQ	(R10)(R11*1), BP
	XORQ	BP, BX
	NOTQ	BX
	POPCNTQ	BX, BX
	ADDQ	BX, R12

	ADDQ	$8, R11
	JMP	scalar_loop

	// 3. Final Block (Mask処理)
final_block:
	MOVQ	(R9)(R11*1), BX
	MOVQ	(R10)(R11*1), BP
	XORQ	BP, BX
	NOTQ	BX
	ANDQ	R8, BX       // Apply Mask
	POPCNTQ	BX, BX
	ADDQ	BX, R12

	// 結果保存
	MOVQ	R12, (CX)
	ADDQ	$8, CX

	// Next B Row
	MOVQ	oData_base+24(FP), BX // Restore Base B
	ADDQ	DX, R10
	INCQ	R15
	JMP	loop_rows_o

next_row_m:
	// Next A Row
	ADDQ	DX, AX
	INCQ	R14
	JMP	loop_rows_m

ret_end:
	RET
