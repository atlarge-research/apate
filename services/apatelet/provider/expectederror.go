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

// An emulation error is an error occurring because the system is emulating an error. This type of error is completely expected.
// This error should never be wrapped
type emulationError string

func (e emulationError) Error() string {
	return string(e)
}

func (e emulationError) Expected() bool {
	return true
}
