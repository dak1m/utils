package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// ConversionFrom is an interface to allow retrieve data from data source
type ConversionFrom interface {
	FromSource(string) error
}

// MapConvertStructByTag 通过tag匹配结构体字段并赋值
func MapConvertStructByTag(input map[string]string, obj interface{}, tag string) (err error) {
	if reflect.TypeOf(obj).Kind() != reflect.Ptr {
		return
	}
	t := reflect.TypeOf(obj).Elem()
	v := reflect.ValueOf(obj).Elem()
	for i := 0; i < t.NumField(); i++ {
		tagName := t.Field(i).Tag.Get(tag)
		field := t.Field(i)
		if tagName == "-" {
			continue
		}
		tagName = strings.Replace(tagName, ",omitempty", "", -1)
		if field.Type.Kind() == reflect.Struct && tagName == "" && field.Anonymous {
			// 匿名嵌套结构体
			AnonymousStructAssignment(t.Field(i), v.Field(i), tag, input)
			continue
		}

		if field.Type.Kind() == reflect.Ptr {
			field.Type = field.Type.Elem()
		}
		flag := 1
		if (field.Type.Kind() == reflect.Struct || field.Type.Kind() == reflect.Slice || field.Type.Kind() == reflect.Array) && tagName != "" {
			if val, ok := input[tagName]; ok {
				newValue := reflect.New(field.Type)
				v.Field(i).Set(newValue)
				if v.Field(i).CanAddr() {
					if structConvert, ok2 := v.Field(i).Addr().Interface().(ConversionFrom); ok2 {
						return structConvert.FromSource(val)
					} else {
						_ = json.Unmarshal([]byte(val), v.Field(i).Addr().Interface())
					}
				}
				if structConvert, ok2 := v.Field(i).Interface().(ConversionFrom); ok2 {
					return structConvert.FromSource(val)
				}

				flag = 2
			}
		}

		if tagName == "" {
			continue
		}

		switch v.Field(i).Interface().(type) {
		case time.Time:
			if val, ok := input[tagName]; ok {
				time, err := time.ParseInLocation(time.RFC3339, val, time.Local)
				if err != nil {
					return err
				}
				v.Field(i).Set(reflect.ValueOf(time))
			}
		case *time.Time:
			if val, ok := input[tagName]; ok {
				time, err := time.ParseInLocation(time.RFC3339, val, time.Local)
				if err != nil {
					return err
				}
				v.Field(i).Set(reflect.ValueOf(&time))
			}
		default:
			if val, ok := input[tagName]; ok && flag == 1 {
				if val == "" {
					val = field.Tag.Get("default")
				}
				var setVal reflect.Value
				setVal, err = decode(val, field.Type.Kind())
				if err != nil {
					return
				}
				v.Field(i).Set(setVal)
			}
		}
	}

	return
}

func decode(input string, intType reflect.Kind) (resVal reflect.Value, err error) {
	switch intType {
	case reflect.Int:
		val, err := strconv.Atoi(input)
		if err != nil {
			return reflect.Value{}, err
		}
		resVal = reflect.ValueOf(val)
	case reflect.Int8:
		val, err := strconv.ParseInt(input, 10, 8)
		if err != nil {
			return reflect.Value{}, err
		}
		resVal = reflect.ValueOf(int8(val))
	case reflect.Int16:
		val, err := strconv.ParseInt(input, 10, 16)
		if err != nil {
			return reflect.Value{}, err
		}
		resVal = reflect.ValueOf(int16(val))
	case reflect.Int32:
		val, err := strconv.ParseInt(input, 10, 32)
		if err != nil {
			return reflect.Value{}, err
		}
		resVal = reflect.ValueOf(int32(val))
	case reflect.Int64:
		val, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			return reflect.Value{}, err
		}
		resVal = reflect.ValueOf(val)
	case reflect.Uint:
		val, err := strconv.ParseUint(input, 10, 0)
		if err != nil {
			return reflect.Value{}, err
		}
		resVal = reflect.ValueOf(uint(val))
	case reflect.Uint8:
		val, err := strconv.ParseUint(input, 10, 8)
		if err != nil {
			return reflect.Value{}, err
		}
		resVal = reflect.ValueOf(uint8(val))
	case reflect.Uint16:
		val, err := strconv.ParseUint(input, 10, 16)
		if err != nil {
			return reflect.Value{}, err
		}
		resVal = reflect.ValueOf(uint16(val))
	case reflect.Uint32:
		val, err := strconv.ParseUint(input, 10, 32)
		if err != nil {
			return reflect.Value{}, err
		}
		resVal = reflect.ValueOf(uint32(val))
	case reflect.Uint64:
		val, err := strconv.ParseUint(input, 10, 64)
		if err != nil {
			return reflect.Value{}, err
		}
		resVal = reflect.ValueOf(val)
	case reflect.String:
		resVal = reflect.ValueOf(input)
	case reflect.Bool:
		val, err := strconv.ParseBool(input)
		if err != nil {
			return reflect.Value{}, err
		}
		resVal = reflect.ValueOf(val)
	case reflect.Float32:
		val, err := strconv.ParseFloat(input, 32)
		if err != nil {
			return reflect.Value{}, err
		}
		resVal = reflect.ValueOf(float32(val))
	case reflect.Float64:
		val, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return reflect.Value{}, err
		}
		resVal = reflect.ValueOf(val)
	default:
		return reflect.Value{}, fmt.Errorf("unsupported type")
	}

	return resVal, nil
}

