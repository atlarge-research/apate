package cluster

// NodeInfo contains all information used for creating an equivalent kubernetes node
type NodeInfo struct {
	NodeType    string
	Role        string
	Name        string
	Version     string
	MetricsPort int
}

// NewNode create a new NodeInfo
func NewNode(nodeType, role, name, version string, metricsPort int) NodeInfo {
	return NodeInfo{
		NodeType:    nodeType,
		Role:        role,
		Name:        name,
		Version:     version,
		MetricsPort: metricsPort,
	}
}
