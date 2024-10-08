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

		fieldType := field.Type
		if field.Type.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}
		flag := 1
		if (fieldType.Kind() == reflect.Struct || fieldType.Kind() == reflect.Slice || fieldType.Kind() == reflect.Array) && tagName != "" {
			if val, ok := input[tagName]; ok {
				if v.Field(i).CanAddr() {
					if structConvert, ok2 := v.Field(i).Addr().Interface().(ConversionFrom); ok2 {
						err = structConvert.FromSource(val)
						if err != nil {
							return
						}
						continue
					}
				}
				if _, ok2 := v.Field(i).Interface().(ConversionFrom); ok2 {
					newValue := reflect.New(fieldType)
					v.Field(i).Set(newValue)
					structConvert := v.Field(i).Interface().(ConversionFrom)
					err = structConvert.FromSource(val)
					if err != nil {
						return
					}
					continue
				}
				if v.Field(i).CanAddr() {
					_ = json.Unmarshal([]byte(val), v.Field(i).Addr().Interface())
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
				timeParse, err := time.ParseInLocation(time.RFC3339, val, time.Local)
				if err != nil {
					timeParse, err = time.ParseInLocation(time.DateTime, val, time.Local)
					if err != nil {
						return err
					}
				}
				v.Field(i).Set(reflect.ValueOf(timeParse))
			}
		case *time.Time:
			if val, ok := input[tagName]; ok {
				timeParse, err := time.ParseInLocation(time.RFC3339, val, time.Local)
				if err != nil {
					timeParse, err = time.ParseInLocation(time.DateTime, val, time.Local)
					if err != nil {
						return err
					}
					return err
				}
				v.Field(i).Set(reflect.ValueOf(&timeParse))
			}
		default:
			if val, ok := input[tagName]; ok && flag == 1 {
				if val == "" {
					val = field.Tag.Get("default")
				}
				var setVal reflect.Value
				setVal, err = decode(val, field.Type, fieldType.Kind())
				if err != nil {
					return
				}
				v.Field(i).Set(setVal)
			}
		}
	}

	return
}

