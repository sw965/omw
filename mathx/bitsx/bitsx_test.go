package bitsx_test

import (
	"github.com/sw965/omw/constraints"
	"github.com/sw965/omw/mathx/bitsx"
	"github.com/sw965/omw/mathx/randx"
	"slices"
	"testing"
	"fmt"
	"strings"
)

// --- 共通ヘルパー ---

func assertIndexError(t *testing.T, gotErr error, wantErrIdx, wantErrBitSize int) {
	t.Helper()
	if gotErr == nil {
		t.Fatalf("エラーを期待したが、エラーが起きなかった")
	}

	errMsg := gotErr.Error()
	wantErrMsgSubs := []string{
		"out of range",
		fmt.Sprintf("index %d", wantErrIdx),
		fmt.Sprintf("[0, %d)", wantErrBitSize),
	}

	for _, sub := range wantErrMsgSubs {
		if !strings.Contains(errMsg, sub) {
			t.Errorf("errMsg: %s, sub: %s", errMsg, sub)
		}
	}
}

// --- FromIndices ---

type fromIndicesTestCase[B constraints.Unsigned] struct {
	name           string
	idxs           []int
	want           B
	wantErr        bool
	wantErrIdx     int
	wantErrBitSize int
}

func runFromIndicesTests[B constraints.Unsigned](t *testing.T, tests []fromIndicesTestCase[B]) {
	t.Helper()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Helper()

			got, err := bitsx.FromIndices[B](tc.idxs)
			if tc.wantErr {
				assertIndexError(t, err, tc.wantErrIdx, tc.wantErrBitSize)
				return
			}

			if err != nil {
				t.Fatalf("予期せぬエラーが発生しました： %v", err)
			}

			if got != tc.want {
				t.Errorf("want: %d, got: %d", tc.want, got)
			}
		})
	}
}

func TestFromIndices(t *testing.T) {
	t.Run("uint8", func(t *testing.T) {
		runFromIndicesTests[uint8](t, []fromIndicesTestCase[uint8]{
			// 正常系
			{
				name: "正常 境界(下限)",
				idxs: []int{0, 3, 5},
				want: 0b00101001,
			},
			{
				name: "正常 境界(上限)",
				idxs: []int{1, 4, 7},
				want: 0b10010010,
			},
			{
				name: "正常 境界(複合)",
				idxs: []int{0, 3, 5, 7},
				want: 0b10101001,
			},
			{
				name: "正常 空スライス",
				idxs: []int{},
				want: 0,
			},
			// 異常系
			{
				name:           "異常 境界(下限越え)",
				idxs:           []int{-1, 3, 5},
				wantErr:        true,
				wantErrIdx:     -1,
				wantErrBitSize: 8,
			},
			{
				name:           "異常 境界(上限越え)",
				idxs:           []int{2, 6, 8},
				wantErr:        true,
				wantErrIdx:     8,
				wantErrBitSize: 8,
			},
			// 準正常系
			{
				name: "準正常 重複インデックス",
				idxs: []int{1, 1, 3, 3, 5, 5},
				want: 0b00101010,
			},
			{
				name: "準正常 nilスライス",
				idxs: nil,
				want: 0,
			},
			{
				name: "準正常 順不同",
				idxs: []int{5, 7, 1, 2, 0},
				want: 0b10100111,
			},
		})
	})

	t.Run("uint16", func(t *testing.T) {
		runFromIndicesTests[uint16](t, []fromIndicesTestCase[uint16]{
			// 正常系
			{
				name: "正常 境界(下限)",
				idxs: []int{0, 5, 10},
				want: 0b00000100_00100001,
			},
			{
				name: "正常 境界(上限)",
				idxs: []int{1, 8, 15},
				want: 0b10000001_00000010,
			},
			{
				name: "正常 境界(複合)",
				idxs: []int{0, 8, 15},
				want: 0b10000001_00000001,
			},
			{
				name: "正常 空スライス",
				idxs: []int{},
				want: 0,
			},
			// 異常系
			{
				name:           "異常 境界(下限越え)",
				idxs:           []int{-1, 5, 10},
				wantErr:        true,
				wantErrIdx:     -1,
				wantErrBitSize: 16,
			},
			{
				name:           "異常 境界(上限越え)",
				idxs:           []int{0, 5, 16},
				wantErr:        true,
				wantErrIdx:     16,
				wantErrBitSize: 16,
			},
			// 準正常系
			{
				name: "準正常 重複インデックス",
				idxs: []int{0, 0, 15, 15},
				want: 0b10000000_00000001,
			},
			{
				name: "準正常 nilスライス",
				idxs: nil,
				want: 0,
			},
			{
				name: "準正常 順不同",
				idxs: []int{15, 0, 8},
				want: 0b10000001_00000001,
			},
		})
	})

	// uint32 テストケース
	t.Run("uint32", func(t *testing.T) {
		runFromIndicesTests[uint32](t, []fromIndicesTestCase[uint32]{
			// 正常系
			{
				name: "正常 境界(下限)",
				idxs: []int{0, 10, 20},
				want: 0b00000000_00010000_00000100_00000001,
			},
			{
				name: "正常 境界(上限)",
				idxs: []int{1, 15, 31},
				want: 0b10000000_00000000_10000000_00000010,
			},
			{
				name: "正常 境界(複合)",
				idxs: []int{0, 16, 31},
				want: 0b10000000_00000001_00000000_00000001,
			},
			{
				name: "正常 空スライス",
				idxs: []int{},
				want: 0,
			},
			// 異常系
			{
				name:           "異常 境界(下限越え)",
				idxs:           []int{-1, 10, 20},
				wantErr:        true,
				wantErrIdx:     -1,
				wantErrBitSize: 32,
			},
			{
				name:           "異常 境界(上限越え)",
				idxs:           []int{0, 10, 32},
				wantErr:        true,
				wantErrIdx:     32,
				wantErrBitSize: 32,
			},
			// 準正常系
			{
				name: "準正常 重複インデックス",
				idxs: []int{0, 0, 31, 31},
				want: 0b10000000_00000000_00000000_00000001,
			},
			{
				name: "準正常 nilスライス",
				idxs: nil,
				want: 0,
			},
			{
				name: "準正常 順不同",
				idxs: []int{31, 0, 16},
				want: 0b10000000_00000001_00000000_00000001,
			},
		})
	})

	// uint64 テストケース
	t.Run("uint64", func(t *testing.T) {
		runFromIndicesTests[uint64](t, []fromIndicesTestCase[uint64]{
			// 正常系
			{
				name: "正常 境界(下限)",
				idxs: []int{0, 32, 48},
				want: 0b00000000_00000001_00000000_00000001_00000000_00000000_00000000_00000001,
			},
			{
				name: "正常 境界(上限)",
				idxs: []int{1, 31, 63},
				want: 0b10000000_00000000_00000000_00000000_10000000_00000000_00000000_00000010,
			},
			{
				name: "正常 境界(複合)",
				idxs: []int{0, 32, 63},
				want: 0b10000000_00000000_00000000_00000001_00000000_00000000_00000000_00000001,
			},
			{
				name: "正常 空スライス",
				idxs: []int{},
				want: 0,
			},
			// 異常系
			{
				name:           "異常 境界(下限越え)",
				idxs:           []int{-1, 32, 63},
				wantErr:        true,
				wantErrIdx:     -1,
				wantErrBitSize: 64,
			},
			{
				name:           "異常 境界(上限越え)",
				idxs:           []int{0, 32, 64},
				wantErr:        true,
				wantErrIdx:     64,
				wantErrBitSize: 64,
			},
			// 準正常系
			{
				name: "準正常 重複インデックス",
				idxs: []int{0, 0, 63, 63},
				want: 0b10000000_00000000_00000000_00000000_00000000_00000000_00000000_00000001,
			},
			{
				name: "準正常 nilスライス",
				idxs: nil,
				want: 0,
			},
			{
				name: "準正常 順不同",
				idxs: []int{63, 0, 32},
				want: 0b10000000_00000000_00000000_00000001_00000000_00000000_00000000_00000001,
			},
		})
	})
}

