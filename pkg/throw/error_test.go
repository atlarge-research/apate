package throw

import (
	"testing"
)

func TestCause(t *testing.T) {
	err := Exception("test")
	if err.Cause() != err {
		t.Fail()
	}
}
