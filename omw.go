package omw

import (
	"golang.org/x/exp/constraints"
	"math/rand"
	"io/ioutil"
	"bytes"
	"encoding/json"
	"golang.org/x/exp/slices"
	"github.com/seehuhn/mt19937"
	"time"
)

func descendingConsecutiveCount[N constraints.Integer](ns ...N) int {
	y := 1
	a := ns[0] - 1
	for _, n := range ns[1:] {
		if n != a {
			break
		}
		y += 1
		a = n - 1
	}
	return y
}

func MapFunc[XS ~[]X, YS ~[]Y, X, Y any](xs XS, f func(X) Y) YS {
	ys := make(YS, len(xs))
	for i, x := range xs {
		ys[i] = f(x)
	}
	return ys
}

func Filter[XS ~[]X, X any](xs XS, f func(X) bool) XS {
	ys := make(XS, 0, len(xs))
	for _, x := range xs {
		if f(x) {
			ys = append(ys, x)
		}
	}
	return ys
}

func Identity[X any](x X) X {
	return x
}

func IsRemainderZero[N constraints.Integer](n N) func(N) bool {
	return func(a N) bool { return a%n == 0 }
}

func StrTildeToStrTilde[X, Y ~string](x X) Y {
	return Y(x)
}

func LoadJson[T any](path string) (T, error) {
	var y T
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return y, err
	}
	file = bytes.TrimPrefix(file, []byte("\xef\xbb\xbf"))
	if err := json.Unmarshal(file, &y); err != nil {
		return y, err
	}
	return y, nil
}

func WriteJson[T any](y *T, path string) error {
	file, err := json.MarshalIndent(y, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, file, 0644)
	return err
}

func PermutationTotalNum(n, r int) int {
	y := 1
	for i := 0; i < r; i++ {
		y *= (n - i)
	}
	return y
}

func PermutationNumberss(n, r int) [][]int {
	yn := PermutationTotalNum(n, r)
	y := make([][]int, 0, yn)
	if r == 0 {
		return y
	}
	var f func(int, []int)
	f = func(nest int, nums []int) {
		if nest == r {
			y = append(y, nums)
			return
		}
		for i := 0; i < n; i++ {
			isContinue := false
			for _, num := range nums {
				if num == i {
					isContinue = true
					break
				}
			}
			if isContinue {
				continue
			}
			clone := slices.Clone(nums)
			f(nest+1, append(clone, i))
		}
	}
	f(0, make([]int, 0, r))
	return y
}

func CombinationTotalNum(n, r int) int {
	a := 1
	for i := 0; i < r; i++ {
		a *= (n - i)
	}

	m := 1
	for i := 0; i < r; i++ {
		m *= (r - i)
	}
	return a / m
}

func CombinationNumberss(n, r int) [][]int {
	nums := make([]int, r)
	for i := 0; i < r; i++ {
		nums[i] = i
	}

	yn := CombinationTotalNum(n, r)
	y := make([][]int, 0, yn)
	if r == 0 {
		return y
	}

	end := r - 1
	for i := 0; i < yn; i++ {
		clone := slices.Clone(nums)
		y = append(y, clone)
		max := Max(nums...)
		if max == (n - 1) {
			reverse := ElementReverse(nums)
			count := descendingConsecutiveCount(reverse...)
			idx := end - count
			if idx < 0 {
				break
			}
			nums[idx] += 1
			for j := idx + 1; j < r; j++ {
				nums[j] = nums[idx] + j - (idx)
			}
		} else {
			nums[end] += 1
		}
	}
	return y
}

func Min[X constraints.Ordered](xs ...X) X {
	y := xs[0]
	for _, x := range xs[1:] {
		if x < y {
			y = x
		}
	}
	return y
}

func Max[X constraints.Ordered](xs ...X) X {
	y := xs[0]
	for _, x := range xs[1:] {
		if x > y {
			y = x
		}
	}
	return y
}

func Sum[X constraints.Ordered](xs ...X) X {
	y := xs[0]
	for _, x := range xs[1:] {
		y += x
	}
	return y
}

