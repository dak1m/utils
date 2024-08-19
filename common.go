package utils

import "reflect"

var DefaultTag = "json"

// IsOnlySet checks if only the X field is set and all other fields are zero values.
func IsOnlySet(obj any, field string) bool {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	for i := 0; i < v.NumField(); i++ {
		reflectField := v.Field(i)
		name := t.Field(i).Tag.Get(DefaultTag)
		if name == field {
			if reflectField.IsZero() {
				return false
			}
		} else {
			if !reflectField.IsZero() {
				return false
			}
		}
	}
	return true
}
