package bitsx

import (
	"math"
	"math/bits"
	"math/rand/v2"
	"testing"
	"testing/quick"
)

func callDotGo(left, right *Matrix) []int {
	results := make([]int, left.rows*right.rows)
	dotGo(left.data, right.data, left.rows, right.rows, left.cols, left.Stride(), results)
	return results
}

func callDotAVX512(left, right *Matrix) []int {
	results := make([]int, left.rows*right.rows)
	dotAVX512(&left.data[0], &right.data[0], left.rows, right.rows, left.cols, left.Stride(), &results[0])
	return results
}

func callDotTernaryGo(value, sign, nonZero *Matrix) []int {
	results := make([]int, value.rows*sign.rows)
	dotTernaryGo(value.data, sign.data, nonZero.data, value.rows, sign.rows, value.Stride(), results)
	return results
}

func callDotTernaryAVX512(value, sign, nonZero *Matrix) []int {
	results := make([]int, value.rows*sign.rows)
	dotTernaryAVX512(&value.data[0], &sign.data[0], &nonZero.data[0], value.rows, sign.rows, value.Stride(), &results[0])
	return results
}

func TestXorPopcntGoExpectedValues(t *testing.T) {
	tests := []struct {
		name string
		a    []uint64
		b    []uint64
		want int
	}{
		{
			name: "空の入力",
			want: 0,
		},
		{
			name: "複数ワード",
			a: []uint64{
				0b11110000_11110000_11110000_11110000_11110000_11110000_11110000_11110000,
				0b10101010_10101010_10101010_10101010_10101010_10101010_10101010_10101010,
			},
			b: []uint64{
				0b00001111_00001111_00001111_00001111_00001111_00001111_00001111_00001111,
				0,
			},
			want: 96,
		},
		{
			name: "同一の入力",
			a: []uint64{
				0b00000001_00100011_01000101_01100111_10001001_10101011_11001101_11101111,
				0b11111110_11011100_10111010_10011000_01110110_01010100_00110010_00010000,
			},
			b: []uint64{
				0b00000001_00100011_01000101_01100111_10001001_10101011_11001101_11101111,
				0b11111110_11011100_10111010_10011000_01110110_01010100_00110010_00010000,
			},
			want: 0,
		},
		{
			name: "全ビット不一致",
			a:    []uint64{0, ^uint64(0)},
			b:    []uint64{^uint64(0), 0},
			want: 128,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := xorPopcntGo(tt.a, tt.b)
			if got != tt.want {
				t.Fatalf("値の不一致: got = %d want = %d", got, tt.want)
			}
		})
	}
}

func TestXorPopcntGoBitwiseAgreement(t *testing.T) {
	// uint8の引数には、0～255のいずれかが代入される
	property := func(words uint8, seed1, seed2 uint64) bool {
		// 0～255を1～32に変換
		wordCount := int(words%32 + 1)

		rng := rand.New(rand.NewPCG(seed1, seed2))
		a := make([]uint64, wordCount)
		b := make([]uint64, wordCount)
		for i := range wordCount {
			a[i] = rng.Uint64()
			b[i] = rng.Uint64()
		}

		// 素朴(愚直)な実装で比較する
		want := 0
		for i := range wordCount {
			x := a[i] ^ b[i]
			for j := range 64 {
				want += int((x >> uint(j)) & 1)
			}
		}

		return xorPopcntGo(a, b) == want
	}

	if err := quick.Check(property, &quick.Config{MaxCount: 512}); err != nil {
		t.Fatalf("%v", err)
	}
}

func assertResults(t *testing.T, name string, got, want []int) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("%s: 長さの不一致: got = %d want %d", name, len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("%s: [%d] = %d want %d (全体: got = %v want %v)", name, i, got[i], want[i], got, want)
		}
	}
}

func TestDotGoExpectedValues(t *testing.T) {
	tests := []struct {
		name  string
		cols  int
		left  [][]int
		right [][]int
		want  []int
	}{
		{
			name: "2×5・3×5",
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

			// leftの各行とrightの各行について、対応する列のビット値が一致する数を数える。(XNOR)
			// left[0]とright[0] = 10101と10001。一致するビットは4
			// left[0]とright[1] = 10101と10101。一致するビットは5
			// left[0]とright[2] = 10101と01110。一致するビットは1
			// left[1]とright[0] = 01010と10001。一致するビットは1
			// left[1]とright[1] = 01010と10101。一致するビットは0
			// left[1]とright[2] = 01010と01110。一致するビットは4
			//
			//              right[0] right[1] right[2]
			// left[0]          4        5        1
			// left[1]          1        0        4
			want: []int{
				4, 5, 1,
				1, 0, 4,
			},
		},

		{
			name: "1×65・3×65_ワード境界",
			cols: 65,
			// 下記のコメントは、8ビットごとに区切り、左端を列0、右端を列64とする。
			left: [][]int{{0, 64}}, // 10000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000_1
			right: [][]int{
				{0},  // 10000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000_0
				{64}, // 00000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000_1
				{},   // 00000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000_0
			},
			want: []int{64, 64, 63},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			left := NewMatrixForTest(t, tt.cols, tt.left)
			right := NewMatrixForTest(t, tt.cols, tt.right)
			assertResults(t, "dotGo", callDotGo(left, right), tt.want)
		})
	}
}

