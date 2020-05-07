package throw

type Exception string

func (e Exception) Error() string {
	return string(e)
}

func (e Exception) Cause() error {
	return e
}
