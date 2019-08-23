package golib

import (
	"fmt"
	"net/url"
	"reflect"
)

// BindQueryParam binding query param from HTTP context to struct model with key in json tag
func BindQueryParam(u *url.URL, target interface{}) error {
	refValue := reflect.ValueOf(target)
	if refValue.Kind() != reflect.Ptr {
		return fmt.Errorf("target is not pointer")
	}
	refValue = refValue.Elem()
	q := u.Query()

	for i := 0; i < refValue.NumField(); i++ {
		field := refValue.Field(i)

		jsonTag := refValue.Type().Field(i).Tag.Get("json")
		field.SetString(q.Get(jsonTag))
	}
	return nil
}
