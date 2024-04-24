package omw

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/seehuhn/mt19937"
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

const (
	JSON_EXTENSION = ".json"
)

var SW965_PATH = os.Getenv("GOPATH") + "sw965/"

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

func CountConsecutiveDecrease[X constraints.Integer](xs ...X) int {
	y := 1
	a := xs[0] - 1
	for _, x := range xs[1:] {
		if x != a {
			break
		}
		y += 1
		a = x - 1
	}
	return y
}

func PermutationCount(n, r int) int {
	y := 1
	for i := 0; i < r; i++ {
		y *= (n - i)
	}
	return y
}

func IntPermutations(n, r int) [][]int {
	yn := PermutationCount(n, r)
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

func CombinationCount(n, r int) int {
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

func IntCombinations(n, r int) [][]int {
	nums := make([]int, r)
	for i := 0; i < r; i++ {
		nums[i] = i
	}

	yn := CombinationCount(n, r)
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
			reversed := Reverse(nums)
			count := CountConsecutiveDecrease(reversed...)
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

func MapFunc[YS ~[]Y, XS ~[]X, X, Y any](xs XS, f func(X) Y) YS {
	ys := make(YS, len(xs))
	for i, x := range xs {
		ys[i] = f(x)
	}
	return ys
}

func MapFuncWithError[YS ~[]Y, XS ~[]X, X, Y any](xs XS, f func(X) (Y, error)) (YS, error) {
	ys := make(YS, len(xs))
	for i, x := range xs {
		y, err := f(x)
		if err != nil {
			return ys, err
		}
		ys[i] = y
	}
	return ys, nil
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

func ConvertToInt[X, Y ~int](x X) Y {
	return Y(x)
}

func ConvertToStr[X, Y ~string](x X) Y {
	return Y(x)
}

func NewMt19937() *rand.Rand {
	r := rand.New(mt19937.New())
	r.Seed(time.Now().UnixNano())
	return r
}

func RandBool(r *rand.Rand) bool {
	return r.Intn(2) == 0
}

func RandFloat64(min, max float64, r *rand.Rand) float64 {
	return r.Float64()*(max-min) + min
}

func RandInt(min, max int, r *rand.Rand) int {
	return r.Intn(max-min) + min
}

func RandIntByWeight(ws []float64, r *rand.Rand) int {
	sum := Sum(ws...)
	if sum == 0.0 {
		return r.Intn(len(ws))
	}

	threshold := RandFloat64(0.0, sum, r)
	total := 0.0
	for i, w := range ws {
		total += w
		if total >= threshold {
			return i
		}
	}
	return len(ws) - 1
}

func RandChoice[XS ~[]X, X any](xs XS, r *rand.Rand) X {
	idx := r.Intn(len(xs))
	return xs[idx]
}

func RandSample[XS ~[]X, X any](n int, xs XS, r *rand.Rand) XS {
	y := make(XS, n)
	for i := range y {
		y[i] = RandChoice(xs, r)
	}
	return y
}

func ShuffleSlice[XS ~[]X, X any](xs XS, r *rand.Rand) {
	r.Shuffle(len(xs), func(i, j int) { xs[i], xs[j] = xs[j], xs[i] })
}

func MakeRangeInteger[IS ~[]I, I constraints.Integer](start, end I) IS {
	n := end - start
	y := make(IS, int(n))
	for i := I(0); i < n; i++ {
		y[i] = i
	}
	return y
}

func Reverse[XS ~[]X, X any](xs XS) XS {
	n := len(xs)
	y := make(XS, 0, n)
	for i := n - 1; i > -1; i-- {
		y = append(y, xs[i])
	}
	return y
}

func Concat[XS ~[]X, X any](xs1 XS, xs2 XS) XS {
	y := make(XS, 0, len(xs1) + len(xs2))
	for _, x1 := range xs1 {
		y = append(y, x1)
	}

	for _, x2 := range xs2 {
		y = append(y, x2)
	}
	return y
}

func CountOccurrences[XS ~[]X, X comparable](xs XS, a X) int {
	y := 0
	for _, x := range xs {
		if x == a {
			y += 1
		}
	}
	return y
}

func CountIf[XS ~[]X, X any](xs XS, f func(x X) bool) int {
	y := 0
	for _, x := range xs {
		if f(x) {
			y += 1
		}
	}
	return y
}

func MinIndices[XS ~[]X, X constraints.Ordered](xs XS) []int {
	min := Min(xs...)
	idxs := make([]int, 0, len(xs))
	for i, x := range xs {
		if x == min {
			idxs = append(idxs, i)
		}
	}
	return idxs
}

func MaxIndices[XS ~[]X, X constraints.Ordered](xs XS) []int {
	max := Max(xs...)
	idxs := make([]int, 0, len(xs))
	for i, x := range xs {
		if x == max {
			idxs = append(idxs, i)
		}
	}
	return idxs
}

func ElementAt[XS ~[]X, X any](xs XS) func(int) X {
	return func(idx int) X {
		return xs[idx]
	}
}

func ElementsAtIndices[XS ~[]X, X any](xs XS) func([]int) XS {
	return func(idxs []int) XS {
		y := make(XS, len(idxs))
		for i, idx := range idxs {
			y[i] = xs[idx]
		}
		return y
	}
}

func FindIndicesOf[XS ~[]X, X comparable](xs XS, e X) []int {
	y := make([]int, 0, len(xs))
	for i, x := range xs {
		if x == e {
			y = append(y, i)
		}
	}
	return y
}

func FindIndicesWhere[XS ~[]X, X any](xs XS, f func(X) bool) []int {
	y := make([]int, 0, len(xs))
	for i, x := range xs {
		if f(x) {
			y = append(y, i)
		}
	}
	return y
}

func Deduplicate[XS ~[]X, X comparable](xs XS) XS {
	y := make(XS, 0, len(xs))
	for _, x := range xs {
		if !slices.Contains(y, x) {
			y = append(y, x)
		}
	}
	return y
}

func IsAllUnique[XS ~[]X, X comparable](xs XS) bool {
	for _, x := range xs {
		if CountOccurrences(xs, x) != 1 {
			return false
		}
	}
	return true
}

func IsSubset[XS ~[]X, X comparable](xs, subs XS) bool {
	for _, sub := range subs {
		if !slices.Contains(xs, sub) {
			return false
		}
	}
	return true
}

func SlicesPermutations[XSS ~[]XS, XS ~[]X, X any](xs XS, r int) XSS {
	n := len(xs)
	idxss := IntPermutations(n, r)
	return MapFunc[XSS](idxss, ElementsAtIndices(xs))
}

func SlicesCombinations[XSS ~[]XS, XS ~[]X, X any](xs XS, r int) XSS {
	n := len(xs)
	idxss := IntCombinations(n, r)
	return MapFunc[XSS](idxss, ElementsAtIndices(xs))
}

func Any(bs []bool) bool {
	for _, b := range bs {
		if b {
			return true
		}
	}
	return false
}

func AnyMatch[XS ~[]X, X any](xs XS, f func(X) bool) bool {
	for _, x := range xs {
		if f(x) {
			return true
		}
	}
	return false
}

func All(bs []bool) bool {
	for _, b := range bs {
		if !b {
			return false
		}
	}
	return true
}

func AllMatch[XS ~[]X, X any](xs XS, f func(X) bool) bool {
	for _, x := range xs {
		if !f(x) {
			return false
		}
	}
	return true
}

func MapKeys[KS ~[]K, M ~map[K]V, K comparable, V any](m M) KS {
	ks := make(KS, 0, len(m))
	for k, _ := range m {
		ks = append(ks, k)
	}
	return ks
}

func MapValues[VS ~[]V, M ~map[K]V, K comparable, V any](m M) VS {
	vs := make(VS, 0, len(m))
	for _, v := range m {
		vs = append(vs, v)
	}
	return vs
}

func InvertMap[YM ~map[V]K, XM ~map[K]V, K, V comparable](xm XM) YM {
	ym := YM{}
	for k, v := range xm {
		ym[v] = k
	}
	return ym
}

type DirEntries []os.DirEntry

func NewDirEntries(path string) (DirEntries, error) {
	return os.ReadDir(path)
}

func (es DirEntries) Names() []string {
	y := make([]string, len(es))
	for i, e := range es {
		y[i] = e.Name()
	}
	return y
}

func LoadJSON[T any](path string) (T, error) {
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

func WriteJSON[T any](y *T, path string) error {
	file, err := json.MarshalIndent(y, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, file, 0644)
	return err
}

func AllocateData(parallel int, totalData int) []int {
	allocations := make([]int, parallel)
	baseCount := totalData / parallel
	remainder := totalData % parallel

	for i := 0; i < parallel; i++ {
		if i < remainder {
			allocations[i] = baseCount + 1
		} else {
			allocations[i] = baseCount
		}
	}
	return allocations
}

func splitData(data []int, allocations []int) [][]int {
	var result [][]int
	start := 0
	for _, count := range allocations {
		if start+count > len(data) {
			count = len(data) - start // データ範囲を超えないように調整
		}
		end := start + count
		result = append(result, data[start:end])
		start = end
	}
	return result
}