// --- IndexOperation(ToggleBit, SetBit, ClearBit) ---

type indexOperationTestCase[B constraints.Unsigned] struct {
	name           string
	b              B
	idx            int
	want           B
	wantErr        bool
	wantErrIdx     int
	wantErrBitSize int
}

func runIndexOperationTests[B constraints.Unsigned](t *testing.T, tests []indexOperationTestCase[B], f func(B, int) (B, error)) {
	t.Helper()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := f(tc.b, tc.idx)
			if tc.wantErr {
				assertIndexError(t, err, tc.wantErrIdx, tc.wantErrBitSize)
				return
			}

			if err != nil {
				t.Fatalf("予期せぬエラー: %v", err)
			}

			if got != tc.want {
				t.Errorf("want: %d, got: %d", tc.want, got)
			}
		})
	}
}

func TestToggleBit(t *testing.T) {
	t.Run("uint8", func(t *testing.T) {
		runIndexOperationTests(t, []indexOperationTestCase[uint8]{
			// 正常系
			{
				name: "正常_境界(下限)",
				b:    0b01100100, // 第0ビットは0
				idx:  0,
				want: 0b01100101, // 第0ビットを反転
			},
			{
				name: "正常_境界(上限)",
				b:    0b10011001, // 第7ビットは1
				idx:  7,
				want: 0b00011001, // 第7ビットを反転
			},
			// 異常系
			{
				name:           "異常_境界(下限越え)",
				b:              0b10000001,
				idx:            -1,
				wantErr:        true,
				wantErrIdx:     -1,
				wantErrBitSize: 8,
			},
			{
				name:           "異常_境界(上限越え)",
				b:              0b10010001,
				idx:            8,
				wantErr:        true,
				wantErrIdx:     8,
				wantErrBitSize: 8,
			},
		}, bitsx.ToggleBit)
	})

	t.Run("uint16", func(t *testing.T) {
		runIndexOperationTests[uint16](t, []indexOperationTestCase[uint16]{
			// 正常系
			{
				name: "正常_境界(下限)",
				b:    0b10101010_11001100, // 第0ビットは0
				idx:  0,
				want: 0b10101010_11001101, // 第0ビットを反転
			},
			{
				name: "正常_境界(上限)",
				b:    0b11110000_00001111, // 第15ビットは1
				idx:  15,
				want: 0b01110000_00001111, // 第15ビットを反転
			},
			// 異常系
			{
				name:           "異常_境界(下限越え)",
				b:              0,
				idx:            -1,
				wantErr:        true,
				wantErrIdx:     -1,
				wantErrBitSize: 16,
			},
			{
				name:           "異常_境界(上限越え)",
				b:              0,
				idx:            16,
				wantErr:        true,
				wantErrIdx:     16,
				wantErrBitSize: 16,
			},
		}, bitsx.ToggleBit)
	})

	t.Run("uint32", func(t *testing.T) {
		runIndexOperationTests[uint32](t, []indexOperationTestCase[uint32]{
			// 正常系
			{
				name: "正常_境界(下限)",
				b:    0b10101010_01010101_11001100_00110010, // 第0ビットは0
				idx:  0,
				want: 0b10101010_01010101_11001100_00110011, // 第0ビットを反転
			},
			{
				name: "正常_境界(上限)",
				b:    0b11000011_11110000_00001111_10101010, // 第31ビットは1
				idx:  31,
				want: 0b01000011_11110000_00001111_10101010, // 第31ビットを反転
			},
			// 異常系
			{
				name:           "異常_境界(下限越え)",
				b:              0,
				idx:            -1,
				wantErr:        true,
				wantErrIdx:     -1,
				wantErrBitSize: 32,
			},
			{
				name:           "異常_境界(上限越え)",
				b:              0,
				idx:            32,
				wantErr:        true,
				wantErrIdx:     32,
				wantErrBitSize: 32,
			},
		}, bitsx.ToggleBit)
	})

	t.Run("uint64", func(t *testing.T) {
		runIndexOperationTests[uint64](t, []indexOperationTestCase[uint64]{
			// 正常系
			{
				name: "正常_境界(下限)",
				b:    0b01010101_00000000_11111111_00000000_10101010_00000000_11001100_00110010, // 第0ビットは0
				idx:  0,
				want: 0b01010101_00000000_11111111_00000000_10101010_00000000_11001100_00110011, // 第0ビットを反転
			},
			{
				name: "正常_境界(上限)",
				b:    0b11110000_11110000_10101010_01010101_00001111_00001111_11001100_00110011, // 第63ビットは1
				idx:  63,
				want: 0b01110000_11110000_10101010_01010101_00001111_00001111_11001100_00110011, // 第63ビットを反転
			},
			// 異常系
			{
				name:           "異常_境界(下限越え)",
				b:              0,
				idx:            -1,
				wantErr:        true,
				wantErrIdx:     -1,
				wantErrBitSize: 64,
			},
			{
				name:           "異常_境界(上限越え)",
				b:              0,
				idx:            64,
				wantErr:        true,
				wantErrIdx:     64,
				wantErrBitSize: 64,
			},
		}, bitsx.ToggleBit)
	})
}

