// Package e2e contains end to end tests for Apate.
package e2e

import (
	"bytes"
	"io"
	"os"
)

type capturer struct {
	contents  chan string
	w         *os.File
	oldStdout *os.File
}

func capture() capturer {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	c := make(chan string)

	go func(ch chan<- string) {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		ch <- buf.String()
	}(c)

	return capturer{
		c,
		w,
		old,
	}
}

func (c capturer) stop() string {
	_ = c.w.Close()
	os.Stdout = c.oldStdout
	return <-c.contents
}
