package mathx

import (
	"github.com/sw965/omw/constraints"
)

func MulOverflowChecked[T constraints.Signed](a, b T) (T, bool) {
	if a == 0 || b == 0 {
		return 0, true
	}

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
