package jsonschema

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/Bhinneka/golib"
	"github.com/xeipuuv/gojsonschema"
)

var jsonSchemaList = map[string]*gojsonschema.Schema{}

// Load all schema
func Load(path string) error {
	return filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		fileName := info.Name()
		if strings.HasSuffix(fileName, ".json") {
			s, err := ioutil.ReadFile(p)
			if err != nil {
				return err
			}

			var data interface{}
			err = json.Unmarshal(s, &data)
			if err != nil {
				return err
			}

			t := reflect.ValueOf(data)
			if t.Kind() == reflect.Slice {
				for i := 0; i < t.Len(); i++ {
					obj := t.Index(i).Interface()
					id, err := getID(obj)
					if err != nil {
						continue
					}
					jsonSchemaList[id], err = gojsonschema.NewSchema(gojsonschema.NewGoLoader(obj))
					if err != nil {
						continue
					}
				}
			} else {
				id, err := getID(data)
				if err != nil {
					return nil
				}
				jsonSchemaList[id], err = gojsonschema.NewSchema(gojsonschema.NewBytesLoader(s))
				if err != nil {
					return nil
				}
			}
		}
		return nil
	})
}

func getID(obj interface{}) (id string, err error) {
	m, ok := obj.(map[string]interface{})
	if !ok {
		err = errors.New("invalid type")
		return
	}
	id, ok = m["id"].(string)
	if !ok {
		err = errors.New("ID not found in schema")
	}
	return
}

// Get json schema by ID
func Get(schemaID string) (*gojsonschema.Schema, error) {
	schema, ok := jsonSchemaList[schemaID]
	if !ok {
		return nil, fmt.Errorf("schema '%s' not found", schemaID)
	}

	return schema, nil
}

// Validate from Go data type
func Validate(schemaID string, input interface{}) error {
	multiError := golib.NewMultiError()

	schema, err := Get(schemaID)
	if err != nil {
		multiError.Append("getSchema", err)
		return multiError
	}

	document := gojsonschema.NewGoLoader(input)
	return validate(schema, document)
}

// ValidateDocument document
func ValidateDocument(schemaID string, jsonByte []byte) error {
	multiError := golib.NewMultiError()

	schema, err := Get(schemaID)
	if err != nil {
		multiError.Append("getSchema", err)
		return multiError
	}

	document := gojsonschema.NewBytesLoader(jsonByte)
	return validate(schema, document)
}

func validate(schema *gojsonschema.Schema, document gojsonschema.JSONLoader) error {
	multiError := golib.NewMultiError()

	result, err := schema.Validate(document)
	if err != nil {
		multiError.Append("validateInput", err)
		return multiError
	}

	if !result.Valid() {
		for _, desc := range result.Errors() {
			multiError.Append(desc.Field(), fmt.Errorf("value '%v' %v", desc.Value(), desc.Description()))
		}
	}

	if multiError.HasError() {
		return multiError
	}

	return nil
}
