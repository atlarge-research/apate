package cluster

// NodeInfo contains all information used for creating an equivalent kubernetes node
type NodeInfo struct {
	NodeType string
	Role     string
	Name     string
	Version  string
}

// NewNode create a new NodeInfo
func NewNode(nodeType string, role string, name string, version string) NodeInfo {
	return NodeInfo{
		NodeType: nodeType,
		Role:     role,
		Name:     name,
		Version:  version,
	}
}
