package throw

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCause(t *testing.T) {
	err := Exception("test")

	assert.EqualValues(t, err.Cause(), errors.New("test"))
}
