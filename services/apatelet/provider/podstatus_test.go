package provider

import (
	"context"
	"testing"
	"time"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes/node"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/types"
	stats "k8s.io/kubernetes/pkg/kubelet/apis/stats/v1alpha1"

	podconfigv1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/podconfiguration/v1"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/provider/podmanager"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store/mock_store"
)

func prepareState(t *testing.T, nodeResources int64, podResources uint64, podMaxResources int64, podStatus scenario.PodStatus, response scenario.Response) (Provider, *gomock.Controller) {
	ctrl := gomock.NewController(t)
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
			Name:  podContainerName,
			Image: podImageName,
			Resources: corev1.ResourceRequirements{
				Limits: corev1.ResourceList{
					corev1.ResourceCPU:              *resource.NewQuantity(podMaxResources, ""),
					corev1.ResourceMemory:           *resource.NewQuantity(podMaxResources, ""),
					corev1.ResourceEphemeralStorage: *resource.NewQuantity(podMaxResources, ""),
				},
			},
		},
	}

	// expect
	ms.EXPECT().GetNodeFlag(events.NodeAddedLatency).Return(time.Duration(0), nil)
	ms.EXPECT().GetPodFlag(&pod, events.PodGetPodStatusResponse).Return(response, nil)
	ms.EXPECT().GetNodeFlag(events.NodeGetPodStatusResponse).Return(scenario.ResponseUnset, nil)

	isNormal := response == scenario.ResponseNormal || response == scenario.ResponseUnset

	expectedResourceGets := 1
	if isNormal {
		expectedResourceGets = 2 // First by updateStatsSummary and then in getPodStatus limitReached
	}

	// Because we compute the resources up front
	ms.EXPECT().GetPodFlag(&pod, events.PodResources).Return(&stats.PodStats{
		CPU: &stats.CPUStats{
			UsageNanoCores: &podResources,
		},
		Memory: &stats.MemoryStats{
			UsageBytes: &podResources,
		},
		EphemeralStorage: &stats.FsStats{
			UsedBytes: &podResources,
		},
	}, nil).Times(expectedResourceGets)

	if isNormal {
		ms.EXPECT().GetPodFlag(&pod, events.PodStatus).Return(podStatus, nil)
	}

	// sot
	var s store.Store = ms
	prov := Provider{
		Store: &s,
		Pods:  podmanager.New(),
		Resources: &scenario.NodeResources{
			CPU:              nodeResources,
			Memory:           nodeResources,
			EphemeralStorage: nodeResources,
		},
		NodeInfo: &node.Info{},
		Stats: &Stats{
			statsSummary: &stats.Summary{},
		},
	}
	prov.Pods.AddPod(&pod)

	prov.updateStatsSummary()
	return prov, ctrl
}

func TestGetPodStatus(t *testing.T) {
	t.Parallel()

	prov, ctrl := prepareState(t, 1000, 1, 1000, scenario.PodStatusRunning, scenario.ResponseNormal)
	defer ctrl.Finish()

	ps, err := prov.GetPodStatus(context.Background(), podNamespace, podName)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, corev1.PodRunning, ps.Phase)

	assert.Len(t, ps.Conditions, 1)
	assert.Equal(t, corev1.PodReady, ps.Conditions[0].Type)
	assert.Equal(t, corev1.ConditionTrue, ps.Conditions[0].Status)

	assert.Len(t, ps.ContainerStatuses, 1)
	assert.EqualValues(t, corev1.ContainerStatus{
		Name: podContainerName,
		State: corev1.ContainerState{
			Running: &corev1.ContainerStateRunning{},
		},
		Ready:       true,
		Image:       podImageName,
		ImageID:     "",
		ContainerID: "",
	}, ps.ContainerStatuses[0])

	assert.Len(t, prov.Pods.GetAllPods(), 1)
}

func TestGetPodStatusPodLimitReached(t *testing.T) {
	t.Parallel()

	prov, ctrl := prepareState(t, 1000, 128, 64, scenario.PodStatusRunning, scenario.ResponseNormal)
	defer ctrl.Finish()

	ps, err := prov.GetPodStatus(context.Background(), podNamespace, podName)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, corev1.PodFailed, ps.Phase)
	assert.Equal(t, corev1.PodReady, ps.Conditions[0].Type)
	assert.Equal(t, corev1.ConditionFalse, ps.Conditions[0].Status)
	assert.Len(t, prov.Pods.GetAllPods(), 1)
}

func TestGetPodStatusNodeLimitReached(t *testing.T) {
	t.Parallel()

	prov, ctrl := prepareState(t, 64, 128, 1000, scenario.PodStatusRunning, scenario.ResponseNormal)
	defer ctrl.Finish()

	ps, err := prov.GetPodStatus(context.Background(), podNamespace, podName)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, corev1.PodFailed, ps.Phase)
	assert.Equal(t, corev1.PodReady, ps.Conditions[0].Type)
	assert.Equal(t, corev1.ConditionFalse, ps.Conditions[0].Status)
	assert.Len(t, prov.Pods.GetAllPods(), 1)

	assert.Len(t, ps.ContainerStatuses, 1)
	assert.EqualValues(t, corev1.ContainerStatus{
		Name: podContainerName,
		State: corev1.ContainerState{
			Terminated: &corev1.ContainerStateTerminated{
				ExitCode: 1,
				Reason:   "Pod status is failed, for reason Pod used too many resources and was then killed",
			},
		},
		Ready:       false,
		Image:       podImageName,
		ImageID:     "",
		ContainerID: "",
	}, ps.ContainerStatuses[0])
}

