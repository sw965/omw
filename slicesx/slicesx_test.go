package slicesx_test

import (
	"github.com/sw965/omw/slicesx"
	"iter"
	"maps"
	"slices"
	"testing"
	"cmp"
)

type selectTestCase struct {
	name string
	s    []string
	r    int
	want [][]string
}

func runSelectTests(t *testing.T, tests []selectTestCase, f func([]string, int) iter.Seq[[]string]) {
	t.Helper()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Helper()
			seq := f(tc.s, tc.r)
			got := slices.Collect(seq)

			gotN := len(got)
			wantN := len(tc.want)

			if len(got) != len(tc.want) {
				t.Fatalf("len(want): %d, len(got): %d", wantN, gotN)
			}

			for i, gv := range got {
				gw := tc.want[i]
				if !slices.Equal(gv, gw) {
					t.Errorf("i = %d, want: %v, got: %v", i, gw, gv)
					break
				}
			}
		})
	}
}

func TestPermutations(t *testing.T) {
	tests := []selectTestCase{
		{
			name: "正常",
			s:    []string{"りんご", "みかん", "ぶどう", "いちご"},
			r:    3,
			want: [][]string{
				{"りんご", "みかん", "ぶどう"},
				{"りんご", "みかん", "いちご"},
				{"りんご", "ぶどう", "みかん"},
				{"りんご", "ぶどう", "いちご"},
				{"りんご", "いちご", "みかん"},
				{"りんご", "いちご", "ぶどう"},
				{"みかん", "りんご", "ぶどう"},
				{"みかん", "りんご", "いちご"},
				{"みかん", "ぶどう", "りんご"},
				{"みかん", "ぶどう", "いちご"},
				{"みかん", "いちご", "りんご"},
				{"みかん", "いちご", "ぶどう"},
				{"ぶどう", "りんご", "みかん"},
				{"ぶどう", "りんご", "いちご"},
				{"ぶどう", "みかん", "りんご"},
				{"ぶどう", "みかん", "いちご"},
				{"ぶどう", "いちご", "りんご"},
				{"ぶどう", "いちご", "みかん"},
				{"いちご", "りんご", "みかん"},
				{"いちご", "りんご", "ぶどう"},
				{"いちご", "みかん", "りんご"},
				{"いちご", "みかん", "ぶどう"},
				{"いちご", "ぶどう", "りんご"},
				{"いちご", "ぶどう", "みかん"},
			},
		},
		{
			name: "正常_rが0",
			s:    []string{"卵", "パン", "ハム"},
			r:    0,
			// rが0の時は、空スライスを1つ持つ
			want: [][]string{{}},
		},
		{
			name: "正常_len(s)とrが同じ",
			s:    []string{"剣", "盾", "兜"},
			r:    3,
			want: [][]string{
				{"剣", "盾", "兜"},
				{"剣", "兜", "盾"},
				{"盾", "剣", "兜"},
				{"盾", "兜", "剣"},
				{"兜", "剣", "盾"},
				{"兜", "盾", "剣"},
			},
		},
		{
			name: "準正常_rが負の値",
			s:    []string{"手", "足", "顔"},
			r:    -1,
			want: [][]string{},
		},
		{
			name: "準正常_rがlen(s)より大きい",
			s:    []string{"洗濯機", "電子レンジ", "掃除機"},
			r:    4,
			want: [][]string{},
		},
		{
			name: "準正常_重複",
			s:    []string{"サメ", "アザラシ", "アザラシ"},
			r:    2,
			want: [][]string{
				{"サメ", "アザラシ"},
				{"サメ", "アザラシ"},
				{"アザラシ", "サメ"},
				{"アザラシ", "アザラシ"},
				{"アザラシ", "サメ"},
				{"アザラシ", "アザラシ"},
			},
		},
	}
	runSelectTests(t, tests, slicesx.Permutations[[]string, string])
}

func TestSequences(t *testing.T) {
	tests := []selectTestCase{
		{
			name: "正常",
			s:    []string{"Python", "Go", "C++"},
			r:    2,
			want: [][]string{
				{"Python", "Python"},
				{"Python", "Go"},
				{"Python", "C++"},
				{"Go", "Python"},
				{"Go", "Go"},
				{"Go", "C++"},
				{"C++", "Python"},
				{"C++", "Go"},
				{"C++", "C++"},
			},
		},
		{
			name: "正常_境界値(下限)_rが0",
			s:    []string{"卵", "パン", "ハム"},
			r:    0,
			// rが0の時は、空スライスを1つ持つ
			want: [][]string{{}},
		},
		{
			name: "正常_len(s)が1",
			s:    []string{"単体"},
			r:    3,
			want: [][]string{
				{"単体", "単体", "単体"},
			},
		},
		{
			name: "準正常_rが負の値",
			s:    []string{"手", "足", "顔"},
			r:    -1,
			want: [][]string{},
		},
		{
			name: "準正常_sが空",
			s:    []string{},
			r:    2,
			want: [][]string{},
		},
		{
			name: "準正常_重複",
			s:    []string{"ダイヤ", "ダイヤ", "ダイヤ"},
			r:    2,
			want: [][]string{
				{"ダイヤ", "ダイヤ"},
				{"ダイヤ", "ダイヤ"},
				{"ダイヤ", "ダイヤ"},
				{"ダイヤ", "ダイヤ"},
				{"ダイヤ", "ダイヤ"},
				{"ダイヤ", "ダイヤ"},
				{"ダイヤ", "ダイヤ"},
				{"ダイヤ", "ダイヤ"},
				{"ダイヤ", "ダイヤ"},
			},
		},
	}

	runSelectTests(t, tests, slicesx.Sequences[[]string, string])
}

