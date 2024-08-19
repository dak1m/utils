package utils

import (
	"reflect"
	"slices"
)

var DefaultTag = "json"

// IsOnlySet checks if only the X field is set and all other fields are zero values.
func IsOnlySet(obj any, field string, exclude ...string) bool {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	for i := 0; i < v.NumField(); i++ {
		reflectField := v.Field(i)
		name := t.Field(i).Tag.Get(DefaultTag)
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
	return true
}
