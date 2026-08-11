package mathx_test

import (
	"math"
	"testing"

	"github.com/sw965/omw/constraints"
	"github.com/sw965/omw/mathx"
)

type mulOverflowCheckedCase[T constraints.Signed] struct {
	name   string
	a, b   T
	want   T
	wantOK bool
}

func runMulOverflowChecked[T constraints.Signed](t *testing.T, cases []mulOverflowCheckedCase[T]) {
	t.Helper()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, ok := mathx.MulOverflowChecked(c.a, c.b)
			if ok != c.wantOK {
				t.Errorf("okの不一致: got = %t want = %t", ok, c.wantOK)
			}
			if c.wantOK && got != c.want {
				t.Errorf("値の不一致: got = %d want = %d", got, c.want)
			}
		})
	}
}

func TestMulOverflowChecked(t *testing.T) {
	t.Run("int8", func(t *testing.T) {
		cases := []mulOverflowCheckedCase[int8]{
			// オーバーフローがギリギリ起きない
			{
				name:   "ギリギリ起きない_正_63*2",
				a:      63,
				b:      2,
				want:   126,
				wantOK: true,
			},
			{
				name:   "ギリギリ起きない_負_-128*1",
				a:      -128,
				b:      1,
				want:   -128,
				wantOK: true,
			},

			// オーバーフローがギリギリ起きる
			{
				name:   "ギリギリ起きる_正の限界値超過_64*2",
				a:      64,
				b:      2,
				wantOK: false,
			},
			{
				name:   "ギリギリ起きる_正の限界値超過_-128*-1",
				a:      -128,
				b:      -1,
				wantOK: false,
			},
			{
				name:   "ギリギリ起きる_負の限界値超過_-65*2",
				a:      -65,
				b:      2,
				wantOK: false,
			},

			// オーバーフローが起きるが、一周回って、同じ値に戻ってくる
			{
				name:   "周回_同じ値に戻る_正の限界値超過_64*5",
				a:      64,
				b:      5,
				wantOK: false,
			},
			{
				name:   "周回_同じ値に戻る_負の限界値超過_-32*9",
				a:      -32,
				b:      9,
				wantOK: false,
			},
		}
		runMulOverflowChecked(t, cases)
	})

	t.Run("int", func(t *testing.T) {
		cases := []mulOverflowCheckedCase[int]{
			{
				name:   "通常_0*0",
				a:      0,
				b:      0,
				want:   0,
				wantOK: true,
			},
			{
				name:   "通常_3*5",
				a:      3,
				b:      5,
				want:   15,
				wantOK: true,
			},

			// オーバーフローがギリギリ起きない
			{
				name:   "ギリギリ起きない_正_MaxInt*1",
				a:      math.MaxInt,
				b:      1,
				want:   math.MaxInt,
				wantOK: true,
			},
			{
				name:   "ギリギリ起きない_正_-1*(MinInt+1)",
				a:      -1,
				b:      math.MinInt + 1,
				want:   math.MaxInt,
				wantOK: true,
			},
			{
				name:   "ギリギリ起きない_負_MinInt*1",
				a:      math.MinInt,
				b:      1,
				want:   math.MinInt,
				wantOK: true,
			},
			{
				name:   "ギリギリ起きない_負_MaxInt*-1",
				a:      math.MaxInt,
				b:      -1,
				want:   math.MinInt + 1,
				wantOK: true,
			},

			// オーバーフローがギリギリ起きる
			{
				name:   "ギリギリ起きる_正の限界値超過_MinInt*-1",
				a:      math.MinInt,
				b:      -1,
				wantOK: false,
			},
			{
				name:   "ギリギリ起きる_正の限界値超過_(MaxInt/2+1)*2",
				a:      math.MaxInt/2 + 1,
				b:      2,
				wantOK: false,
			},
			{
				name:   "ギリギリ起きる_負の限界値超過_MinInt*2",
				a:      math.MinInt,
				b:      2,
				wantOK: false,
			},
			{
				name:   "ギリギリ起きる_負の限界値超過_MaxInt*-2",
				a:      math.MaxInt,
				b:      -2,
				wantOK: false,
			},

			// 3. オーバーフローが起きるが、一周回って、同じ値に戻ってくる
			{
				name:   "周回_同じ値に戻る_正の限界値超過_(1<<62)*5",
				a:      1 << 62,
				b:      5,
				wantOK: false,
			},
			{
				name:   "周回_同じ値に戻る_負の限界値超過_-(1<<62)*5",
				a:      -(1 << 62),
				b:      5,
				wantOK: false,
			},
		}
		runMulOverflowChecked(t, cases)
	})
}

// 総当たりで、より広い型での計算結果と一致することを確認する。
func TestMatchesWiderType(t *testing.T) {
	// aとbはint型
	for a := math.MinInt8; a <= math.MaxInt8; a++ {
		for b := math.MinInt8; b <= math.MaxInt8; b++ {
			a8, b8 := int8(a), int8(b)
			got, ok := mathx.MulOverflowChecked(a8, b8)

			wantMul := a * b
			wantOK := wantMul >= math.MinInt8 && wantMul <= math.MaxInt8

			if ok != wantOK {
				t.Fatalf("okの不一致: got = %t, want = %t", ok, wantOK)
			}
			if wantOK && int(got) != wantMul {
				t.Fatalf("値の不一致: got = %d, want = %d", got, wantMul)
			}
		}
	}
}
