package hw09structvalidator

import (
	"reflect"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	panic("implement me")
}

func vaildateInt(v int, rule string) error {
	rules := strings.Split(rule, "|")
	if len(rules) < 1 {
		return errBadValidator
	}

	for _, entry := range rules {
		values := strings.Split(entry, ":")
		if len(values) < 2 {
			return errBadValidator
		}

		var err error
		switch values[0] {
		case "in":
			err = validateIntIn(v, values[1])
		case "min":
			err = validateIntMin(v, values[1])
		case "max":
			err = validateIntMax(v, values[1])
		default:
			return errBadValidator
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func vaildateString(v string, rule string) error {
	rules := strings.Split(rule, "|")
	if len(rules) < 1 {
		return errBadValidator
	}

	for _, entry := range rules {
		values := strings.Split(entry, ":")
		if len(values) < 2 {
			return errBadValidator
		}

		var err error
		switch values[0] {
		case "in":
			err = validateStringIn(v, values[1])
		case "len":
			err = validateStringLen(v, values[1])
		case "regexp":
			err = validateStringRegexp(v, values[1])
		default:
			return errBadValidator
		}
		if err != nil {
			return err
		}
	}
	return nil
}

//nolint:exhaustive
func Validate(v interface{}) error {
	reflectValue := reflect.ValueOf(v)
	if reflectValue.Kind() != reflect.Struct {
		return errNotAStruct
	}
	vType := reflect.TypeOf(v)
	for i := 0; i < vType.NumField(); i++ {
		field := vType.Field(i)

		tag := field.Tag
		validate := tag.Get("validate")
		if validate == "" {
			continue
		}

		fieldKind := field.Type.Kind()
		switch fieldKind {
		case reflect.Slice:
			switch field.Type.Elem().Kind() {
			case reflect.Int:
				intSlice := reflectValue.Field(i).Interface().([]int)
				for _, entry := range intSlice {
					err := vaildateInt(entry, validate)
					if err != nil {
						return err
					}
				}
			case reflect.String:
				stringSlice := reflectValue.Field(i).Interface().([]string)
				for _, entry := range stringSlice {
					err := vaildateString(entry, validate)
					if err != nil {
						return err
					}
				}
			default:
				continue
			}
		case reflect.Int:
			underlyingInt := reflectValue.Field(i).Interface().(int)
			return vaildateInt(underlyingInt, validate)
		case reflect.String:
			underlyingString := reflectValue.Field(i).Interface().(string)
			return vaildateString(underlyingString, validate)
		default:
			continue
		}
	}
	return nil
}
