package golib

import (
	"fmt"
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
