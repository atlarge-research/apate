package service

// ConnectionInfo contains all information required for connecting to a services
type ConnectionInfo struct {
	Address string
	Port    int
	Tls     bool
}

// NewConnectionInfo creates new connection information struct
func NewConnectionInfo(address string, port int, tls bool) *ConnectionInfo {
	return &ConnectionInfo{
		Address: address,
		Port:    port,
		Tls:     tls,
	}
}
