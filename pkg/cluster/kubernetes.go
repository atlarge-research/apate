package cluster

// NodeInfo contains all information used for creating an equivalent kubernetes node
type NodeInfo struct {
	NodeType, Role, Name, Version, Spec string
	MetricsPort                         int
}

// NewNodeInfo creates a new NodeInfo
func NewNodeInfo(nodeType, role, name, version, spec string, metricsPort int) NodeInfo {
	return NodeInfo{
		NodeType:    nodeType,
		Role:        role,
		Name:        name,
		Version:     version,
		Spec:        spec,
		MetricsPort: metricsPort,
	}
}
