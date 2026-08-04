package bitsx

import "testing"

// newTestMatrix は、行ごとに立てる列を指定して Matrix を作る。
func newTestMatrix(t *testing.T, cols int, oneColIdxsPerRow [][]int) *Matrix {
	t.Helper()
	m, err := NewZerosMatrix(len(oneColIdxsPerRow), cols)
	if err != nil {
		t.Fatal(err)
	}

	for r, oneColIdxs := range oneColIdxsPerRow {
		for _, c := range oneColIdxs {
			if err := m.Set(r, c); err != nil {
				t.Fatal(err)
			}
		}
	}
	return m
}

func callDotGo(left, right *Matrix) []int {
	results := make([]int, left.Rows*right.Rows)
	dotGo(left.Data, right.Data, left.Rows, right.Rows, left.Cols, left.Stride(), results)
	return results
}

func assertResults(t *testing.T, name string, got, want []int) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("%s: 長さの不一致: got = %d, want %d", name, len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("%s: [%d] = %d, want %d (全体: got = %v, want %v)", name, i, got[i], want[i], got, want)
		}
	}
}

// テスト要件：TST-001
// dotGo は「left行rとright行cで、値が一致する列の数」を results[r*rightRows+c] に書く。
func TestDotGoExpectedValues(t *testing.T) {
	tests := []struct {
		name  string
		cols  int
		left  [][]int
		right [][]int
		want  []int
	}{
		{
			name: "2x5と3x5",
			cols: 5,
			left: [][]int{
				{0, 2, 4}, // 10101
				{1, 3},    // 01010
			},
			right: [][]int{
				{0, 4},    // 10001
				{0, 2, 4}, // 10101
				{1, 2, 3}, // 01110
			},

			// 次の順番で、一致するビットを数える。
			// left[0]とright[0] = 10101と10001。一致するビットは4
			// left[0]とright[1] = 10101と10101。一致するビットは5
			// left[0]とright[2] = 10101と01110。一致するビットは1
			// left[1]とright[0] = 01010と10001。一致するビットは1
			// left[1]とright[1] = 01010と10101。一致するビットは0
			// left[1]とright[2] = 01010と01110。一致するビットは4
			//
			// 二次元表の結果は次の通り
			//              right[0] right[1] right[2]
			// left[0]          4        5        1
			// left[1]          1        0        4
			want: []int{
				4, 5, 1,
				1, 0, 4,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			left := newTestMatrix(t, tt.cols, tt.left)
			right := newTestMatrix(t, tt.cols, tt.right)
			assertResults(t, "dotGo", callDotGo(left, right), tt.want)
		})
	}
}
