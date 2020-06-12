package channel

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStopChannel(t *testing.T) {
	sc := NewStopChannel()

	ch := sc.GetChannel()
	select {
	case <-ch:
		assert.Fail(t, "channel either contains a value or is closed!")
	default:
		//
	}

	go func() {
		sc.ch <- struct{}{}
	}()
	select {
	case x, ok := <-ch:
		assert.True(t, ok)
		assert.EqualValues(t, struct{}{}, x)
	case <-time.After(100 * time.Millisecond): // It should be written by then
		assert.Fail(t, "channel should contain a value!")
	}

	sc.Close()
	sc.Close()
	sc.Close()
	sc.Close()
	sc.Close()

	select {
	case x, ok := <-ch:
		assert.False(t, ok)
		assert.EqualValues(t, struct{}{}, x)
	default:
		assert.Fail(t, "channel should contain a value!")
	}
}
