package kubernetes

import (
	"strings"

	"github.com/pkg/errors"
)

// NodeInfo contains all information used for creating an equivalent kubernetes node
type NodeInfo struct {
	NodeType, Role, Name, Version, Namespace, Selector string

	MetricsPort int
}

// NewNodeInfo creates a new NodeInfo
func NewNodeInfo(nodeType string, role string, name string, version string, selector string, metricsPort int) (NodeInfo, error) {
	selectorParts := strings.Split(selector, "/")

	if len(selectorParts) != 2 {
		return NodeInfo{}, errors.Errorf("invalid selector %s", selector)
	}

	return NodeInfo{
		NodeType:    nodeType,
		Role:        role,
		Name:        name,
		Version:     version,
		MetricsPort: metricsPort,
		Namespace:   selectorParts[0],
		Selector:    selectorParts[1],
	}, nil
}
