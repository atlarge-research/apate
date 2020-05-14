package provider

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/virtual-kubelet/node-cli/provider"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	stats "k8s.io/kubernetes/pkg/kubelet/apis/stats/v1alpha1"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/normalization"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/provider/podmanager"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store/mock_store"
)

const (
	name  = "name"
	port  = 42
	flag  = "yes"
	event = events.PodResources
)

func createProvider(t *testing.T, cpu, mem int64) (provider.PodMetricsProvider, *gomock.Controller, *mock_store.MockStore, podmanager.PodManager) {
	ctrl := gomock.NewController(t)
	ms := mock_store.NewMockStore(ctrl)
	pm := podmanager.New() // TODO mock?
	var s store.Store = ms

	res := normalization.NodeResources{CPU: cpu, Memory: mem}
	info := cluster.NewNodeInfo("", "", name, "", port)
	prov := NewProvider(pm, NewStats(), &res, provider.InitConfig{}, info, &s)

	return prov.(provider.PodMetricsProvider), ctrl, ms, pm
}

func TestEmpty(t *testing.T) {
	mem := int64(34)
	prov, ctrl, _, _ := createProvider(t, 12, mem)

	result, err := prov.GetStatsSummary(context.Background())
	assert.NoError(t, err)

	// Verify node
	zero := uint64(0)
	assert.Equal(t, name, result.Node.NodeName)
	assert.Equal(t, zero, *result.Node.CPU.UsageNanoCores)
	assert.Equal(t, zero, *result.Node.Memory.UsageBytes)
	assert.Equal(t, uint64(mem), *result.Node.Memory.AvailableBytes)

	// Verify pods
	var statistics []stats.PodStats
	assert.EqualValues(t, statistics, result.Pods)

	ctrl.Finish()
}

func TestSinglePod(t *testing.T) {
	cpu := int64(123124)
	mem := int64(52562)
	memUsage := uint64(15)
	cpuUsage := uint64(16)
	prov, ctrl, ms, pm := createProvider(t, cpu, mem)

	// Create pod
	lbl := make(map[string]string)
	lbl["apate"] = flag
	pod := v1.Pod{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{Labels: lbl},
		Spec:       v1.PodSpec{},
		Status:     v1.PodStatus{},
	}
	pm.AddPod(pod) //TODO mock?

	// Create stats
	statistics := stats.PodStats{
		PodRef:           stats.PodReference{},
		StartTime:        metav1.Time{},
		Containers:       nil,
		CPU:              &stats.CPUStats{UsageNanoCores: &cpuUsage},
		Memory:           &stats.MemoryStats{UsageBytes: &memUsage},
		Network:          nil,
		VolumeStats:      nil,
		EphemeralStorage: nil,
	}

	// Setup store
	ms.EXPECT().GetPodFlag(flag, event).Return(statistics, nil)

	result, err := prov.GetStatsSummary(context.Background())
	assert.NoError(t, err)

	// Verify node
	left := uint64(mem) - memUsage
	assert.Equal(t, name, result.Node.NodeName)
	assert.Equal(t, cpuUsage, *result.Node.CPU.UsageNanoCores)
	assert.Equal(t, memUsage, *result.Node.Memory.UsageBytes)
	assert.Equal(t, left, *result.Node.Memory.AvailableBytes)

	// Verify pod
	podStats := []stats.PodStats{statistics}
	assert.EqualValues(t, podStats, result.Pods)

	ctrl.Finish()
}

func TestUnspecifiedPods(t *testing.T) {
	cpu := int64(2)
	mem := int64(2)
	memUsage := uint64(1)
	cpuUsage := uint64(1)
	prov, ctrl, ms, pm := createProvider(t, cpu, mem)

	// Create pods
	lbl := make(map[string]string)
	lbl["apate"] = flag
	pod := v1.Pod{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{Labels: lbl, UID: flag},
		Spec:       v1.PodSpec{},
		Status:     v1.PodStatus{},
	}
	pm.AddPod(pod) //TODO mock?

	lbl2 := make(map[string]string)
	lbl2["apate"] = flag + "2"
	pod2 := v1.Pod{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{Labels: lbl2, UID: flag + "2"},
		Spec:       v1.PodSpec{},
		Status:     v1.PodStatus{},
	}
	pm.AddPod(pod2) //TODO mock?

	lbl3 := make(map[string]string)
	lbl3["apate"] = flag + "3"
	pod3 := v1.Pod{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{Labels: lbl3, UID: flag + "3"},
		Spec:       v1.PodSpec{},
		Status:     v1.PodStatus{},
	}
	pm.AddPod(pod3) //TODO mock?
	pod4 := v1.Pod{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{UID: flag + "4"},
		Spec:       v1.PodSpec{},
		Status:     v1.PodStatus{},
	}
	pm.AddPod(pod4) //TODO mock?

	// Create stats
	statistics := stats.PodStats{
		PodRef:           stats.PodReference{UID: flag},
		StartTime:        metav1.Time{},
		Containers:       nil,
		CPU:              nil,
		Memory:           nil,
		Network:          nil,
		VolumeStats:      nil,
		EphemeralStorage: nil,
	}

	statistics2 := stats.PodStats{
		PodRef:           stats.PodReference{UID: flag + "2"},
		StartTime:        metav1.Time{},
		Containers:       nil,
		CPU:              &stats.CPUStats{},
		Memory:           &stats.MemoryStats{},
		Network:          nil,
		VolumeStats:      nil,
		EphemeralStorage: nil,
	}

	statistics3 := stats.PodStats{
		PodRef:           stats.PodReference{UID: flag + "3"},
		StartTime:        metav1.Time{},
		Containers:       nil,
		CPU:              &stats.CPUStats{UsageNanoCores: &cpuUsage},
		Memory:           &stats.MemoryStats{UsageBytes: &memUsage},
		Network:          nil,
		VolumeStats:      nil,
		EphemeralStorage: nil,
	}

	// Setup store
	ms.EXPECT().GetPodFlag(flag, event).Return(statistics, nil)
	ms.EXPECT().GetPodFlag(flag+"2", event).Return(statistics2, nil)
	ms.EXPECT().GetPodFlag(flag+"3", event).Return(statistics3, nil)

	result, err := prov.GetStatsSummary(context.Background())
	assert.NoError(t, err)

	// Verify node
	left := uint64(mem) - memUsage
	assert.Equal(t, name, result.Node.NodeName)
	assert.Equal(t, cpuUsage, *result.Node.CPU.UsageNanoCores)
	assert.Equal(t, memUsage, *result.Node.Memory.UsageBytes)
	assert.Equal(t, left, *result.Node.Memory.AvailableBytes)

	// Verify pod
	podStats := []stats.PodStats{statistics, statistics2, statistics3, {PodRef: stats.PodReference{UID: flag + "4"}}}
	assert.EqualValues(t, podStats, result.Pods)

	ctrl.Finish()
}