func decode(input string, fieldType reflect.Type, inType reflect.Kind) (resVal reflect.Value, err error) {
	fieldVal := func(val any) reflect.Value {
		reflectVal := reflect.ValueOf(val)
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
			ptr := reflect.New(fieldType)
			ptr.Elem().Set(reflect.ValueOf(val))

			reflectVal = ptr
		}
		if fieldType.PkgPath() != "" {
			return reflectVal.Convert(fieldType)
		}
		return reflectVal
	}
	switch inType {
	case reflect.Int:
		val, err := strconv.Atoi(input)
		if err != nil {
			return reflect.Value{}, err
		}
		resVal = fieldVal(val)
	case reflect.Int8:
		val, err := strconv.ParseInt(input, 10, 8)
		if err != nil {
			return reflect.Value{}, err
		}
		resVal = fieldVal(int8(val))
	case reflect.Int16:
		val, err := strconv.ParseInt(input, 10, 16)
		if err != nil {
			return reflect.Value{}, err
		}
		resVal = fieldVal(int16(val))
	case reflect.Int32:
		val, err := strconv.ParseInt(input, 10, 32)
		if err != nil {
			return reflect.Value{}, err
		}
		resVal = fieldVal(int32(val))
	case reflect.Int64:
		val, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			return reflect.Value{}, err
		}
		resVal = fieldVal(val)
	case reflect.Uint:
		val, err := strconv.ParseUint(input, 10, 0)
		if err != nil {
			return reflect.Value{}, err
		}
		resVal = fieldVal(uint(val))
	case reflect.Uint8:
		val, err := strconv.ParseUint(input, 10, 8)
		if err != nil {
			return reflect.Value{}, err
		}
		resVal = fieldVal(uint8(val))
	case reflect.Uint16:
		val, err := strconv.ParseUint(input, 10, 16)
		if err != nil {
			return reflect.Value{}, err
		}
		resVal = fieldVal(uint16(val))
	case reflect.Uint32:
		val, err := strconv.ParseUint(input, 10, 32)
		if err != nil {
			return reflect.Value{}, err
		}
		resVal = fieldVal(uint32(val))
	case reflect.Uint64:
		val, err := strconv.ParseUint(input, 10, 64)
		if err != nil {
			return reflect.Value{}, err
		}
		resVal = fieldVal(val)
	case reflect.String:
		resVal = fieldVal(input)
	case reflect.Bool:
		val, err := strconv.ParseBool(input)
		if err != nil {
			return reflect.Value{}, err
		}
		resVal = fieldVal(val)
	case reflect.Float32:
		val, err := strconv.ParseFloat(input, 32)
		if err != nil {
			return reflect.Value{}, err
		}
		resVal = fieldVal(float32(val))
	case reflect.Float64:
		val, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return reflect.Value{}, err
		}
		resVal = fieldVal(val)
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
				if ptr {
					if v.Field(i).IsNil() {
						continue
					} else {
						data[tagName] = reflect.Indirect(v.Field(i)).Interface()
						continue
					}
				}
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
			b, _ := json.Marshal(fv.Field(j).Interface())
			data[tagName] = string(b)
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
				if ptr {
					if fv.Field(j).IsNil() {
						continue
					} else {
						data[tagName] = reflect.Indirect(fv.Field(j)).Interface()
						continue
					}
				}
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

		fieldType := field.Type
		if field.Type.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}
		flag := 1
		if (fieldType.Kind() == reflect.Struct || fieldType.Kind() == reflect.Slice || fieldType.Kind() == reflect.Array) && tagName != "" {
			if val, ok := input[tagName]; ok {
				if fv.Field(j).CanAddr() {
					if structConvert, ok2 := fv.Field(j).Addr().Interface().(ConversionFrom); ok2 {
						err = structConvert.FromSource(val)
						if err != nil {
							return
						}
						continue
					}
				}
				if _, ok2 := fv.Field(j).Interface().(ConversionFrom); ok2 {
					newValue := reflect.New(fieldType)
					fv.Field(j).Set(newValue)
					structConvert := fv.Field(j).Interface().(ConversionFrom)
					err = structConvert.FromSource(val)
					if err != nil {
						return
					}
					continue
				}
				if fv.Field(j).CanAddr() {
					_ = json.Unmarshal([]byte(val), fv.Field(j).Addr().Interface())
				}

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
					t, err = time.ParseInLocation(time.DateTime, val, time.Local)
					if err != nil {
						return err
					}
				}
				fv.Field(j).Set(reflect.ValueOf(t))
			}
		case *time.Time:
			if val, ok := input[tagName]; ok && val != "" {
				t, err := time.ParseInLocation(time.RFC3339, val, time.Local)
				if err != nil {
					t, err = time.ParseInLocation(time.DateTime, val, time.Local)
					if err != nil {
						return err
					}
				}
				fv.Field(j).Set(reflect.ValueOf(&t))
			}
		default:
			if val, ok := input[tagName]; ok && flag == 1 {
				if val == "" {
					val = field.Tag.Get("default")
				}
				var setVal reflect.Value
				setVal, err := decode(val, field.Type, fieldType.Kind())
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
		if t == nil {
			continue
		}
		flag := 1
		if t.Kind() == reflect.Struct || t.Kind() == reflect.Slice {
			b, _ := json.Marshal(v)
			newFields[k] = string(b)
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

	var fnc func(reflect.Type, reflect.Value)
	fnc = func(t reflect.Type, v reflect.Value) {
		for i := 0; i < t.NumField(); i++ {
			tagName := t.Field(i).Tag.Get(tag)
			if tagName == "-" {
				continue
			}

			tagName = strings.Replace(tagName, ",omitempty", "", -1)
			field := t.Field(i)
			if field.Type.Kind() == reflect.Struct && tagName == "" && field.Anonymous {
				// 匿名嵌套结构体
				fnc(t.Field(i).Type, v.Field(i))
			} else {
				if tagName == "" {
					continue
				}
				data = append(data, tagName)
			}
		}
	}

	fnc(t, v)

	return data
}

func GenerateTypeMapping(obj any, tagName string) map[string]any {
	typeMapping := make(map[string]interface{})
	fields := reflect.TypeOf(obj)

	for i := 0; i < fields.NumField(); i++ {
		field := fields.Field(i)
		fieldName := field.Tag.Get(tagName)
		fieldType := field.Type

		switch fieldType.Kind() {
		case reflect.Slice:
			elemType := fieldType.Elem()
			switch elemType.Kind() {
			case reflect.Struct:
				typeMapping[fieldName] = reflect.New(reflect.SliceOf(elemType)).Interface()
			case reflect.String:
				typeMapping[fieldName] = &[]string{}
			case reflect.Uint:
				typeMapping[fieldName] = &[]uint{}
			case reflect.Uint8:
				typeMapping[fieldName] = &[]uint8{}
			case reflect.Uint16:
				typeMapping[fieldName] = &[]uint16{}
			case reflect.Uint32:
				typeMapping[fieldName] = &[]uint32{}
			case reflect.Uint64:
				typeMapping[fieldName] = &[]uint64{}
			case reflect.Int:
				typeMapping[fieldName] = &[]int{}
			case reflect.Int8:
				typeMapping[fieldName] = &[]int8{}
			case reflect.Int16:
				typeMapping[fieldName] = &[]int16{}
			case reflect.Int32:
				typeMapping[fieldName] = &[]int32{}
			case reflect.Int64:
				typeMapping[fieldName] = &[]int64{}
			case reflect.Float64:
				typeMapping[fieldName] = &[]float64{}
			case reflect.Float32:
				typeMapping[fieldName] = &[]float32{}
			case reflect.Bool:
				typeMapping[fieldName] = &[]bool{}
			case reflect.Interface:
				typeMapping[fieldName] = &[]any{}
			default:
				continue
			}
		case reflect.Struct:
			if !field.Anonymous {
				typeMapping[fieldName] = reflect.New(fieldType).Interface()
			}
		default:
			continue
		}
	}

	return typeMapping
}
