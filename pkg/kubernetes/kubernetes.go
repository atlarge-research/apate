package kubernetes

import (
	"strings"

	"github.com/pkg/errors"
)

// NodeInfo contains all information used for creating an equivalent kubernetes node
type NodeInfo struct {
	NodeType, Role, Name, Version, Namespace, Label string

	MetricsPort int
}

// NewNodeInfo creates a new NodeInfo
func NewNodeInfo(nodeType string, role string, name string, version string, label string, metricsPort int) (NodeInfo, error) {
	labelParts := strings.Split(label, "/")

	if len(labelParts) != 2 {
		return NodeInfo{}, errors.Errorf("invalid label %s", label)
	}

	return NodeInfo{
		NodeType:    nodeType,
		Role:        role,
		Name:        name,
		Version:     version,
		MetricsPort: metricsPort,
		Namespace:   labelParts[0],
		Label:       labelParts[1],
	}, nil
}
