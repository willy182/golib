package golib

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiError(t *testing.T) {
	multiError := NewMultiError()

	multiError.Append("err1", fmt.Errorf("error 1"))
	assert.True(t, multiError.HasError())
	errMap := multiError.ToMap()
	assert.Equal(t, 1, len(errMap))
	assert.Equal(t, "err1: error 1", multiError.Error())
}

func TestAppendMultiError(t *testing.T) {
	multiError1 := NewMultiError()
	multiError1.Append("err1", fmt.Errorf("error 1"))

	multiError2 := NewMultiError()
	multiError2.Append("err2", fmt.Errorf("error 2"))

	multiErrorAll := NewMultiError()
	multiErrorAll.Append("err1", fmt.Errorf("error 1"))
	multiErrorAll.Append("err2", fmt.Errorf("error 2"))

	type args struct {
		map1 *MultiError
		map2 *MultiError
	}
	tests := []struct {
		name string
		args args
		want *MultiError
	}{
		{
			name: "test case 1",
			args: args{
				map1: multiError1,
				map2: multiError2,
			},
			want: multiErrorAll,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AppendMultiError(tt.args.map1, tt.args.map2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppendMultiError() = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("APPEND EXISTING KEY", func(t *testing.T) {
		m := NewMultiError()
		m.errs["test"] = "test"
		m.Append("test", errors.New("testing"))
		assert.Equal(t, "test; testing", m.errs["test"])
	})
}

func TestMultiErrorNotNill(t *testing.T) {
	type args struct {
		multiError *MultiError
	}
	tests := []struct {
		name string
		args args
		want *MultiError
	}{
		{
			name: "test 1",
			args: args{},
			want: NewMultiError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MultiErrorNotNill(tt.args.multiError); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MultiErrorNotNill() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMultiError_Clear(t *testing.T) {
	t.Run("CLEAR MULTI ERROR", func(t *testing.T) {
		m := NewMultiError()
		m.Clear()
		assert.Equal(t, 0, len(m.errs))
	})
}
