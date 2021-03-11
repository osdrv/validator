package validator

import (
	"errors"
	"fmt"
	"reflect"
)

type Equality uint8

const (
	Continue = true
	Break    = false
)

const (
	CompareEqual Equality = 1 << iota
	CompareLessThan
	CompareGreaterThan
)

const (
	ValidateTagName = "validate"
)

var validators map[string]func(interface{}, ...interface{}) (bool, error)

func init() {
	validators = make(map[string]func(interface{}, ...interface{}) (bool, error))

	Register("empty", StdEmpty)
	Register("enum", StdEnum)
	Register("eq", StdEq)
	Register("gt", StdGt)
	Register("gte", StdGte)
	Register("len", StdLen)
	Register("lt", StdLt)
	Register("lte", StdLte)
	Register("maxlen", StdMaxLen)
	Register("ne", StdNe)
	Register("none", StdNone)
	Register("nonempty", StdNonEmpty)
	Register("optional", StdOptional)
	Register("range", StdRange)
}

func Register(handle string, check interface{}) error {
	if _, ok := validators[handle]; ok {
		return duplicateValidatorDefErr(handle)
	}
	checkV := reflect.ValueOf(check)
	checkT := reflect.TypeOf(check)
	isVariadic := checkT.IsVariadic()
	types := make([]reflect.Type, 0, checkT.NumIn())
	for i := 0; i < checkT.NumIn(); i++ {
		types = append(types, checkT.In(i))
	}

	validators[handle] = func(v interface{}, args ...interface{}) (bool, error) {
		if len(types) > 0 {
			args = append([]interface{}{v}, args...)
		}
		argV, err := convArgV(types, args, isVariadic)
		if err != nil {
			return Break, fmt.Errorf("argument conversion failed: %s", err)
		}
		resV := checkV.Call(argV)

		match := resV[0].Bool()
		reason := "constraint mismatch"
		cont := Continue
		if len(resV) > 1 {
			reason = resV[1].String()
		}
		if len(resV) > 2 {
			cont = resV[2].Bool()
		}
		if !match {
			if len(resV) > 1 {
				reason = resV[1].String()
			}
			return Break, errors.New(reason)
		}
		return cont, nil
	}

	return nil
}

func Validate(datum interface{}) error {
	datumT := reflect.TypeOf(datum)
	datumV := reflect.ValueOf(datum)

	for datumT.Kind() == reflect.Ptr {
		datumT = datumT.Elem()
		datumV = datumV.Elem()
	}

	if datumT.Kind() != reflect.Struct {
		return fmt.Errorf("Validate accepts a struct, %#v %T given", datum, datum)
	}

Datum:
	for i := 0; i < datumT.NumField(); i++ {
		field := datumT.Field(i)
		v := datumV.FieldByIndex([]int{i})
		if tagDef, ok := field.Tag.Lookup(ValidateTagName); ok {
			tags := parseValidateTags(tagDef)
			for _, tag := range tags {
				check, ok := validators[tag.Op]
				if !ok {
					return fmt.Errorf("Validator %q is unknown", tag.Op)
				}
				cont, err := check(v.Interface(), tag.Args...)
				if err != nil {
					return fmt.Errorf("Validation failed for field %q: %s", field.Name, err)
				}
				if cont == Break {
					return nil
				}
			}
		}

		p := v

	Deref:
		switch p.Kind() {
		case reflect.Struct:
			if err := Validate(p.Interface()); err != nil {
				return err
			}
		case reflect.Ptr:
			if p.IsNil() {
				continue Datum
			}
			if p.Kind() == reflect.Ptr {
				p = p.Elem()
				goto Deref
			}
			if err := Validate(p.Interface()); err != nil {
				return err
			}
		}
	}

	return nil
}
