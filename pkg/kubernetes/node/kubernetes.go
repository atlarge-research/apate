// Package node contains some utilities to describe Kubernetes nodes
package node

import (
	"strings"

	"github.com/pkg/errors"
)

// Info contains all information used for creating an equivalent kubernetes node
type Info struct {
	NodeType, Role, Name, Version, Namespace, Label string

	MetricsPort int
}

// NewInfo creates a new Info
func NewInfo(nodeType, role, name, version, label string) (Info, error) {
	labelParts := strings.Split(label, "/")

	if len(labelParts) != 2 {
		return Info{}, errors.Errorf("invalid label %s", label)
	}

	return Info{
		NodeType:  nodeType,
		Role:      role,
		Name:      name,
		Version:   version,
		Namespace: labelParts[0],
		Label:     labelParts[1],
	}, nil
}
