package bitsx_test

import (
	"slices"
	"testing"

	"github.com/sw965/omw/constraints"
	"github.com/sw965/omw/mathx/bitsx"
)

type fromIndicesCase[B constraints.Unsigned] struct {
	name    string
	idxs    []int
	want    B
	wantErr bool
}

func runFromIndices[B constraints.Unsigned](t *testing.T, cases []fromIndicesCase[B]) {
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := bitsx.FromIndices[B](c.idxs)

			if err != nil && !c.wantErr {
				t.Errorf("予期せぬエラー: err = %v", err)
			}

			if err == nil && c.wantErr {
				t.Errorf("想定外の非エラー")
			}

			if got != c.want {
				t.Errorf("値の不一致: got = %b, want = %b", got, c.want)
			}
		})
	}
}

func TestFromIndices(t *testing.T) {
	t.Run("uint8", func(t *testing.T) {
		cases := []fromIndicesCase[uint8]{
			{
				name: "正常 通常",
				idxs: []int{2, 5},
				want: 0b0010_0100,
			},

			{
				name: "正常 境界 下限",
				idxs: []int{0},
				want: 0b0000_0001,
			},

			{
				name: "正常 境界 上限",
				idxs: []int{7},
				want: 0b1000_0000,
			},

			{
				name: "正常 空値",
				idxs: []int{},
				want: 0,
			},

			{
				name: "正常 nil",
				idxs: nil,
				want: 0,
			},

			{
				name:    "異常 境界 下限未満",
				idxs:    []int{-1},
				wantErr: true,
			},

			{
				name:    "異常 境界 上限超過",
				idxs:    []int{8},
				wantErr: true,
			},
		}

		runFromIndices(t, cases)
	})

	t.Run("uint64", func(t *testing.T) {
		cases := []fromIndicesCase[uint64]{
			{
				name: "正常 通常",
				idxs: []int{8, 16, 32, 48},
				want: 0b00000000_00000001_00000000_00000001_00000000_00000001_00000001_00000000,
			},

			{
				name: "正常 境界 下限",
				idxs: []int{0},
				want: 0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_00000001,
			},

			{
				name: "正常 境界 上限",
				idxs: []int{63},
				want: 0b10000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000,
			},

			{
				name: "正常 空値",
				idxs: []int{},
				want: 0,
			},

			{
				name: "正常 nil",
				idxs: nil,
				want: 0,
			},

			{
				name:    "異常 境界 下限未満",
				idxs:    []int{-1},
				wantErr: true,
			},

			{
				name:    "異常 境界 上限超過",
				idxs:    []int{64},
				wantErr: true,
			},
		}

		runFromIndices(t, cases)
	})
}

func TestSize(t *testing.T) {
	got8 := bitsx.Size[uint8]()
	want8 := 8
	if got8 != want8 {
		t.Errorf("値の不一致: got = %v, want = %v", got8, want8)
	}

	got64 := bitsx.Size[uint64]()
	want64 := 64
	if got64 != want64 {
		t.Errorf("値の不一致: got = %v, want = %v", got64, want64)
	}
}

// Bit, Set, Toggle, Clear 等の「インデックスを指定して操作する関数のテストケース
type indexOperationCase[B constraints.Unsigned] struct {
	name    string
	b       B
	idx     int
	want    B
	wantErr bool
}

func runIndexOperation[B constraints.Unsigned](t *testing.T, f func(B, int) (B, error), cases []indexOperationCase[B]) {
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := f(c.b, c.idx)

			if err != nil && !c.wantErr {
				t.Errorf("%s 予期せぬエラー: err = %v", c.name, err)
			}

			if err == nil && c.wantErr {
				t.Errorf("%s 想定外の非エラー", c.name)
			}

			if got != c.want {
				t.Errorf("%s 値の不一致: got = %b, want = %b", c.name, got, c.want)
			}
		})
	}
}

