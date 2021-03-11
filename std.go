package validator

import (
	"fmt"
	"reflect"
)

func StdNone() (bool, string, bool) {
	return true, "", Break
}

func StdOptional(v interface{}) (bool, string, bool) {
	rv := reflect.ValueOf(v)
	zv := reflect.Zero(rv.Type())
	if reflect.DeepEqual(rv.Interface(), zv.Interface()) {
		// it's a zero value, break the validation chain
		return true, "", Break
	}
	return true, "", Continue
}

func StdEmpty(v interface{}) (bool, string) {
	rv := reflect.ValueOf(v)
	zv := reflect.Zero(rv.Type())
	return reflect.DeepEqual(rv.Interface(), zv.Interface()), "should be empty"
}

func StdNonEmpty(v interface{}) (bool, string) {
	rv := reflect.ValueOf(v)
	zv := reflect.Zero(rv.Type())
	return !reflect.DeepEqual(rv.Interface(), zv.Interface()), "should not be empty"
}

func StdEq(v interface{}, cmp string) (bool, string) {
	eq, err := compare(v, cmp)
	if err != nil {
		return false, err.Error()
	}
	return eq == CompareEqual, fmt.Sprintf("should be equal to %s", cmp)
}

func StdNe(v interface{}, cmp string) (bool, string) {
	eq, err := compare(v, cmp)
	if err != nil {
		return false, err.Error()
	}
	return eq != CompareEqual, fmt.Sprintf("should not be equal to %s", cmp)
}

func StdGt(v interface{}, cmp string) (bool, string) {
	eq, err := compare(v, cmp)
	if err != nil {
		return false, err.Error()
	}
	return eq == CompareGreaterThan, fmt.Sprintf("should be greater than %s", cmp)
}

func StdGte(v interface{}, cmp string) (bool, string) {
	eq, err := compare(v, cmp)
	if err != nil {
		return false, err.Error()
	}
	return (CompareEqual|CompareGreaterThan)&eq > 0, fmt.Sprintf("should be greater or equal to %s", cmp)
}

func StdLt(v interface{}, cmp string) (bool, string) {
	eq, err := compare(v, cmp)
	if err != nil {
		return false, err.Error()
	}
	return eq == CompareLessThan, fmt.Sprintf("should be less than %s", cmp)
}

func StdLte(v interface{}, cmp string) (bool, string) {
	eq, err := compare(v, cmp)
	if err != nil {
		return false, err.Error()
	}
	return (CompareEqual|CompareLessThan)&eq > 0, fmt.Sprintf("should be less or equal to %s", cmp)
}

func StdRange(v interface{}, low, high string) (bool, string) {
	eq, err := compare(v, low)
	if err != nil {
		return false, err.Error()
	}
	if (CompareEqual|CompareGreaterThan)&eq > 0 {
		eq, err = compare(v, high)
		if err != nil {
			return false, err.Error()
		}
		if (CompareEqual|CompareLessThan)&eq > 0 {
			return true, ""
		}
	}
	return false, fmt.Sprintf("should be in the range [%s, %s]", low, high)
}

func StdEnum(v interface{}, opts ...string) (bool, string) {
	for _, opt := range opts {
		eq, err := compare(v, opt)
		if err != nil {
			return false, err.Error()
		}
		if CompareEqual&eq > 0 {
			return true, ""
		}
	}
	return false, fmt.Sprintf("should be in range %+v", opts)
}

func StdLen(v interface{}, maxlen int) (bool, string) {
	err := fmt.Sprintf("length must be exactly %d", maxlen)
	if s, ok := v.(string); ok {
		return len(s) == maxlen, err
	} else if s, ok := v.(stringer); ok {
		return len(s.String()) == maxlen, err
	}
	return false, fmt.Sprintf("unexpected string type: %t", v)
}

func StdMaxLen(v interface{}, maxlen int) (bool, string) {
	err := fmt.Sprintf("length must be up to %d", maxlen)
	if s, ok := v.(string); ok {
		return len(s) <= maxlen, err
	} else if s, ok := v.(stringer); ok {
		return len(s.String()) <= maxlen, err
	}
	return false, fmt.Sprintf("unexpected string type: %t", v)
}
