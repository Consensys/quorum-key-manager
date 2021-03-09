package utils

import (
	"fmt"
	"reflect"
)

func ExtractType(element interface{}) string {
	v := reflect.ValueOf(element).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.Ptr && field.Elem().Kind() == reflect.Struct {
			if !field.IsNil() {
				return v.Type().Field(i).Name
			}
		}
	}
	return ""
}

func ExtractField(cfg interface{}) (interface{}, error) {
	v := reflect.ValueOf(cfg).Elem()
	nonNilFieldsCount := 0
	var rv interface{}
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.Ptr && field.Elem().Kind() == reflect.Struct {
			if !field.IsNil() {
				nonNilFieldsCount++
				rv = field.Interface()
			}
		}
	}

	if nonNilFieldsCount == 0 {
		return nil, fmt.Errorf("invalid configuration: empty")
	}

	if nonNilFieldsCount > 1 {
		return nil, fmt.Errorf("invalid configuration: multiple fields provided")
	}

	return rv, nil
}
