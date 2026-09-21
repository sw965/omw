package parallel_test

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/sw965/omw/parallel"
)

func TestFor_TableDriven(t *testing.T) {
	tests := []struct {
		name      string
		n         int
		p         int
		isNilFunc bool
		want      []string
		wantErr   bool
	}{
		// 正常系
		{
			name: "正常_余りあり",
			n:    11,
			p:    3,
			// n = 11, p = 3の時、余り(r) = 2。余った量は、worker0とworker1にそれぞれに1つずつ割り当てられる
			want: []string{
				// worker0 (4 items)
				"w0: i0", "w0: i1", "w0: i2", "w0: i3",
				// worker1 (4 items)
				"w1: i4", "w1: i5", "w1: i6", "w1: i7",
				// worker2 (3 items)
				"w2: i8", "w2: i9", "w2: i10",
			},
		},
		{
			name: "正常_余りなし",
			n:    6,
			p:    2,
			want: []string{
				"w0: i0", "w0: i1", "w0: i2",
				"w1: i3", "w1: i4", "w1: i5",
			},
		},
		{
			name: "正常_pがnより大きい",
			n:    3,
			p:    5,
			// p > n なので p = n = 3 に正則化される
			want: []string{
				"w0: i0",
				"w1: i1",
				"w2: i2",
			},
		},
		{
			name: "正常_nが0",
			n:    0,
			p:    4,
			want: []string{},
		},

		// 異常系
		{
			name:    "異常_nが負の値",
			n:       -1,
			p:       4,
			wantErr: true,
		},
		{
			name:    "異常_pが0以下",
			n:       16,
			p:       0,
			wantErr: true,
		},
		{
			name:      "異常_fがnil",
			n:         10,
			p:         2,
			isNilFunc: true,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()

			var got []string
			// 異常系テストの時、tt.n < 0 を渡す事があるため、ここでガードしておく。
			if tt.n >= 0 {
				got = make([]string, tt.n)
			}

			var f func(workerIdx, itemIdx int) error
			if !tt.isNilFunc {
				f = func(workerIdx, itemIdx int) error {
					got[itemIdx] = fmt.Sprintf("w%d: i%d", workerIdx, itemIdx)
					return nil
				}
			}

			err := parallel.For(tt.n, tt.p, f)

			if tt.wantErr {
				if err == nil {
					t.Fatal("エラーを期待したが、nilが返された")
				}
				return
			}

			if err != nil {
				t.Fatalf("nilを期待したが、エラーが返された: %v", err)
			}

			if !slices.Equal(got, tt.want) {
				t.Errorf("値の不一致: got = %v want = %v", got, tt.want)
			}
		})
	}
}

func TestFor_CallbackError(t *testing.T) {
	t.Run("異常_1つ", func(t *testing.T) {
		t.Helper()
		// worker0 に割り当てられるインデックス 0, 1, 2, 3, 4
		// worker1 に割り当てられるインデックス 5, 6, 7, 8, 9
		const n = 10
		const p = 2

		// worker0に割り割り当てられるインデックスが2の時にエラーが起きる想定
		errIdx := 2
		failErr := errors.New("boom")
		gotSucceeded := make([]bool, n)

		gotErr := parallel.For(n, p, func(workerIdx, itemIdx int) error {
			if itemIdx == errIdx {
				return failErr
			}
			gotSucceeded[itemIdx] = true
			return nil
		})

		if gotErr == nil {
			t.Fatal("エラーを期待したが、nilが返された")
		}

		if !errors.Is(gotErr, failErr) {
			t.Fatalf("failErr を拾えない: err: %v", gotErr)
		}

		wantSucceeded := []bool{
			// worker0は、インデックスが2番目の時に、エラーが起きるので、それ以降はfalse
			true, true, false, false, false,
			// worker1は、エラーが起きない為、担当するインデックスは全てtrue
			true, true, true, true, true,
		}

		if !slices.Equal(gotSucceeded, wantSucceeded) {
			t.Errorf("gotSucceeded: %v, wantSucceeded: %v", gotSucceeded, wantSucceeded)
		}
	})

	t.Run("異常_2つ", func(t *testing.T) {
		// worker0に割り当てられるインデックス 0, 1, 2, 3
		// worker1に割り当てられるインデックス 4, 5, 6, 7
		// worker2に割り当てられるインデックス 8, 9, 10, 11
		const n = 12
		const p = 3

		// worker0とworker2がエラーを起こす想定
		worker0Err := errors.New("boom0")
		worker2Err := errors.New("boom2")

		gotSucceeded := make([]bool, n)
		gotErr := parallel.For(n, p, func(workerIdx, itemIdx int) error {
			switch itemIdx {
			case 1:
				// worker0がインデックス1でエラーを返す
				return worker0Err
			case 9:
				// worker2がインデックス9でエラーを返す
				return worker2Err
			default:
				gotSucceeded[itemIdx] = true
				return nil
			}
		})

		if gotErr == nil {
			t.Fatal("エラーを期待したが、nilが返された")
		}

		if !errors.Is(gotErr, worker0Err) || !errors.Is(gotErr, worker2Err) {
			t.Fatalf("両方のエラーを拾えない: gotErr: %v", gotErr)
		}

		wantSucceeded := []bool{
			// worker0は、インデックスが1の時にエラーが起きるので、それ以降はfalse
			true, false, false, false,
			// worker1は、エラーが起きない為、完走
			true, true, true, true,
			// worker2は、インデックスが9の時にエラーが起きるので、それ以降はfalse
			true, false, false, false,
		}

		if !slices.Equal(gotSucceeded, wantSucceeded) {
			t.Errorf("gotSucceeded: %v, wantSucceeded: %v", gotSucceeded, wantSucceeded)
		}
	})
}

func TestForContext_Cancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := parallel.ForContext(ctx, 100, 4, func(w, i int) error {
		return nil
	})

	if !errors.Is(err, context.Canceled) {
		t.Errorf("context.Canceled を期待したが、異なるエラーが返された: %v", err)
	}
}

func TestForContext_Timeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	err := parallel.ForContext(ctx, 10000, 4, func(w, i int) error {
		time.Sleep(5 * time.Millisecond)
		return nil
	})

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("context.DeadlineExceeded を期待したが、異なるエラーが返された: %v", err)
	}
}
