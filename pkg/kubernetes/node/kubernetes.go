// Package node contains some utilities to describe Kubernetes nodes
package node

import (
	"strings"

	"github.com/pkg/errors"
)

// Info contains all information used for creating an equivalent kubernetes node
type Info struct {
	NodeType, Role, Name, Version, Namespace, Selector string

	MetricsPort int
}

// NewInfo creates a new Info
func NewInfo(nodeType string, role string, name string, version string, selector string, metricsPort int) (Info, error) {
	selectorParts := strings.Split(selector, "/")

	if len(selectorParts) != 2 {
		return Info{}, errors.Errorf("invalid selector %s", selector)
	}

	return Info{
		NodeType:    nodeType,
		Role:        role,
		Name:        name,
		Version:     version,
		MetricsPort: metricsPort,
		Namespace:   selectorParts[0],
		Selector:    selectorParts[1],
	}, nil
}
