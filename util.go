package validator

import (
	"fmt"
	"reflect"
	"strconv"
)

type stringer interface {
	String() string
}

func copyTypes(types []reflect.Type) []reflect.Type {
	typescp := make([]reflect.Type, len(types))
	copy(typescp, types)
	return typescp
}

func convArgV(types []reflect.Type, args []interface{}, isVariadic bool) ([]reflect.Value, error) {
	types = copyTypes(types)
	if isVariadic {
		var vt reflect.Type
		types, vt = types[:len(types)-1], types[len(types)-1]
		et := vt.Elem()
		for len(types) < len(args) {
			types = append(types, et)
		}
	}
	if len(types) != len(args) {
		return nil, fmt.Errorf("number of factual parameters does not match with validator definition (want: %d, got: %d)", len(types), len(args))
	}
	argV := make([]reflect.Value, 0, len(args))
	for i, arg := range args {
		kind := types[i].Kind()
		var val reflect.Value
		var err error
		switch reflect.ValueOf(arg).Kind() {
		case reflect.String:
			var s string
			if strngr, ok := arg.(stringer); ok {
				s = strngr.String()
			} else {
				s = arg.(string)
			}
			val, err = convStringVal(s, kind)
		default:
			val, err = convDirectCast(arg, kind)
		}
		if err != nil {
			return nil, err
		}
		argV = append(argV, val)
	}

	return argV, nil
}

func convDirectCast(arg interface{}, kind reflect.Kind) (reflect.Value, error) {
	var val reflect.Value
	switch kind {
	case reflect.Interface:
		val = reflect.ValueOf(arg)
	case reflect.String:
		val = reflect.ValueOf(arg.(string))
	case reflect.Int:
		val = reflect.ValueOf(arg.(int))
	case reflect.Int8:
		val = reflect.ValueOf(arg.(int8))
	case reflect.Int16:
		val = reflect.ValueOf(arg.(int16))
	case reflect.Int32:
		val = reflect.ValueOf(arg.(int32))
	case reflect.Int64:
		val = reflect.ValueOf(arg.(int64))
	case reflect.Uint:
		val = reflect.ValueOf(arg.(uint))
	case reflect.Uint8:
		val = reflect.ValueOf(arg.(uint8))
	case reflect.Uint16:
		val = reflect.ValueOf(arg.(uint16))
	case reflect.Uint32:
		val = reflect.ValueOf(arg.(uint32))
	case reflect.Uint64:
		val = reflect.ValueOf(arg.(uint64))
	case reflect.Uintptr:
		val = reflect.ValueOf(arg.(uintptr))
	case reflect.Float32:
		val = reflect.ValueOf(arg.(float32))
	case reflect.Float64:
		val = reflect.ValueOf(arg.(float64))
	case reflect.Bool:
		val = reflect.ValueOf(arg.(bool))
	default:
		return val, fmt.Errorf("unsupported kind: %v", kind)
	}
	return val, nil
}

func convStringVal(arg string, kind reflect.Kind) (reflect.Value, error) {
	var val reflect.Value
	switch kind {
	case reflect.Interface, reflect.String:
		val = reflect.ValueOf(arg)
	case reflect.Int:
		if i, err := strconv.ParseInt(arg, 10, 64); err != nil {
			return val, err
		} else {
			val = reflect.ValueOf(int(i))
		}
	case reflect.Int8:
		if i, err := strconv.ParseInt(arg, 10, 8); err != nil {
			return val, err
		} else {
			val = reflect.ValueOf(int8(i))
		}
	case reflect.Int16:
		if i, err := strconv.ParseInt(arg, 10, 16); err != nil {
			return val, err
		} else {
			val = reflect.ValueOf(int16(i))
		}
	case reflect.Int32:
		if i, err := strconv.ParseInt(arg, 10, 32); err != nil {
			return val, err
		} else {
			val = reflect.ValueOf(int32(i))
		}
	case reflect.Int64:
		if i, err := strconv.ParseInt(arg, 10, 64); err != nil {
			return val, err
		} else {
			val = reflect.ValueOf(int64(i))
		}
	case reflect.Uint:
		if u, err := strconv.ParseUint(arg, 10, 64); err != nil {
			return val, err
		} else {
			val = reflect.ValueOf(uint(u))
		}
	case reflect.Uint8:
		if u, err := strconv.ParseUint(arg, 10, 8); err != nil {
			return val, err
		} else {
			val = reflect.ValueOf(uint8(u))
		}
	case reflect.Uint16:
		if u, err := strconv.ParseUint(arg, 10, 16); err != nil {
			return val, err
		} else {
			val = reflect.ValueOf(uint16(u))
		}
	case reflect.Uint32:
		if u, err := strconv.ParseUint(arg, 10, 32); err != nil {
			return val, err
		} else {
			val = reflect.ValueOf(uint32(u))
		}
	case reflect.Uint64:
		if u, err := strconv.ParseUint(arg, 10, 64); err != nil {
			return val, err
		} else {
			val = reflect.ValueOf(uint64(u))
		}
	case reflect.Uintptr:
		if u, err := strconv.ParseUint(arg, 10, 64); err != nil {
			return val, err
		} else {
			val = reflect.ValueOf(uintptr(u))
		}
	case reflect.Float32:
		if f, err := strconv.ParseFloat(arg, 32); err != nil {
			return val, err
		} else {
			val = reflect.ValueOf(f)
		}
	case reflect.Float64:
		if f, err := strconv.ParseFloat(arg, 64); err != nil {
			return val, err
		} else {
			val = reflect.ValueOf(f)
		}
	case reflect.Bool:
		if b, err := strconv.ParseBool(arg); err != nil {
			return val, err
		} else {
			val = reflect.ValueOf(b)
		}
	default:
		return val, fmt.Errorf("unsupported kind: %v", kind)
	}
	return val, nil
}

