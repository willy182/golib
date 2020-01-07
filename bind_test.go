package golib

import (
	"net/http"
	"net/url"
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
		{
			name:       "Testcase #4: Positive",
			queryParam: "orderId=0001&soNumber=2019:08:28T00:00:00+07:00",
			target:     new(exampleParam),
			wantResult: exampleParam{
				OrderID: "0001", SoNumber: "2019:08:28T00:00:00+07:00",
			},
			wantEqual: true,
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

func TestParseFromQueryParam(t *testing.T) {
	type Embed struct {
		Page   int    `json:"page"`
		Offset int    `json:"-"`
		Sort   string `json:"sort,omitempty" default:"desc" lower:"true"`
	}
	type params struct {
		Embed
		IsActive bool    `json:"isActive"`
		Ptr      *string `json:"ptr"`
	}

	t.Run("Testcase #1: Positive", func(t *testing.T) {
		urlVal, err := url.ParseQuery("page=1&ptr=val&isActive=true")
		assert.NoError(t, err)

		var p params
		err = ParseFromQueryParam(urlVal, &p)
		assert.NoError(t, err)
		assert.Equal(t, p.Page, 1)
		assert.Equal(t, *p.Ptr, "val")
		assert.Equal(t, p.IsActive, true)
	})
	t.Run("Testcase #2: Negative, invalid data type (string to int in struct)", func(t *testing.T) {
		urlVal, err := url.ParseQuery("page=undefined")
		assert.NoError(t, err)

		var p params
		err = ParseFromQueryParam(urlVal, &p)
		assert.Error(t, err)
	})
	t.Run("Testcase #3: Negative, invalid data type (not boolean)", func(t *testing.T) {
		urlVal, err := url.ParseQuery("isActive=terue")
		assert.NoError(t, err)

		var p params
		err = ParseFromQueryParam(urlVal, &p)
		assert.Error(t, err)
	})
	t.Run("Testcase #4: Negative, invalid target type (not pointer)", func(t *testing.T) {
		urlVal, err := url.ParseQuery("isActive=true")
		assert.NoError(t, err)

		var p params
		err = ParseFromQueryParam(urlVal, p)
		assert.Error(t, err)
	})
}
