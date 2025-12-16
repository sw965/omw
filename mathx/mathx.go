// Package mathx provides generic utility functions for mathematical operations and aggregations.
// Package mathx は、数学的な操作や集計を行うためのユーティリティ関数を提供します。
package mathx

import (
	"cmp"
)

// Sum returns the sum of the provided arguments.
//
// For numeric types (integers, floats), it computes the arithmetic sum.
// For strings, it concatenates the values in order.
// If no arguments are provided, it returns the zero value of type T.
//
// Sum は、渡されたすべての引数の合計（または連結結果）を返します。
//
// 数値型（整数、浮動小数点数）の場合、算術的な和を計算します。
// 文字列の場合、順に連結します。
// 引数が指定されなかった場合、その型 T のゼロ値を返します。
func Sum[T cmp.Ordered](xs ...T) T {
	var sum T
	for _, x := range xs {
		sum += x
	}
	return sum
}