func TestBit(t *testing.T) {
	t.Run("uint8", func(t *testing.T) {
		cases := []indexOperationCase[uint8]{
			{
				name: "正常 0取得",
				b:    0b0000_1010,
				idx:  2,
				want: 0,
			},

			{
				name: "正常 1を取得",
				b:    0b0010_0000,
				idx:  5,
				want: 1,
			},

			{
				name: "正常 境界 下限",
				b:    0b0000_0001,
				idx:  0,
				want: 1,
			},

			{
				name: "正常 境界 上限",
				b:    0b1000_0000,
				idx:  7,
				want: 1,
			},

			{
				name:    "異常 境界 下限未満",
				idx:     -1,
				wantErr: true,
			},

			{
				name:    "異常 境界 上限超過",
				idx:     8,
				wantErr: true,
			},
		}

		runIndexOperation(t, bitsx.Bit[uint8], cases)
	})

	t.Run("uint64", func(t *testing.T) {
		cases := []indexOperationCase[uint64]{
			{
				name: "正常 0を取得",
				b:    0b11111111_11111111_11111111_11101111_11111111_11111111_11111111_11111111,
				idx:  36,
				want: 0,
			},

			{
				name: "正常 1を取得",
				b:    0b00000000_00000000_00000000_00000000_00000000_00000001_00000000_00000000,
				idx:  16,
				want: 1,
			},

			{
				name: "正常 境界 下限",
				b:    0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_00000001,
				idx:  0,
				want: 1,
			},

			{
				name: "正常 境界 上限",
				b:    0b10000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000,
				idx:  63,
				want: 1,
			},

			{
				name:    "異常 境界 下限未満",
				idx:     -1,
				wantErr: true,
			},

			{
				name:    "異常 境界 上限超過",
				idx:     64,
				wantErr: true,
			},
		}

		runIndexOperation(t, bitsx.Bit[uint64], cases)
	})
}

func TestSet(t *testing.T) {
	t.Run("uint8", func(t *testing.T) {
		cases := []indexOperationCase[uint8]{
			{
				name: "正常 通常",
				b:    0b0010_1000,
				idx:  4,
				want: 0b0011_1000,
			},

			{
				name: "正常 境界 下限",
				b:    0b1010_1010,
				idx:  0,
				want: 0b1010_1011,
			},

			{
				name: "正常 境界 上限",
				b:    0b0101_0101,
				idx:  7,
				want: 0b1101_0101,
			},

			{
				name: "準正常 セット済み",
				b:    0b0000_1000,
				idx:  3,
				want: 0b0000_1000,
			},

			{
				name:    "異常 境界 下限未満",
				idx:     -1,
				wantErr: true,
			},

			{
				name:    "異常 境界 上限超過",
				idx:     8,
				wantErr: true,
			},
		}

		runIndexOperation(t, bitsx.Set[uint8], cases)
	})

	t.Run("uint64", func(t *testing.T) {
		cases := []indexOperationCase[uint64]{
			{
				name: "正常 通常",
				b:    0b11111111_01111111_11111111_11111111_11111111_11111111_11111111_11111111,
				idx:  55,
				want: ^uint64(0),
			},

			{
				name: "正常 境界 下限",
				b:    0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_10101010,
				idx:  0,
				want: 0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_10101011,
			},

			{
				name: "正常 境界 上限",
				b:    0b00000000_11111111_00000000_11111111_00000000_11111111_00000000_11111111,
				idx:  63,
				want: 0b10000000_11111111_00000000_11111111_00000000_11111111_00000000_11111111,
			},

			{
				name: "準正常 セット済み",
				b:    ^uint64(0),
				idx:  32,
				want: ^uint64(0),
			},

			{
				name:    "異常 境界 下限未満",
				idx:     -1,
				wantErr: true,
			},

			{
				name:    "異常 境界 上限超過",
				idx:     64,
				wantErr: true,
			},
		}

		runIndexOperation(t, bitsx.Set[uint64], cases)
	})
}

