package jhttp

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/Freelance-launchpad/backend-interview/common/jdate"
	"github.com/google/uuid"
)

// isEmptyValue check whether the value is empty or not.
// shamefully copied from encoding/json
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.Interface, reflect.Pointer:
		return v.IsZero()
	}
	return false
}

// AddStructInValues adds fields in the data struct to values using the `form` tag a key.
// If the `form` tag is "-" or the field is a nil pointer, it will be ignored.
// Options omitempty is accepted
// It will return an error if data is not a struct.
// Fields can be array or slice, in which case each element of the list will be added to values.
// Following encoding/json package system, the name of the field must be the first elem in the tag.
// This function allows nested struct/composition
func AddStructInValues(values url.Values, data any) error {
	t := reflect.TypeOf(data)
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("not a struct: %s", t.Kind().String())
	}

	for field := range t.Fields() {
		if !field.IsExported() {
			continue
		}

		key, opts, _ := strings.Cut(field.Tag.Get("form"), ",")
		if key == "" {
			key = strings.ToLower(field.Name)
		} else if key == "-" {
			continue
		}

		fieldVal := reflect.ValueOf(data).FieldByIndex(field.Index)
		if fieldVal.Kind() == reflect.Pointer {
			if fieldVal.IsNil() {
				continue
			}
			fieldVal = fieldVal.Elem()
		}

		if isEmptyValue(fieldVal) && strings.Contains(opts, "omitempty") {
			continue
		}

		var vals []reflect.Value
		if fieldVal.Type() != reflect.TypeFor[uuid.UUID]() &&
			(fieldVal.Kind() == reflect.Slice || fieldVal.Kind() == reflect.Array) {
			for j := 0; j < fieldVal.Len(); j++ {
				vals = append(vals, fieldVal.Index(j))
			}
		} else {
			vals = append(vals, fieldVal)
		}

		for _, val := range vals {
			var s string
			switch {
			case val.Type() == reflect.TypeFor[time.Time]():
				s = val.Interface().(time.Time).Format(time.RFC3339)
			case val.Type() == reflect.TypeFor[jdate.Date]():
				s = val.Interface().(jdate.Date).Format(time.DateOnly)
			case val.CanInterface() && reflect.TypeOf(val.Interface()).Implements(reflect.TypeFor[fmt.Stringer]()):
				s = val.Interface().(fmt.Stringer).String()
			case val.CanInt():
				s = strconv.FormatInt(val.Int(), 10)
			case val.CanUint():
				s = strconv.FormatUint(val.Uint(), 10)
			case val.CanFloat():
				s = strconv.FormatFloat(val.Float(), 'g', 5, 64)
			case val.Kind() == reflect.String:
				s = val.String()
			case val.Kind() == reflect.Bool:
				s = strconv.FormatBool(val.Bool())
			default:
				if val.Kind() == reflect.Struct {
					if err := AddStructInValues(values, val.Interface()); err != nil {
						return err
					}
					continue
				} else {
					return fmt.Errorf("unsupported type: %s", val.Type())
				}
			}
			values.Add(key, s)
		}
	}

	return nil
}