func StructConvertMapByTag(obj interface{}, tag string) map[string]any {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	data := make(map[string]any)

	for i := 0; i < t.NumField(); i++ {
		tagName := t.Field(i).Tag.Get(tag)
		if tagName == "-" {
			continue
		}

		tagName = strings.Replace(tagName, ",omitempty", "", -1)
		field := t.Field(i)
		if field.Type.Kind() == reflect.Struct && tagName == "" && field.Anonymous {
			// 匿名嵌套结构体
			AnonymousStructToMap(t.Field(i), v.Field(i), tag, data)
			continue
		}

		ptr := false
		if field.Type.Kind() == reflect.Ptr {
			field.Type = field.Type.Elem()
			ptr = true
		}
		flag := 1
		if (field.Type.Kind() == reflect.Struct || field.Type.Kind() == reflect.Slice || field.Type.Kind() == reflect.Array) && tagName != "" {
			if v.Field(i).IsZero() {
				continue
			}

			if ptr {
				if v.Field(i).IsNil() {
					continue
				}
			}

			b, _ := json.Marshal(v.Field(i).Interface())
			data[tagName] = string(b)
			flag = 2
		}

		if tagName == "" {
			continue
		}

		switch value := v.Field(i).Interface().(type) {
		case time.Time:
			if v.Field(i).IsZero() {
				continue
			} else {
				data[tagName] = value.Format(time.RFC3339)
			}
		case *time.Time:
			if v.Field(i).IsNil() {
				continue
			} else {
				data[tagName] = value.Format(time.RFC3339)
			}
		default:
			if flag == 1 {
				data[tagName] = v.Field(i).Interface()
			}
		}

	}
	return data
}

func AnonymousStructToMap(ft reflect.StructField, fv reflect.Value, tag string, data map[string]any) {
	for j := 0; j < fv.NumField(); j++ {
		field := ft.Type.Field(j)
		tagName := field.Tag.Get(tag)
		if tagName == "-" {
			continue
		}
		tagName = strings.Replace(tagName, ",omitempty", "", -1)
		if field.Type.Kind() == reflect.Struct && tagName == "" && field.Anonymous {
			AnonymousStructToMap(field, fv.Field(j), tag, data)
			continue
		}

		ptr := false
		if field.Type.Kind() == reflect.Ptr {
			field.Type = field.Type.Elem()
			ptr = true
		}
		flag := 1
		if (field.Type.Kind() == reflect.Struct || field.Type.Kind() == reflect.Slice || field.Type.Kind() == reflect.Array) && tagName != "" {
			if fv.Field(j).IsZero() {
				continue
			}

			if ptr {
				if fv.Field(j).IsNil() {
					continue
				}
			}

			data[tagName], _ = json.Marshal(fv.Field(j).Interface())
			flag = 2
		}

		if tagName == "" {
			continue
		}

		switch value := fv.Field(j).Interface().(type) {
		case time.Time:
			if fv.Field(j).IsZero() {
				continue
			} else {
				data[tagName] = value.Format(time.RFC3339)
			}
		case *time.Time:
			if fv.Field(j).IsNil() {
				continue
			} else {
				data[tagName] = value.Format(time.RFC3339)
			}
		default:
			if flag == 1 {
				data[tagName] = fv.Field(j).Interface()
			}
		}
	}
}

