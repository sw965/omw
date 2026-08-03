package mathx_test

import (
	"math"
	"testing"

	"github.com/sw965/omw/mathx"
)

func TestAdd(t *testing.T) {
	t.Run("正常", func(t *testing.T) {
		cases := []struct{ a, b, want int }{
			{0, 0, 0},
			{2, 3, 5},
			{-2, -3, -5},
			{-2, 3, 1},
			{math.MaxInt, 0, math.MaxInt},
			{math.MaxInt - 1, 1, math.MaxInt},
			{math.MinInt, 0, math.MinInt},
			{math.MinInt + 1, -1, math.MinInt},
		}
		for _, c := range cases {
			got, ok := mathx.Add(c.a, c.b)
			if !ok {
				t.Errorf("Add(%d, %d): 桁あふれと判定された", c.a, c.b)
			}
			if got != c.want {
				t.Errorf("Add(%d, %d) = %d, want %d", c.a, c.b, got, c.want)
			}
		}
	})

	t.Run("異常_桁あふれ", func(t *testing.T) {
		cases := []struct{ a, b int }{
			{math.MaxInt, 1},
			{1, math.MaxInt},
			{math.MaxInt, math.MaxInt},
			{math.MinInt, -1},
			{-1, math.MinInt},
			{math.MinInt, math.MinInt},
		}
		for _, c := range cases {
			if _, ok := mathx.Add(c.a, c.b); ok {
				t.Errorf("Add(%d, %d): 桁あふれを検出できなかった", c.a, c.b)
			}
		}
	})

	t.Run("正常_int8", func(t *testing.T) {
		if _, ok := mathx.Add[int8](127, 1); ok {
			t.Error("int8の桁あふれを検出できなかった")
		}
		if got, ok := mathx.Add[int8](126, 1); !ok || got != 127 {
			t.Errorf("Add[int8](126, 1) = %d, %v", got, ok)
		}
	})
}

func TestSub(t *testing.T) {
	t.Run("正常", func(t *testing.T) {
		cases := []struct{ a, b, want int }{
			{0, 0, 0},
			{5, 3, 2},
			{3, 5, -2},
			{-3, -5, 2},
			{math.MaxInt, 0, math.MaxInt},
			{math.MaxInt, 1, math.MaxInt - 1},
			{math.MinInt, 0, math.MinInt},
			{-1, math.MinInt, math.MaxInt},
		}
		for _, c := range cases {
			got, ok := mathx.Sub(c.a, c.b)
			if !ok {
				t.Errorf("Sub(%d, %d): 桁あふれと判定された", c.a, c.b)
			}
			if got != c.want {
				t.Errorf("Sub(%d, %d) = %d, want %d", c.a, c.b, got, c.want)
			}
		}
	})

	t.Run("異常_桁あふれ", func(t *testing.T) {
		cases := []struct{ a, b int }{
			{math.MinInt, 1},
			{math.MaxInt, -1},
			{0, math.MinInt},
			{math.MinInt, math.MaxInt},
			{math.MaxInt, math.MinInt},
		}
		for _, c := range cases {
			if _, ok := mathx.Sub(c.a, c.b); ok {
				t.Errorf("Sub(%d, %d): 桁あふれを検出できなかった", c.a, c.b)
			}
		}
	})

	t.Run("正常_int8", func(t *testing.T) {
		if _, ok := mathx.Sub[int8](-128, 1); ok {
			t.Error("int8の桁あふれを検出できなかった")
		}
		if got, ok := mathx.Sub[int8](-127, 1); !ok || got != -128 {
			t.Errorf("Sub[int8](-127, 1) = %d, %v", got, ok)
		}
	})
}

