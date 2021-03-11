package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseValidateTags(t *testing.T) {
	tests := []struct {
		input string
		want  []ValidateTag
	}{
		{
			input: "",
			want:  []ValidateTag{},
		},
		{
			input: "optional",
			want: []ValidateTag{
				{
					Op:   "optional",
					Args: []interface{}{},
				},
			},
		},
		{
			input: "gt(-10),lt(10)",
			want: []ValidateTag{
				{
					Op:   "gt",
					Args: []interface{}{"-10"},
				},
				{
					Op:   "lt",
					Args: []interface{}{"10"},
				},
			},
		},
		{
			input: "oneOf(foo, bar, baz)",
			want: []ValidateTag{
				{
					Op:   "oneOf",
					Args: []interface{}{"foo", "bar", "baz"},
				},
			},
		},
		{
			input: "foo1(+1, -2, 6.02e+23),foo_bar(this, that)",
			want: []ValidateTag{
				{
					Op:   "foo1",
					Args: []interface{}{"+1", "-2", "6.02e+23"},
				},
				{
					Op:   "foo_bar",
					Args: []interface{}{"this", "that"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			tags := parseValidateTags(tt.input)
			assert.Equal(t, tt.want, tags)
		})
	}
}