func TestDotGoResultLength(t *testing.T) {
	// uint8の引数には、0～255のいずれかが代入される
	property := func(lRows, rRows, cols uint8) bool {
		// 0～255を1～16に変換
		leftRows := int(lRows%16 + 1)
		rightRows := int(rRows%16 + 1)
		// 0～255を1～65に変換
		columns := int(cols%65 + 1)

		left := NewMatrixForTest(t, columns, make([][]int, leftRows))
		right := NewMatrixForTest(t, columns, make([][]int, rightRows))
		got := callDotGo(left, right)

		return len(got) == leftRows*rightRows
	}

	if err := quick.Check(property, &quick.Config{MaxCount: 512}); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestDotGoValueRange(t *testing.T) {
	// uint8の引数には、0～255のいずれかが代入される
	property := func(lRows, rRows, cols uint8, seed1, seed2 uint64) bool {
		// 0～255を1～16に変換
		leftRows := int(lRows%16 + 1)
		rightRows := int(rRows%16 + 1)
		// 0～255を1～65に変換
		columns := int(cols%65 + 1)

		rng := rand.New(rand.NewPCG(seed1, seed2))
		left, err := NewRandMatrix(leftRows, columns, 0, rng)
		if err != nil {
			return false
		}
		right, err := NewRandMatrix(rightRows, columns, 0, rng)
		if err != nil {
			return false
		}

		got := callDotGo(left, right)
		for _, value := range got {
			if value < 0 || value > columns {
				return false
			}
		}
		return true
	}

	if err := quick.Check(property, &quick.Config{MaxCount: 512}); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestDotGoTranspose(t *testing.T) {
	// uint8の引数には、0～255のいずれかが代入される
	property := func(lRows, rRows, cols uint8, seed1, seed2 uint64) bool {
		// 0～255を1～16に変換
		leftRows := int(lRows%16 + 1)
		rightRows := int(rRows%16 + 1)
		// 0～255を1～65に変換
		columns := int(cols%65 + 1)

		rng := rand.New(rand.NewPCG(seed1, seed2))
		left, err := NewRandMatrix(leftRows, columns, 0, rng)
		if err != nil {
			return false
		}
		right, err := NewRandMatrix(rightRows, columns, 0, rng)
		if err != nil {
			return false
		}

		gotLeftRight := callDotGo(left, right)
		gotRightLeft := callDotGo(right, left)

		// leftとrightを入れ替えた結果は、元の結果を転置したものになる
		for r := range left.rows {
			for c := range right.rows {
				leftRightIdx := r*rightRows + c
				rightLeftIdx := c*leftRows + r
				if gotLeftRight[leftRightIdx] != gotRightLeft[rightLeftIdx] {
					return false
				}
			}
		}
		return true
	}

	if err := quick.Check(property, &quick.Config{MaxCount: 512}); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestDotGoBitwiseAgreement(t *testing.T) {
	// uint8の引数には、0～255のいずれかが代入される
	property := func(lRows, rRows, cols uint8, seed1, seed2 uint64) bool {
		// 0～255を1～4に変換
		leftRows := int(lRows%4 + 1)
		rightRows := int(rRows%4 + 1)
		// 0～255を1～256に変換
		columns := int(cols) + 1

		rng := rand.New(rand.NewPCG(seed1, seed2))
		left, err := NewRandMatrix(leftRows, columns, 0, rng)
		if err != nil {
			return false
		}
		right, err := NewRandMatrix(rightRows, columns, 0, rng)
		if err != nil {
			return false
		}

		got := callDotGo(left, right)

		// 素朴(愚直)な実装で比較する
		for r := range leftRows {
			for c := range rightRows {
				match := 0
				for k := range columns {
					lb, err := left.Bit(r, k)
					if err != nil {
						return false
					}
					rb, err := right.Bit(c, k)
					if err != nil {
						return false
					}
					if lb == rb {
						match++
					}
				}
				if got[r*rightRows+c] != match {
					return false
				}
			}
		}
		return true
	}

	if err := quick.Check(property, &quick.Config{MaxCount: 512}); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestDotTernaryGoExpectedValues(t *testing.T) {
	tests := []struct {
		name    string
		cols    int
		value   [][]int
		sign    [][]int
		nonZero [][]int
		want    []int
	}{
		{
			name: "2×5・3×5",
			cols: 5,
			value: [][]int{
				{0, 2, 4}, // 10101
				{1, 3},    // 01010
			},
			sign: [][]int{
				{0, 4},    // 10001
				{0, 2, 4}, // 10101
				{1, 2, 3}, // 01110
			},
			nonZero: [][]int{
				{0, 1, 3, 4}, // 11011
				{1, 2, 3},    // 01110
				{0, 2, 4},    // 10101
			},

			// valueの各行とsignの各行について、nonZeroが1の列だけを対象に、valueとsignが一致すれば+1、不一致なら-1を加算する。
			// value[0]とsign[0]とnonZero[0] = 10101と10001と11011は、+1+1+1+1=4。
			// value[0]とsign[1]とnonZero[1] = 10101と10101と01110は、+1+1+1=3。
			// value[0]とsign[2]とnonZero[2] = 10101と01110と10101は、-1+1-1=-1。
			// value[1]とsign[0]とnonZero[0] = 01010と10001と11011は、-1-1-1-1=-4。
			// value[1]とsign[1]とnonZero[1] = 01010と10101と01110は、-1-1-1=-3。
			// value[1]とsign[2]とnonZero[2] = 01010と01110と10101は、+1-1+1=1。
			//
			//              sign[0] sign[1] sign[2]
			// value[0]         4       3     -1
			// value[1]        -4      -3      1
			want: []int{
				4, 3, -1,
				-4, -3, 1,
			},
		},

		{
			name: "1×65・3×65_ワード境界",
			cols: 65,
			// 下記のコメントは、8ビットごとに区切り、左端を列0、右端を列64とする。
			value: [][]int{{0, 64}}, // 10000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000_1
			sign: [][]int{
				{0},     // 10000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000_0
				{64},    // 00000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000_1
				{0, 64}, // 10000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000_1
			},
			nonZero: [][]int{
				{0, 64}, // 10000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000_1
				{0, 64}, // 10000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000_1
				{0, 64}, // 10000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000_1
			},
			want: []int{0, 0, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value := NewMatrixForTest(t, tt.cols, tt.value)
			sign := NewMatrixForTest(t, tt.cols, tt.sign)
			nonZero := NewMatrixForTest(t, tt.cols, tt.nonZero)
			assertResults(t, "dotTernaryGo", callDotTernaryGo(value, sign, nonZero), tt.want)
		})
	}
}

func TestDotTernaryGoResultLength(t *testing.T) {
	// uint8の引数には、0～255のいずれかが代入される
	property := func(vRows, sRows, cols uint8) bool {
		// 0～255を1～16に変換
		valueRows := int(vRows%16 + 1)
		signRows := int(sRows%16 + 1)
		// 0～255を1～65に変換
		columns := int(cols%65 + 1)

		value := NewMatrixForTest(t, columns, make([][]int, valueRows))
		sign := NewMatrixForTest(t, columns, make([][]int, signRows))
		nonZero := NewMatrixForTest(t, columns, make([][]int, signRows))
		got := callDotTernaryGo(value, sign, nonZero)

		return len(got) == valueRows*signRows
	}

	if err := quick.Check(property, &quick.Config{MaxCount: 512}); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestDotTernaryGoValueRange(t *testing.T) {
	// uint8の引数には、0～255のいずれかが代入される
	property := func(vRows, sRows, cols uint8, seed1, seed2 uint64) bool {
		// 0～255を1～16に変換
		valueRows := int(vRows%16 + 1)
		signRows := int(sRows%16 + 1)
		// 0～255を1～65に変換
		columns := int(cols%65 + 1)

		rng := rand.New(rand.NewPCG(seed1, seed2))
		value, err := NewRandMatrix(valueRows, columns, 0, rng)
		if err != nil {
			return false
		}
		sign, err := NewRandMatrix(signRows, columns, 0, rng)
		if err != nil {
			return false
		}
		nonZero, err := NewRandMatrix(signRows, columns, 0, rng)
		if err != nil {
			return false
		}

		got := callDotTernaryGo(value, sign, nonZero)
		stride := nonZero.Stride()

		nonZeroCounts := make([]int, signRows)
		for c := range signRows {
			rowWords := nonZero.data[c*stride : (c+1)*stride]
			for _, w := range rowWords {
				nonZeroCounts[c] += bits.OnesCount64(w)
			}
		}

		for i, gotValue := range got {
			c := i % signRows
			limit := nonZeroCounts[c]
			if gotValue < -limit || gotValue > limit {
				return false
			}
		}
		return true
	}

	if err := quick.Check(property, &quick.Config{MaxCount: 512}); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestDotTernaryGoBitwiseAgreement(t *testing.T) {
	// uint8の引数には、0～255のいずれかが代入される
	property := func(vRows, sRows, cols uint8, seed1, seed2 uint64) bool {
		// 0～255を1～4に変換
		valueRows := int(vRows%4 + 1)
		signRows := int(sRows%4 + 1)
		// 0～255を1～256に変換
		columns := int(cols) + 1

		rng := rand.New(rand.NewPCG(seed1, seed2))
		value, err := NewRandMatrix(valueRows, columns, 0, rng)
		if err != nil {
			return false
		}
		sign, err := NewRandMatrix(signRows, columns, 0, rng)
		if err != nil {
			return false
		}
		nonZero, err := NewRandMatrix(signRows, columns, 0, rng)
		if err != nil {
			return false
		}

		got := callDotTernaryGo(value, sign, nonZero)

		// 素朴(愚直)な実装で比較する
		for r := range valueRows {
			for c := range signRows {
				sum := 0
				for k := range columns {
					nz, err := nonZero.Bit(c, k)
					if err != nil {
						return false
					}
					if nz == 0 {
						continue
					}
					vb, err := value.Bit(r, k)
					if err != nil {
						return false
					}
					sb, err := sign.Bit(c, k)
					if err != nil {
						return false
					}
					if vb == sb {
						sum++
					} else {
						sum--
					}
				}
				if got[r*signRows+c] != sum {
					return false
				}
			}
		}
		return true
	}

	if err := quick.Check(property, &quick.Config{MaxCount: 512}); err != nil {
		t.Fatalf("%v", err)
	}
}

func FuzzXorPopcntAVX512VsGo(f *testing.F) {
	if !useAVX512 {
		f.Skip("AVX512命令は非対応の環境")
	}

	seeds := []struct {
		words        uint8
		seed1, seed2 uint64
	}{
		{0, 0, 0}, // 1ワード
		{6, 0x123456789ABCDEF0, 0x0FEDCBA987654321},   // 7ワード
		{7, 0x5555555555555555, 0xAAAAAAAAAAAAAAAA},   // 8ワード
		{8, 0x0123456789ABCDEF, 0xFEDCBA9876543210},   // 9ワード
		{11, 0x2222222222222222, 0xDDDDDDDDDDDDDDDD},  // 12ワード(1ループ+端数4)
		{15, 0, ^uint64(0)},                           // 16ワード
		{31, ^uint64(0), ^uint64(0)},                  // 32ワード
		{63, 0x13579BDF2468ACE0, 0x0ECA8642FDB97531},  // 64ワード
		{127, 0xAAAAAAAAAAAAAAAA, 0x5555555555555555}, // 128ワード
		{255, ^uint64(0), 0},                          // 256ワード
	}
	for _, s := range seeds {
		f.Add(s.words, s.seed1, s.seed2)
	}

	f.Fuzz(func(t *testing.T, words uint8, seed1, seed2 uint64) {
		wordCount := int(words) + 1
		rng := rand.New(rand.NewPCG(seed1, seed2))
		a := make([]uint64, wordCount)
		b := make([]uint64, wordCount)
		for i := range wordCount {
			a[i] = rng.Uint64()
			b[i] = rng.Uint64()
		}

		gotGo := xorPopcntGo(a, b)
		gotAVX512 := xorPopcntAVX512(&a[0], &b[0], wordCount)
		if gotAVX512 != gotGo {
			t.Fatalf("xorPopcntAVX512とxorPopcntGoの不一致: gotAVX512 = %d, gotGo = %d", gotAVX512, gotGo)
		}
	})
}

func FuzzDotAVX512VsGo(f *testing.F) {
	if !useAVX512 {
		f.Skip("AVX512命令は非対応の環境")
	}

	seeds := []struct {
		lRows, rRows uint8
		cols         uint16
		seed1, seed2 uint64
	}{
		{0, 0, 0, 0, 0}, // 極小ケース (1x1行列, 1列, 0シード)
		{0, 15, 63, 0x123456789ABCDEF0, 0x0FEDCBA987654321}, // 不均衡ケース (1x16行列, 64列/1ワード)
		{1, 2, 64, 0x5555555555555555, 0xAAAAAAAAAAAAAAAA},  // ワード跨ぎ境界 (2x3行列, 65列/1ワード+1)
		{2, 3, 127, 0x0123456789ABCDEF, 0xFEDCBA9876543210}, // マルチワード境界 (3x4行列, 128列/2ワード)
		{3, 4, 128, 0x3333333333333333, 0xCCCCCCCCCCCCCCCC}, // ワード数境界 (4x5行列, 129列/3ワード)
		{4, 5, 448, 0x0F0F0F0F0F0F0F0F, 0xF0F0F0F0F0F0F0F0}, // メインループ境界(stride=8, 端数無し) (5x6行列, 449列/8ワード)
		{5, 6, 512, 0x00FF00FF00FF00FF, 0xFF00FF00FF00FF00}, // メインループ+端数(stride=9) (6x7行列, 513列/8ワード+1)
		{6, 7, 960, 0x3C3C3C3C3C3C3C3C, 0xC3C3C3C3C3C3C3C3}, // メインループ2回転(stride=16, 端数無し) (7x8行列, 961列/16ワード)
		{15, 15, 255, ^uint64(0), ^uint64(0)},               // 最大ケース (16x16行列, 256列/4ワード, 最大値シード)
	}
	for _, s := range seeds {
		f.Add(s.lRows, s.rRows, s.cols, s.seed1, s.seed2)
	}

	// uint8の引数には、0～255のいずれかが代入される
	f.Fuzz(func(t *testing.T, lRows, rRows uint8, cols uint16, seed1, seed2 uint64) {
		// 0～255を1～16に変換
		leftRows := int(lRows%16 + 1)
		rightRows := int(rRows%16 + 1)
		// uint16の引数には、0～65535のいずれかが代入される
		// 0～65535を1～65536に変換
		columns := int(cols) + 1

		rng := rand.New(rand.NewPCG(seed1, seed2))
		left, err := NewRandMatrix(leftRows, columns, 0, rng)
		if err != nil {
			t.Fatalf("%v", err)
		}
		right, err := NewRandMatrix(rightRows, columns, 0, rng)
		if err != nil {
			t.Fatalf("%v", err)
		}

		gotGo := callDotGo(left, right)
		gotAVX512 := callDotAVX512(left, right)

		assertResults(t, "dotAVX512 vs dotGo", gotAVX512, gotGo)
	})
}

func FuzzDotTernaryAVX512VsGo(f *testing.F) {
	if !useAVX512 {
		f.Skip("AVX512命令は非対応の環境")
	}

	seeds := []struct {
		vRows, sRows uint8
		cols         uint16
		seed1, seed2 uint64
	}{
		{0, 0, 0, 0, 0}, // 極小ケース (1x1行列, 1列, 0シード)
		{0, 15, 63, 0x123456789ABCDEF0, 0x0FEDCBA987654321}, // 不均衡ケース (1x16行列, 64列/1ワード)
		{1, 2, 64, 0x5555555555555555, 0xAAAAAAAAAAAAAAAA},  // ワード跨ぎ境界 (2x3行列, 65列/1ワード+1)
		{2, 3, 127, 0x0123456789ABCDEF, 0xFEDCBA9876543210}, // マルチワード境界 (3x4行列, 128列/2ワード)
		{3, 8, 128, 0x6666666666666666, 0x9999999999999999}, // ブロック境界(8フル+1端数)+ワード数境界 (4x9行列, 129列/3ワード)
		{4, 5, 448, 0x0F0F0F0F0F0F0F0F, 0xF0F0F0F0F0F0F0F0}, // メインループ境界(stride=8, 端数無し) (5x6行列, 449列/8ワード)
		{5, 6, 512, 0x00FF00FF00FF00FF, 0xFF00FF00FF00FF00}, // メインループ+端数(stride=9) (6x7行列, 513列/8ワード+1)
		{6, 7, 960, 0x3C3C3C3C3C3C3C3C, 0xC3C3C3C3C3C3C3C3}, // メインループ2回転(stride=16, 端数無し) (7x8行列, 961列/16ワード)
		{15, 15, 255, ^uint64(0), ^uint64(0)},               // 最大ケース (16x16行列, 256列/4ワード, 最大値シード)
	}
	for _, s := range seeds {
		f.Add(s.vRows, s.sRows, s.cols, s.seed1, s.seed2)
	}

	// uint8の引数には、0～255のいずれかが代入される
	f.Fuzz(func(t *testing.T, vRows, sRows uint8, cols uint16, seed1, seed2 uint64) {
		// 0～255を1～16に変換
		valueRows := int(vRows%16 + 1)
		signRows := int(sRows%16 + 1)
		// uint16の引数には、0～65535のいずれかが代入される
		// 0～65535を1～65536に変換
		columns := int(cols) + 1

		rng := rand.New(rand.NewPCG(seed1, seed2))
		value, err := NewRandMatrix(valueRows, columns, 0, rng)
		if err != nil {
			t.Fatalf("%v", err)
		}
		sign, err := NewRandMatrix(signRows, columns, 0, rng)
		if err != nil {
			t.Fatalf("%v", err)
		}
		nonZero, err := NewRandMatrix(signRows, columns, 0, rng)
		if err != nil {
			t.Fatalf("%v", err)
		}

		gotGo := callDotTernaryGo(value, sign, nonZero)
		gotAVX512 := callDotTernaryAVX512(value, sign, nonZero)

		assertResults(t, "dotTernaryAVX512 vs dotTernaryGo", gotAVX512, gotGo)
	})
}

const (
	benchXorPopcntCols = 8192

	benchDotRows = 768
	benchDotCols = 768

	benchDotTernaryRows = 768
	benchDotTernaryCols = 768
)

func newBenchMatrix(b *testing.B, rows, cols int, rng *rand.Rand) *Matrix {
	b.Helper()
	m, err := NewRandMatrix(rows, cols, 0, rng)
	if err != nil {
		b.Fatalf("%v", err)
	}
	return m
}

func skipIfNoAVX512(b *testing.B) {
	b.Helper()
	if !useAVX512 {
		b.Skip("AVX512命令は非対応の環境")
	}
}

func BenchmarkXorPopcntGo(b *testing.B) {
	rng := rand.New(rand.NewPCG(1, 2))
	a := newBenchMatrix(b, 1, benchXorPopcntCols, rng)
	c := newBenchMatrix(b, 1, benchXorPopcntCols, rng)

	for b.Loop() {
		xorPopcntGo(a.data, c.data)
	}
}

// 環境: AMD Ryzen 5 8600G
// 形状: 1行 × 8192列
// 期待: Go版に対して、おおよそ5〜7倍速い
// ※ 上記のコメントは勝手に変更・削除してはならない
func BenchmarkXorPopcntAVX512(b *testing.B) {
	skipIfNoAVX512(b)
	rng := rand.New(rand.NewPCG(1, 2))
	a := newBenchMatrix(b, 1, benchXorPopcntCols, rng)
	c := newBenchMatrix(b, 1, benchXorPopcntCols, rng)

	for b.Loop() {
		xorPopcntAVX512(&a.data[0], &c.data[0], len(a.data))
	}
}

func BenchmarkDotGo(b *testing.B) {
	rng := rand.New(rand.NewPCG(3, 4))
	left := newBenchMatrix(b, benchDotRows, benchDotCols, rng)
	right := newBenchMatrix(b, benchDotRows, benchDotCols, rng)
	results := make([]int, left.rows*right.rows)

	for b.Loop() {
		dotGo(left.data, right.data, left.rows, right.rows, left.cols, left.Stride(), results)
	}
}

// 環境: AMD Ryzen 5 8600G
// 形状: 768行 × 768列
// 期待: Go版に対して、おおよそ4〜5倍速い
// ※ 上記のコメントは勝手に変更・削除してはならない
func BenchmarkDotAVX512(b *testing.B) {
	skipIfNoAVX512(b)
	rng := rand.New(rand.NewPCG(3, 4))
	left := newBenchMatrix(b, benchDotRows, benchDotCols, rng)
	right := newBenchMatrix(b, benchDotRows, benchDotCols, rng)
	results := make([]int, left.rows*right.rows)

	for b.Loop() {
		dotAVX512(&left.data[0], &right.data[0], left.rows, right.rows, left.cols, left.Stride(), &results[0])
	}
}

func BenchmarkDotTernaryGo(b *testing.B) {
	rng := rand.New(rand.NewPCG(5, 6))
	value := newBenchMatrix(b, benchDotTernaryRows, benchDotTernaryCols, rng)
	sign := newBenchMatrix(b, benchDotTernaryRows, benchDotTernaryCols, rng)
	nonZero := newBenchMatrix(b, benchDotTernaryRows, benchDotTernaryCols, rng)
	results := make([]int, value.rows*sign.rows)

	for b.Loop() {
		dotTernaryGo(value.data, sign.data, nonZero.data, value.rows, sign.rows, value.Stride(), results)
	}
}

// 環境: AMD Ryzen 5 8600G
// 形状: 768行 × 768列
// 期待: Go版に対して、おおよそ4〜5倍速い
// ※ 上記のコメントは勝手に変更・削除してはならない
func BenchmarkDotTernaryAVX512(b *testing.B) {
	skipIfNoAVX512(b)
	rng := rand.New(rand.NewPCG(5, 6))
	value := newBenchMatrix(b, benchDotTernaryRows, benchDotTernaryCols, rng)
	sign := newBenchMatrix(b, benchDotTernaryRows, benchDotTernaryCols, rng)
	nonZero := newBenchMatrix(b, benchDotTernaryRows, benchDotTernaryCols, rng)
	results := make([]int, value.rows*sign.rows)

	for b.Loop() {
		dotTernaryAVX512(&value.data[0], &sign.data[0], &nonZero.data[0], value.rows, sign.rows, value.Stride(), &results[0])
	}
}

// 正しいバリデーションが書かれていない場合、このテストは失敗ではなくクラッシュする可能性がある。詳細は「kernels_amd64.go」を参照
func TestValidate(t *testing.T) {
	// 意図的に不整合な内部長を持つ Matrix を組み立てる
	// data は非公開フィールドの為、このテストは package bitsx に置く必要がある。

	newMatrix := func(rows, cols, dataLen int) *Matrix {
		return &Matrix{rows: rows, cols: cols, data: make([]uint64, dataLen)}
	}

	tests := []struct {
		name    string
		rows    int
		cols    int
		dataLen int
		wantErr bool
	}{
		// cols = 100のとき、stride = 2
		// 1行あたり、2ワード(dataLen)が必要
		// よって、rows = 3のとき、dataLen = 6
		{
			name:    "正常",
			rows:    3,
			cols:    100,
			dataLen: 6,
			wantErr: false,
		},
		{
			name:    "異常_1ワード不足",
			rows:    3,
			cols:    100,
			dataLen: 5,
			wantErr: true,
		},
		// 1. 正しいバリデーション: rows * stride == len(data)
		//    (rows = 3, stride = 2 の場合、len(data) = 6 でのみ通る)
		// 2. 商へ式変形した場合: len(data) / stride == rows
		//    (7 / 2 = 3 と余りが切り捨てられる為、len(data) = 7 でも通る)
		// 3. 検証内容: 上記の変形バグを防ぎ、len(data) = 7 (1ワード過剰) を正しく弾けるか確認する。
		{
			name:    "異常_1ワード過剰",
			rows:    3,
			cols:    100,
			dataLen: 7,
			wantErr: true,
		},
		{
			name:    "異常_1行分過剰",
			rows:    3,
			cols:    100,
			dataLen: 8,
			wantErr: true,
		},
		{
			name:    "異常_dataが空",
			rows:    1,
			cols:    64,
			dataLen: 0,
			wantErr: true,
		},
		{
			name:    "異常_Rowsが0",
			rows:    0,
			cols:    100,
			dataLen: 6,
			wantErr: true,
		},
		{
			name:    "異常_Rowsが負",
			rows:    -1,
			cols:    100,
			dataLen: 6,
			wantErr: true,
		},
		{
			name:    "異常_Colsが0",
			rows:    3,
			cols:    0,
			dataLen: 6,
			wantErr: true,
		},
		{
			name:    "異常_Colsが負",
			rows:    3,
			cols:    -1,
			dataLen: 6,
			wantErr: true,
		},
		{
			name:    "異常_Cols+63が桁あふれしてStrideが負になる",
			rows:    1,
			cols:    math.MaxInt,
			dataLen: 8,
			wantErr: true,
		},
		{
			name:    "異常_Cols+63がギリギリ桁あふれしてStrideが負になる",
			rows:    1,
			cols:    math.MaxInt - 62,
			dataLen: 1,
			wantErr: true,
		},
		{
			name:    "異常_Rows*Strideの桁あふれが負になる",
			rows:    math.MaxInt,
			cols:    100,
			dataLen: 6,
			wantErr: true,
		},
		{
			name:    "異常_Rows*Strideの桁あふれが小さい正の値になる",
			rows:    1<<62 + 1,
			cols:    200,
			dataLen: 8,
			wantErr: true,
		},
		{
			name:    "異常_Rows*Strideの桁あふれが周回して元の値に戻る",
			rows:    1 << 62,
			cols:    257,
			dataLen: 6,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := newMatrix(tt.rows, tt.cols, tt.dataLen).validateDotFamily()
			if !tt.wantErr && err != nil {
				t.Fatalf("nilを期待したが、エラーが返された: %v", err)
			}
			if tt.wantErr && err == nil {
				t.Fatal("エラーを期待したが、nilが返された")
			}
		})
	}
}

// Matrix.Dot() と validateDotArgs() の判定が一致することを検証する。
// 結果配列の桁あふれは入力値域外のため、TestValidate で検証する。
func FuzzDotValidationAgreement(f *testing.F) {
	seeds := []struct {
		validRows, validCols                         uint8
		candidateRows, candidateCols, candidateWords int16
	}{
		// 基準は3行×100列（stride=2、必要データ長=6）。
		{3, 100, 3, 100, 5}, // 内部データが1ワード不足
		{3, 100, 3, 100, 7}, // 内部データが1ワード過剰
		{3, 100, 0, 100, 0}, // 行数が不正
		{3, 100, 3, 128, 6}, // 同じstrideでも列数が不一致
		{3, 100, 1, 100, 2}, // 行数だけ異なる正常系
	}
	for _, s := range seeds {
		f.Add(s.validRows, s.validCols, s.candidateRows, s.candidateCols, s.candidateWords)
	}

	f.Fuzz(func(t *testing.T, vRows, vCols uint8, cRows, cCols, cWords int16) {
		// valid側をuint8の範囲に制限し、判定一致時の計算量を抑える。
		if vRows == 0 || vCols == 0 {
			return
		}
		valid, err := NewZerosMatrix(int(vRows), int(vCols))
		if err != nil {
			t.Fatalf("%v", err)
		}

		// fuzz入力による過大な確保を避ける。
		if cWords < 0 || cWords > 1024 {
			return
		}
		// 公開APIでは作れない不整合な内部状態を検証するため、直接組み立てる。
		candidate := &Matrix{rows: int(cRows), cols: int(cCols), data: make([]uint64, cWords)}

		cases := []struct {
			name        string
			left, right *Matrix
		}{
			{"candidateがレシーバー", candidate, valid},
			{"candidateが第1引数", valid, candidate},
		}

		for _, c := range cases {
			_, wantErr := validateDotArgs(c.left, c.right)
			_, gotErr := c.left.Dot(c.right)
			if (wantErr != nil) != (gotErr != nil) {
				t.Errorf("%s: バリデーション判定が不一致: gotErr = %v wantErr = %v", c.name, gotErr, wantErr)
			}
		}
	})
}

// Matrix.DotTernary() と validateDotTernaryArgs() の判定が一致することを検証する。
// 結果配列の桁あふれは入力値域外のため、TestValidate で検証する。
func FuzzDotTernaryValidationAgreement(f *testing.F) {
	seeds := []struct {
		aRows, aCols, bRows, bCols                   uint8
		candidateRows, candidateCols, candidateWords int16
	}{
		// 基準は3つとも3行×100列（stride=2、必要データ長=6）。
		{3, 100, 3, 100, 3, 100, 5}, // 内部データが1ワード不足
		{3, 100, 3, 100, 3, 100, 7}, // 内部データが1ワード過剰
		{3, 100, 3, 100, 0, 100, 0}, // 行数が不正
		{3, 100, 3, 100, 3, 128, 6}, // 同じstrideでも列数が不一致
		{3, 100, 3, 100, 2, 100, 4}, // valueだけなら行数が異なっても正常
		{3, 100, 3, 100, 3, 100, 6}, // 3つのどの位置でも正常
	}
	for _, s := range seeds {
		f.Add(s.aRows, s.aCols, s.bRows, s.bCols, s.candidateRows, s.candidateCols, s.candidateWords)
	}

	f.Fuzz(func(t *testing.T, aRows, aCols, bRows, bCols uint8, cRows, cCols, cWords int16) {
		// valid側をuint8の範囲に制限し、判定一致時の計算量を抑える。
		if aRows == 0 || aCols == 0 || bRows == 0 || bCols == 0 {
			return
		}
		validA, err := NewZerosMatrix(int(aRows), int(aCols))
		if err != nil {
			t.Fatalf("%v", err)
		}
		validB, err := NewZerosMatrix(int(bRows), int(bCols))
		if err != nil {
			t.Fatalf("%v", err)
		}

		// fuzz入力による過大な確保を避ける。
		if cWords < 0 || cWords > 1024 {
			return
		}
		// 公開APIでは作れない不整合な内部状態を検証するため、直接組み立てる。
		candidate := &Matrix{rows: int(cRows), cols: int(cCols), data: make([]uint64, cWords)}

		cases := []struct {
			name                 string
			value, sign, nonZero *Matrix
		}{
			{"candidateがvalue(レシーバー)", candidate, validA, validB},
			{"candidateがsign(第1引数)", validA, candidate, validB},
			{"candidateがnonZero(第2引数)", validA, validB, candidate},
		}

		for _, c := range cases {
			_, wantErr := validateDotTernaryArgs(c.value, c.sign, c.nonZero)
			_, gotErr := c.value.DotTernary(c.sign, c.nonZero)
			if (wantErr != nil) != (gotErr != nil) {
				t.Errorf("%s: バリデーション判定が不一致: gotErr = %v wantErr = %v", c.name, gotErr, wantErr)
			}
		}
	})
}
