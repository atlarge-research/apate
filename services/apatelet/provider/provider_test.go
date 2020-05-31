package provider

import (
	"context"
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/api/resource"
	stats "k8s.io/kubernetes/pkg/kubelet/apis/stats/v1alpha1"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/virtual-kubelet/node-cli/provider"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	podconfigv1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/podconfiguration/v1"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/provider/podmanager"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store/mock_store"
)

const podNamespace = "podnamespace"
const podName = "pod"
const podLabel = "label"

func TestConfigureNodeWithCreate(t *testing.T) {
	t.Parallel()

	resources := scenario.NodeResources{
		UUID:    uuid.New(),
		Memory:  42,
		CPU:     1337,
		MaxPods: 1001,
	}

	ctrl := gomock.NewController(t)
	st := store.NewStore()

	prov := NewProvider(podmanager.New(), NewStats(), &resources, provider.InitConfig{}, kubernetes.NodeInfo{}, &st)

	fakeNode := corev1.Node{}

	// Run the method
	prov.ConfigureNode(context.Background(), &fakeNode)

	assert.EqualValues(t, resources.CPU, fakeNode.Status.Capacity.Cpu().Value())
	assert.EqualValues(t, resources.Memory, fakeNode.Status.Capacity.Memory().Value())
	assert.EqualValues(t, resources.MaxPods, fakeNode.Status.Capacity.Pods().Value())

	ctrl.Finish()
}

func TestCreatePod(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	ms := mock_store.NewMockStore(ctrl)

	/// vars
	pod := corev1.Pod{}
	pod.Namespace = podNamespace
	pod.Name = podName
	pod.Labels = map[string]string{
		podconfigv1.PodConfigurationLabel: podLabel,
	}
	pod.UID = types.UID(uuid.New().String())
	PCPRF := events.PodCreatePodResponse

	// expect
	ms.EXPECT().GetNodeFlag(events.NodeAddedLatency).Return(time.Duration(0), nil)
	ms.EXPECT().GetPodFlag(podNamespace+"/"+podLabel, PCPRF).Return(scenario.ResponseNormal, nil)

	// sot
	var s store.Store = ms

	p := Provider{
		Store: &s,
		Pods:  podmanager.New(),
	}

	err := p.CreatePod(context.Background(), &pod)

	// assert
	assert.NoError(t, err)

	uid, ok := p.Pods.GetPodByUID(pod.UID)
	assert.True(t, ok)
	assert.Equal(t, &pod, uid)
	ctrl.Finish()
}

func TestUpdatePod(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	ms := mock_store.NewMockStore(ctrl)

	/// vars
	pod := corev1.Pod{}
	pod.Namespace = podNamespace
	pod.Name = podName
	pod.Labels = map[string]string{
		podconfigv1.PodConfigurationLabel: podLabel,
	}
	pod.UID = types.UID(uuid.New().String())
	PCPRF := events.PodUpdatePodResponse

	// expect
	ms.EXPECT().GetNodeFlag(events.NodeAddedLatency).Return(time.Duration(0), nil)
	ms.EXPECT().GetPodFlag(podNamespace+"/"+podLabel, PCPRF).Return(scenario.ResponseNormal, nil)

	// sot
	var s store.Store = ms
	p := Provider{
		Store: &s,
		Pods:  podmanager.New(),
	}

	err := p.UpdatePod(context.Background(), &pod)

	// assert
	assert.NoError(t, err)
	uid, ok := p.Pods.GetPodByUID(pod.UID)
	assert.True(t, ok)
	assert.Equal(t, &pod, uid)
	ctrl.Finish()
}

func TestDeletePod(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	ms := mock_store.NewMockStore(ctrl)

	/// vars
	pod := corev1.Pod{}
	pod.Namespace = podNamespace
	pod.Name = podName
	pod.Labels = map[string]string{
		podconfigv1.PodConfigurationLabel: podLabel,
	}
	pod.UID = types.UID(uuid.New().String())
	PCPRF := events.PodDeletePodResponse

	// expect
	ms.EXPECT().GetNodeFlag(events.NodeAddedLatency).Return(time.Duration(0), nil)
	ms.EXPECT().GetPodFlag(podNamespace+"/"+podLabel, PCPRF).Return(scenario.ResponseNormal, nil)

	// sot
	var s store.Store = ms
	p := Provider{
		Store: &s,
		Pods:  podmanager.New(),
	}

	err := p.DeletePod(context.Background(), &pod)

	// assert
	assert.NoError(t, err)
	assert.NotContains(t, p.Pods.GetAllPods(), &pod)
	ctrl.Finish()
}