func compare(v interface{}, cmp string) (Equality, error) {
	rv := reflect.ValueOf(v)
	cmpv, err := convStringVal(cmp, rv.Kind())
	var eq Equality
	if err != nil {
		return 0, err
	}
	switch rv.Kind() {
	case reflect.Bool:
		bv := v.(bool)
		if rv.Interface() == cmpv.Interface() {
			eq = CompareEqual
		} else {
			if bv == true { // true is greater than false
				eq = CompareGreaterThan
			} else {
				eq = CompareLessThan
			}
		}
	case reflect.Int:
		iv := v.(int)
		if rv.Interface() == cmpv.Interface() {
			eq = CompareEqual
		} else {
			if iv < int(cmpv.Int()) {
				eq = CompareLessThan
			} else {
				eq = CompareGreaterThan
			}
		}
	case reflect.Int8:
		iv := v.(int8)
		if rv.Interface() == cmpv.Interface() {
			eq = CompareEqual
		} else {
			if iv < int8(cmpv.Int()) {
				eq = CompareLessThan
			} else {
				eq = CompareGreaterThan
			}
		}
	case reflect.Int16:
		iv := v.(int16)
		if rv.Interface() == cmpv.Interface() {
			eq = CompareEqual
		} else {
			if iv < int16(cmpv.Int()) {
				eq = CompareLessThan
			} else {
				eq = CompareGreaterThan
			}
		}
	case reflect.Int32:
		iv := v.(int32)
		if rv.Interface() == cmpv.Interface() {
			eq = CompareEqual
		} else {
			if iv < int32(cmpv.Int()) {
				eq = CompareLessThan
			} else {
				eq = CompareGreaterThan
			}
		}
	case reflect.Int64:
		iv := v.(int64)
		if rv.Interface() == cmpv.Interface() {
			eq = CompareEqual
		} else {
			if iv < int64(cmpv.Int()) {
				eq = CompareLessThan
			} else {
				eq = CompareGreaterThan
			}
		}
	case reflect.Uint:
		iv := v.(uint)
		if rv.Interface() == cmpv.Interface() {
			eq = CompareEqual
		} else {
			if iv < uint(cmpv.Uint()) {
				eq = CompareLessThan
			} else {
				eq = CompareGreaterThan
			}
		}
	case reflect.Uint8:
		iv := v.(uint8)
		if rv.Interface() == cmpv.Interface() {
			eq = CompareEqual
		} else {
			if iv < uint8(cmpv.Uint()) {
				eq = CompareLessThan
			} else {
				eq = CompareGreaterThan
			}
		}
	case reflect.Uint16:
		iv := v.(uint16)
		if rv.Interface() == cmpv.Interface() {
			eq = CompareEqual
		} else {
			if iv < uint16(cmpv.Uint()) {
				eq = CompareLessThan
			} else {
				eq = CompareGreaterThan
			}
		}
	case reflect.Uint32:
		iv := v.(uint32)
		if rv.Interface() == cmpv.Interface() {
			eq = CompareEqual
		} else {
			if iv < uint32(cmpv.Uint()) {
				eq = CompareLessThan
			} else {
				eq = CompareGreaterThan
			}
		}
	case reflect.Uint64:
		iv := v.(uint64)
		if rv.Interface() == cmpv.Interface() {
			eq = CompareEqual
		} else {
			if iv < uint64(cmpv.Uint()) {
				eq = CompareLessThan
			} else {
				eq = CompareGreaterThan
			}
		}
	case reflect.Uintptr:
		iv := v.(uintptr)
		if rv.Interface() == cmpv.Interface() {
			eq = CompareEqual
		} else {
			if iv < uintptr(cmpv.Uint()) {
				eq = CompareLessThan
			} else {
				eq = CompareGreaterThan
			}
		}
	case reflect.String:
		sv := v.(string)
		if rv.Interface() == cmpv.Interface() {
			eq = CompareEqual
		} else {
			if sv < cmpv.String() {
				eq = CompareLessThan
			} else {
				eq = CompareGreaterThan
			}
		}
	default:
		return 0, fmt.Errorf("kind %v is not comparable", rv.Kind())
	}

	return eq, nil
}

func duplicateValidatorDefErr(handle string) error {
	return fmt.Errorf("Duplicate validator definition: %s", handle)
}
