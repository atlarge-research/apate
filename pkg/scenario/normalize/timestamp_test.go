package normalize

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
	assert.Equal(t, time, 42 * 60 * 1000)
}

func TestTimestampHourPostfix(t *testing.T) {
	time, err := desugarTimestamp("42h")
	assert.NoError(t, err)
	assert.Equal(t, time, 42 * 3600 * 1000)
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

func TestTimestampInvalidPostfix(t *testing.T) {
	_, err := desugarTimestamp("42db")
	assert.Error(t, err)
}