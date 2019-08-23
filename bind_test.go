package golib

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBindQueryParam(t *testing.T) {
	type exampleParam struct {
		OrderID  string `json:"orderId"`
		SoNumber string `json:"soNumber"`
	}
	tests := []struct {
		name, queryParam     string
		target, wantResult   interface{}
		wantError, wantEqual bool
	}{
		{
			name:       "Testcase #1: Positive",
			queryParam: "orderId=0001&soNumber=SO001",
			target:     new(exampleParam),
			wantResult: exampleParam{
				OrderID: "0001", SoNumber: "SO001",
			},
			wantEqual: true,
		},
		{
			name:       "Testcase #2: Negative, query param not match with json tag in struct model target",
			queryParam: "orderID=0001&soNumber=SO001",
			target:     new(exampleParam),
			wantResult: exampleParam{
				OrderID: "0001", SoNumber: "SO001",
			},
			wantEqual: false,
		},
		{
			name:       "Testcase #3: Negative, target is not pointer",
			queryParam: "orderId=0001&soNumber=SO001",
			target:     exampleParam{},
			wantResult: exampleParam{
				OrderID: "0001", SoNumber: "SO001",
			},
			wantError: true, wantEqual: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/api/order?"+tt.queryParam, nil)
			assert.NoError(t, err)

			err = BindQueryParam(req.URL, tt.target)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			targetValue := reflect.ValueOf(tt.target)
			if targetValue.Kind() == reflect.Ptr {
				targetValue = targetValue.Elem()
			}
			if tt.wantEqual {
				assert.Equal(t, tt.wantResult, targetValue.Interface())
			} else {
				assert.NotEqual(t, tt.wantResult, targetValue.Interface())
			}
		})
	}
}