func TestSetBit(t *testing.T) {
	t.Run("uint8", func(t *testing.T) {
		runIndexOperationTests(t, []indexOperationTestCase[uint8]{
			// 正常系
			{
				name: "正常_境界(下限)",
				b:    0b00000110, // 第0ビットは0
				idx:  0,
				want: 0b00000111, // 第0ビットを1にセット
			},
			{
				name: "正常_境界(上限)",
				b:    0b01100000, // 第7ビットは0
				idx:  7,
				want: 0b11100000, // 第7ビットを1にセット
			},
			// 異常系
			{
				name:           "異常_境界(下限越え)",
				b:              0,
				idx:            -1,
				wantErr:        true,
				wantErrIdx:     -1,
				wantErrBitSize: 8,
			},
			{
				name:           "異常_境界(上限越え)",
				b:              0,
				idx:            8,
				wantErr:        true,
				wantErrIdx:     8,
				wantErrBitSize: 8,
			},
			// 準正常系
			{
				name: "準正常_重複",
				b:    0b00011000, // 第3ビットは既に1
				idx:  3,
				want: 0b00011000, // 変化なし
			},
		}, bitsx.SetBit)
	})

	t.Run("uint16", func(t *testing.T) {
		runIndexOperationTests[uint16](t, []indexOperationTestCase[uint16]{
			// 正常系
			{
				name: "正常_境界(下限)",
				b:    0b10101010_11001100, // 第0ビットは0
				idx:  0,
				want: 0b10101010_11001101, // 第0ビットを1にセット
			},
			{
				name: "正常_境界(上限)",
				b:    0b01110000_00001111, // 第15ビットは0
				idx:  15,
				want: 0b11110000_00001111, // 第15ビットを1にセット
			},
			// 異常系
			{
				name:           "異常_境界(下限越え)",
				b:              0,
				idx:            -1,
				wantErr:        true,
				wantErrIdx:     -1,
				wantErrBitSize: 16,
			},
			{
				name:           "異常_境界(上限越え)",
				b:              0,
				idx:            16,
				wantErr:        true,
				wantErrIdx:     16,
				wantErrBitSize: 16,
			},
			// 準正常系
			{
				name: "準正常_重複",
				b:    0b00000000_00000001, // 第0ビットは既に1
				idx:  0,
				want: 0b00000000_00000001, // 変化なし
			},
		}, bitsx.SetBit)
	})

	t.Run("uint32", func(t *testing.T) {
		runIndexOperationTests[uint32](t, []indexOperationTestCase[uint32]{
			// 正常系
			{
				name: "正常_境界(下限)",
				b:    0b10101010_01010101_11001100_00110010, // 第0ビットは0
				idx:  0,
				want: 0b10101010_01010101_11001100_00110011, // 第0ビットを1にセット
			},
			{
				name: "正常_境界(上限)",
				b:    0b01000011_11110000_00001111_10101010, // 第31ビットは0
				idx:  31,
				want: 0b11000011_11110000_00001111_10101010, // 第31ビットを1にセット
			},
			// 異常系
			{
				name:           "異常_境界(下限越え)",
				b:              0,
				idx:            -1,
				wantErr:        true,
				wantErrIdx:     -1,
				wantErrBitSize: 32,
			},
			{
				name:           "異常_境界(上限越え)",
				b:              0,
				idx:            32,
				wantErr:        true,
				wantErrIdx:     32,
				wantErrBitSize: 32,
			},
			// 準正常系
			{
				name: "準正常_重複",
				b:    0b10000000_00000000_00000000_00000000, // 第31ビットは既に1
				idx:  31,
				want: 0b10000000_00000000_00000000_00000000, // 変化なし
			},
		}, bitsx.SetBit)
	})

	t.Run("uint64", func(t *testing.T) {
		runIndexOperationTests[uint64](t, []indexOperationTestCase[uint64]{
			// 正常系
			{
				name: "正常_境界(下限)",
				b:    0b01010101_00000000_11111111_00000000_10101010_00000000_11001100_00110010, // 第0ビットは0
				idx:  0,
				want: 0b01010101_00000000_11111111_00000000_10101010_00000000_11001100_00110011, // 第0ビットを1にセット
			},
			{
				name: "正常_境界(上限)",
				b:    0b01110000_11110000_10101010_01010101_00001111_00001111_11001100_00110011, // 第63ビットは0
				idx:  63,
				want: 0b11110000_11110000_10101010_01010101_00001111_00001111_11001100_00110011, // 第63ビットを1にセット
			},
			// 異常系
			{
				name:           "異常_境界(下限越え)",
				b:              0,
				idx:            -1,
				wantErr:        true,
				wantErrIdx:     -1,
				wantErrBitSize: 64,
			},
			{
				name:           "異常_境界(上限越え)",
				b:              0,
				idx:            64,
				wantErr:        true,
				wantErrIdx:     64,
				wantErrBitSize: 64,
			},
			// 準正常系
			{
				name: "準正常_重複",
				b:    0b10000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000, // 第63ビットは既に1
				idx:  63,
				want: 0b10000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000, // 変化なし
			},
		}, bitsx.SetBit)
	})
}