func TestMul(t *testing.T) {
	t.Run("正常", func(t *testing.T) {
		cases := []struct{ a, b, want int }{
			{0, 0, 0},
			{0, math.MaxInt, 0},
			{math.MinInt, 0, 0},
			{3, 5, 15},
			{-3, 5, -15},
			{3, -5, -15},
			{-3, -5, 15},
			{1, math.MaxInt, math.MaxInt},
			{math.MaxInt, 1, math.MaxInt},
			{1, math.MinInt, math.MinInt},
			{math.MinInt, 1, math.MinInt},
			{-1, math.MaxInt, math.MinInt + 1},
			{math.MaxInt, -1, math.MinInt + 1},
			{-1, -1, 1},
		}
		for _, c := range cases {
			got, ok := mathx.Mul(c.a, c.b)
			if !ok {
				t.Errorf("Mul(%d, %d): 桁あふれと判定された", c.a, c.b)
			}
			if got != c.want {
				t.Errorf("Mul(%d, %d) = %d, want %d", c.a, c.b, got, c.want)
			}
		}
	})

	t.Run("異常_桁あふれ", func(t *testing.T) {
		cases := []struct {
			name string
			a, b int
		}{
			// 折り返した結果が負になる
			{"MaxInt*2", math.MaxInt, 2},
			{"2*MaxInt", 2, math.MaxInt},
			// 折り返した結果が0になる
			{"2^62*4", 1 << 62, 4},
			// 折り返した結果が小さい正の値になる
			{"(2^62+1)*4", 1<<62 + 1, 4},
			// -1での除算が桁あふれする為、c/bによる判定が使えない組み合わせ
			{"MinInt*-1", math.MinInt, -1},
			{"-1*MinInt", -1, math.MinInt},
			{"MinInt*2", math.MinInt, 2},
			{"MaxInt*MaxInt", math.MaxInt, math.MaxInt},
			{"MinInt*MinInt", math.MinInt, math.MinInt},
		}
		for _, c := range cases {
			if _, ok := mathx.Mul(c.a, c.b); ok {
				t.Errorf("%s: Mul(%d, %d) が桁あふれを検出できなかった", c.name, c.a, c.b)
			}
		}
	})

	t.Run("正常_int8", func(t *testing.T) {
		if _, ok := mathx.Mul[int8](64, 2); ok {
			t.Error("int8の桁あふれを検出できなかった")
		}
		if _, ok := mathx.Mul[int8](-128, -1); ok {
			t.Error("int8の最小値の符号反転を検出できなかった")
		}
		if got, ok := mathx.Mul[int8](63, 2); !ok || got != 126 {
			t.Errorf("Mul[int8](63, 2) = %d, %v", got, ok)
		}
	})
}

// 総当たりで、より広い型での計算結果と一致することを確認する。
func TestMatchesWiderType(t *testing.T) {
	const lo, hi = -130, 130
	for a := lo; a <= hi; a++ {
		for b := lo; b <= hi; b++ {
			a8, b8 := int8(a), int8(b)
			if int(a8) != a || int(b8) != b {
				continue // int8で表せない値は対象外
			}

			gotAdd, okAdd := mathx.Add(a8, b8)
			wantAdd := a + b
			if okAdd != (wantAdd >= math.MinInt8 && wantAdd <= math.MaxInt8) {
				t.Fatalf("Add(%d, %d): ok = %v, 実際の和 = %d", a8, b8, okAdd, wantAdd)
			}
			if okAdd && int(gotAdd) != wantAdd {
				t.Fatalf("Add(%d, %d) = %d, want %d", a8, b8, gotAdd, wantAdd)
			}

			gotSub, okSub := mathx.Sub(a8, b8)
			wantSub := a - b
			if okSub != (wantSub >= math.MinInt8 && wantSub <= math.MaxInt8) {
				t.Fatalf("Sub(%d, %d): ok = %v, 実際の差 = %d", a8, b8, okSub, wantSub)
			}
			if okSub && int(gotSub) != wantSub {
				t.Fatalf("Sub(%d, %d) = %d, want %d", a8, b8, gotSub, wantSub)
			}

			gotMul, okMul := mathx.Mul(a8, b8)
			wantMul := a * b
			if okMul != (wantMul >= math.MinInt8 && wantMul <= math.MaxInt8) {
				t.Fatalf("Mul(%d, %d): ok = %v, 実際の積 = %d", a8, b8, okMul, wantMul)
			}
			if okMul && int(gotMul) != wantMul {
				t.Fatalf("Mul(%d, %d) = %d, want %d", a8, b8, gotMul, wantMul)
			}
		}
	}
}
