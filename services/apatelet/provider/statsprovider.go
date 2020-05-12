package provider

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	stats "k8s.io/kubernetes/pkg/kubelet/apis/stats/v1alpha1"
)

type Stats struct {
	startTime metav1.Time
}

// NewStats creates a new Stats instance
func NewStats() *Stats {
	return &Stats{startTime: metav1.NewTime(time.Now())}
}

func (p *Provider) now() metav1.Time {
	return metav1.NewTime(time.Now())
}

// GetStatsSummary should return a node level statistic report
func (p *Provider) GetStatsSummary(context.Context) (*stats.Summary, error) {
	return &stats.Summary{
		Node: p.getNodeStats(),
		Pods: p.getAggregatePodStats(),
	}, nil
}

// TODO: Decide which we want to implement

// Node statistics
func (p *Provider) getNodeStats() stats.NodeStats {
	return stats.NodeStats{
		NodeName:         p.nodeInfo.Name,
		SystemContainers: []stats.ContainerStats{},
		StartTime:        p.stats.startTime,
		CPU:              p.cpuStats(),
		Memory:           p.memoryStats(),
		Network:          nil,
		Fs:               nil,
		Runtime:          nil,
		Rlimit:           nil,
	}
}

func (p *Provider) cpuStats() *stats.CPUStats {
	cnt := uint64(532425124124)
	return &stats.CPUStats{
		Time:                 p.now(),
		UsageNanoCores:       &cnt,
		UsageCoreNanoSeconds: &cnt,
	}
}

func (p *Provider) memoryStats() *stats.MemoryStats {
	cnt := uint64(532425124124)
	return &stats.MemoryStats{
		Time:            p.now(),
		AvailableBytes:  &cnt,
		UsageBytes:      &cnt,
		WorkingSetBytes: &cnt,
		RSSBytes:        &cnt,
		PageFaults:      &cnt,
		MajorPageFaults: &cnt,
	}
}

// Pod statistics
func (p *Provider) getAggregatePodStats() []stats.PodStats {
	return []stats.PodStats{}
}

func (p *Provider) getPodStats() stats.PodStats {
	return stats.PodStats{
		PodRef: stats.PodReference{
			Name:      "",
			Namespace: "",
			UID:       "",
		},
		StartTime: metav1.NewTime(time.Now()),
		Containers: []stats.ContainerStats{
			{
				Name:               "",
				StartTime:          metav1.Time{},
				CPU:                nil,
				Memory:             nil,
				Accelerators:       nil,
				Rootfs:             nil,
				Logs:               nil,
				UserDefinedMetrics: nil,
			},
		},
		CPU:     nil,
		Memory:  nil,
		Network: nil,
		VolumeStats: []stats.VolumeStats{
			{
				FsStats: stats.FsStats{
					Time:           metav1.Time{},
					AvailableBytes: nil,
					CapacityBytes:  nil,
					UsedBytes:      nil,
					InodesFree:     nil,
					Inodes:         nil,
					InodesUsed:     nil,
				},
				Name:   "",
				PVCRef: nil,
			},
		},
		EphemeralStorage: nil,
	}
}
