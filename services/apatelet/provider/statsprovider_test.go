package provider

import (
	"testing"

	"github.com/finitum/node-cli/provider"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes/node"

	podconfigv1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/podconfiguration/v1"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario"

	"github.com/finitum/node-cli/stats"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/provider/podmanager"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store/mock_store"
)

const (
	name      = "name"
	label     = "label"
	namespace = "namespace"
	event     = events.PodResources
)

func createProvider(t *testing.T, cpu, mem, fs int64) (*Provider, *gomock.Controller, *mock_store.MockStore, podmanager.PodManager) {
	ctrl := gomock.NewController(t)
	// no defer ctrl.Finish() here because this function returns the ctrl

	ms := mock_store.NewMockStore(ctrl)
	var s store.Store = ms

	pm := podmanager.New() // TODO mock?

	res := scenario.NodeResources{CPU: cpu, Memory: mem, EphemeralStorage: fs}
	info, err := node.NewInfo("", "", name, "", "a/b")
	assert.NoError(t, err)

	ms.EXPECT().AddPodFlagListener(events.PodResources, gomock.Any())

	e, err := env.ApateletEnv()
	assert.NoError(t, err)
	prov := NewProvider(pm, NewStats(), &res, &provider.InitConfig{}, &info, &s, true, e)
	p := prov.(*Provider)
	return p, ctrl, ms, pm
}

func TestEmpty(t *testing.T) {
	t.Parallel()

	mem := int64(34)

	prov, ctrl, _, _ := createProvider(t, 12, mem, 0)
	defer ctrl.Finish()

	result, err := prov.GetStatsSummary()
	assert.NoError(t, err)

	// Verify node
	zero := uint64(0)
	assert.Equal(t, name, result.Node.Name)
	assert.Equal(t, zero, result.Node.UsageNanoCores)
	assert.Equal(t, zero, result.Node.UsageBytesMemory)
	assert.Equal(t, uint64(mem), result.Node.AvailableBytesMemory)

	// Verify pods
	assert.Empty(t, result.Pods)
}

func TestSinglePod(t *testing.T) {
	t.Parallel()

	cpu := int64(123124)
	mem := int64(52562)
	memUsage := uint64(15)
	cpuUsage := uint64(16)

	prov, ctrl, ms, pm := createProvider(t, cpu, mem, 0)
	defer ctrl.Finish()

	// Create pod
	lbl := make(map[string]string)
	lbl[podconfigv1.PodConfigurationLabel] = label
	pod := &corev1.Pod{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{Labels: lbl, Namespace: namespace},
		Spec:       corev1.PodSpec{},
		Status:     corev1.PodStatus{},
	}
	pm.AddPod(pod)

	// Create stats
	statistics := &stats.PodStats{
		PodRef:           stats.PodReference{Namespace: namespace},
		UsageNanoCores:   cpuUsage,
		UsageBytesMemory: memUsage,
	}

	// Setup store
	ms.EXPECT().GetPodFlag(pod, event).Return(statistics, nil)

	prov.updateStatsSummary()

	result, err := prov.GetStatsSummary()
	assert.NoError(t, err)

	// Verify node
	left := uint64(mem) - memUsage
	assert.Equal(t, name, result.Node.Name)
	assert.Equal(t, cpuUsage, result.Node.UsageNanoCores)
	assert.Equal(t, memUsage, result.Node.UsageBytesMemory)
	assert.Equal(t, left, result.Node.AvailableBytesMemory)

	// Verify pod
	podStats := []stats.PodStats{*statistics}
	assert.EqualValues(t, podStats, result.Pods)
}

func TestUnspecifiedPods(t *testing.T) {
	t.Parallel()

	cpu := int64(2)
	mem := int64(2)
	fs := int64(15)
	memUsage := uint64(1)
	cpuUsage := uint64(1)
	fsUsage := uint64(12)

	prov, ctrl, ms, pm := createProvider(t, cpu, mem, fs)
	defer ctrl.Finish()

	// Create pods
	lbl := make(map[string]string)
	lbl[podconfigv1.PodConfigurationLabel] = label
	pod := &corev1.Pod{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{Labels: lbl, UID: label, Namespace: namespace},
		Spec:       corev1.PodSpec{},
		Status:     corev1.PodStatus{},
	}
	pm.AddPod(pod)

	lbl2 := make(map[string]string)
	lbl2[podconfigv1.PodConfigurationLabel] = label + "2"
	pod2 := &corev1.Pod{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{Labels: lbl2, UID: label + "2", Namespace: namespace},
		Spec:       corev1.PodSpec{},
		Status:     corev1.PodStatus{},
	}
	pm.AddPod(pod2)

	lbl3 := make(map[string]string)
	lbl3[podconfigv1.PodConfigurationLabel] = label + "3"
	pod3 := &corev1.Pod{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{Labels: lbl3, UID: label + "3", Namespace: namespace},
		Spec:       corev1.PodSpec{},
		Status:     corev1.PodStatus{},
	}
	pm.AddPod(pod3)

	lbl4 := make(map[string]string)
	lbl4[podconfigv1.PodConfigurationLabel] = label + "4"
	pod4 := &corev1.Pod{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{Labels: lbl4, UID: label + "4", Namespace: namespace},
		Spec:       corev1.PodSpec{},
		Status:     corev1.PodStatus{},
	}
	pm.AddPod(pod4)

	// Create stats
	statistics := &stats.PodStats{
		PodRef: stats.PodReference{UID: label, Namespace: namespace},
	}

	statistics2 := &stats.PodStats{
		PodRef:             stats.PodReference{UID: label + "2", Namespace: namespace},
		UsedBytesEphemeral: fsUsage,
	}

	statistics3 := &stats.PodStats{
		PodRef:           stats.PodReference{UID: label + "3", Namespace: namespace},
		UsageNanoCores:   cpuUsage,
		UsageBytesMemory: memUsage,
	}

	statistics4 := &stats.PodStats{
		PodRef: stats.PodReference{UID: label + "4", Namespace: namespace},
	}

	statMap := make(map[string]*stats.PodStats)
	statMap[label] = statistics
	statMap[label+"2"] = statistics2
	statMap[label+"3"] = statistics3
	statMap[label+"4"] = statistics4

	// Setup store
	ms.EXPECT().GetPodFlag(pod, event).Return(statistics, nil)
	ms.EXPECT().GetPodFlag(pod2, event).Return(statistics2, nil)
	ms.EXPECT().GetPodFlag(pod3, event).Return(statistics3, nil)
	ms.EXPECT().GetPodFlag(pod4, event).Return(statistics4, nil)

	prov.updateStatsSummary()
	result, err := prov.GetStatsSummary()
	assert.NoError(t, err)

	// Verify node
	memLeft := uint64(mem) - memUsage
	fsLeft := uint64(fs) - fsUsage
	assert.Equal(t, name, result.Node.Name)
	assert.Equal(t, cpuUsage, result.Node.UsageNanoCores)
	assert.Equal(t, memUsage, result.Node.UsageBytesMemory)
	assert.Equal(t, memLeft, result.Node.AvailableBytesMemory)
	assert.Equal(t, fsUsage, result.Node.UsedBytesEphemeral)
	assert.Equal(t, fsLeft, result.Node.AvailableBytesEphemeral)
	assert.Equal(t, uint64(fs), result.Node.CapacityBytesEphemeral)

	// Verify pod
	for _, podStat := range result.Pods {
		assert.Equal(t, *statMap[podStat.PodRef.UID], podStat)
	}
}
