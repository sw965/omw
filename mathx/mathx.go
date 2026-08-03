package mathx

import (
	"github.com/sw965/omw/constraints"
)

// Goの符号付き整数の演算は、桁あふれしても未定義動作にはならず2の補数で折り返す。
// その結果は負にも0にも小さい正の値にもなり得る為、あふれたことを結果から見分けられない。
// 以下の関数は、あふれの有無をokで返す。

// Add は a + b を返す。桁あふれした場合、okはfalseになる。
func Add[T constraints.Signed](a, b T) (T, bool) {
	c := a + b
	// bが正なら和はaより大きく、bが負なら和はaより小さくなるはず
	if (b > 0 && c < a) || (b < 0 && c > a) {
		return 0, false
	}
	return c, true
}

// Sub は a - b を返す。桁あふれした場合、okはfalseになる。
func Sub[T constraints.Signed](a, b T) (T, bool) {
	c := a - b
	// bが正なら差はaより小さく、bが負なら差はaより大きくなるはず
	if (b > 0 && c > a) || (b < 0 && c < a) {
		return 0, false
	}
	return c, true
}

// Mul は a * b を返す。桁あふれした場合、okはfalseになる。
func Mul[T constraints.Signed](a, b T) (T, bool) {
	if a == 0 || b == 0 {
		return 0, true
	}

	// -1での除算は最小値に対して桁あふれする為、下のc/bによる判定を使えない。
	// a*b は -a に等しく、最小値だけが符号を反転できない(-a == a となる)。
	if b == -1 {
		if -a == a {
			return 0, false
		}
		return -a, true
	}

	// 桁あふれしていなければ、積をbで割るとaに戻る
	c := a * b
	if c/b != a {
		return 0, false
	}
	return c, true
}
