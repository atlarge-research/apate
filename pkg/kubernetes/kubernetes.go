package kubernetes

// NodeInfo contains all information used for creating an equivalent kubernetes node
type NodeInfo struct {
	NodeType, Role, Name, Version, Selector string
	MetricsPort                             int
}

// NewNodeInfo creates a new NodeInfo
func NewNodeInfo(nodeType string, role string, name string, version string, selector string, metricsPort int) NodeInfo {
	return NodeInfo{
		NodeType:    nodeType,
		Role:        role,
		Name:        name,
		Version:     version,
		MetricsPort: metricsPort,
		Selector:    selector,
	}
}