func TestToggle(t *testing.T) {
	t.Run("uint8", func(t *testing.T) {
		cases := []indexOperationCase[uint8]{
			{
				name: "正常 0を1に",
				b:    0b0000_1010,
				idx:  2,
				want: 0b0000_1110,
			},

			{
				name: "正常 1を0に",
				b:    0b0111_0000,
				idx:  5,
				want: 0b0101_0000,
			},

			{
				name: "正常 境界 下限",
				b:    0b0000_0000,
				idx:  0,
				want: 0b0000_0001,
			},

			{
				name: "正常 境界 上限",
				b:    0b1000_0000,
				idx:  7,
				want: 0,
			},

			{
				name:    "異常 境界 下限未満",
				idx:     -1,
				wantErr: true,
			},

			{
				name:    "異常 境界 上限超過",
				idx:     8,
				wantErr: true,
			},
		}

		runIndexOperation(t, bitsx.Toggle[uint8], cases)
	})

	t.Run("uint64", func(t *testing.T) {
		cases := []indexOperationCase[uint64]{
			{
				name: "正常 0を1に",
				b:    0b00000000_00000000_10100000_00000000_00000000_00000000_00000000_00000000,
				idx:  46,
				want: 0b00000000_00000000_11100000_00000000_00000000_00000000_00000000_00000000,
			},

			{
				name: "正常 1を0に",
				b:    0b00000000_00000000_00000000_00000000_00000000_11111111_00000000_00000000,
				idx:  22,
				want: 0b00000000_00000000_00000000_00000000_00000000_10111111_00000000_00000000,
			},

			{
				name: "正常 境界 下限",
				b:    0b10000000_00000000_00000000_00000000_00000000_00000000_00000000_00000001,
				idx:  0,
				want: 0b10000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000,
			},

			{
				name: "正常 境界 上限",
				b:    0b11111111_00000000_00000000_00000000_00000000_00000000_00000000_11111111,
				idx:  63,
				want: 0b01111111_00000000_00000000_00000000_00000000_00000000_00000000_11111111,
			},

			{
				name:    "異常 境界 下限未満",
				idx:     -1,
				wantErr: true,
			},

			{
				name:    "異常 境界 上限超過",
				idx:     64,
				wantErr: true,
			},
		}

		runIndexOperation(t, bitsx.Toggle[uint64], cases)
	})
}

func TestClear(t *testing.T) {
	t.Run("uint8", func(t *testing.T) {
		cases := []indexOperationCase[uint8]{
			{
				name: "正常 通常",
				b:    0b1010_1010,
				idx:  3,
				want: 0b1010_0010,
			},

			{
				name: "正常 境界 下限",
				b:    0b1000_0001,
				idx:  0,
				want: 0b1000_0000,
			},

			{
				name: "正常 境界 上限",
				b:    0b1100_0010,
				idx:  7,
				want: 0b0100_0010,
			},

			{
				name: "準正常 クリア済み",
				b:    0b1111_0111,
				idx:  3,
				want: 0b1111_0111,
			},

			{
				name:    "異常 境界 下限未満",
				idx:     -1,
				wantErr: true,
			},

			{
				name:    "異常 境界 上限超過",
				idx:     8,
				wantErr: true,
			},
		}

		runIndexOperation(t, bitsx.Clear[uint8], cases)
	})

	t.Run("uint64", func(t *testing.T) {
		cases := []indexOperationCase[uint64]{
			{
				name: "正常 通常",
				b:    0b10000001_10000001_10000001_10000001_10000001_10000001_10000001_10000001,
				idx:  24,
				want: 0b10000001_10000001_10000001_10000001_10000000_10000001_10000001_10000001,
			},

			{
				name: "正常 境界 下限",
				b:    0b10000001_10000001_10000001_10000001_10000001_10000001_10000001_10000001,
				idx:  0,
				want: 0b10000001_10000001_10000001_10000001_10000001_10000001_10000001_10000000,
			},

			{
				name: "正常 境界 上限",
				b:    0b10000001_10000001_10000001_10000001_10000001_10000001_10000001_10000001,
				idx:  63,
				want: 0b00000001_10000001_10000001_10000001_10000001_10000001_10000001_10000001,
			},

			{
				name: "準正常 クリア済み",
				b:    0b10000001_10000001_10010001_10001001_10000001_10000001_10000001_10000001,
				idx:  62,
				want: 0b10000001_10000001_10010001_10001001_10000001_10000001_10000001_10000001,
			},

			{
				name:    "異常 境界 下限未満",
				idx:     -1,
				wantErr: true,
			},

			{
				name:    "異常 境界 上限超過",
				idx:     64,
				wantErr: true,
			},
		}

		runIndexOperation(t, bitsx.Clear[uint64], cases)
	})
}

