package constraints

// Signed is a constraint that permits any signed integer type.
// Signed は、任意の符号付き整数型を許可する制約です。
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Unsigned is a constraint that permits any unsigned integer type.
// Unsigned は、任意の符号なし整数型を許可する制約です。
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Integer is a constraint that permits any integer type.
// Integer は、任意の整数型（符号付き・符号なし）を許可する制約です。
type Integer interface {
	Signed | Unsigned
}

// Float is a constraint that permits any floating-point type.
// Float は、任意の浮動小数点型を許可する制約です。
type Float interface {
	~float32 | ~float64
}

// Number is a constraint that permits any numeric type (excluding complex numbers).
// Number は、任意の数値型（複素数を除く）を許可する制約です。
type Number interface {
	Integer | Float
}