func TestGetPod(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	ms := mock_store.NewMockStore(ctrl)

	/// vars
	pod := corev1.Pod{}
	pod.Namespace = podNamespace
	pod.Name = podName
	pod.Labels = map[string]string{
		podconfigv1.PodConfigurationLabel: podLabel,
	}
	pod.UID = types.UID(uuid.New().String())
	PCPRF := events.PodGetPodResponse

	// expect
	ms.EXPECT().GetNodeFlag(events.NodeAddedLatency).Return(time.Duration(0), nil)
	ms.EXPECT().GetPodFlag(podNamespace+"/"+podLabel, PCPRF).Return(scenario.ResponseNormal, nil)

	// sot
	var s store.Store = ms
	prov := Provider{
		Store: &s,
		Pods:  podmanager.New(),
	}

	prov.Pods.AddPod(&pod)

	np, err := prov.GetPod(context.Background(), podNamespace, podName)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, &pod, np)
	ctrl.Finish()
}

func TestGetPods(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	ms := mock_store.NewMockStore(ctrl)

	/// vars
	pod := corev1.Pod{}
	pod.Namespace = podNamespace
	pod.Name = podName
	pod.Labels = map[string]string{
		podconfigv1.PodConfigurationLabel: podLabel,
	}
	pod.UID = types.UID(uuid.New().String())
	PCPRF := events.NodeGetPodsResponse

	// expect
	ms.EXPECT().GetNodeFlag(PCPRF).Return(scenario.ResponseNormal, nil)
	ms.EXPECT().GetNodeFlag(events.NodeAddedLatency).Return(time.Duration(0), nil)

	// sot
	var s store.Store = ms
	prov := Provider{
		Store: &s,
		Pods:  podmanager.New(),
	}
	prov.Pods.AddPod(&pod)

	ps, err := prov.GetPods(context.Background())

	// assert
	assert.NoError(t, err)
	uid, ok := prov.Pods.GetPodByUID(pod.UID)
	assert.True(t, ok)
	assert.Contains(t, ps, uid)
	ctrl.Finish()
}

func TestGetPodStatus(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := mock_store.NewMockStore(ctrl)

	// vars
	pod := corev1.Pod{}
	pod.Namespace = podNamespace
	pod.Name = podName
	pod.Labels = map[string]string{
		podconfigv1.PodConfigurationLabel: podLabel,
	}
	pod.UID = types.UID(uuid.New().String())
	pod.Spec.Containers = []corev1.Container{
		{
			Resources: corev1.ResourceRequirements{
				Limits: corev1.ResourceList{
					corev1.ResourceCPU:              *resource.NewQuantity(1000, ""),
					corev1.ResourceMemory:           *resource.NewQuantity(1000, ""),
					corev1.ResourceEphemeralStorage: *resource.NewQuantity(1000, ""),
				},
			},
		},
	}

	// expect
	one := uint64(1)

	ms.EXPECT().GetNodeFlag(events.NodeAddedLatency).Return(time.Duration(0), nil)
	ms.EXPECT().GetPodFlag(podNamespace+"/"+podLabel, events.PodResources).Return(stats.PodStats{
		CPU: &stats.CPUStats{
			UsageNanoCores: &one,
		},
		Memory: &stats.MemoryStats{
			UsageBytes: &one,
		},
		EphemeralStorage: &stats.FsStats{
			UsedBytes: &one,
		},
	}, nil).Times(2)
	ms.EXPECT().GetPodFlag(podNamespace+"/"+podLabel, events.PodGetPodStatusResponse).Return(scenario.ResponseNormal, nil)
	ms.EXPECT().GetPodFlag(podNamespace+"/"+podLabel, events.PodStatus).Return(scenario.PodStatusSucceeded, nil)

	// sot
	var s store.Store = ms
	prov := Provider{
		Store: &s,
		Pods:  podmanager.New(),
		Resources: &scenario.NodeResources{
			CPU:              1000,
			Memory:           1000,
			EphemeralStorage: 1000,
		},
		Stats: &Stats{
			statsSummary: &stats.Summary{},
		},
	}
	prov.Pods.AddPod(&pod)

	prov.updateStatsSummary()
	ps, err := prov.GetPodStatus(context.Background(), podNamespace, podName)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, corev1.PodSucceeded, ps.Phase)
	assert.Equal(t, corev1.PodReady, ps.Conditions[0].Type)
	assert.Equal(t, corev1.ConditionTrue, ps.Conditions[0].Status)
}