func TestCombinations(t *testing.T) {
	tests := []selectTestCase{
		{
			name: "正常",
			s:    []string{"ChatGPT", "Gemini", "DeepSeek", "Grok"},
			r:    3,
			want: [][]string{
				{"ChatGPT", "Gemini", "DeepSeek"},
				{"ChatGPT", "Gemini", "Grok"},
				{"ChatGPT", "DeepSeek", "Grok"},
				{"Gemini", "DeepSeek", "Grok"},
			},
		},
		{
			name: "正常_境界値(下限)_rが0",
			s:    []string{"剣", "盾", "兜"},
			r:    0,
			// rが0の時は、空スライスを1つ持つ
			want: [][]string{{}},
		},
		{
			name: "正常_len(s)とrが同じ",
			s:    []string{"赤", "青", "緑"},
			r:    3,
			want: [][]string{
				{"赤", "青", "緑"},
			},
		},
		{
			name: "準正常_rが負の値",
			s:    []string{"手", "足", "顔"},
			r:    -1,
			want: [][]string{},
		},
		{
			name: "準正常_rがlen(s)より大きい",
			s:    []string{"洗濯機", "電子レンジ", "掃除機"},
			r:    4,
			want: [][]string{},
		},
		{
			name: "準正常_重複",
			s:    []string{"日本", "韓国", "中国", "日本"},
			r:    2,
			want: [][]string{
				{"日本", "韓国"},
				{"日本", "中国"},
				{"日本", "日本"},
				{"韓国", "中国"},
				{"韓国", "日本"},
				{"中国", "日本"},
			},
		},
	}

	runSelectTests(t, tests, slicesx.Combinations[[]string, string])
}

func TestCartesianProducts(t *testing.T) {
	tests := []struct {
		name string
		args [][]string
		want [][]string
	}{
		{
			name: "正常",
			args: [][]string{
				[]string{"りんご", "みかん", "ぶどう", "いちご"},
				[]string{"Python", "Go", "C++"},
				[]string{"ChatGPT", "Gemini", "DeepSeek", "Grok"},
			},
			want: [][]string{
				{"りんご", "Python", "ChatGPT"},
				{"りんご", "Python", "Gemini"},
				{"りんご", "Python", "DeepSeek"},
				{"りんご", "Python", "Grok"},
				{"りんご", "Go", "ChatGPT"},
				{"りんご", "Go", "Gemini"},
				{"りんご", "Go", "DeepSeek"},
				{"りんご", "Go", "Grok"},
				{"りんご", "C++", "ChatGPT"},
				{"りんご", "C++", "Gemini"},
				{"りんご", "C++", "DeepSeek"},
				{"りんご", "C++", "Grok"},

				{"みかん", "Python", "ChatGPT"},
				{"みかん", "Python", "Gemini"},
				{"みかん", "Python", "DeepSeek"},
				{"みかん", "Python", "Grok"},
				{"みかん", "Go", "ChatGPT"},
				{"みかん", "Go", "Gemini"},
				{"みかん", "Go", "DeepSeek"},
				{"みかん", "Go", "Grok"},
				{"みかん", "C++", "ChatGPT"},
				{"みかん", "C++", "Gemini"},
				{"みかん", "C++", "DeepSeek"},
				{"みかん", "C++", "Grok"},

				{"ぶどう", "Python", "ChatGPT"},
				{"ぶどう", "Python", "Gemini"},
				{"ぶどう", "Python", "DeepSeek"},
				{"ぶどう", "Python", "Grok"},
				{"ぶどう", "Go", "ChatGPT"},
				{"ぶどう", "Go", "Gemini"},
				{"ぶどう", "Go", "DeepSeek"},
				{"ぶどう", "Go", "Grok"},
				{"ぶどう", "C++", "ChatGPT"},
				{"ぶどう", "C++", "Gemini"},
				{"ぶどう", "C++", "DeepSeek"},
				{"ぶどう", "C++", "Grok"},

				{"いちご", "Python", "ChatGPT"},
				{"いちご", "Python", "Gemini"},
				{"いちご", "Python", "DeepSeek"},
				{"いちご", "Python", "Grok"},
				{"いちご", "Go", "ChatGPT"},
				{"いちご", "Go", "Gemini"},
				{"いちご", "Go", "DeepSeek"},
				{"いちご", "Go", "Grok"},
				{"いちご", "C++", "ChatGPT"},
				{"いちご", "C++", "Gemini"},
				{"いちご", "C++", "DeepSeek"},
				{"いちご", "C++", "Grok"},
			},
		},
		{
			name: "正常_入力そのものが空集合",
			args: [][]string{},
			// 空集合を1つ持つ
			want: [][]string{{}},
		},
		{
			name: "正常_空集合を含む",
			args: [][]string{
				[]string{"マグロ", "サーモン", "カツオ"},
				[]string{},
				[]string{"犬", "猫"},
			},
			// 空集合を「含む」場合は戻り値も空集合
			want: [][]string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Helper()
			seq := slicesx.CartesianProducts(tc.args...)
			got := slices.Collect(seq)

			if len(got) != len(tc.want) {
				t.Fatalf("len(want): %d, len(got): %d", len(tc.want), len(got))
			}

			for i, gv := range got {
				wv := tc.want[i]
				if !slices.Equal(gv, wv) {
					t.Errorf("i = %d, want: %v, got: %v", i, wv, gv)
					break
				}
			}
		})
	}
}

