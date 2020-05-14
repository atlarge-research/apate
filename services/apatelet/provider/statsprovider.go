package provider

import (
	"context"
	"time"

	v1 "k8s.io/api/core/v1"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	stats "k8s.io/kubernetes/pkg/kubelet/apis/stats/v1alpha1"
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

		// TODO: Do we want these? They have 0 added value afaik
		Network: nil,
		Fs:      nil,
		Runtime: nil,
		Rlimit:  nil,
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

func (p *Provider) getAggregatePodStats() []stats.PodStats {
	var statistics []stats.PodStats

	for _, pod := range p.pods.GetAllPods() {
		statistics = append(statistics, p.getPodStats(pod))
	}

	return statistics
}

func (p *Provider) getPodStats(pod *v1.Pod) stats.PodStats {
	for k, label := range pod.Labels {
		if k == "apate" { //TODO: Const or something
			event := events.PodResources

			unconvertedStats, err := (*p.store).GetPodFlag(label, event)
			statistics := unconvertedStats.(stats.PodStats)

			if err != nil {
				return stats.PodStats{PodRef: stats.PodReference{Name: pod.Name, Namespace: pod.Namespace, UID: string(pod.UID)}}
			}

			addPodSpecificStats(pod, &statistics)
			return statistics
		}
	}

	return stats.PodStats{PodRef: stats.PodReference{Name: pod.Name, Namespace: pod.Namespace, UID: string(pod.UID)}}
}

func addPodSpecificStats(pod *v1.Pod, statistics *stats.PodStats) {
	statistics.PodRef = stats.PodReference{
		Name:      pod.Name,
		Namespace: pod.Namespace,
		UID:       string(pod.UID),
	}

	if pod.Status.StartTime != nil {
		statistics.StartTime = *pod.Status.StartTime
	}
}
