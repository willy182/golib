package golib

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"reflect"
)

// HTTPResponse abstract interface
type HTTPResponse interface {
	JSON(w http.ResponseWriter) error
	XML(w http.ResponseWriter) error
}

type (
	// ResponseV2 model
	ResponseV2 struct {
		Success bool        `json:"success"`
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Meta    interface{} `json:"meta,omitempty"`
		Data    interface{} `json:"data,omitempty"`
		Errors  interface{} `json:"errors,omitempty"`
	}

	// Meta model
	Meta struct {
		Page         int `json:"page"`
		Limit        int `json:"limit"`
		TotalRecords int `json:"totalRecords"`
		TotalPages   int `json:"totalPages"`
	}
)

// NewHTTPResponseV2 for create common response, data must in first params and meta in second params
func NewHTTPResponseV2(code int, message string, params ...interface{}) HTTPResponse {
	commonResponse := new(ResponseV2)

	for _, param := range params {
		// get value param if type is pointer
		refValue := reflect.ValueOf(param)
		if refValue.Kind() == reflect.Ptr {
			refValue = refValue.Elem()
		}
		param = refValue.Interface()

		switch val := param.(type) {
		case Meta:
			commonResponse.Meta = val
		case MultiError:
			commonResponse.Errors = val.ToMap()
		default:
			commonResponse.Data = param
		}
	}

	if code < http.StatusBadRequest && message != ErrorDataNotFound {
		commonResponse.Success = true
	} else {
		commonResponse.Success = false
	}
	commonResponse.Code = code
	commonResponse.Message = message
	return commonResponse
}

// JSON for set http JSON response (Content-Type: application/json) with parameter is http response writer
func (resp *ResponseV2) JSON(w http.ResponseWriter) error {
	if resp.Data == nil {
		resp.Data = struct{}{}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Code)
	return json.NewEncoder(w).Encode(resp)
}

// XML for set http XML response (Content-Type: application/xml)
func (resp *ResponseV2) XML(w http.ResponseWriter) error {
	if resp.Data == nil {
		resp.Data = struct{}{}
	}
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(resp.Code)
	return xml.NewEncoder(w).Encode(resp)
}
