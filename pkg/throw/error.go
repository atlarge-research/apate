// Package throw provides error helpers specifically for const sentinel errors
package throw

// TODO: Replace with github.com/pkg/errors

// Exception is a wrapped error
type Exception struct {
	inner   error
	message string
}

// Error returns the error message
func (e Exception) Error() string {
	return e.message
}

// Cause returns the inner error
func (e Exception) Cause() error {
	return e.inner
}

// NewException creates a new exception
func NewException(inner error, message string) Exception {
	return Exception{inner, message}
}

// ConstException is a valid error type which is constant
type ConstException string

// Error returns the error message
func (e ConstException) Error() string {
	return string(e)
}