func Mean[X constraints.Integer | constraints.Float](xs ...X) X {
	return Sum(xs...) / X(len(xs))
}

func RandBool(r *rand.Rand) bool {
	return r.Intn(2) == 0
}

func RandFloat64(min, max float64, r *rand.Rand) float64 {
	return r.Float64()*(max-min) + min
}

func RandInt(start, end int, r *rand.Rand) int {
	return r.Intn(end-start) + start
}

func RandIntWithWeight(ws []float64, r *rand.Rand) int {
	sum := Sum(ws...)
	if sum == 0.0 {
		return r.Intn(len(ws))
	}

	threshold := RandFloat64(0.0, sum, r)
	accum := 0.0
	for i, w := range ws {
		accum += w
		if accum >= threshold {
			return i
		}
	}
	return len(ws) - 1
}

func NewMt19937() *rand.Rand {
	y := rand.New(mt19937.New())
	y.Seed(time.Now().UnixNano())
	return y
}

func RandChoice[XS ~[]X, X any](xs XS, r *rand.Rand) X {
	idx := r.Intn(len(xs))
	return xs[idx]
}

func MakeSliceFunc[XS ~[]X, X any](n int, f func(int) X) XS {
	y := make(XS, n)
	for i := 0; i < n; i++ {
		y[i] = f(i)
	}
	return y
}

func MakeIntegerRange[XS ~[]X, X constraints.Integer](start, end, step X) XS {
	n := int((end - 1 - start) / step) + 1
	y := make(XS, n)
	for i := 0; i < n; i++ {
		y[i] = start + (step * X(i))
	}
	return y
}

func Permutation[XSS ~[]XS, XS ~[]X, X any](xs XS, r int) XSS {
	n := len(xs)
	numss := PermutationNumberss(n, r)
	access := func(nums []int) XS {return IndicesAccess(xs, nums...) }
	return MapFunc[[][]int, XSS](numss, access)
}

func Combination[XSS ~[]XS, XS ~[]X, X any](xs XS, r int) XSS {
	n := len(xs)
	numss := CombinationNumberss(n, r)
	access := func(nums []int) XS { return IndicesAccess(xs, nums...) }
	return MapFunc[[][]int, XSS](numss, access)
}

func IsSubset[XS ~[]X, X comparable](xs, subs XS) bool {
	for _, sub := range subs {
		if !slices.Contains(xs, sub) {
			return false
		}
	}
	return true
}

func IndicesAccess[XS ~[]X, X any](xs XS , indices ...int) XS {
	y := make(XS, len(indices))
	for i, idx := range indices {
		y[i] = xs[idx]
	}
	return y
}

func ElementCount[XS ~[]X, X comparable](xs XS, a X) int {
	y := 0
	for _, x := range xs {
		if x == a {
			y += 1
		}
	}
	return y
}

func ElementReverse[XS ~[]X, X any](xs XS) XS {
	n := len(xs)
	y := make(XS, 0, n)
	for i := n - 1; i > -1; i-- {
		y = append(y, xs[i])
	}
	return y
}

func ElementIndices[XS ~[]X, X comparable](xs XS, a X) []int {
	y := make([]int, 0, len(xs))
	for i, x := range xs {
		if x == a {
			y = append(y, i)
		}
	}
	return y
}

func ToUnique[XS ~[]X, X comparable](xs XS) XS {
	y := make(XS, 0, len(xs))
	for _, x := range xs {
		if !slices.Contains(y, x) {
			y = append(y, x)
		}
	}
	return y
}

func IsUnique[XS ~[]X, X comparable](xs XS) bool {
	for _, x := range xs {
		if ElementCount(xs, x) != 1 {
			return false
		}
	}
	return true
}

func DirNames(path string) ([]string, error) {
	dirs, err := ioutil.ReadDir(path)
	if err != nil {
		return []string{}, err
	}
	y := make([]string, len(dirs))
	for i, dir := range dirs {
		y[i] = dir.Name()
	}
	return y, nil
}