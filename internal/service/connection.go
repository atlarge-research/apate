package service

// ConnectionInfo contains all information required for connecting to a services
type ConnectionInfo struct {
	Address string
	Port    int
}

// NewConnectionInfo creates new connection information struct
func NewConnectionInfo(address string, port int) *ConnectionInfo {
	return &ConnectionInfo{
		Address: address,
		Port:    port,
	}
}
