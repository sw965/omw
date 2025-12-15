package parallel_test

import (
	"errors"
	"fmt"
	"slices"
	"testing"

	"github.com/sw965/omw/parallel"
)

func TestFor(t *testing.T) {
	tests := []struct {
		name    string
		n       int
		p       int
		want    []string
		wantErr bool
		wantErrIs error
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
		// 異常系
		{
			name:"異常_nが負の値",
			n:-1,
			p:4,
			wantErr:true,
			wantErrIs:parallel.ErrNegativeN,
		},
		{
			name:"異常_pが0以下",
			n:16,
			p:0,
			wantErr:true,
			wantErrIs:parallel.ErrInvalidP,
		},
		//準正常
		{
			name:"準正常_nが0",
			n:0,
			p:4,
			want:[]string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
        	var got []string
        	if tc.n >= 0 {
            	got = make([]string, tc.n)
        	}

            gotErr := parallel.For(tc.n, tc.p, func(workerId, idx int) error {
                got[idx] = fmt.Sprintf("w%d: i%d", workerId, idx)
                return nil
            })

			if tc.wantErr {
				if gotErr == nil {
					t.Fatal("エラーを期待したが、nilが返された")
				}

				if !errors.Is(gotErr, tc.wantErrIs) {
					t.Errorf("期待されたエラー型と異なります: want: %T, got :%T", tc.wantErrIs, gotErr)
				}
				return
			}

			if gotErr != nil {
				t.Fatalf("予期せぬエラー: %v", gotErr)
			}

			if !slices.Equal(got, tc.want) {
				t.Errorf("want: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestFor_CallbackError(t *testing.T) {
	t.Run("異常_1つ", func(t *testing.T) {
		// worker0 に割り当てられるインデックス 0, 1, 2, 3, 4
		// worker1 に割り当てられるインデックス 5, 6, 7, 8, 9
		const n = 10
		const p = 2

		// worker0に割り割り当てられるインデックスが2の時にエラーが起きる想定
		errIdx := 2
		failErr := errors.New("boom")
		succeeded := make([]bool, n)

		gotErr := parallel.For(n, p, func(workerId, idx int) error {
			if idx == errIdx {
				return failErr
			}
			succeeded[idx] = true
			return nil
		})

		if gotErr == nil {
			t.Fatal("エラーを期待したが、nilが返された")
		}
		if !errors.Is(gotErr, failErr) {
			t.Fatalf("errors.Is が failErr を拾えない: gotErr: %v", gotErr)
		}

		wantSucceeded := []bool{
			// worker0は、インデックスが2番目の時に、エラーが起きるので、それ以降はfalse
			true, true, false, false, false,
			// worker1は、エラーが起きない為、担当するインデックスは全てtrue
			true, true, true, true, true,
		}

		if !slices.Equal(succeeded, wantSucceeded) {
			t.Errorf("wantSucceeded: %v, succeeded: %v", wantSucceeded, succeeded)
		}
	})

	t.Run("異常_2つ", func(t *testing.T) {
		//worker0に割り当てられるインデックス 0, 1, 2, 3
		//worker1に割り当てられるインデックス 4, 5, 6, 7
		//worker2に割り当てられるインデックス 8, 9, 10, 11
		const n = 12
		const p = 3

		//worker0とworker2がエラーを起こす想定
		err0 := errors.New("boom0")
		err2 := errors.New("boom2")

		succeeded := make([]bool, n)

		gotErr := parallel.For(n, p, func(workerId, idx int) error {
			switch idx {
			case 1:
				// worker0がインデックス1でエラーを返す
				return err0
			case 9:
				// worker2がインデックス9でエラーを返す
				return err2
			default:
				succeeded[idx] = true
				return nil
			}
		})

		if gotErr == nil {
			t.Fatal("エラーを期待したが、nilが返された")
		}

		if !errors.Is(gotErr, err0) || !errors.Is(gotErr, err2) {
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

		if !slices.Equal(succeeded, wantSucceeded) {
			t.Errorf("wantSucceeded: %v, succeeded: %v", wantSucceeded, succeeded)
		}
	})
}