package ternary

import (
	"fmt"
	"math/bits"
	"github.com/sw965/omw/mathx/bitsx"
	"math/rand/v2"
)

type Matrix struct {
	// スイッチがオンの場合の情報
    Sign    bitsx.Matrix
	// オンオフのスイッチ
    NonZero bitsx.Matrix
}

func NewZerosMatrix(rows, cols int) (Matrix, error) {
	sign, err := bitsx.NewZerosMatrix(rows, cols)
	if err != nil {
		return Matrix{}, err
	}

	nonZero, err := bitsx.NewZerosMatrix(rows, cols)
	if err != nil {
		return Matrix{}, err
	}

	return Matrix{
		Sign:sign,
		NonZero:nonZero,
	}, nil
}

func NewOnesMatrix(rows, cols int) (Matrix, error) {
	sign, err := bitsx.NewOnesMatrix(rows, cols)
	if err != nil {
		return Matrix{}, err
	}

	nonZero, err := bitsx.NewOnesMatrix(rows, cols)
	if err != nil {
		return Matrix{}, err
	}

	return Matrix{
		Sign:sign,
		NonZero:nonZero,
	}, nil
}

func NewRandMatrix(rows, cols int, kSign, kNonZero int, rng *rand.Rand) (Matrix, error) {
	sign, err := bitsx.NewRandMatrix(rows, cols, kSign, rng)
	if err != nil {
		return Matrix{}, err
	}

	nonZero, err := bitsx.NewRandMatrix(rows, cols, kNonZero, rng)
	if err != nil {
		return Matrix{}, err
	}

	return Matrix{
		Sign:    sign,
		NonZero: nonZero,
	}, nil
}

func (m *Matrix) SetZero(r, c int) error {
	if err := m.NonZero.Clear(r, c); err != nil {
		return err
	}
	return m.Sign.Clear(r, c)
}

func (m *Matrix) SetPlus(r, c int) error {
	if err := m.NonZero.Set(r, c); err != nil {
		return err
	}
	return m.Sign.Set(r, c)
}

func (m *Matrix) SetMinus(r, c int) error {
	if err := m.NonZero.Set(r, c); err != nil {
		return err
	}
	return m.Sign.Clear(r, c)
}

func (m Matrix) Dot(other Matrix) (DotResult, error) {
	if m.Sign.Cols != other.Sign.Cols {
		return DotResult{}, fmt.Errorf("dimension mismatch: m.Cols %d != other.Cols %d", m.Sign.Cols, other.Sign.Cols)
	}
	if m.Sign.Stride != other.Sign.Stride {
		return DotResult{}, fmt.Errorf("stride mismatch")
	}

	if m.NonZero.Stride != other.NonZero.Stride {
		return DotResult{}, fmt.Errorf("nonzero stride mismatch")
	}

	outRows := m.Sign.Rows
	outCols := other.Sign.Rows // otherは転置されている（列ベクトルが並んでいる）前提
	size := outRows * outCols

	res := DotResult{
		Rows:          outRows,
		Cols:          outCols,
		MatchCounts:   make([]int, size),
		NonZeroCounts: make([]int, size),
	}

	// ポインタと定数のキャッシュ
	mSignData := m.Sign.Data
	mNzData := m.NonZero.Data
	oSignData := other.Sign.Data
	oNzData := other.NonZero.Data
	stride := m.Sign.Stride
	mask := m.Sign.RowMask

	// 2. 行列積のループ
	for r := 0; r < outRows; r++ {
		// m (左側) の行オフセット
		mRowOffset := r * stride
		// 結果配列 (res) の行オフセット
		resRowOffset := r * outCols

		for c := 0; c < outCols; c++ {
			// other (右側) の行オフセット
			oRowOffset := c * stride

			matchCount := 0
			nzCount := 0

			for k := 0; k < stride; k++ {
				// データをロード
				ms := mSignData[mRowOffset+k]
				mn := mNzData[mRowOffset+k]
				os := oSignData[oRowOffset+k]
				on := oNzData[oRowOffset+k]

				// A. 双方が非ゼロ(有効)であるビット
				// 0 * 1 = 0, 0 * 0 = 0, 1 * 1 = 1 (有効)
				commonNonZero := mn & on

				// B. 符号が一致しているビット (XNOR)
				// 1(Pos) ^ 1(Pos) = 0 -> ^0 = 1
				sameSign := ^(ms ^ os)

				// 最後のブロックのみマスク処理
				if k == stride-1 {
					// 有効ビット判定の方にマスクを掛ければ、以降の計算も安全になる
					commonNonZero &= mask
				}

				// C. 「有効」かつ「符号一致」しているビット
				validMatch := sameSign & commonNonZero

				matchCount += bits.OnesCount64(validMatch)
				nzCount += bits.OnesCount64(commonNonZero)
			}

			// 結果の格納
			idx := resRowOffset + c
			res.MatchCounts[idx] = matchCount
			res.NonZeroCounts[idx] = nzCount
		}
	}

	return res, nil
}

