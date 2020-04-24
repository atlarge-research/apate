package service

// ConnectionInfo contains all information required for connecting to a services
type ConnectionInfo struct {
	address string
	port    int
	tls     bool
}

// NewConnectionInfo creates new connection information struct
func NewConnectionInfo(address string, port int, tls bool) *ConnectionInfo {
	return &ConnectionInfo{
		address: address,
		port:    port,
		tls:     tls,
	}
}