func TestClearBit(t *testing.T) {
	t.Run("uint8", func(t *testing.T) {
		runIndexOperationTests(t, []indexOperationTestCase[uint8]{
			// 正常系
			{
				name: "正常_境界(下限)",
				b:    0b00000111, // 第0ビットは1
				idx:  0,
				want: 0b00000110, // 第0ビットを0にクリア
			},
			{
				name: "正常_境界(上限)",
				b:    0b11100000, // 第7ビットは1
				idx:  7,
				want: 0b01100000, // 第7ビットを0にクリア
			},
			// 異常系
			{
				name:           "異常_境界(下限越え)",
				b:              0,
				idx:            -1,
				wantErr:        true,
				wantErrIdx:     -1,
				wantErrBitSize: 8,
			},
			{
				name:           "異常_境界(上限越え)",
				b:              0,
				idx:            8,
				wantErr:        true,
				wantErrIdx:     8,
				wantErrBitSize: 8,
			},
			// 準正常系
			{
				name: "準正常_重複",
				b:    0b11110111, // 第3ビットは既に0
				idx:  3,
				want: 0b11110111, // 変化なし
			},
		}, bitsx.ClearBit)
	})

	t.Run("uint16", func(t *testing.T) {
		runIndexOperationTests[uint16](t, []indexOperationTestCase[uint16]{
			// 正常系
			{
				name: "正常_境界(下限)",
				b:    0b10101010_11001101, // 第0ビットは1
				idx:  0,
				want: 0b10101010_11001100, // 第0ビットを0にクリア
			},
			{
				name: "正常_境界(上限)",
				b:    0b11110000_00001111, // 第15ビットは1
				idx:  15,
				want: 0b01110000_00001111, // 第15ビットを0にクリア
			},
			// 異常系
			{
				name:           "異常_境界(下限越え)",
				b:              0,
				idx:            -1,
				wantErr:        true,
				wantErrIdx:     -1,
				wantErrBitSize: 16,
			},
			{
				name:           "異常_境界(上限越え)",
				b:              0,
				idx:            16,
				wantErr:        true,
				wantErrIdx:     16,
				wantErrBitSize: 16,
			},
			// 準正常系
			{
				name: "準正常_重複",
				b:    0b11111111_11111110, // 第0ビットは既に0
				idx:  0,
				want: 0b11111111_11111110, // 変化なし
			},
		}, bitsx.ClearBit)
	})

	t.Run("uint32", func(t *testing.T) {
		runIndexOperationTests[uint32](t, []indexOperationTestCase[uint32]{
			// 正常系
			{
				name: "正常_境界(下限)",
				b:    0b10101010_01010101_11001100_00110011, // 第0ビットは1
				idx:  0,
				want: 0b10101010_01010101_11001100_00110010, // 第0ビットを0にクリア
			},
			{
				name: "正常_境界(上限)",
				b:    0b11000011_11110000_00001111_10101010, // 第31ビットは1
				idx:  31,
				want: 0b01000011_11110000_00001111_10101010, // 第31ビットを0にクリア
			},
			// 異常系
			{
				name:           "異常_境界(下限越え)",
				b:              0,
				idx:            -1,
				wantErr:        true,
				wantErrIdx:     -1,
				wantErrBitSize: 32,
			},
			{
				name:           "異常_境界(上限越え)",
				b:              0,
				idx:            32,
				wantErr:        true,
				wantErrIdx:     32,
				wantErrBitSize: 32,
			},
			// 準正常系
			{
				name: "準正常_重複",
				b:    0b01111111_11111111_11111111_11111111, // 第31ビットは既に0
				idx:  31,
				want: 0b01111111_11111111_11111111_11111111, // 変化なし
			},
		}, bitsx.ClearBit)
	})

	t.Run("uint64", func(t *testing.T) {
		runIndexOperationTests[uint64](t, []indexOperationTestCase[uint64]{
			// 正常系
			{
				name: "正常_境界(下限)",
				b:    0b01010101_00000000_11111111_00000000_10101010_00000000_11001100_00110011, // 第0ビットは1
				idx:  0,
				want: 0b01010101_00000000_11111111_00000000_10101010_00000000_11001100_00110010, // 第0ビットを0にクリア
			},
			{
				name: "正常_境界(上限)",
				b:    0b11110000_11110000_10101010_01010101_00001111_00001111_11001100_00110011, // 第63ビットは1
				idx:  63,
				want: 0b01110000_11110000_10101010_01010101_00001111_00001111_11001100_00110011, // 第63ビットを0にクリア
			},
			// 異常系
			{
				name:           "異常_境界(下限越え)",
				b:              0,
				idx:            -1,
				wantErr:        true,
				wantErrIdx:     -1,
				wantErrBitSize: 64,
			},
			{
				name:           "異常_境界(上限越え)",
				b:              0,
				idx:            64,
				wantErr:        true,
				wantErrIdx:     64,
				wantErrBitSize: 64,
			},
			// 準正常系
			{
				name: "準正常_重複",
				b:    0b01111111_11111111_11111111_11111111_11111111_11111111_11111111_11111111, // 第63ビットは既に0
				idx:  63,
				want: 0b01111111_11111111_11111111_11111111_11111111_11111111_11111111_11111111, // 変化なし
			},
		}, bitsx.ClearBit)
	})
}

// --- ExtractLowestBit ---

type extractLowestBitTestCase[B constraints.Unsigned] struct {
	name string
	b    B
	want B
}

func runExtractLowestBitTests[B constraints.Unsigned](t *testing.T, tests []extractLowestBitTestCase[B]) {
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := bitsx.ExtractLowestBit[B](tc.b)
			if got != tc.want {
				t.Errorf("want: %d, got: %d", tc.want, got)
			}
		})
	}
}

func TestExtractLowestBit(t *testing.T) {
	t.Run("uint8", func(t *testing.T) {
		runExtractLowestBitTests[uint8](t, []extractLowestBitTestCase[uint8]{
			{
				name: "正常_最下位ビット",
				b:    0b00000001,
				want: 0b00000001,
			},
			{
				name: "正常_最上位ビット",
				b:    0b10000000,
				want: 0b10000000,
			},
			{
				name: "正常_両端",
				b:    0b10000001,
				want: 0b00000001,
			},
			{
				name: "正常_連続",
				b:    0b00111000,
				want: 0b00001000,
			},
			{
				name: "準正常_ゼロ値",
				b:    0,
				want: 0,
			},
		})
	})

	t.Run("uint16", func(t *testing.T) {
		runExtractLowestBitTests[uint16](t, []extractLowestBitTestCase[uint16]{
			{
				name: "正常_最下位ビット",
				b:    0b00000000_00000001,
				want: 0b00000000_00000001,
			},
			{
				name: "正常_最上位ビット",
				b:    0b10000000_00000000,
				want: 0b10000000_00000000,
			},
			{
				name: "正常_両端",
				b:    0b10000000_00000001,
				want: 0b00000000_00000001,
			},
			{
				name: "正常_連続",
				b:    0b00000111_00000000,
				want: 0b00000001_00000000,
			},
			{
				name: "準正常_ゼロ値",
				b:    0,
				want: 0,
			},
		})
	})

	t.Run("uint32", func(t *testing.T) {
		runExtractLowestBitTests[uint32](t, []extractLowestBitTestCase[uint32]{
			{
				name: "正常_最下位ビット",
				b:    0b00000000_00000000_00000000_00000001,
				want: 0b00000000_00000000_00000000_00000001,
			},
			{
				name: "正常_最上位ビット",
				b:    0b10000000_00000000_00000000_00000000,
				want: 0b10000000_00000000_00000000_00000000,
			},
			{
				name: "正常_両端",
				b:    0b10000000_00000000_00000000_00000001,
				want: 0b00000000_00000000_00000000_00000001,
			},
			{
				name: "正常_連続",
				b:    0b00000000_00000111_00000000_00000000,
				want: 0b00000000_00000001_00000000_00000000,
			},
			{
				name: "準正常_ゼロ値",
				b:    0,
				want: 0,
			},
		})
	})

	t.Run("uint64", func(t *testing.T) {
		runExtractLowestBitTests[uint64](t, []extractLowestBitTestCase[uint64]{
			{
				name: "正常_最下位ビット",
				b:    0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_00000001,
				want: 0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_00000001,
			},
			{
				name: "正常_最上位ビット",
				b:    0b10000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000,
				want: 0b10000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000,
			},
			{
				name: "正常_両端",
				b:    0b10000000_00000000_00000000_00000000_00000000_00000000_00000000_00000001,
				want: 0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_00000001,
			},
			{
				name: "正常_連続",
				b:    0b00000000_00000000_00000000_00000111_00000000_00000000_00000000_00000000,
				want: 0b00000000_00000000_00000000_00000001_00000000_00000000_00000000_00000000,
			},
			{
				name: "準正常_ゼロ値",
				b:    0,
				want: 0,
			},
		})
	})
}

// --- Indices ---

type indicesTestCase[B constraints.Unsigned] struct {
	name string
	b    B
	want []int
}