type clearLowestCase[B constraints.Unsigned] struct {
	name string
	b    B
	want B
}

func runClearLowestCase[B constraints.Unsigned](t *testing.T, cases []clearLowestCase[B]) {
	for _, c := range cases {
		got := bitsx.ClearLowest(c.b)
		if got != c.want {
			t.Errorf("値の不一致: got = %v, want = %v", got, c.want)
		}
	}
}

func TestClearLowest(t *testing.T) {
	t.Run("uint8", func(t *testing.T) {
		cases := []clearLowestCase[uint8]{
			{
				name: "正常 通常",
				b:    0b1010_1000,
				want: 0b1010_0000,
			},

			{
				name: "準正常 0",
				b:    0,
				want: 0,
			},
		}

		runClearLowestCase(t, cases)
	})

	t.Run("uint64", func(t *testing.T) {
		cases := []clearLowestCase[uint64]{
			{
				name: "正常 通常",
				b:    0b00000000_10000000_10000000_10000000_10000000_10000000_00000000_00000000,
				want: 0b00000000_10000000_10000000_10000000_10000000_00000000_00000000_00000000,
			},

			{
				name: "準正常 0",
				b:    0,
				want: 0,
			},
		}

		runClearLowestCase(t, cases)
	})
}

type extractLowestCase[B constraints.Unsigned] struct {
	name string
	b    B
	want B
}

func runExtractLowestCase[B constraints.Unsigned](t *testing.T, cases []extractLowestCase[B]) {
	for _, c := range cases {
		got := bitsx.ExtractLowest(c.b)
		if got != c.want {
			t.Errorf("値の不一致: got = %v, want = %v", got, c.want)
		}
	}
}

func TestExtractLowest(t *testing.T) {
	t.Run("uint8", func(t *testing.T) {
		cases := []extractLowestCase[uint8]{
			{
				name: "正常 通常",
				b:    0b1010_1000,
				want: 0b0000_1000,
			},

			{
				name: "準正常 0",
				b:    0,
				want: 0,
			},
		}

		runExtractLowestCase(t, cases)
	})

	t.Run("uint64", func(t *testing.T) {
		cases := []extractLowestCase[uint64]{
			{
				name: "正常 通常",
				b:    0b00000000_01000000_10000000_00100000_00010000_00000000_00000000_00000000,
				want: 0b00000000_00000000_00000000_00000000_00010000_00000000_00000000_00000000,
			},

			{
				name: "準正常 0",
				b:    0,
				want: 0,
			},
		}

		runExtractLowestCase(t, cases)
	})
}

type indicesCase[B constraints.Unsigned] struct {
	name string
	b    B
	want []int
}

func runIndicesCase[B constraints.Unsigned](t *testing.T, cases []indicesCase[B]) {
	for _, c := range cases {
		got := bitsx.Indices(c.b)
		if !slices.Equal(got, c.want) {
			t.Errorf("値の不一致: got = %v, want = %v", got, c.want)
		}
	}
}

func TestIndices(t *testing.T) {
	t.Run("uint8", func(t *testing.T) {
		cases := []indicesCase[uint8]{
			{
				name: "正常 通常",
				b:    0b0101_0101,
				want: []int{0, 2, 4, 6},
			},

			{
				name: "正常 0",
				b:    0,
				want: []int{},
			},
		}

		runIndicesCase(t, cases)
	})

	t.Run("uint64", func(t *testing.T) {
		cases := []indicesCase[uint64]{
			{
				name: "正常 通常",
				b:    0b00000000_00000000_11110000_00000000_00001111_00000000_00000000_00000000,
				want: []int{24, 25, 26, 27, 44, 45, 46, 47},
			},

			{
				name: "正常 0",
				b:    0,
				want: []int{},
			},
		}

		runIndicesCase(t, cases)
	})
}

