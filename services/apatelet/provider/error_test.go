package provider

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWError_Cause(t *testing.T) {
	erra := errors.New("yeet")
	errb := wError(erra.Error())

	assert.EqualError(t, erra, errb.Cause().Error())
}
