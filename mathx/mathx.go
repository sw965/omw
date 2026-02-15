package mathx

import (
	"math"

	"github.com/sw965/omw/constraints"
	"slices"
	"fmt"
)

func IsInf[F constraints.Float](f F, sign int) bool {
	return math.IsInf(float64(f), sign)
}

func IsNaN[F constraints.Float](f F) bool {
	return math.IsNaN(float64(f))
}

func Abs[N constraints.Number](x N) N {
    if x < 0 {
        return -x
    }
    return x
}

func Log[F constraints.Float](x F) F {
	return F(math.Log(float64(x)))
}

func Sqrt[F constraints.Float](x F) F {
	return F(math.Sqrt(float64(x)))
}

func Sum[N constraints.Number](xs ...N) N {
	var sum N
	for _, x := range xs {
		sum += x
	}
	return sum
}

// テストはまだ 後で消す？
func Median[N constraints.Number](xs ...N) (float64, error) {
	n := len(xs)
	if n == 0 {
		return 0, fmt.Errorf("len(xs) = 0: len(xs) > 0 であるべき")
	}

	c := slices.Clone(xs)
	slices.Sort(c)

	// 3. 中央値の計算
	if n%2 != 0 {
		// 要素数が奇数の場合、真ん中の値を float64 で返す
		return float64(c[n/2]), nil
	}

	// 要素数が偶数の場合、真ん中2つの値の平均を返す
	mid1 := c[n/2-1]
	mid2 := c[n/2]
	return float64(mid1+mid2) / 2.0, nil
}