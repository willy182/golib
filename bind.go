package golib

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// BindQueryParam binding query param from HTTP context to struct model with key in json tag
func BindQueryParam(u *url.URL, target interface{}) error {
	refValue := reflect.ValueOf(target)
	if refValue.Kind() != reflect.Ptr {
		return fmt.Errorf("target is not pointer")
	}
	refValue = refValue.Elem()
	q := make(map[string]string)
	urlPathUnscape, _ := url.PathUnescape(u.RawQuery)
	for _, r := range strings.Split(urlPathUnscape, "&") {
		sp := strings.Split(r, "=")
		if len(sp) > 1 {
			q[sp[0]] = sp[1]
		}
	}

	for i := 0; i < refValue.NumField(); i++ {
		field := refValue.Field(i)

		jsonTag := refValue.Type().Field(i).Tag.Get("json")
		jsonTag = strings.TrimSuffix(jsonTag, ",omitempty")
		field.SetString(q[jsonTag])
	}
	return nil
}

// ParseFromQueryParam parse url query string to struct target (with multiple data type in struct field), target must in pointer
func ParseFromQueryParam(query url.Values, target interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	var parseDataTypeValue func(typ reflect.Type, val reflect.Value)

	var errs = NewMultiError()

	pValue := reflect.ValueOf(target)
	if pValue.Kind() != reflect.Ptr {
		panic(fmt.Errorf("%v is not pointer", pValue.Kind()))
	}
	pValue = pValue.Elem()
	pType := reflect.TypeOf(target).Elem()
	for i := 0; i < pValue.NumField(); i++ {
		field := pValue.Field(i)
		typ := pType.Field(i)
		if typ.Anonymous { // embedded struct
			err = ParseFromQueryParam(query, field.Addr().Interface())
		}

		key := strings.TrimSuffix(typ.Tag.Get("json"), ",omitempty")
		if key == "-" {
			continue
		}

		var v string
		if val := query[key]; len(val) > 0 && val[0] != "" {
			v = val[0]
		} else {
			v = typ.Tag.Get("default")
		}

		parseDataTypeValue = func(tp reflect.Type, targetField reflect.Value) {
			switch tp.Kind() {
			case reflect.String:
				if ok, _ := strconv.ParseBool(typ.Tag.Get("lower")); ok {
					v = strings.ToLower(v)
				}
				targetField.SetString(v)
			case reflect.Int32, reflect.Int, reflect.Int64:
				vInt, err := strconv.Atoi(v)
				if v != "" && err != nil {
					errs.Append(key, fmt.Errorf("Cannot parse '%s' (%T) to type number", v, v))
				}
				targetField.SetInt(int64(vInt))
			case reflect.Bool:
				vBool, err := strconv.ParseBool(v)
				if v != "" && err != nil {
					errs.Append(key, fmt.Errorf("Cannot parse '%s' (%T) to type boolean", v, v))
				}
				targetField.SetBool(vBool)
			case reflect.Ptr:
				if v != "" {
					// allocate new value to pointer targetField
					targetField.Set(reflect.ValueOf(reflect.New(tp.Elem()).Interface()))
					parseDataTypeValue(tp.Elem(), targetField.Elem())
				}
			}
		}

		parseDataTypeValue(field.Type(), field)
	}

	if errs.HasError() {
		return errs
	}

	return
}
