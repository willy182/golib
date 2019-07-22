package golib

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

// ParseToFormValue convert struct to form values
func ParseToFormValue(source interface{}) (form url.Values, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	form = url.Values{}
	val := reflect.ValueOf(source)
	if val.Kind() == reflect.Ptr {
		val = val.Elem() // take element if source type is pointer
	}

	// must struct
	if val.Kind() != reflect.Struct {
		err = fmt.Errorf("invalid source type %v: must struct", val.Kind())
		return
	}

	objType := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		jsonTag := objType.Field(i).Tag.Get("json")
		if jsonTags := strings.Split(jsonTag, ","); len(jsonTags) > 0 {
			jsonTag = jsonTags[0]
		}
		if jsonTag == "" {
			jsonTag = objType.Field(i).Name
		}

		var value string
		if field.Kind() != reflect.String {
			value = fmt.Sprint(field.Interface())
		} else {
			value = field.String()
		}
		form.Set(jsonTag, value)
	}
	return
}
