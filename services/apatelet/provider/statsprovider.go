package provider

import (
	"errors"
	"log"
	"time"

	"github.com/finitum/node-cli/stats"

	corev1 "k8s.io/api/core/v1"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type resources struct {
	cpu              uint64
	memory           uint64
	ephemeralStorage uint64
}

// Stats is a simple wrapper for statistics fields
type Stats struct {
	statsSummary *stats.Summary
}

// NewStats creates a new Stats instance
func NewStats() *Stats {
	return &Stats{}
}

func (p *Provider) now() metav1.Time {
	return metav1.NewTime(time.Now())
}

// GetStatsSummary should return a node level statistic report
func (p *Provider) GetStatsSummary() (*stats.Summary, error) {
	if p.Stats.statsSummary == nil {
		return nil, errors.New("statsSummary is nil, please call updateStatsSummary first")
	}
	return p.Stats.statsSummary, nil
}

func (p *Provider) updateStatsSummary() {
	pods := p.getAggregatePodStats()

	p.Stats.statsSummary = &stats.Summary{
		Node: p.getNodeStats(pods),
		Pods: pods,
	}
}

// Node statistics
func (p *Provider) getNodeStats(pods []stats.PodStats) stats.NodeStats {
	aMem, uMem := p.memoryStats(pods)
	aEph, cEph, uEph := p.ephemeralStats(pods)
	aSto, cSto, uSto := p.storageStats(pods)

	return stats.NodeStats{
		Name:                 p.NodeInfo.Name,
		UsageNanoCores:       p.cpuStats(pods),
		AvailableBytesMemory: aMem,
		UsageBytesMemory:     uMem,

		AvailableBytesEphemeral: aEph,
		CapacityBytesEphemeral:  cEph,
		UsedBytesEphemeral:      uEph,

		AvailableBytesStorage: aSto,
		CapacityBytesStorage:  cSto,
		UsedBytesStorage:      uSto,
	}
}

func (p *Provider) cpuStats(pods []stats.PodStats) uint64 {
	usage := uint64(0)

	for _, pod := range pods {
		usage += pod.UsageNanoCores
	}

	return usage
}

func (p *Provider) memoryStats(pods []stats.PodStats) (uint64, uint64) {
	usage := uint64(0)

	for _, pod := range pods {
		usage += pod.UsageBytesMemory
	}

	available := uint64(p.Resources.Memory) - usage

	return available, usage
}

func (p *Provider) ephemeralStats(pods []stats.PodStats) (uint64, uint64, uint64) {
	capacity := uint64(p.Resources.EphemeralStorage)
	usage := uint64(0)

	for _, pod := range pods {
		usage += pod.UsedBytesEphemeral
	}

	free := capacity - usage
	return free, capacity, usage
}

func (p *Provider) storageStats(pods []stats.PodStats) (uint64, uint64, uint64) {
	capacity := uint64(p.Resources.Storage)
	usage := uint64(0)

	for _, pod := range pods {
		usage += pod.UsedBytesStorage
	}

	free := capacity - usage
	return free, capacity, usage
}

func (p *Provider) getAggregatePodStats() []stats.PodStats {
	var statistics []stats.PodStats

	for _, pod := range p.Pods.GetAllPods() {
		statistics = append(statistics, *p.getPodStats(pod))
	}

	return statistics
}

func (p *Provider) getPodStats(pod *corev1.Pod) *stats.PodStats {
	unconvertedStats, err := (*p.Store).GetPodFlag(pod, events.PodResources)
	if err != nil {
		log.Printf("error while retrieving pod flag for resources: %v\n", err)
		return addPodSpecificStats(pod, &stats.PodStats{})
	}

	statistics, ok := unconvertedStats.(*stats.PodStats)
	if !ok {
		log.Printf("unable to convert '%v' to PodStats\n", unconvertedStats)
		return addPodSpecificStats(pod, &stats.PodStats{})
	}

	return addPodSpecificStats(pod, statistics)
}

func addPodSpecificStats(pod *corev1.Pod, statistics *stats.PodStats) *stats.PodStats {
	statistics.PodRef = stats.PodReference{
		Name:      pod.Name,
		Namespace: pod.Namespace,
		UID:       string(pod.UID),
	}

	return statistics
}
