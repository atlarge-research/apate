package provider

import (
	"context"
	"log"
	"time"

	v1 "k8s.io/api/core/v1"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	stats "k8s.io/kubernetes/pkg/kubelet/apis/stats/v1alpha1"
)

const (
	podConfigurationLabel = "apate" // TODO: Change to constant somewhere else
)

// Stats is a simple wrapper for statistics fields
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
	pods := p.getAggregatePodStats()

	return &stats.Summary{
		Node: p.getNodeStats(&pods),
		Pods: pods,
	}, nil
}

// Node statistics
func (p *Provider) getNodeStats(pods *[]stats.PodStats) stats.NodeStats {
	return stats.NodeStats{
		NodeName:         p.nodeInfo.Name,
		SystemContainers: []stats.ContainerStats{},
		StartTime:        p.stats.startTime,
		CPU:              p.cpuStats(pods),
		Memory:           p.memoryStats(pods),
		Fs:               p.filesystemStats(pods),
	}
}

func (p *Provider) cpuStats(pods *[]stats.PodStats) *stats.CPUStats {
	zero := uint64(0)
	usage := uint64(0)

	for _, pod := range *pods {
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

func (p *Provider) memoryStats(pods *[]stats.PodStats) *stats.MemoryStats {
	zero := uint64(0)
	usage := uint64(0)

	for _, pod := range *pods {
		if pod.Memory != nil && pod.Memory.UsageBytes != nil {
			usage += *pod.Memory.UsageBytes
		}
	}

	available := uint64(p.resources.Memory) - usage

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

func (p *Provider) filesystemStats(pods *[]stats.PodStats) *stats.FsStats {
	zero := uint64(0)
	capacity := uint64(p.resources.EphemeralStorage)
	usage := uint64(0)

	for _, pod := range *pods {
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

func (p *Provider) getAggregatePodStats() []stats.PodStats {
	var statistics []stats.PodStats

	for _, pod := range p.pods.GetAllPods() {
		statistics = append(statistics, *p.getPodStats(pod))
	}

	return statistics
}

func (p *Provider) getPodStats(pod *v1.Pod) *stats.PodStats {
	for k, label := range pod.Labels {
		if k == podConfigurationLabel {
			unconvertedStats, err := (*p.store).GetPodFlag(label, events.PodResources)
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
	}

	return addPodSpecificStats(pod, &stats.PodStats{})
}

func addPodSpecificStats(pod *v1.Pod, statistics *stats.PodStats) *stats.PodStats {
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