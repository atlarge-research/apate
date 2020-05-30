package provider

import (
	"context"
	"log"
	"time"

	corev1 "k8s.io/api/core/v1"

	podconfigv1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/podconfiguration/v1"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	stats "k8s.io/kubernetes/pkg/kubelet/apis/stats/v1alpha1"
)

type resources struct {
	cpu              uint64
	memory           uint64
	ephemeralStorage uint64
}

// Stats is a simple wrapper for statistics fields
type Stats struct {
	startTime         metav1.Time
	podStats          *[]stats.PodStats
	podTotalResources *resources
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
		Pods: *p.Stats.podStats,
	}, nil
}

// Node statistics
func (p *Provider) getNodeStats() stats.NodeStats {
	return stats.NodeStats{
		NodeName:         p.NodeInfo.Name,
		SystemContainers: []stats.ContainerStats{},
		StartTime:        p.Stats.startTime,
		CPU:              p.cpuStats(),
		Memory:           p.memoryStats(),
		Fs:               p.filesystemStats(),
	}
}

func (p *Provider) cpuStats() *stats.CPUStats {
	zero := uint64(0)
	usage := uint64(0)

	for _, pod := range *p.Stats.podStats {
		if pod.CPU != nil && pod.CPU.UsageNanoCores != nil {
			usage += *pod.CPU.UsageNanoCores
		}
	}

	return &stats.CPUStats{
		Time:                 p.now(),
		UsageNanoCores:       &usage,
		UsageCoreNanoSeconds: &zero,
	}
}

func (p *Provider) memoryStats() *stats.MemoryStats {
	zero := uint64(0)
	usage := uint64(0)

	for _, pod := range *p.Stats.podStats {
		if pod.Memory != nil && pod.Memory.UsageBytes != nil {
			usage += *pod.Memory.UsageBytes
		}
	}

	available := uint64(p.Resources.Memory) - usage

	return &stats.MemoryStats{
		Time:            p.now(),
		AvailableBytes:  &available,
		UsageBytes:      &usage,
		WorkingSetBytes: &zero,
		RSSBytes:        &zero,
		PageFaults:      &zero,
		MajorPageFaults: &zero,
	}
}

func (p *Provider) filesystemStats() *stats.FsStats {
	zero := uint64(0)
	capacity := uint64(p.Resources.EphemeralStorage)
	usage := uint64(0)

	for _, pod := range *p.Stats.podStats {
		if pod.EphemeralStorage != nil && pod.EphemeralStorage.UsedBytes != nil {
			usage += *pod.EphemeralStorage.UsedBytes
		}
	}

	free := capacity - usage
	return &stats.FsStats{
		Time:           p.now(),
		AvailableBytes: &free,
		CapacityBytes:  &capacity,
		UsedBytes:      &usage,
		InodesFree:     &zero,
		Inodes:         &zero,
		InodesUsed:     &zero,
	}
}

func (p *Provider) updateAggregatePodStats() {
	var statistics []stats.PodStats
	var totalResources = resources{}

	for _, pod := range p.Pods.GetAllPods() {
		podStats := *p.getPodStats(pod)
		totalResources.cpu += *podStats.CPU.UsageNanoCores
		totalResources.memory += *podStats.Memory.UsageBytes
		totalResources.ephemeralStorage += *podStats.EphemeralStorage.UsedBytes
		statistics = append(statistics, podStats)
	}

	p.Stats.podStats = &statistics
	p.Stats.podTotalResources = &totalResources
}

func (p *Provider) getPodStats(pod *corev1.Pod) *stats.PodStats {
	if label, ok := pod.Labels[podconfigv1.PodConfigurationLabel]; ok {
		unconvertedStats, err := (*p.Store).GetPodFlag(label, events.PodResources)
		if err != nil {
			log.Printf("error while retrieving pod flag for resources: %v\n", err)
			return addPodSpecificStats(pod, &stats.PodStats{})
		}

		statistics, ok := unconvertedStats.(stats.PodStats)
		if !ok {
			log.Printf("unable to convert '%v' to PodStats\n", unconvertedStats)
			return addPodSpecificStats(pod, &stats.PodStats{})
		}

		return addPodSpecificStats(pod, &statistics)
	}

	return addPodSpecificStats(pod, &stats.PodStats{})
}

func addPodSpecificStats(pod *corev1.Pod, statistics *stats.PodStats) *stats.PodStats {
	statistics.PodRef = stats.PodReference{
		Name:      pod.Name,
		Namespace: pod.Namespace,
		UID:       string(pod.UID),
	}

	if pod.Status.StartTime != nil {
		statistics.StartTime = *pod.Status.StartTime
	}

	return statistics
}
