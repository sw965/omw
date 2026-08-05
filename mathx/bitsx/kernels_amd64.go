//go:build amd64

package bitsx

import "golang.org/x/sys/cpu"

var useAVX512 = cpu.X86.HasAVX512F && cpu.X86.HasAVX512VPOPCNTDQ

// 実装は kernels_amd64.s
//
//  1. Goアセンブラを呼び出す場合、各関数に書かれた制約を守らなければ、範囲外の
//     メモリの読み書きをする恐れがある。
//  2. 範囲外のメモリの読み書きをしてしまうと、まったく無関係の変数が書き変わったり、
//     クラッシュしたりする恐れがある。
//  3. 各関数に書かれたコメントの制約以外に、範囲外のメモリの読み書きをしてしまう
//     可能性に気付いたら、ただちに報告する事。
//  4. 引数名が○○DataFirstElemとなっているものは、&Matrix.data[0]として渡される。

// 1. n >= 0
// 2. n <= aDataFirstElemが格納されたスライスの長さ
// 3. n <= bDataFirstElemが格納されたスライスの長さ
func xorPopcntAVX512(aDataFirstElem, bDataFirstElem *uint64, n int) int

// 1. leftRows, rightRows, stride >= 0
// 2. leftRows*stride == leftDataFirstElemが格納されたスライスの長さ
// 3. rightRows*stride == rightDataFirstElemが格納されたスライスの長さ
// 4. leftRows*rightRows == resultsFirstElemが格納されたスライスの長さ
func dotAVX512(leftDataFirstElem, rightDataFirstElem *uint64, leftRows, rightRows, cols, stride int, resultsFirstElem *int)

// 1. valueRows, signRows, stride >= 0
// 2. valueRows*stride == valueDataFirstElemが格納されたスライスの長さ
// 3. signRows*stride == signDataFirstElemが格納されたスライスの長さ
// 4. signRows*stride == nonZeroDataFirstElemが格納されたスライスの長さ
// 5. valueRows*signRows == resultsFirstElemが格納されたスライスの長さ
func dotTernaryAVX512(valueDataFirstElem, signDataFirstElem, nonZeroDataFirstElem *uint64, valueRows, signRows, stride int, resultsFirstElem *int)