func TestGetPodStatusSucceeded(t *testing.T) {
	t.Parallel()

	prov, ctrl := prepareState(t, 1000, 128, 1000, scenario.PodStatusSucceeded, scenario.ResponseNormal)
	defer ctrl.Finish()

	ps, err := prov.GetPodStatus(context.Background(), podNamespace, podName)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, corev1.PodSucceeded, ps.Phase)
	assert.Empty(t, ps.Conditions)
	assert.Len(t, prov.Pods.GetAllPods(), 1)

	assert.Len(t, ps.ContainerStatuses, 1)
	assert.EqualValues(t, corev1.ContainerStatus{
		Name: podContainerName,
		State: corev1.ContainerState{
			Terminated: &corev1.ContainerStateTerminated{
				ExitCode: 0,
				Reason:   "Pod status is succeeded",
			},
		},
		Ready:       false,
		Image:       podImageName,
		ImageID:     "",
		ContainerID: "",
	}, ps.ContainerStatuses[0])
}

func TestGetPodStatusPending(t *testing.T) {
	t.Parallel()

	prov, ctrl := prepareState(t, 1000, 128, 1000, scenario.PodStatusPending, scenario.ResponseNormal)
	defer ctrl.Finish()

	ps, err := prov.GetPodStatus(context.Background(), podNamespace, podName)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, corev1.PodPending, ps.Phase)
	assert.Equal(t, corev1.PodScheduled, ps.Conditions[0].Type)
	assert.Equal(t, corev1.ConditionTrue, ps.Conditions[0].Status)
	assert.Len(t, prov.Pods.GetAllPods(), 1)

	assert.Len(t, ps.ContainerStatuses, 1)
	assert.EqualValues(t, corev1.ContainerStatus{
		Name: podContainerName,
		State: corev1.ContainerState{
			Waiting: &corev1.ContainerStateWaiting{
				Reason: "Pod status is pending",
			},
		},
		Ready:       false,
		Image:       podImageName,
		ImageID:     "",
		ContainerID: "",
	}, ps.ContainerStatuses[0])
}

func TestGetPodStatusUnknown(t *testing.T) {
	t.Parallel()

	prov, ctrl := prepareState(t, 1000, 128, 1000, scenario.PodStatusUnknown, scenario.ResponseNormal)
	defer ctrl.Finish()

	ps, err := prov.GetPodStatus(context.Background(), podNamespace, podName)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, corev1.PodUnknown, ps.Phase)
	assert.Empty(t, ps.Conditions)

	assert.Len(t, ps.ContainerStatuses, 1)
	assert.EqualValues(t, corev1.ContainerStatus{
		Name: podContainerName,
		State: corev1.ContainerState{
			Waiting: &corev1.ContainerStateWaiting{
				Reason: "Pod status is unknown",
			},
		},
		Ready:       false,
		Image:       podImageName,
		ImageID:     "",
		ContainerID: "",
	}, ps.ContainerStatuses[0])
}

func TestGetPodStatusEmulationError(t *testing.T) {
	t.Parallel()

	prov, ctrl := prepareState(t, 1000, 128, 1000, scenario.PodStatusUnknown, scenario.ResponseError)
	defer ctrl.Finish()

	_, err := prov.GetPodStatus(context.Background(), podNamespace, podName)

	// assert
	assert.Error(t, err)
	assert.True(t, IsExpected(err))
}

func TestGetPodStatusDefault(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mock_store.NewMockStore(ctrl)
	pod := corev1.Pod{}
	pod.Namespace = podNamespace
	pod.Name = podName
	pod.Labels = map[string]string{
		podconfigv1.PodConfigurationLabel: podLabel,
	}

	ms.EXPECT().GetNodeFlag(events.NodeAddedLatency).Return(time.Duration(0), nil)
	ms.EXPECT().GetPodFlag(&pod, events.PodGetPodStatusResponse).Return(scenario.ResponseUnset, nil)
	ms.EXPECT().GetNodeFlag(events.NodeGetPodStatusResponse).Return(scenario.ResponseUnset, nil)

	ms.EXPECT().GetPodFlag(&pod, events.PodResources).Return(&stats.PodStats{}, nil).Times(2)
	ms.EXPECT().GetPodFlag(&pod, events.PodStatus).Return(scenario.PodStatusUnset, nil)

	var s store.Store = ms
	nodeResources := int64(1000)
	prov := Provider{
		Store: &s,
		Pods:  podmanager.New(),
		Resources: &scenario.NodeResources{
			CPU:              nodeResources,
			Memory:           nodeResources,
			EphemeralStorage: nodeResources,
		},
		NodeInfo: &node.Info{},
		Stats: &Stats{
			statsSummary: &stats.Summary{},
		},
	}
	prov.Pods.AddPod(&pod)

	prov.updateStatsSummary()

	ps, err := prov.GetPodStatus(context.Background(), podNamespace, podName)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, corev1.PodRunning, ps.Phase)
	assert.Equal(t, corev1.PodReady, ps.Conditions[0].Type)
	assert.Equal(t, corev1.ConditionTrue, ps.Conditions[0].Status)
	assert.Len(t, prov.Pods.GetAllPods(), 1)
}
