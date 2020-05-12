package provider

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	stats "k8s.io/kubernetes/pkg/kubelet/apis/stats/v1alpha1"
)

// GetStatsSummary should return a node level statistic report
func (p *Provider) GetStatsSummary(context.Context) (*stats.Summary, error) {
	// TODO: Implement
	mem := uint64(421251256)
	return &stats.Summary{Node: stats.NodeStats{Memory: &stats.MemoryStats{Time: metav1.NewTime(time.Now()), AvailableBytes: &mem}}}, nil
}
