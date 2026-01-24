package ternary

import (
	"fmt"
	"math/bits"
	"github.com/sw965/omw/mathx/bitsx"
)

type Matrix struct {
    Sign    bitsx.Matrix
    NonZero bitsx.Matrix
}

type DotVecResult struct {
	MatchCounts   []int
	NonZeroCounts []int
}

// 後でコメントを書く
func (d DotVecResult) Zs() ([]int, error) {
	mn := len(d.MatchCounts)
	nn := len(d.NonZeroCounts)
	if mn != nn {
		return nil, fmt.Errorf("len(d.MatchCounts) != len(d.NonZeroCounts): len(d.MatchCounts) = %d, len(d.NonZeroCounts) = %d: 同じ要素数であるべき", nn, mn)
	}

	zs := make([]int, mn)
	for i, mc := range d.MatchCounts {
		nc := d.NonZeroCounts[i]
		// 後でコメントを書く
		zs[i] = 2*mc - nc
	}
	return zs, nil
}

// この関数の中のローカル変数の変更の余地あり
func (m Matrix) DotVec(v Matrix) (DotVecResult, error) {
    // 1. 形状チェック
    if v.Sign.Rows != 1 || v.NonZero.Rows != 1 {
        return DotVecResult{}, fmt.Errorf("v must be a vector (Rows=1)")
    }
    if m.Sign.Cols != v.Sign.Cols {
        return DotVecResult{}, fmt.Errorf("dimension mismatch: %d != %d", m.Sign.Cols, v.Sign.Cols)
    }

    stride := m.Sign.Stride
    rowMask := m.Sign.RowMask
	ret := DotVecResult{
		MatchCounts:   make([]int, m.Sign.Rows),
		NonZeroCounts: make([]int, m.Sign.Rows),
	}

    for r := 0; r < m.Sign.Rows; r++ {
        offset := r * stride
		matchCount := 0
		nonZeroCount := 0

        for k := 0; k < stride; k++ {
            // 符号が一致しているか (XNOR)
            sign := ^(m.Sign.Data[offset+k] ^ v.Sign.Data[k])
            
            // 両方が 0 でない場所 (AND)
            nonZero := m.NonZero.Data[offset+k] & v.NonZero.Data[k]

            // 最後のブロックならパディング部分をマスク
            if k == stride-1 {
                nonZero &= rowMask
            }

            // 符号が一致し、かつ両方が0でない場所を1にする
            valid := sign & nonZero
			matchCount += bits.OnesCount64(valid)
			nonZeroCount += bits.OnesCount64(nonZero)
        }
		ret.MatchCounts[r] = matchCount
		ret.NonZeroCounts[r] = nonZeroCount
    }
    return ret, nil
}