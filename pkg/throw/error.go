// Package throw provides error helpers specifically for const sentinel errors
package throw

import "errors"

// Exception is a type alias of string with Error implemented on it
type Exception string

func (e Exception) Error() string {
	return string(e)
}

// Cause returns the inner error
func (e Exception) Cause() error {
	return errors.New(string(e))
}
