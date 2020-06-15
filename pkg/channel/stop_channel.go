// Package channel contains utilities related to channels
package channel

import "sync"

// StopChannel contains a channel that is safe to close twice
type StopChannel struct {
	ch   chan struct{}
	once sync.Once
}

// NewStopChannel creates a new stop channel
func NewStopChannel() *StopChannel {
	return &StopChannel{ch: make(chan struct{})}
}

// Close closes the internal channel. Can be executed an arbitrary amount of times.
func (sc *StopChannel) Close() {
	sc.once.Do(func() {
		close(sc.ch)
	})
}

// GetChannel gets the read channel.
func (sc *StopChannel) GetChannel() <-chan struct{} {
	return sc.ch
}
