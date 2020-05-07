package provider

import (
	"context"

	corev1 "k8s.io/api/core/v1"
)

// Ping TODO
func (p *Provider) Ping(ctx context.Context) error {
	return ctx.Err()
}

// NotifyNodeStatus TODO
func (p *Provider) NotifyNodeStatus(_ context.Context, _ func(*corev1.Node)) {
	// TODO
}
