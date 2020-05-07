package provider

type wError string

func (e wError) Error() string {
	return string(e)
}

func (e wError) Cause() error {
	return e
}
