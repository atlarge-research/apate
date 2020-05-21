package provider

// ExpectedError is an error which is thrown during emulation.
type ExpectedError interface {
	Expected() bool
}

// IsExpected returns true if err is expected.
func IsExpected(err error) bool {
	te, ok := err.(ExpectedError)
	return ok && te.Expected()
}
