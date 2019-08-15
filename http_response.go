package golib

import (
	"encoding/xml"
	"net/http"
	"reflect"

	"github.com/labstack/echo"
)

// HTTPResponse abstract interface
type HTTPResponse interface {
	JSON(c echo.Context) error
	XML(c echo.Context) error
}

type (
	// httpResponse model
	httpResponse struct {
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

	// NavResponse xml model for nav response
	NavResponse struct {
		XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ ReturnValue"`

		Message string `xml:"Message"`
	}
)

// EchoHTTPResponseV2 for create common response, data must in first params and meta in second params
func EchoHTTPResponseV2(code int, message string, params ...interface{}) HTTPResponse {
	commonResponse := new(httpResponse)

	for _, param := range params {
		// get value param if type is pointer
		refValue := reflect.ValueOf(param)
		if refValue.Kind() == reflect.Ptr {
			refValue = refValue.Elem()
		}
		param = refValue.Interface()

		switch param.(type) {
		case Meta:
			commonResponse.Meta = param
		case MultiError:
			multiError := param.(MultiError)
			commonResponse.Errors = multiError.ToMap()
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

// JSON for set http JSON response (Content-Type: application/json)
func (resp *httpResponse) JSON(c echo.Context) error {
	return c.JSON(resp.Code, resp)
}

// XML for set http XML response (Content-Type: application/xml)
func (resp *httpResponse) XML(c echo.Context) error {
	return c.XML(resp.Code, resp)
}
