package constraints_test

import (
	"testing"

	omwconstraints "github.com/sw965/omw/constraints"
)

type customInt int
type customUint8 uint8
type customUint64 uint64
type customFloat32 float32
type customFloat64 float64

// コンパイルが通るかだけをテストするので、中身は空
func compileAsSigned[T omwconstraints.Signed]()     {}
func compileAsUnsigned[T omwconstraints.Unsigned]() {}
func compileAsInteger[T omwconstraints.Integer]()   {}
func compileAsFloat[T omwconstraints.Float]()       {}
func compileAsNumber[T omwconstraints.Number]()     {}

// テスト要件：TST-001
func TestInt(t *testing.T) {
	compileAsSigned[int]()
	compileAsSigned[customInt]()
	compileAsInteger[int]()
	compileAsInteger[customInt]()
	compileAsNumber[int]()
	compileAsNumber[customInt]()
}

// テスト要件：TST-002
func TestUint8(t *testing.T) {
	compileAsUnsigned[uint8]()
	compileAsUnsigned[customUint8]()
	compileAsInteger[uint8]()
	compileAsInteger[customUint8]()
	compileAsNumber[uint8]()
	compileAsNumber[customUint8]()
}

// テスト要件：TST-003
func TestUint64(t *testing.T) {
	compileAsUnsigned[uint64]()
	compileAsUnsigned[customUint64]()
	compileAsInteger[uint64]()
	compileAsInteger[customUint64]()
	compileAsNumber[uint64]()
	compileAsNumber[customUint64]()
}

// テスト要件：TST-004
func TestFloat32(t *testing.T) {
	compileAsFloat[float32]()
	compileAsFloat[customFloat32]()
	compileAsNumber[float32]()
	compileAsNumber[customFloat32]()
}

// テスト要件：TST-005
func TestFloat64(t *testing.T) {
	compileAsFloat[float64]()
	compileAsFloat[customFloat64]()
	compileAsNumber[float64]()
	compileAsNumber[customFloat64]()
}
