package validator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPointerAttr(t *testing.T) {
	type Inner struct {
		Val int `validate:"range(2,5)"`
	}
	type Outer struct {
		Inner *Inner
	}

	tests := []struct {
		input   Outer
		wantErr error
	}{
		{
			input: Outer{
				Inner: &Inner{
					Val: 4,
				},
			},
		},
		{
			input: Outer{
				Inner: &Inner{
					Val: 1,
				},
			},
			wantErr: fmt.Errorf(`Validation failed for field "Val": should be in the range [2, 5]`),
		},
		{
			input: Outer{
				Inner: nil,
			},
		},
	}

	for _, tt := range tests {
		err := Validate(tt.input)
		if tt.wantErr != nil {
			assert.Error(t, err)
			assert.Equal(t, tt.wantErr.Error(), err.Error())
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestValidate(t *testing.T) {
	//Register("range", func(v int, low int, high int) (bool, string) {
	//	return v >= low && v <= high, fmt.Sprintf("value should be in range between %d and %d", low, high)
	//})

	type TestStruct struct {
		Attr int `validate:"range(-10, 10)"`
	}

	tests := []struct {
		name    string
		val     int
		wantErr error
	}{
		{
			name:    "out of range",
			val:     42,
			wantErr: fmt.Errorf("Validation failed for field %q: should be in the range [%d, %d]", "Attr", -10, 10),
		},
		{
			name:    "out of range",
			val:     -42,
			wantErr: fmt.Errorf("Validation failed for field %q: should be in the range [%d, %d]", "Attr", -10, 10),
		},
		{
			name: "in the range",
			val:  0,
		},
		{
			name: "in the range",
			val:  -10,
		},
		{
			name: "in the range",
			val:  10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ts TestStruct
			ts.Attr = tt.val
			err := Validate(ts)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
