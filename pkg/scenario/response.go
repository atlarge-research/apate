package scenario

// Response specifies the response type
type Response int

const (
	// ResponseUnset No response
	ResponseUnset Response = iota
	// ResponseNormal will respond normally to the request
	ResponseNormal
	// ResponseTimeout will timeout any request
	ResponseTimeout
	// ResponseError will error on any request
	ResponseError
)