type singlesCase[B constraints.Unsigned] struct {
	name string
	b    B
	want []B
}

func runSingles[B constraints.Unsigned](t *testing.T, cases []singlesCase[B]) {
	for _, c := range cases {
		got := bitsx.Singles(c.b)
		if !slices.Equal(got, c.want) {
			t.Errorf("値の不一致: got = %v, want = %v", got, c.want)
		}
	}
}

func TestSingles(t *testing.T) {
	t.Run("uint8", func(t *testing.T) {
		cases := []singlesCase[uint8]{
			{
				name: "正常 通常",
				b:    0b1010_0100,
				want: []uint8{
					0b0000_0100,
					0b0010_0000,
					0b1000_0000,
				},
			},

			{
				name: "正常 0",
				b:    0,
				want: []uint8{},
			},
		}

		runSingles(t, cases)
	})

	t.Run("uint64", func(t *testing.T) {
		cases := []singlesCase[uint64]{
			{
				name: "正常 通常",
				b:    0b10000000_01000000_00100000_00010000_00001000_00000100_00000010_00000001,
				want: []uint64{
					0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_00000001,
					0b00000000_00000000_00000000_00000000_00000000_00000000_00000010_00000000,
					0b00000000_00000000_00000000_00000000_00000000_00000100_00000000_00000000,
					0b00000000_00000000_00000000_00000000_00001000_00000000_00000000_00000000,
					0b00000000_00000000_00000000_00010000_00000000_00000000_00000000_00000000,
					0b00000000_00000000_00100000_00000000_00000000_00000000_00000000_00000000,
					0b00000000_01000000_00000000_00000000_00000000_00000000_00000000_00000000,
					0b10000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000,
				},
			},

			{
				name: "正常 0",
				b:    0,
				want: []uint64{},
			},
		}

		runSingles(t, cases)
	})
}

type isSubsetCase[B constraints.Unsigned] struct {
	name  string
	super B
	sub   B
	want  bool
}

func runIsSubset[B constraints.Unsigned](t *testing.T, cases []isSubsetCase[B]) {
	for _, c := range cases {
		got := bitsx.IsSubset(c.super, c.sub)
		if got != c.want {
			t.Errorf("値の不一致: got = %t, want = %t", got, c.want)
		}
	}
}

func TestIsSubset(t *testing.T) {
	t.Run("uint8", func(t *testing.T) {
		cases := []isSubsetCase[uint8]{
			{
				name:  "正常 true",
				super: 0b0010_1100,
				sub:   0b0010_0100,
				want:  true,
			},

			{
				name:  "正常 false",
				super: 0b0101_0010,
				sub:   0b0100_1000,
				want:  false,
			},

			{
				name:  "正常 等しい",
				super: 0b0101_0100,
				sub:   0b0101_0100,
				want:  true,
			},

			{
				name:  "正常 0",
				super: 0,
				sub:   0,
				want:  true,
			},
		}

		runIsSubset(t, cases)
	})

	t.Run("uint64", func(t *testing.T) {
		cases := []isSubsetCase[uint64]{
			{
				name:  "正常 true",
				super: 0b00000000_00000000_00000000_11111111_00000000_11111111_00000000_00000000,
				sub:   0b00000000_00000000_00000000_00011100_00000000_00011100_00000000_00000000,
				want:  true,
			},

			{
				name:  "正常 false",
				super: 0b00000000_00000000_00000000_11111111_00000000_11111111_00000000_00000000,
				sub:   0b00000000_00000000_00000000_00011100_00001000_00011100_00000000_00000000,
				want:  false,
			},

			{
				name:  "正常 等しい",
				super: 0b00000001_00000010_00000100_00001000_00010000_00100000_01000000_10000000,
				sub:   0b00000001_00000010_00000100_00001000_00010000_00100000_01000000_10000000,
				want:  true,
			},

			{
				name:  "正常 0",
				super: 0,
				sub:   0,
				want:  true,
			},
		}

		runIsSubset(t, cases)
	})
}
