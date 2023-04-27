package omw_test

import (
	"testing"
	"fmt"
	"github.com/sw965/omw"
	"golang.org/x/exp/slices"
	"os"
)

func Test1MapFunc(t *testing.T) {
	xs := omw.MakeIntegerRange[[]int](0, 6, 1)
	f := func(x int) int {return x * x}
	result := omw.MapFunc[[]int, []int](xs, f)
	expected := []int{0, 1, 4, 9, 16, 25}
	if !slices.Equal(result, expected) {
		t.Errorf("テスト失敗")
	}
}

func Test2MapFunc(t *testing.T) {
	xs := []bool{true, false, true, false, true}
	f := func(x bool) int {
		if x {
			return 1
		} else {
			return 0
		}
	}
	result := omw.MapFunc[[]bool, []int](xs, f)
	expected := []int{1, 0, 1, 0, 1}
	if !slices.Equal(result, expected) {
		t.Errorf("テスト失敗")
	}
}

func Test1Filter(t *testing.T) {
	xs := omw.MakeIntegerRange[[]int](0, 11, 1)
	f := func(x int) bool {
		return omw.IsRemainderZero(2)(x)
	}
	result := omw.Filter[[]int](xs, f)
	expected := []int{0, 2, 4, 6, 8, 10}
	if !slices.Equal(result, expected) {
		t.Errorf("テスト失敗")
	}
}

func Test2Filter(t *testing.T) {
	xs := omw.MakeIntegerRange[[]int](0, 100, 1)
	f := func(x int) bool {
		return omw.IsRemainderZero(16)(x)
	}
	result := omw.Filter[[]int](xs, f)
	expected := []int{0, 16, 32, 48, 64, 80, 96}
	if !slices.Equal(result, expected) {
		t.Errorf("テスト失敗")
	}
}

func TestPermutationTotalNum(t *testing.T) {
	result := omw.PermutationTotalNum(10, 5)
	expected := 30240
	if result != expected {
		t.Errorf("テスト失敗")
	}
}

func TestPermutationNumberss(t *testing.T) {
	numss := omw.PermutationNumberss(5, 3)
	for _, nums := range numss {
		fmt.Println(nums)
	}
}

func TestCombinationTotalNum(t *testing.T) {
	result := omw.CombinationTotalNum(10, 5)
	expected := 252
	if result != expected {
		t.Errorf("テスト失敗")
	}
}

func TestCombinationNumberss(t *testing.T) {
	numss := omw.CombinationNumberss(5, 3)
	for _, nums := range numss {
		fmt.Println(nums)
	}
}

func TestMin(t *testing.T) {
	xs := []float64{0.0, 1.0, -0.5, 2.0, 0.01}
	result := omw.Min(xs...)
	expected := -0.5
	if result != expected {
		t.Errorf("テスト失敗")
	}
}

func TestMax(t *testing.T) {
	xs := []float64{0.0, 1.0, -0.5, 2.0, 0.01}
	result := omw.Max(xs...)
	expected := 2.0
	if result != expected {
		t.Errorf("テスト失敗")
	}
}

func TestSum(t *testing.T) {
	xs := omw.MakeIntegerRange[[]int](0, 11, 1)
	result := omw.Sum(xs...)
	expected := 55
	if result != expected {
		t.Errorf("テスト失敗")
	}
}

func TestMean(t *testing.T) {
	xs := omw.MakeIntegerRange[[]int](0, 101, 2)
	result := omw.Mean(xs...)
	expected := 50
	if result != expected {
		t.Errorf("テスト失敗")
	}
}

type Fruit string

const (
	APPLE = Fruit("りんご")
	ORANGE = Fruit("オレンジ")
	GRAPE = Fruit("グレープ")
)

type Fruits []Fruit

var ALL_FRUITS = Fruits{APPLE, ORANGE, GRAPE}

func TestPermutation(t *testing.T) {
	result := omw.Permutation[[]Fruits](ALL_FRUITS, 2)
	expected := []Fruits{
		Fruits{"りんご", "オレンジ"},
		Fruits{"りんご", "グレープ"},
		Fruits{"オレンジ", "りんご"},
		Fruits{"オレンジ", "グレープ"},
		Fruits{"グレープ", "りんご"},
		Fruits{"グレープ", "オレンジ"},
	}

	eq := func(fruits1, fruits2 Fruits) bool {
		return slices.Equal(fruits1, fruits2)
	}

	if !slices.EqualFunc(result, expected, eq) {
		t.Errorf("テスト失敗")
	}
}

