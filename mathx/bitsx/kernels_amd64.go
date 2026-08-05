//go:build amd64

package bitsx

import "golang.org/x/sys/cpu"

var useAVX512 = cpu.X86.HasAVX512F && cpu.X86.HasAVX512VPOPCNTDQ

// 実装は kernels_amd64.s

// Goアセンブラで書かれた生ポインタ関数で、境界チェックを持たない。
// *uint64はポインタなので長さの情報を持たない(len()も使えない)為、
// n として別途渡している。呼び出し側は、n が a・b それぞれの元のスライスの
// 長さ(語数)以下である事を保証しなければならない。掛け算を伴わない為、
// 桁あふれの心配はない。
func xorPopcntAVX512(a, b *uint64, n int) int

// Goアセンブラで書かれた生ポインタ関数で、境界チェックを持たない。
// leftData・rightData・resultsはいずれもポインタで長さの情報を持たない為、
// 呼び出し側は、呼び出し前に以下を保証しなければならない。満たさない場合、
// 範囲外のメモリを読み書きする(validateDotAVX512Argsが検証する)。
//
//  1. leftRows, rightRows, stride がいずれも0以上である事
//  2. leftRows*stride, rightRows*stride, leftRows*rightRows のいずれも桁あふれせず、
//     それぞれ leftData, rightData, results が指す元のスライスの長さと一致する事
//  3. stride が leftData・rightData 双方の元のスライスに対して正しい値である事
//     (両者の列数(cols)が異なると、共通の stride ではどちらかの実際の
//     バッファ長と食い違う)
func dotAVX512(leftData, rightData *uint64, leftRows, rightRows, cols, stride int, results *int)

// Goアセンブラで書かれた生ポインタ関数で、境界チェックを持たない。
// valueData・signData・nonZeroData・resultsはいずれもポインタで長さの情報を
// 持たない為、呼び出し側は、呼び出し前に以下を保証しなければならない。
// 満たさない場合、範囲外のメモリを読み書きする
// (validateDotTernaryAVX512Argsが検証する)。
//
//  1. valueRows, signRows, stride がいずれも0以上である事
//  2. valueRows*stride, signRows*stride, valueRows*signRows のいずれも桁あふれせず、
//     それぞれ valueData、signData(及びnonZeroData)、results が指す元のスライスの
//     長さと一致する事
//  3. stride が valueData・signData・nonZeroData 全ての元のスライスに対して
//     正しい値である事(いずれかの列数(cols)が異なると、共通の stride では
//     実際のバッファ長と食い違う)
func dotTernaryAVX512(valueData, signData, nonZeroData *uint64, valueRows, signRows, stride int, results *int)