func (m Matrix) DotTernary(otherSign, otherNonZero Matrix) ([]int, error) {
	if m.Cols != otherSign.Cols {
		return nil, fmt.Errorf("dimension mismatch: m.Cols %d != otherSign.Cols %d", m.Cols, otherSign.Cols)
	}

	if otherSign.Cols != otherNonZero.Cols || otherSign.Rows != otherNonZero.Rows {
		return nil, fmt.Errorf("otherSign and otherNonZero dimension mismatch")
	}
	
	if m.Stride != otherSign.Stride || otherSign.Stride != otherNonZero.Stride {
		return nil, fmt.Errorf("stride mismatch")
	}
	if m.RowMask != otherSign.RowMask || otherSign.RowMask != otherNonZero.RowMask {
		return nil, fmt.Errorf("mask mismatch")
	}

	outRows := m.Rows
	outCols := otherSign.Rows // otherは転置されている前提
	counts := make([]int, outRows*outCols)

	mData := m.Data
	sData := otherSign.Data
	nData := otherNonZero.Data
	stride := m.Stride
	mask := m.RowMask

	for r := 0; r < outRows; r++ {
		mRowOffset := r * stride
		resRowOffset := r * outCols

		for c := 0; c < outCols; c++ {
			oRowOffset := c * stride

			matchCount := 0   // 符号が一致し、かつNonZeroであるビット数
			nzCount := 0      // NonZeroであるビット数

			for k := 0; k < stride; k++ {
				// データをロード
				mWord := mData[mRowOffset+k]
				sWord := sData[oRowOffset+k]
				nWord := nData[oRowOffset+k]

				// A. 符号が一致しているビット (XNOR)
				// 1(Pos) ^ 1(Pos) = 0 -> ^0 = 1 (Match)
				// 0(Neg) ^ 0(Neg) = 0 -> ^0 = 1 (Match)
				sameSign := ^(mWord ^ sWord)

				// B. Ternary側が非ゼロ(有効)であるビット
				// nWord そのもの

				// C. 有効かつ符号一致
				validMatch := sameSign & nWord

				// 最後のブロックのみマスク処理
				if k == stride-1 {
					validMatch &= mask
					nWord &= mask
				}

				matchCount += bits.OnesCount64(validMatch)
				nzCount += bits.OnesCount64(nWord)
			}

			// 結果の計算
			// z = Match - Mismatch
			// Mismatch = TotalNonZero - Match
			// z = Match - (TotalNonZero - Match)
			// z = 2 * Match - TotalNonZero
			counts[resRowOffset+c] = 2*matchCount - nzCount
		}
	}
	return counts, nil
}

func (m Matrix) Transpose() (Matrix, error) {
	signT, err := m.Sign.Transpose()
	if err != nil {
		return Matrix{}, err
	}

	nonZeroT, err := m.NonZero.Transpose()
	if err != nil {
		return Matrix{}, err
	}

	return Matrix{
		Sign:signT,
		NonZero:nonZeroT,
	}, nil
}

type DotResult struct {
	Rows int
	Cols int
	MatchCounts   []int
	NonZeroCounts []int
}

// 後でコメントを書く
func (d DotResult) Zs() ([]int, error) {
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