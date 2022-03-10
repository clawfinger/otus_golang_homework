package hw09structvalidator

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	IntRange struct {
		Code int `validate:"in:200,404,500"`
	}
	IntMin struct {
		Code int `validate:"min:200"`
	}
	IntMax struct {
		Code int `validate:"max:200"`
	}
	StringLen struct {
		Version string `validate:"len:5"`
	}
	StringIn struct {
		Version string `validate:"in:abc,def,hig"`
	}
	StringRegexp struct {
		Version string `validate:"regexp:\\d+a"`
	}
	StringSlice struct {
		Version []string `validate:"len:5"`
	}
	IntSlice struct {
		Code []int `validate:"in:200,250,500"`
	}
	IntMulti struct {
		Code int `validate:"min:0|max:10"`
	}
	StringMulti struct {
		Version string `validate:"len:5|regexp:\\d+a"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in:          1,
			expectedErr: errNotAStruct,
		},
		{
			in:          IntRange{Code: 1},
			expectedErr: errIntNotInSet,
		},
		{
			in:          IntRange{Code: 404},
			expectedErr: nil,
		},
		{
			in:          IntMin{Code: 404},
			expectedErr: nil,
		},
		{
			in:          IntMin{Code: 104},
			expectedErr: errIntLess,
		},
		{
			in:          IntMax{Code: 104},
			expectedErr: nil,
		},
		{
			in:          IntMax{Code: 404},
			expectedErr: errIntGreater,
		},
		{
			in:          StringLen{Version: "123"},
			expectedErr: errStringLen,
		},
		{
			in:          StringLen{Version: "12345"},
			expectedErr: nil,
		},
		{
			in:          StringIn{Version: "12345"},
			expectedErr: errStringNotInSet,
		},
		{
			in:          StringIn{Version: "def"},
			expectedErr: nil,
		},
		{
			in:          StringRegexp{Version: "456a"},
			expectedErr: nil,
		},
		{
			in:          StringRegexp{Version: "456"},
			expectedErr: errStringRegexp,
		},
		{
			in:          IntSlice{Code: []int{200, 250, 3}},
			expectedErr: errIntNotInSet,
		},
		{
			in:          IntSlice{Code: []int{200, 250, 200}},
			expectedErr: nil,
		},
		{
			in:          StringSlice{Version: []string{"12345", "12", "asdfg"}},
			expectedErr: errStringLen,
		},
		{
			in:          StringSlice{Version: []string{"12345", "12543", "asdfg"}},
			expectedErr: nil,
		},
		{
			in:          IntMulti{Code: 5},
			expectedErr: nil,
		},
		{
			in:          IntMulti{Code: 11},
			expectedErr: errIntGreater,
		},
		{
			in:          IntMulti{Code: -6},
			expectedErr: errIntLess,
		},
		{
			in:          StringMulti{Version: "4564a"},
			expectedErr: nil,
		},
		{
			in:          StringMulti{Version: "454a"},
			expectedErr: errStringLen,
		},
		{
			in:          StringMulti{Version: "454wa"},
			expectedErr: errStringRegexp,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			require.Truef(t, errors.Is(err, tt.expectedErr), "Wrong error.")
		})
	}
}
