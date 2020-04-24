package normalise

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTimestampSecondsPostfix(t *testing.T) {
	time, err := desugarTimestamp("42s")
	assert.NoError(t, err)
	assert.Equal(t, time, 42000)
}

func TestTimestampMilliSecondsPostfix(t *testing.T) {
	time, err := desugarTimestamp("42ms")
	assert.NoError(t, err)
	assert.Equal(t, time, 42)
}

func TestTimestampMinutePostfix(t *testing.T) {
	time, err := desugarTimestamp("42m")
	assert.NoError(t, err)
	assert.Equal(t, time, 42*60*1000)
}

func TestTimestampHourPostfix(t *testing.T) {
	time, err := desugarTimestamp("42h")
	assert.NoError(t, err)
	assert.Equal(t, time, 42*3600*1000)
}

func TestTimestampNoPostfix(t *testing.T) {
	time, err := desugarTimestamp("42")
	assert.NoError(t, err)
	assert.Equal(t, time, 42000)
}

func TestTimestampInvalidPostfixSecond(t *testing.T) {
	_, err := desugarTimestamp("42ns")
	assert.Error(t, err)
}

func TestTimestampInvalidPostfixInvalid(t *testing.T) {
	_, err := desugarTimestamp("42db")
	assert.Error(t, err)
}

func TestTimestampSecondsPostfixInvalid(t *testing.T) {
	_, err := desugarTimestamp("a42s")
	assert.Error(t, err)
}

func TestTimestampMilliSecondsPostfixInvalid(t *testing.T) {
	_, err := desugarTimestamp("a42ms")
	assert.Error(t, err)
}

func TestTimestampMinutePostfixInvalid(t *testing.T) {
	_, err := desugarTimestamp("a42m")
	assert.Error(t, err)
}

func TestTimestampHourPostfixInvalid(t *testing.T) {
	_, err := desugarTimestamp("a42h")
	assert.Error(t, err)
}

func TestTimestampNoPostfixInvalid(t *testing.T) {
	_, err := desugarTimestamp("a42")
	assert.Error(t, err)
}
