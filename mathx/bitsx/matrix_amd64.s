#include "textflag.h"

// func mulVecPopCountsAVX512Asm(mat []uint64, vec []uint64, res []int, stride int, mask uint64)
TEXT ·mulVecPopCountsAVX512Asm(SB), NOSPLIT, $0
    // ----------------------------------------------------------------
    // 引数のロード
    // ----------------------------------------------------------------
    MOVQ mat_base+0(FP), SI    // SI = mat.Data ポインタ
    MOVQ vec_base+24(FP), DX   // DX = vec.Data ポインタ
    MOVQ res_base+48(FP), DI   // DI = res スライスのポインタ
    MOVQ res_len+56(FP), R10   // R10 = rows (処理する行数)
    MOVQ stride+72(FP), R8     // R8 = stride
    MOVQ mask+80(FP), R9       // R9 = RowMask

    // 行数が0以下の場合は即座にリターン（安全策）
    TESTQ R10, R10
    JLE ret_end

    // 定数・ゼロの準備
    VPXORD Z0, Z0, Z0          // Z0 = ゼロ（クリア用）

loop_rows:
    // ----------------------------------------------------------------
    // 行ごとの初期化
    // ----------------------------------------------------------------
    // PopCount累積用レジスタをクリア
    VPXORD Z1, Z1, Z1          
    XORQ R11, R11              // R11 = 現在の列インデックス (k)
    XORQ R12, R12              // R12 = スカラー（端数）用カウンタ

loop_stride_512:
    // ----------------------------------------------------------------
    // AVX-512 ループ (8 uint64 ずつ処理)
    // ----------------------------------------------------------------
    // 【修正点】: 残りが8個ちょうどの場合も含めてシリアル処理に回す。
    // これにより、最後のブロック(Stride末尾)が常にシリアルループで処理され、
    // 確実にRowMaskが適用されるようになる。
    MOVQ R8, AX
    SUBQ R11, AX
    CMPQ AX, $8
    JLE loop_stride_serial     // JL (Less) ではなく JLE (Less Equal) に変更

    // データロード
    // mat[k...k+7] -> Z2
    VMOVDQU64 0(SI)(R11*8), Z2
    // vec[k...k+7] -> Z3 (ベクトルは行によってアドレスが変わらないのでDX固定)
    VMOVDQU64 0(DX)(R11*8), Z3

    // XNOR計算: Z2 = ~(Z2 ^ Z3)
    // VPTERNLOGQ $0x99, src1(Z3), src2(Z2), dest(Z2)
    VPTERNLOGQ $0x99, Z3, Z2, Z2 

    // ポップカウント実行
    VPOPCNTQ Z2, Z2
    // 累積加算
    VPADDQ Z2, Z1, Z1

    ADDQ $8, R11
    JMP loop_stride_512

loop_stride_serial:
    // ----------------------------------------------------------------
    // シリアル処理 & マスク処理 (残りの端数および末尾ブロック)
    // ----------------------------------------------------------------
    CMPQ R11, R8
    JGE row_done

    MOVQ 0(SI)(R11*8), AX
    MOVQ 0(DX)(R11*8), BX
    
    // XNOR
    XORQ BX, AX
    NOTQ AX

    // 【重要】最後のワードであれば RowMask を適用
    MOVQ R8, CX
    DECQ CX            // CX = Stride - 1
    CMPQ R11, CX       // 現在位置が末尾か？
    JNE skip_mask
    ANDQ R9, AX        // 末尾ならマスク適用

skip_mask:
    POPCNTQ AX, AX
    ADDQ AX, R12

    INCQ R11
    JMP loop_stride_serial

row_done:
    // ----------------------------------------------------------------
    // 行の集計
    // ----------------------------------------------------------------
    // Z1 (8x uint64) の水平加算
    // Z1: [ A | B | C | D | E | F | G | H ]
    
    // 上位256bitを下位へ加算
    VEXTRACTI64X4 $1, Z1, Y2
    VPADDQ Y2, Y1, Y1   // Y1 = [A+E | B+F | C+G | D+H]

    // 上位128bitを下位へ加算
    VEXTRACTI128 $1, Y1, X2
    VPADDQ X2, X1, X1   // X1 = [(A+E)+(C+G) | (B+F)+(D+H)]

    // X1 の上位64bitを下位へ加算
    MOVQ X1, AX         // AX = 下位64bit
    PSRLDQ $8, X1       // X1を右シフトして上位64bitをずらす
    MOVQ X1, BX         // BX = 元の上位64bit
    ADDQ BX, AX         // 合計

    // スカラー計算分(R12)を足す
    ADDQ R12, AX

    // 結果を res[r] に格納
    MOVQ AX, 0(DI)
    ADDQ $8, DI         // resポインタを進める
    
    // 行ポインタを進める (mat += stride * 8bytes)
    LEAQ (SI)(R8*8), SI
    
    DECQ R10
    JNZ loop_rows

ret_end:
    VZEROUPPER
    RET
