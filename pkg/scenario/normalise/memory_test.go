package normalise

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDesugarMemoryLower(t *testing.T) {
	time, err := desugarMemory("42b")
	assert.NoError(t, err)
	assert.Equal(t, time, int64(42))
}

func TestDesugarMemoryBytes(t *testing.T) {
	time, err := desugarMemory("42B")
	assert.NoError(t, err)
	assert.Equal(t, time, int64(42))
}

func TestDesugarMemoryNoExt(t *testing.T) {
	time, err := desugarMemory("42")
	assert.NoError(t, err)
	assert.Equal(t, time, int64(42))
}

func TestDesugarMemoryKB(t *testing.T) {
	time, err := desugarMemory("42KB")
	assert.NoError(t, err)
	assert.Equal(t, time, int64(42000))
}

func TestDesugarMemoryMB(t *testing.T) {
	time, err := desugarMemory("42MB")
	assert.NoError(t, err)
	assert.Equal(t, time, int64(42000000))
}

func TestDesugarMemoryGB(t *testing.T) {
	time, err := desugarMemory("42GB")
	assert.NoError(t, err)
	assert.Equal(t, time, int64(42000000000))
}

func TestDesugarMemoryKiB(t *testing.T) {
	time, err := desugarMemory("42KiB")
	assert.NoError(t, err)
	assert.Equal(t, time, int64(42 * 1024))
}

func TestDesugarMemoryMiB(t *testing.T) {
	time, err := desugarMemory("42MiB")
	assert.NoError(t, err)
	assert.Equal(t, time, int64(42 * 1024 * 1024))
}

func TestDesugarMemoryGiB(t *testing.T) {
	time, err := desugarMemory("42GiB")
	assert.NoError(t, err)
	assert.Equal(t, time, int64(42 * 1024 * 1024 * 1024))
}

func TestDesugarMemoryK(t *testing.T) {
	time, err := desugarMemory("42K")
	assert.NoError(t, err)
	assert.Equal(t, time, int64(42 * 1000))
}

func TestDesugarMemoryM(t *testing.T) {
	time, err := desugarMemory("42M")
	assert.NoError(t, err)
	assert.Equal(t, time, int64(42 * 1000 * 1000))
}

func TestDesugarMemoryG(t *testing.T) {
	time, err := desugarMemory("42G")
	assert.NoError(t, err)
	assert.Equal(t, time, int64(42 * 1000 * 1000 * 1000))
}


func TestDesugarMemoryBytesErr(t *testing.T) {
	_, err := desugarMemory("a42B")
	assert.Error(t, err)
}

func TestDesugarMemoryNoExtErr(t *testing.T) {
	_, err := desugarMemory("a42")
	assert.Error(t, err)
}

func TestDesugarMemoryKBErr(t *testing.T) {
	_, err := desugarMemory("a42KB")
	assert.Error(t, err)
}

func TestDesugarMemoryMBErr(t *testing.T) {
	_, err := desugarMemory("a42MB")
	assert.Error(t, err)
}

func TestDesugarMemoryGBErr(t *testing.T) {
	_, err := desugarMemory("a42GB")
	assert.Error(t, err)
}

func TestDesugarMemoryKiBErr(t *testing.T) {
	_, err := desugarMemory("a42KiB")
	assert.Error(t, err)
}

func TestDesugarMemoryMiBErr(t *testing.T) {
	_, err := desugarMemory("a42MiB")
	assert.Error(t, err)
}

func TestDesugarMemoryGiBErr(t *testing.T) {
	_, err := desugarMemory("a42GiB")
	assert.Error(t, err)
}

func TestDesugarMemoryKErr(t *testing.T) {
	_, err := desugarMemory("a42K")
	assert.Error(t, err)
}

func TestDesugarMemoryMErr(t *testing.T) {
	_, err := desugarMemory("a42M")
	assert.Error(t, err)
}

func TestDesugarMemoryGErr(t *testing.T) {
	_, err := desugarMemory("a42G")
	assert.Error(t, err)
}