func TestCounts(t *testing.T) {
	tests := []struct {
		name string
		s    []string
		want map[string]int
	}{
		{
			name: "正常",
			s:    []string{"りんご", "バナナ", "りんご", "さくらんぼ", "バナナ", "りんご"},
			want: map[string]int{
				"りんご":   3,
				"バナナ":   2,
				"さくらんぼ": 1,
			},
		},
		{
			name: "正常 重複要素なし",
			s:    []string{"a", "b", "c"},
			want: map[string]int{
				"a": 1, "b": 1, "c": 1,
			},
		},
		{
			name: "正常 空スライスの入力",
			s:    []string{},
			want: map[string]int{},
		},
		{
			name: "準正常 nilの入力",
			s:    nil,
			want: map[string]int{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := slicesx.Counts(tc.s)
			if !maps.Equal(got, tc.want) {
				t.Errorf("want: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestArgsort(t *testing.T) {
	tests := []struct{
		name string
		s    []int
		want []int
	}{
		{
			name:"正常_重複なし",
			s:[]int{8, 0, 5, 7},
			want:[]int{1, 2, 3, 0},
		},
		{
			name:"正常_重複あり",
			s:[]int{5, 8, 7, 5},
			want:[]int{0, 3, 2, 1},
		},
		{
			name:"正常_空スライス",
			s:[]int{},
			want:[]int{},
		},
		{
			name:"正常_同じ要素のみ",
			s:[]int{10, 10, 10, 10, 10},
			want:[]int{0, 1, 2, 3, 4},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Helper()
			got := slicesx.Argsort(tc.s)
			if !slices.Equal(got, tc.want) {
				t.Errorf("want: %v, got: %v", got, tc.want)
			}
		})
	}
}

func TestArgsortFunc(t *testing.T) {
	tests := []struct{
		name string
		s    []int
		want []int
		f func(a, b int) int
	}{
		{
			name:"正常_降順",
			s:[]int{8, 0, 5, 7},
			want:[]int{0, 3, 2, 1},
			f:func(a, b int) int {
				return cmp.Compare(b, a)
			},
		},
		{
			name:"正常_恒等",
			s:[]int{10, 5, 20, 15, 30, 25},
			want:[]int{0, 1, 2, 3, 4, 5},
			f:func(a, b int) int {
				return 0
			},
		},
		{
			name:"正常_空スライス",
			s:[]int{},
			want:[]int{},
			f:func(a, b int) int {
				return 0
			},
		},
		{
			name:"正常_特殊関数",
			// sumDigitsの結果 {4, 12, 1, 6}
			s:[]int{121, 66, 100, 33},
			want:[]int{2, 0, 3, 1},
			f:func(a, b int) int {
				a = sumDigits(a)
				b = sumDigits(b)
				return cmp.Compare(a, b)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := slicesx.ArgsortFunc(tc.s, tc.f)
			if !slices.Equal(got, tc.want) {
				t.Errorf("want: %v, got: %v", got, tc.want)
			}
		})
	}
}

func sumDigits(n int) int {
	sum := 0

	// 負の数が入力された場合、正の数として扱います
	if n < 0 {
		n = -n
	}

	// 数値が0になるまで繰り返す
	for n > 0 {
		// 10で割った余り（1の位）を足す 
		// 例: 432 ならば、10で割った余りは2
		sum += n % 10

		// 10で割って、桁を1つずらす（整数の割り算なので小数は切り捨て）
		// 例: 432 ならば、43.2 → 43 になる
		n /= 10
	}

	return sum
}