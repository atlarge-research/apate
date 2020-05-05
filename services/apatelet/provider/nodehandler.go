package provider

import (
	"context"
	"github.com/virtual-kubelet/virtual-kubelet/node"
	corev1 "k8s.io/api/core/v1"
)

type NodeHandler struct {
}

// NewNodeProvider creates a new node provider
func NewNodeProvider() node.NodeProvider {
	return &NodeHandler{}
}

func (n NodeHandler) Ping(ctx context.Context) error {
	return ctx.Err()
}

func (n NodeHandler) NotifyNodeStatus(ctx context.Context, cb func(*corev1.Node)) {
	// TODO
}