func TestCombination(t *testing.T) {
	result := omw.Combination[[]Fruits](ALL_FRUITS, 2)
	expected := []Fruits{
		Fruits{APPLE, ORANGE},
		Fruits{APPLE, GRAPE},
		Fruits{ORANGE, GRAPE},
	}

	eq := func(fruits1, fruits2 Fruits) bool {
		return slices.Equal(fruits1, fruits2)
	}

	if !slices.EqualFunc(result, expected, eq) {
		t.Errorf("テスト失敗")
	}
}

type Animal string

const (
	DOG = Animal("犬")
	CAT = Animal("猫")
	BIRD = Animal("鳥")
	COW = Animal("牛")
	PIG = Animal("豚")
)

type Animals []Animal

var ALL_ANIMALS = Animals{DOG, CAT, BIRD, COW, COW, PIG}

func Test1IsSubset(t *testing.T) {
	a := Animals{CAT, BIRD, COW}
	if !omw.IsSubset(ALL_ANIMALS, a) {
		t.Errorf("テスト失敗")
	}
}

func Test2IsSubset(t *testing.T) {
	a := Animals{CAT, BIRD, COW}
	if omw.IsSubset(a, ALL_ANIMALS) {
		t.Errorf("テスト失敗")
	}
}

func TestIndicesAccess(t *testing.T) {
	result := omw.IndicesAccess[Animals](ALL_ANIMALS, []int{0, 2, 4}...)
	expected := Animals{DOG, BIRD, COW}
	if !slices.Equal(result, expected) {
		t.Errorf("テスト失敗")
	}
}

func Test1ElementCount(t *testing.T) {
	xs := []bool{true, false, true, true, true, true, false, false, true}
	result := omw.ElementCount(xs, true)
	expected := 6
	if result != expected {
		t.Errorf("テスト失敗")
	}
}

func Test2ElementCount(t *testing.T) {
	xs := []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5, 5}
	r := omw.NewMt19937()
	r.Shuffle(len(xs), func(i, j int) {
		xs[i], xs[j] = xs[j], xs[i]
	})
	result := omw.ElementCount(xs, 5)
	expected := 5
	if result != expected {
		t.Errorf("テスト失敗")
	}
}

func Test1ElementReverse(t *testing.T) {
	xs := []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}
	result := omw.ElementReverse(xs)
	expected := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	if !slices.Equal(result, expected) {
		t.Errorf("テスト失敗")
	}
}

func TestToUnique(t *testing.T) {
	xs := []int{0, 1, 0, 2, 4, 5, 2, 3, 9, 7, 5, 4, 1, 4, 8, 9}
	result := omw.ToUnique(xs)
	expected := []int{0, 1, 2, 4, 5, 3, 9, 7, 8}
	if !slices.Equal(result, expected) {
		t.Errorf("テスト失敗")
	}
}

func Test1IsUnique(t *testing.T) {
	xs := []string{"あ", "い", "う", "え", "お"}
	if !omw.IsUnique(xs) {
		t.Errorf("テスト失敗")
	}
}

func Test2IsUnique(t *testing.T) {
	xs := []string{"あ", "い", "う", "え", "お", "あ"}
	if omw.IsUnique(xs) {
		t.Errorf("テスト失敗")
	}
}

type Human struct {
	Name string
	Age int
	Gender string
	Height float64
}

func TestLoadJson(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	result, err := omw.LoadJson[Human](dir + "/test_load.json")
	if err != nil {
		panic(err)
	}
	expected := Human{
		Name:"山田 太郎",
		Age:12,
		Gender:"♂",
		Height:170.5,
	}
	if result != expected {
		t.Errorf("テスト失敗")
	}
}

func TestWriteJson(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	human := Human{
		Name:"鈴木 ウド",
		Age:50,
		Gender:"♀",
		Height:217.7,
	}

	err = omw.WriteJson(&human, dir + "/test_write.json")
	if err != nil {
		panic(err)
	}
}

func TestDirNames(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	result, err := omw.DirNames(dir)
	if err != nil {
		panic(err)
	}
	expected := []string{".git", "go.mod", "go.sum", "omw_test.go", "omw.go", "README.md", "test_load.json", "test_write.json"}

	slices.Sort(result)
	slices.Sort(expected)


	if !slices.Equal(result, expected) {
		t.Errorf("テスト失敗")
	}
}

func TestRandSample(t *testing.T) {
	r := omw.NewMt19937()
	xs := omw.MakeIntegerRange[[]int](0, 10, 1)
	testNum := 16
	for i := 0; i < testNum; i++ {
		result := omw.RandSample(xs, 5, r)
		if !omw.IsUnique(result) {
			t.Errorf("テスト失敗")
			break
		}
		fmt.Println(result)
	}
}