package validator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStdOptional_String(t *testing.T) {
	Register("post_optional_str", func(v string, cmp string) (bool, string) {
		return v == cmp, "should equal to " + cmp
	})

	type TestStruct struct {
		Attr string `validate:"optional,post_optional_str(foobar)"`
	}

	var err error

	var ts TestStruct
	err = Validate(ts)
	assert.NoError(t, err)

	ts.Attr = "foobar"
	err = Validate(ts)
	assert.NoError(t, err)

	ts.Attr = "barbaz"
	err = Validate(ts)
	assert.Error(t, err)
	assert.Equal(t, "Validation failed for field \"Attr\": should equal to foobar", err.Error())
}

func TestStdOptional_Uint32(t *testing.T) {
	Register("post_optional_uint32", func(v uint32, cmp uint32) (bool, string) {
		return v == cmp, fmt.Sprintf("should equal to %d", cmp)
	})

	type TestStruct struct {
		Attr uint32 `validate:"optional,post_optional_uint32(42)"`
	}

	var err error

	var ts TestStruct
	err = Validate(ts)
	assert.NoError(t, err)

	ts.Attr = 42
	err = Validate(ts)
	assert.NoError(t, err)

	ts.Attr = 123
	err = Validate(ts)
	assert.Error(t, err)
	assert.Equal(t, "Validation failed for field \"Attr\": should equal to 42", err.Error())
}

func TestStdOptional_Pointer(t *testing.T) {
	type TestSubStruct struct {
		Attr string `validate:"nonempty"`
	}
	type TestStruct struct {
		SubStruct *TestSubStruct `validate:"optional"`
	}
	var err error

	var ts TestStruct
	err = Validate(ts)
	assert.NoError(t, err)

	ts.SubStruct = &TestSubStruct{}
	err = Validate(ts)
	assert.Error(t, err)
	assert.Equal(t, "Validation failed for field \"Attr\": should not be empty", err.Error())

	ts.SubStruct.Attr = "foobar"
	err = Validate(ts)
	assert.NoError(t, err)
}

func TestStdEmpty(t *testing.T) {
	type TestStruct struct {
		Attr int `validate:"empty"`
	}

	var err error

	var ts TestStruct
	err = Validate(ts)
	assert.NoError(t, err)

	ts.Attr = 42
	err = Validate(ts)
	assert.Error(t, err)
	assert.Equal(t, "Validation failed for field \"Attr\": should be empty", err.Error())
}

func TestStdEnum(t *testing.T) {
	type TestStruct struct {
		Attr string `validate:"enum(foo, bar, baz, boo)"`
	}

	var err error

	legal := []string{"foo", "bar", "baz", "boo"}
	illegal := []string{"bak", "bazz", "123"}

	for _, v := range legal {
		var ts TestStruct
		ts.Attr = v
		err = Validate(ts)
		assert.NoError(t, err)
	}

	for _, v := range illegal {
		var ts TestStruct
		ts.Attr = v
		err = Validate(ts)
		assert.Error(t, err)
		assert.Equal(t, "Validation failed for field \"Attr\": should be in range [foo bar baz boo]", err.Error())
	}
}

func TestStdLen(t *testing.T) {
	type TestStruct struct {
		Str5 string `validate:"len(5)"`
	}

	var err error
	valid := []string{"hello", "world"}
	invalid := []string{"", "hello!", "worl"}

	for _, v := range valid {
		var ts TestStruct
		ts.Str5 = v
		err = Validate(ts)
		assert.NoError(t, err)
	}

	for _, v := range invalid {
		var ts TestStruct
		ts.Str5 = v
		err = Validate(ts)
		assert.Error(t, err)
		assert.Equal(t, "Validation failed for field \"Str5\": length must be exactly 5", err.Error())
	}
}

func TestStdMaxLen(t *testing.T) {
	type TestStruct struct {
		Str5 string `validate:"maxlen(5)"`
	}

	var err error
	valid := []string{"hello", "", "hey"}
	invalid := []string{"hello!", "world123"}

	for _, v := range valid {
		var ts TestStruct
		ts.Str5 = v
		err = Validate(ts)
		assert.NoError(t, err)
	}

	for _, v := range invalid {
		var ts TestStruct
		ts.Str5 = v
		err = Validate(ts)
		assert.Error(t, err)
		assert.Equal(t, "Validation failed for field \"Str5\": length must be up to 5", err.Error())
	}
}
