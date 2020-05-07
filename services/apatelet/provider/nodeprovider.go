package provider

import (
	"context"

	"github.com/virtual-kubelet/virtual-kubelet/node"
	corev1 "k8s.io/api/core/v1"
)

// TODO: Handle node provider in virtual kubelet

// NodeProvider implements the virtual kubelet NodeProvider interface to be able to alter responses to ping
type NodeProvider struct {
}

// NewNodeProvider creates a new node provider
func NewNodeProvider() node.NodeProvider {
	return &NodeProvider{}
}

// Ping TODO
func (n NodeProvider) Ping(ctx context.Context) error {
	return ctx.Err()
}

// NotifyNodeStatus TODO
func (n NodeProvider) NotifyNodeStatus(ctx context.Context, cb func(*corev1.Node)) {
	// TODO
}
