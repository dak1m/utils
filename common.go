package utils

import (
	"fmt"
	"math"
	"reflect"
	"slices"
	"strconv"
)

var DefaultTag = "json"

// IsOnlySet checks if only the X field is set and all other fields are zero values.
func IsOnlySet(obj any, field string, exclude ...string) bool {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	var fnc func(reflect.Type, reflect.Value) bool
	fnc = func(t reflect.Type, v reflect.Value) bool {
		for i := 0; i < v.NumField(); i++ {
			reflectField := v.Field(i)
			name := t.Field(i).Tag.Get(DefaultTag)

			if t.Field(i).Type.Kind() == reflect.Struct && name == "" && t.Field(i).Anonymous {
				// 匿名嵌套结构体
				fnc(t.Field(i).Type, v.Field(i))
			} else {
				if name == field {
					if reflectField.IsZero() {
						return false
					}
				} else if !slices.Contains(exclude, name) {
					if !reflectField.IsZero() {
						return false

					}
				}
			}

		}

		return true
	}

	return fnc(t, v)
}

func FloatStrTrim(s string) string {
	f, _ := strconv.ParseFloat(s, 64)
	d1 := int(f*math.Pow10(2)) % 10
	d2 := int(f*math.Pow10(1)) % 10
	if d1 == 0 && d2 == 0 {
		return fmt.Sprintf("%.0f", f)
	} else if d1 == 0 {
		return fmt.Sprintf("%.1f", f)
	} else {
		return fmt.Sprintf("%.2f", f)
	}
}