func AnonymousStructAssignment(ft reflect.StructField, fv reflect.Value, tag string, input map[string]string) (err error) {
	for j := 0; j < fv.NumField(); j++ {
		field := ft.Type.Field(j)
		tagName := field.Tag.Get(tag)
		if tagName == "-" {
			continue
		}
		tagName = strings.Replace(tagName, ",omitempty", "", -1)
		if field.Type.Kind() == reflect.Struct && tagName == "" && field.Anonymous {
			AnonymousStructAssignment(field, fv.Field(j), tag, input)
			continue
		}

		if field.Type.Kind() == reflect.Ptr {
			field.Type = field.Type.Elem()
		}
		flag := 1
		if (field.Type.Kind() == reflect.Struct || field.Type.Kind() == reflect.Slice || field.Type.Kind() == reflect.Array) && tagName != "" {
			if val, ok := input[tagName]; ok {
				_ = json.Unmarshal([]byte(val), fv.Field(j).Addr().Interface())
				flag = 2
			}
		}

		if tagName == "" {
			continue
		}

		switch fv.Field(j).Interface().(type) {
		case time.Time:
			if val, ok := input[tagName]; ok && val != "" {
				t, err := time.ParseInLocation(time.RFC3339, val, time.Local)
				if err != nil {
					return err
				}
				fv.Field(j).Set(reflect.ValueOf(t))
			}
		case *time.Time:
			if val, ok := input[tagName]; ok && val != "" {
				t, err := time.ParseInLocation(time.RFC3339, val, time.Local)
				if err != nil {
					return err
				}
				fv.Field(j).Set(reflect.ValueOf(&t))
			}
		default:
			if val, ok := input[tagName]; ok && flag == 1 {
				if val == "" {
					val = field.Tag.Get("default")
				}
				var setVal reflect.Value
				setVal, err := decode(val, field.Type.Kind())
				if err != nil {
					return err
				}
				fv.Field(j).Set(setVal)
			}
		}
	}

	return nil
}

func HandleFields(fields map[string]any) map[string]any {
	newFields := make(map[string]any, len(fields))
	for k, v := range fields {
		t := reflect.TypeOf(v)
		flag := 1
		if t.Kind() == reflect.Struct || t.Kind() == reflect.Slice {
			newFields[k], _ = json.Marshal(v)
			flag = 2
		}
		switch value := v.(type) {
		case time.Time:
			if v.(time.Time).IsZero() {
				continue
			} else {
				newFields[k] = value.Format(time.RFC3339)
			}
		case *time.Time:
			if v.(*time.Time).IsZero() {
				continue
			} else {
				newFields[k] = value.Format(time.RFC3339)
			}
		default:
			if flag == 1 {
				newFields[k] = v
			}
		}
	}

	return newFields
}

// AllFieldsByTag 根据tag获取对象中的所有字段
func AllFieldsByTag(obj interface{}, tag string) []string {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	data := make([]string, 0)

	for i := 0; i < t.NumField(); i++ {
		tagName := t.Field(i).Tag.Get(tag)
		if tagName == "-" {
			continue
		}

		tagName = strings.Replace(tagName, ",omitempty", "", -1)
		field := t.Field(i)
		if field.Type.Kind() == reflect.Struct && tagName == "" && field.Anonymous {
			// 匿名嵌套结构体
			data = anonymousStructToField(t.Field(i), v.Field(i), tag, data)
			continue
		} else {
			data = append(data, tagName)
		}
	}
	return data
}

func anonymousStructToField(field reflect.StructField, value reflect.Value, tag string, data []string) []string {
	t := field.Type
	for i := 0; i < t.NumField(); i++ {
		tagName := t.Field(i).Tag.Get(tag)
		if tagName == "-" {
			continue
		}

		tagName = strings.Replace(tagName, ",omitempty", "", -1)
		field := t.Field(i)
		if field.Type.Kind() == reflect.Struct && tagName == "" && field.Anonymous {
			// 匿名嵌套结构体
			anonymousStructToField(t.Field(i), value.Field(i), tag, data)
			continue
		} else {
			data = append(data, tagName)
		}
	}

	return data
}
