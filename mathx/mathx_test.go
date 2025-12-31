package mathx_test

import (
	"cmp"
	"github.com/sw965/omw/mathx"
	"testing"
)

type sumTestCase[T cmp.Ordered] struct {
	name string
	args []T
	want T
}

func runSumTests[T cmp.Ordered](t *testing.T, tests []sumTestCase[T]) {
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := mathx.Sum[T](tc.args...)
			if got != tc.want {
				t.Errorf("want: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestSum(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		runSumTests(t, []sumTestCase[int]{
			{
				name: "正常",
				args: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				want: 55,
			},
			{
				name: "準正常_空スライス",
				args: []int{},
				want: 0,
			},
		})
	})

	t.Run("float64", func(t *testing.T) {
		runSumTests(t, []sumTestCase[float64]{
			{
				name: "正常",
				args: []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0},
				want: 21.0,
			},
			{
				name: "準正常_空スライス",
				args: []float64{},
				want: 0.0,
			},
		})
	})

	t.Run("string", func(t *testing.T) {
		runSumTests(t, []sumTestCase[string]{
			{
				name: "正常",
				args: []string{"a", "b", "c", "d", "e"},
				want: "abcde",
			},
			{
				name: "準正常_空スライス",
				args: []string{},
				want: "",
			},
		})
	})
}