func runIndicesTests[B constraints.Unsigned](t *testing.T, tests []indicesTestCase[B]) {
	t.Helper()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := bitsx.Indices[B](tc.b)
			if !slices.Equal(got, tc.want) {
				t.Errorf("want: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestIndices(t *testing.T) {
	t.Run("uint8", func(t *testing.T) {
		runIndicesTests[uint8](t, []indicesTestCase[uint8]{
			{
				name: "正常_最下位ビット",
				b:    0b00000001,
				want: []int{0},
			},
			{
				name: "正常_最上位ビット",
				b:    0b10000000,
				want: []int{7},
			},
			{
				name: "正常_両端",
				b:    0b10000001,
				want: []int{0, 7},
			},
			{
				name: "正常_連続",
				b:    0b00111000,
				want: []int{3, 4, 5},
			},
			{
				name: "正常_ゼロ値",
				b:    0,
				want: []int{},
			},
		})
	})

	t.Run("uint16", func(t *testing.T) {
		runIndicesTests[uint16](t, []indicesTestCase[uint16]{
			{
				name: "正常_最下位ビット",
				b:    0b00000000_00000001,
				want: []int{0},
			},
			{
				name: "正常_最上位ビット",
				b:    0b10000000_00000000,
				want: []int{15},
			},
			{
				name: "正常_両端",
				b:    0b10000000_00000001,
				want: []int{0, 15},
			},
			{
				name: "正常_連続",
				b:    0b00000011_10000000,
				want: []int{7, 8, 9},
			},
			{
				name: "正常_ゼロ値",
				b:    0,
				want: []int{},
			},
		})
	})

	t.Run("uint32", func(t *testing.T) {
		runIndicesTests[uint32](t, []indicesTestCase[uint32]{
			{
				name: "正常_最下位ビット",
				b:    0b00000000_00000000_00000000_00000001,
				want: []int{0},
			},
			{
				name: "正常_最上位ビット",
				b:    0b10000000_00000000_00000000_00000000,
				want: []int{31},
			},
			{
				name: "正常_両端",
				b:    0b10000000_00000000_00000000_00000001,
				want: []int{0, 31},
			},
			{
				name: "正常_連続",
				b:    0b00000000_00000011_10000000_00000000,
				want: []int{15, 16, 17},
			},
			{
				name: "正常_ゼロ値",
				b:    0,
				want: []int{},
			},
		})
	})

	t.Run("uint64", func(t *testing.T) {
		runIndicesTests[uint64](t, []indicesTestCase[uint64]{
			{
				name: "正常_最下位ビット",
				b:    0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_00000001,
				want: []int{0},
			},
			{
				name: "正常_最上位ビット",
				b:    0b10000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000,
				want: []int{63},
			},
			{
				name: "正常_両端",
				b:    0b10000000_00000000_00000000_00000000_00000000_00000000_00000000_00000001,
				want: []int{0, 63},
			},
			{
				name: "正常_連続",
				b:    0b00000000_00000000_00000000_00000011_10000000_00000000_00000000_00000000,
				want: []int{31, 32, 33},
			},
			{
				name: "正常_ゼロ値",
				b:    0,
				want: []int{},
			},
		})
	})
}

// --- Singles ---

type singlesTestCase[B constraints.Unsigned] struct {
	name string
	b    B
	want []B
}

func runSinglesTests[B constraints.Unsigned](t *testing.T, tests []singlesTestCase[B]) {
	t.Helper()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := bitsx.Singles[B](tc.b)
			if !slices.Equal(got, tc.want) {
				t.Errorf("want: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestSingles(t *testing.T) {
	t.Run("uint8", func(t *testing.T) {
		runSinglesTests[uint8](t, []singlesTestCase[uint8]{
			{
				name: "正常_最下位ビット",
				b:    0b00000001,
				want: []uint8{0b00000001},
			},
			{
				name: "正常_最上位ビット",
				b:    0b10000000,
				want: []uint8{0b10000000},
			},
			{
				name: "正常_両端",
				b:    0b10000001,
				want: []uint8{0b00000001, 0b10000000},
			},
			{
				name: "正常_連続",
				b:    0b00111000,
				want: []uint8{0b00001000, 0b00010000, 0b00100000},
			},
			{
				name: "正常_ゼロ値",
				b:    0,
				want: []uint8{},
			},
		})
	})

	t.Run("uint16", func(t *testing.T) {
		runSinglesTests[uint16](t, []singlesTestCase[uint16]{
			{
				name: "正常_最下位ビット",
				b:    0b00000000_00000001,
				want: []uint16{0b00000000_00000001},
			},
			{
				name: "正常_最上位ビット",
				b:    0b10000000_00000000,
				want: []uint16{0b10000000_00000000},
			},
			{
				name: "正常_両端",
				b:    0b10000000_00000001,
				want: []uint16{0b00000000_00000001, 0b10000000_00000000},
			},
			{
				name: "正常_連続",
				b:    0b00000011_10000000,
				want: []uint16{0b00000000_10000000, 0b00000001_00000000, 0b00000010_00000000},
			},
			{
				name: "正常_ゼロ値",
				b:    0,
				want: []uint16{},
			},
		})
	})

	t.Run("uint32", func(t *testing.T) {
		runSinglesTests[uint32](t, []singlesTestCase[uint32]{
			{
				name: "正常_最下位ビット",
				b:    0b00000000_00000000_00000000_00000001,
				want: []uint32{0b00000000_00000000_00000000_00000001},
			},
			{
				name: "正常_最上位ビット",
				b:    0b10000000_00000000_00000000_00000000,
				want: []uint32{0b10000000_00000000_00000000_00000000},
			},
			{
				name: "正常_両端",
				b:    0b10000000_00000000_00000000_00000001,
				want: []uint32{0b00000000_00000000_00000000_00000001, 0b10000000_00000000_00000000_00000000},
			},
			{
				name: "正常_連続",
				b:    0b00000000_00000011_10000000_00000000,
				want: []uint32{0b00000000_00000000_10000000_00000000, 0b00000000_00000001_00000000_00000000, 0b00000000_00000010_00000000_00000000},
			},
			{
				name: "正常_ゼロ値",
				b:    0,
				want: []uint32{},
			},
		})
	})

	t.Run("uint64", func(t *testing.T) {
		runSinglesTests[uint64](t, []singlesTestCase[uint64]{
			{
				name: "正常_最下位ビット",
				b:    0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_00000001,
				want: []uint64{0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_00000001},
			},
			{
				name: "正常_最上位ビット",
				b:    0b10000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000,
				want: []uint64{0b10000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000},
			},
			{
				name: "正常_両端",
				b:    0b10000000_00000000_00000000_00000000_00000000_00000000_00000000_00000001,
				want: []uint64{0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_00000001, 0b10000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000},
			},
			{
				name: "正常_連続",
				b:    0b00000000_00000000_00000000_00000011_10000000_00000000_00000000_00000000,
				want: []uint64{0b00000000_00000000_00000000_00000000_10000000_00000000_00000000_00000000, 0b00000000_00000000_00000000_00000001_00000000_00000000_00000000_00000000, 0b00000000_00000000_00000000_00000010_00000000_00000000_00000000_00000000},
			},
			{
				name: "正常_ゼロ値",
				b:    0,
				want: []uint64{},
			},
		})
	})
}

// --- IsSubset ---

type isSubsetTestCase[B constraints.Unsigned] struct {
	name  string
	super B
	sub   B
	want  bool
}

func runIsSubsetTests[B constraints.Unsigned](t *testing.T, tests []isSubsetTestCase[B]) {
	t.Helper()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := bitsx.IsSubset[B](tc.super, tc.sub)
			if got != tc.want {
				t.Errorf("want: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestIsSubset(t *testing.T) {
	t.Run("uint8", func(t *testing.T) {
		runIsSubsetTests(t, []isSubsetTestCase[uint8]{
			// superのビットが全て1
			{
				name:  "正常_全部_全部",
				super: 0b11111111,
				sub:   0b11111111,
				want:  true,
			},
			{
				name:  "正常_全部_右半分",
				super: 0b11111111,
				sub:   0b00001111,
				want:  true,
			},
			{
				name:  "正常_全部_左半分",
				super: 0b11111111,
				sub:   0b11110000,
				want:  true,
			},
			{
				name:  "正常_全部_交互",
				super: 0b11111111,
				sub:   0b10101010,
				want:  true,
			},
			{
				name:  "正常_全部_0",
				super: 0b11111111,
				sub:   0,
				want:  true,
			},

			// superの右半分のビットが1
			{
				name:  "正常_右半分_全体",
				super: 0b00001111,
				sub:   0b11111111,
				want:  false,
			},
			{
				name:  "正常_右半分_右半分",
				super: 0b00001111,
				sub:   0b00001111,
				want:  true,
			},
			{
				name:  "正常_右半分_左半分",
				super: 0b00001111,
				sub:   0b11110000,
				want:  false,
			},
			{
				name:  "正常_右半分_交互",
				super: 0b00001111,
				sub:   0b00001010,
				want:  true,
			},
			{
				name:  "正常_右半分_0",
				super: 0b00001111,
				sub:   0,
				want:  true,
			},

			// superの左半分のビットが1
			{
				name:  "正常_左半分_全部",
				super: 0b11110000,
				sub:   0b11111111,
				want:  false,
			},
			{
				name:  "正常_左半分_右半分",
				super: 0b11110000,
				sub:   0b00001111,
				want:  false,
			},
			{
				name:  "正常_左半分_左半分",
				super: 0b11110000,
				sub:   0b11110000,
				want:  true,
			},
			{
				name:  "正常_左半分_交互",
				super: 0b11110000,
				sub:   0b10100000,
				want:  true,
			},
			{
				name:  "正常_左半分_0",
				super: 0b11110000,
				sub:   0,
				want:  true,
			},

			// superが0
			{
				name:  "正常_0_全部",
				super: 0,
				sub:   0b11111111,
				want:  false,
			},
			{
				name:  "正常_0_右半分",
				super: 0,
				sub:   0b00001111,
				want:  false,
			},
			{
				name:  "正常_0_左半分",
				super: 0,
				sub:   0b11110000,
				want:  false,
			},
			{
				name:  "正常_0_0",
				super: 0,
				sub:   0,
				want:  true,
			},

			// その他
			{
				name:  "正常_交互_逆交互",
				super: 0b10101010,
				sub:   0b01010101,
				want:  false,
			},
			{
				name:  "正常_逆交互_交互",
				super: 0b01010101,
				sub:   0b10101010,
				want:  false,
			},
			{
				name:  "正常_最下位_最下位",
				super: 0b00000001,
				sub:   0b00000001,
				want:  true,
			},
			{
				name:  "正常_最下位_最上位",
				super: 0b00000001,
				sub:   0b10000000,
				want:  false,
			},
			{
				name:  "正常_最上位_最下位",
				super: 0b10000000,
				sub:   0b00000001,
				want:  false,
			},
			{
				name:  "正常_最上位_最上位",
				super: 0b10000000,
				sub:   0b10000000,
				want:  true,
			},
			{
				name:  "正常_欠け1bit_全部",
				super: 0b11101111,
				sub:   0b11111111,
				want:  false,
			},
		})
	})

	t.Run("uint16", func(t *testing.T) {
		runIsSubsetTests(t, []isSubsetTestCase[uint16]{
			// superのビットが全て1
			{
				name:  "正常_全部_全部",
				super: 0b11111111_11111111,
				sub:   0b11111111_11111111,
				want:  true,
			},
			{
				name:  "正常_全部_右半分",
				super: 0b11111111_11111111,
				sub:   0b00000000_11111111,
				want:  true,
			},
			{
				name:  "正常_全部_左半分",
				super: 0b11111111_11111111,
				sub:   0b11111111_00000000,
				want:  true,
			},
			{
				name:  "正常_全部_交互",
				super: 0b11111111_11111111,
				sub:   0b10101010_10101010,
				want:  true,
			},
			{
				name:  "正常_全部_0",
				super: 0b11111111_11111111,
				sub:   0,
				want:  true,
			},

			// superの右半分のビットが1
			{
				name:  "正常_右半分_全体",
				super: 0b00000000_11111111,
				sub:   0b11111111_11111111,
				want:  false,
			},
			{
				name:  "正常_右半分_右半分",
				super: 0b00000000_11111111,
				sub:   0b00000000_11111111,
				want:  true,
			},
			{
				name:  "正常_右半分_左半分",
				super: 0b00000000_11111111,
				sub:   0b11111111_00000000,
				want:  false,
			},
			{
				name:  "正常_右半分_交互",
				super: 0b00000000_11111111,
				sub:   0b00000000_10101010,
				want:  true,
			},
			{
				name:  "正常_右半分_0",
				super: 0b00000000_11111111,
				sub:   0,
				want:  true,
			},

			// superの左半分のビットが1
			{
				name:  "正常_左半分_全部",
				super: 0b11111111_00000000,
				sub:   0b11111111_11111111,
				want:  false,
			},
			{
				name:  "正常_左半分_右半分",
				super: 0b11111111_00000000,
				sub:   0b00000000_11111111,
				want:  false,
			},
			{
				name:  "正常_左半分_左半分",
				super: 0b11111111_00000000,
				sub:   0b11111111_00000000,
				want:  true,
			},
			{
				name:  "正常_左半分_交互",
				super: 0b11111111_00000000,
				sub:   0b10101010_00000000,
				want:  true,
			},
			{
				name:  "正常_左半分_0",
				super: 0b11111111_00000000,
				sub:   0,
				want:  true,
			},

			// superが0
			{
				name:  "正常_0_全部",
				super: 0,
				sub:   0b11111111_11111111,
				want:  false,
			},
			{
				name:  "正常_0_右半分",
				super: 0,
				sub:   0b00000000_11111111,
				want:  false,
			},
			{
				name:  "正常_0_左半分",
				super: 0,
				sub:   0b11111111_00000000,
				want:  false,
			},
			{
				name:  "正常_0_0",
				super: 0,
				sub:   0,
				want:  true,
			},

			// その他
			{
				name:  "正常_交互_逆交互",
				super: 0b10101010_10101010,
				sub:   0b01010101_01010101,
				want:  false,
			},
			{
				name:  "正常_逆交互_交互",
				super: 0b01010101_01010101,
				sub:   0b10101010_10101010,
				want:  false,
			},
			{
				name:  "正常_最下位_最下位",
				super: 0b00000000_00000001,
				sub:   0b00000000_00000001,
				want:  true,
			},
			{
				name:  "正常_最下位_最上位",
				super: 0b00000000_00000001,
				sub:   0b10000000_00000000,
				want:  false,
			},
			{
				name:  "正常_最上位_最下位",
				super: 0b10000000_00000000,
				sub:   0b00000000_00000001,
				want:  false,
			},
			{
				name:  "正常_最上位_最上位",
				super: 0b10000000_00000000,
				sub:   0b10000000_00000000,
				want:  true,
			},
			{
				name:  "正常_欠け1bit_全部",
				super: 0b11111111_11101111,
				sub:   0b11111111_11111111,
				want:  false,
			},
		})
	})

	t.Run("uint32", func(t *testing.T) {
		runIsSubsetTests(t, []isSubsetTestCase[uint32]{
			// superのビットが全て1
			{
				name:  "正常_全部_全部",
				super: 0b11111111_11111111_11111111_11111111,
				sub:   0b11111111_11111111_11111111_11111111,
				want:  true,
			},
			{
				name:  "正常_全部_右半分",
				super: 0b11111111_11111111_11111111_11111111,
				sub:   0b00000000_00000000_11111111_11111111,
				want:  true,
			},
			{
				name:  "正常_全部_左半分",
				super: 0b11111111_11111111_11111111_11111111,
				sub:   0b11111111_11111111_00000000_00000000,
				want:  true,
			},
			{
				name:  "正常_全部_交互",
				super: 0b11111111_11111111_11111111_11111111,
				sub:   0b10101010_10101010_10101010_10101010,
				want:  true,
			},
			{
				name:  "正常_全部_0",
				super: 0b11111111_11111111_11111111_11111111,
				sub:   0,
				want:  true,
			},

			// superの右半分のビットが1
			{
				name:  "正常_右半分_全体",
				super: 0b00000000_00000000_11111111_11111111,
				sub:   0b11111111_11111111_11111111_11111111,
				want:  false,
			},
			{
				name:  "正常_右半分_右半分",
				super: 0b00000000_00000000_11111111_11111111,
				sub:   0b00000000_00000000_11111111_11111111,
				want:  true,
			},
			{
				name:  "正常_右半分_左半分",
				super: 0b00000000_00000000_11111111_11111111,
				sub:   0b11111111_11111111_00000000_00000000,
				want:  false,
			},
			{
				name:  "正常_右半分_交互",
				super: 0b00000000_00000000_11111111_11111111,
				sub:   0b00000000_00000000_10101010_10101010,
				want:  true,
			},
			{
				name:  "正常_右半分_0",
				super: 0b00000000_00000000_11111111_11111111,
				sub:   0,
				want:  true,
			},

			// superの左半分のビットが1
			{
				name:  "正常_左半分_全部",
				super: 0b11111111_11111111_00000000_00000000,
				sub:   0b11111111_11111111_11111111_11111111,
				want:  false,
			},
			{
				name:  "正常_左半分_右半分",
				super: 0b11111111_11111111_00000000_00000000,
				sub:   0b00000000_00000000_11111111_11111111,
				want:  false,
			},
			{
				name:  "正常_左半分_左半分",
				super: 0b11111111_11111111_00000000_00000000,
				sub:   0b11111111_11111111_00000000_00000000,
				want:  true,
			},
			{
				name:  "正常_左半分_交互",
				super: 0b11111111_11111111_00000000_00000000,
				sub:   0b10101010_10101010_00000000_00000000,
				want:  true,
			},
			{
				name:  "正常_左半分_0",
				super: 0b11111111_11111111_00000000_00000000,
				sub:   0,
				want:  true,
			},

			// superが0
			{
				name:  "正常_0_全部",
				super: 0,
				sub:   0b11111111_11111111_11111111_11111111,
				want:  false,
			},
			{
				name:  "正常_0_右半分",
				super: 0,
				sub:   0b00000000_00000000_11111111_11111111,
				want:  false,
			},
			{
				name:  "正常_0_左半分",
				super: 0,
				sub:   0b11111111_11111111_00000000_00000000,
				want:  false,
			},
			{
				name:  "正常_0_0",
				super: 0,
				sub:   0,
				want:  true,
			},

			// その他
			{
				name:  "正常_交互_逆交互",
				super: 0b10101010_10101010_10101010_10101010,
				sub:   0b01010101_01010101_01010101_01010101,
				want:  false,
			},
			{
				name:  "正常_逆交互_交互",
				super: 0b01010101_01010101_01010101_01010101,
				sub:   0b10101010_10101010_10101010_10101010,
				want:  false,
			},
			{
				name:  "正常_最下位_最下位",
				super: 0b00000000_00000000_00000000_00000001,
				sub:   0b00000000_00000000_00000000_00000001,
				want:  true,
			},
			{
				name:  "正常_最下位_最上位",
				super: 0b00000000_00000000_00000000_00000001,
				sub:   0b10000000_00000000_00000000_00000000,
				want:  false,
			},
			{
				name:  "正常_最上位_最下位",
				super: 0b10000000_00000000_00000000_00000000,
				sub:   0b00000000_00000000_00000000_00000001,
				want:  false,
			},
			{
				name:  "正常_最上位_最上位",
				super: 0b10000000_00000000_00000000_00000000,
				sub:   0b10000000_00000000_00000000_00000000,
				want:  true,
			},
			{
				name:  "正常_欠け1bit_全部",
				super: 0b11111111_11111111_11111111_11101111,
				sub:   0b11111111_11111111_11111111_11111111,
				want:  false,
			},
		})
	})

	t.Run("uint64", func(t *testing.T) {
		runIsSubsetTests(t, []isSubsetTestCase[uint64]{
			// superのビットが全て1
			{
				name:  "正常_全部_全部",
				super: 0b11111111_11111111_11111111_11111111_11111111_11111111_11111111_11111111,
				sub:   0b11111111_11111111_11111111_11111111_11111111_11111111_11111111_11111111,
				want:  true,
			},
			{
				name:  "正常_全部_右半分",
				super: 0b11111111_11111111_11111111_11111111_11111111_11111111_11111111_11111111,
				sub:   0b00000000_00000000_00000000_00000000_11111111_11111111_11111111_11111111,
				want:  true,
			},
			{
				name:  "正常_全部_左半分",
				super: 0b11111111_11111111_11111111_11111111_11111111_11111111_11111111_11111111,
				sub:   0b11111111_11111111_11111111_11111111_00000000_00000000_00000000_00000000,
				want:  true,
			},
			{
				name:  "正常_全部_交互",
				super: 0b11111111_11111111_11111111_11111111_11111111_11111111_11111111_11111111,
				sub:   0b10101010_10101010_10101010_10101010_10101010_10101010_10101010_10101010,
				want:  true,
			},
			{
				name:  "正常_全部_0",
				super: 0b11111111_11111111_11111111_11111111_11111111_11111111_11111111_11111111,
				sub:   0,
				want:  true,
			},

			// superの右半分のビットが1
			{
				name:  "正常_右半分_全体",
				super: 0b00000000_00000000_00000000_00000000_11111111_11111111_11111111_11111111,
				sub:   0b11111111_11111111_11111111_11111111_11111111_11111111_11111111_11111111,
				want:  false,
			},
			{
				name:  "正常_右半分_右半分",
				super: 0b00000000_00000000_00000000_00000000_11111111_11111111_11111111_11111111,
				sub:   0b00000000_00000000_00000000_00000000_11111111_11111111_11111111_11111111,
				want:  true,
			},
			{
				name:  "正常_右半分_左半分",
				super: 0b00000000_00000000_00000000_00000000_11111111_11111111_11111111_11111111,
				sub:   0b11111111_11111111_11111111_11111111_00000000_00000000_00000000_00000000,
				want:  false,
			},
			{
				name:  "正常_右半分_交互",
				super: 0b00000000_00000000_00000000_00000000_11111111_11111111_11111111_11111111,
				sub:   0b00000000_00000000_00000000_00000000_10101010_10101010_10101010_10101010,
				want:  true,
			},
			{
				name:  "正常_右半分_0",
				super: 0b00000000_00000000_00000000_00000000_11111111_11111111_11111111_11111111,
				sub:   0,
				want:  true,
			},

			// superの左半分のビットが1
			{
				name:  "正常_左半分_全部",
				super: 0b11111111_11111111_11111111_11111111_00000000_00000000_00000000_00000000,
				sub:   0b11111111_11111111_11111111_11111111_11111111_11111111_11111111_11111111,
				want:  false,
			},
			{
				name:  "正常_左半分_右半分",
				super: 0b11111111_11111111_11111111_11111111_00000000_00000000_00000000_00000000,
				sub:   0b00000000_00000000_00000000_00000000_11111111_11111111_11111111_11111111,
				want:  false,
			},
			{
				name:  "正常_左半分_左半分",
				super: 0b11111111_11111111_11111111_11111111_00000000_00000000_00000000_00000000,
				sub:   0b11111111_11111111_11111111_11111111_00000000_00000000_00000000_00000000,
				want:  true,
			},
			{
				name:  "正常_左半分_交互",
				super: 0b11111111_11111111_11111111_11111111_00000000_00000000_00000000_00000000,
				sub:   0b10101010_10101010_10101010_10101010_00000000_00000000_00000000_00000000,
				want:  true,
			},
			{
				name:  "正常_左半分_0",
				super: 0b11111111_11111111_11111111_11111111_00000000_00000000_00000000_00000000,
				sub:   0,
				want:  true,
			},

			// superが0
			{
				name:  "正常_0_全部",
				super: 0,
				sub:   0b11111111_11111111_11111111_11111111_11111111_11111111_11111111_11111111,
				want:  false,
			},
			{
				name:  "正常_0_右半分",
				super: 0,
				sub:   0b00000000_00000000_00000000_00000000_11111111_11111111_11111111_11111111,
				want:  false,
			},
			{
				name:  "正常_0_左半分",
				super: 0,
				sub:   0b11111111_11111111_11111111_11111111_00000000_00000000_00000000_00000000,
				want:  false,
			},
			{
				name:  "正常_0_0",
				super: 0,
				sub:   0,
				want:  true,
			},

			// その他
			{
				name:  "正常_交互_逆交互",
				super: 0b10101010_10101010_10101010_10101010_10101010_10101010_10101010_10101010,
				sub:   0b01010101_01010101_01010101_01010101_01010101_01010101_01010101_01010101,
				want:  false,
			},
			{
				name:  "正常_逆交互_交互",
				super: 0b01010101_01010101_01010101_01010101_01010101_01010101_01010101_01010101,
				sub:   0b10101010_10101010_10101010_10101010_10101010_10101010_10101010_10101010,
				want:  false,
			},
			{
				name:  "正常_最下位_最下位",
				super: 0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_00000001,
				sub:   0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_00000001,
				want:  true,
			},
			{
				name:  "正常_最下位_最上位",
				super: 0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_00000001,
				sub:   0b10000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000,
				want:  false,
			},
			{
				name:  "正常_最上位_最下位",
				super: 0b10000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000,
				sub:   0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_00000001,
				want:  false,
			},
			{
				name:  "正常_最上位_最上位",
				super: 0b10000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000,
				sub:   0b10000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000,
				want:  true,
			},
			{
				name:  "正常_欠け1bit_全部",
				super: 0b11111111_11111111_11111111_11111111_11111111_11111111_11111111_11101111,
				sub:   0b11111111_11111111_11111111_11111111_11111111_11111111_11111111_11111111,
				want:  false,
			},
		})
	})
}

func BenchmarkMulVecAndPopCount(b *testing.B) {
	// テストする行列のサイズ（行数 x 列数）
	// AVX-512の性能を引き出すため、列数を変えて測定します。
	sizes := []struct {
		rows, cols int
	}{
		{100, 64},    // 1ブロック (Small)
		{100, 1024},  // 16ブロック (Medium)
		{1000, 8192}, // 128ブロック (Large)
	}

	for _, sz := range sizes {
		// データ準備
		rng := randx.NewPCGFromGlobalSeed()
		m, _ := bitsx.NewRandMatrix(sz.rows, sz.cols, rng)
		v, _ := bitsx.NewRandMatrix(1, sz.cols, rng)

		// 1. Generic版のベンチマーク
		b.Run(fmt.Sprintf("Generic/R%d-C%d", sz.rows, sz.cols), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = m.MulVecAndPopCount(v)
			}
		})

		// 2. AVX-512版（アセンブラ）のベンチマーク
		b.Run(fmt.Sprintf("AVX512/R%d-C%d", sz.rows, sz.cols), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = m.MulVecAndPopCountAVX512(v)
			}
		})
	}
}