func TestGetPodStatusLimitReached(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := mock_store.NewMockStore(ctrl)

	// vars
	pod := corev1.Pod{}
	pod.Namespace = podNamespace
	pod.Name = podName
	pod.Labels = map[string]string{
		podconfigv1.PodConfigurationLabel: podLabel,
	}
	pod.UID = types.UID(uuid.New().String())
	pod.Spec.Containers = []corev1.Container{
		{
			Resources: corev1.ResourceRequirements{
				Limits: corev1.ResourceList{
					corev1.ResourceCPU:              *resource.NewQuantity(64, ""),
					corev1.ResourceMemory:           *resource.NewQuantity(64, ""),
					corev1.ResourceEphemeralStorage: *resource.NewQuantity(64, ""),
				},
			},
		},
	}

	// expect
	moreThan64 := uint64(128)

	ms.EXPECT().GetNodeFlag(events.NodeAddedLatency).Return(time.Duration(0), nil)
	ms.EXPECT().GetPodFlag(podNamespace+"/"+podLabel, events.PodResources).Return(stats.PodStats{
		CPU: &stats.CPUStats{
			UsageNanoCores: &moreThan64,
		},
		Memory: &stats.MemoryStats{
			UsageBytes: &moreThan64,
		},
		EphemeralStorage: &stats.FsStats{
			UsedBytes: &moreThan64,
		},
	}, nil).Times(2)
	ms.EXPECT().GetPodFlag(podNamespace+"/"+podLabel, events.PodGetPodStatusResponse).Return(scenario.ResponseNormal, nil)
	ms.EXPECT().GetPodFlag(podNamespace+"/"+podLabel, events.PodStatus).Return(scenario.PodStatusSucceeded, nil)

	// sot
	var s store.Store = ms
	prov := Provider{
		Store:     &s,
		Pods:      podmanager.New(),
		Resources: &scenario.NodeResources{},
		Stats: &Stats{
			statsSummary: &stats.Summary{},
		},
	}
	prov.Pods.AddPod(&pod)

	prov.updateStatsSummary()
	ps, err := prov.GetPodStatus(context.Background(), podNamespace, podName)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, corev1.PodFailed, ps.Phase)
	assert.Equal(t, corev1.PodReady, ps.Conditions[0].Type)
	assert.Equal(t, corev1.ConditionFalse, ps.Conditions[0].Status)
}

func TestNewProvider(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := mock_store.NewMockStore(ctrl)

	pm := podmanager.New()
	stats := NewStats()
	resources := scenario.NodeResources{
		UUID:             uuid.New(),
		Memory:           0,
		CPU:              0,
		Storage:          0,
		EphemeralStorage: 0,
		MaxPods:          0,
		Selector:         "",
	}

	cfg := provider.InitConfig{}
	ni, err := kubernetes.NewNodeInfo("a", "b", "c", "d", "e/f", 4242)
	assert.NoError(t, err)

	var s store.Store = ms

	ms.EXPECT().AddPodListener(events.PodResources, gomock.Any())

	p, ok := NewProvider(pm, stats, &resources, cfg, ni, &s).(*Provider)

	assert.True(t, ok)

	assert.EqualValues(t, p.Conditions.ready.Get().Status, metav1.ConditionTrue)
	assert.EqualValues(t, p.Conditions.outOfDisk.Get().Status, metav1.ConditionFalse)
	assert.EqualValues(t, p.Conditions.memoryPressure.Get().Status, metav1.ConditionFalse)
	assert.EqualValues(t, p.Conditions.diskPressure.Get().Status, metav1.ConditionFalse)
	assert.EqualValues(t, p.Conditions.networkUnavailable.Get().Status, metav1.ConditionFalse)
	assert.EqualValues(t, p.Conditions.pidPressure.Get().Status, metav1.ConditionFalse)
}
