package golib

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// writer implement http.ResponseWriter
type writer struct {
	http.ResponseWriter
}

func (w *writer) Header() http.Header {
	return http.Header{}
}

func (w *writer) WriteHeader(code int) {
}

func (w *writer) Write(b []byte) (int, error) {
	return len(b), nil
}

type ExampleModel struct {
	OrderID string `json:"orderId"`
}

func TestNewHTTPResponseV2(t *testing.T) {
	multiError := NewMultiError()
	multiError.Append("test", fmt.Errorf("error test"))
	type args struct {
		code    int
		message string
		params  []interface{}
	}
	tests := []struct {
		name string
		args args
		want *httpResponse
	}{
		{
			name: "Testcase #1: Response data list (include meta)",
			args: args{
				code:    http.StatusOK,
				message: "Fetch all data",
				params: []interface{}{
					[]ExampleModel{{OrderID: "061499700032"}, {OrderID: "061499700033"}},
					Meta{Page: 1, Limit: 10, TotalPages: 10, TotalRecords: 100},
				},
			},
			want: &httpResponse{
				Success: true,
				Code:    200,
				Message: "Fetch all data",
				Meta:    Meta{Page: 1, Limit: 10, TotalPages: 10, TotalRecords: 100},
				Data:    []ExampleModel{{OrderID: "061499700032"}, {OrderID: "061499700033"}},
			},
		},
		{
			name: "Testcase #2: Response data detail",
			args: args{
				code:    http.StatusOK,
				message: "Get detail data",
				params: []interface{}{
					ExampleModel{OrderID: "061499700032"},
				},
			},
			want: &httpResponse{
				Success: true,
				Code:    200,
				Message: "Get detail data",
				Data:    ExampleModel{OrderID: "061499700032"},
			},
		},
		{
			name: "Testcase #3: Response only message (without data)",
			args: args{
				code:    http.StatusOK,
				message: "list data empty",
			},
			want: &httpResponse{
				Success: true,
				Code:    200,
				Message: "list data empty",
			},
		},
		{
			name: "Testcase #4: Response failed (ex: Bad Request)",
			args: args{
				code:    http.StatusBadRequest,
				message: "id cannot be empty",
				params:  []interface{}{multiError},
			},
			want: &httpResponse{
				Success: false,
				Code:    400,
				Message: "id cannot be empty",
				Errors:  map[string]string{"test": "error test"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewHTTPResponseV2(tt.args.code, tt.args.message, tt.args.params...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\x1b[31;1mNewHTTPResponseV2() = %v, \nwant => %v\x1b[0m", got, tt.want)
			}
		})
	}
}

func TestHTTPResponse_JSON(t *testing.T) {
	resp := NewHTTPResponseV2(200, "success")
	w := new(writer)
	assert.NoError(t, resp.JSON(w))
}

func TestHTTPResponse_XML(t *testing.T) {
	resp := NewHTTPResponseV2(200, "success")
	w := new(writer)
	assert.NoError(t, resp.XML(w))
}
