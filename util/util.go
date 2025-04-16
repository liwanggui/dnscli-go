package util

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func GetStructFieldNamesAndValues(s interface{}, tagName string, exclField []string) (fieldNameList, fieldValueList []string) {
	set := make(map[string]bool)
	for _, v := range exclField {
		set[v] = true
	}
	// 获取传入参数的反射值和类型
	value := reflect.ValueOf(s)

	// 如果传入的是指针，需要解引用
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	// 获取结构体的类型
	typeOf := value.Type()

	// 遍历结构体的所有字段
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldType := typeOf.Field(i)

		// 排除字段
		if set[fieldType.Name] {
			continue
		}

		// 获取 tag 名
		fieldName := fieldType.Tag.Get(tagName)

		if fieldName == "" {
			fieldName = fieldType.Name
		}
		fieldNameList = append(fieldNameList, strings.Split(fieldName, ",")[0])

		// 获取字段值
		switch field.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fieldValueList = append(fieldValueList, strconv.FormatInt(field.Int(), 10))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			fieldValueList = append(fieldValueList, strconv.FormatUint(field.Uint(), 10))
		case reflect.Float32, reflect.Float64:
			fieldValueList = append(fieldValueList, strconv.FormatFloat(field.Float(), 'g', -1, 64))
		case reflect.String:
			fieldValueList = append(fieldValueList, field.String())
		case reflect.Bool:
			fieldValueList = append(fieldValueList, fmt.Sprintf("%t", field.Bool()))
		default:
			fieldValueList = append(fieldValueList, fmt.Sprintf("%v", field.Interface()))
		}
	}
	return
}
