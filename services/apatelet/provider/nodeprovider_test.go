package provider

import (
	"context"
	"testing"
	"time"

	"github.com/finitum/node-cli/provider"

	"github.com/atlarge-research/apate/pkg/kubernetes/node"

	nodeconfigv1 "github.com/atlarge-research/apate/pkg/apis/nodeconfiguration/v1"
	podconfigv1 "github.com/atlarge-research/apate/pkg/apis/podconfiguration/v1"

	"github.com/finitum/node-cli/stats"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/atlarge-research/apate/pkg/scenario"
	"github.com/atlarge-research/apate/pkg/scenario/events"
	"github.com/atlarge-research/apate/services/apatelet/provider/condition"
	"github.com/atlarge-research/apate/services/apatelet/provider/podmanager"
	"github.com/atlarge-research/apate/services/apatelet/provider/podmanager/mock_podmanager"
	"github.com/atlarge-research/apate/services/apatelet/store"
	"github.com/atlarge-research/apate/services/apatelet/store/mock_store"
)

func TestPing(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mock_store.NewMockStore(ctrl)
	var st store.Store = ms

	pmm := mock_podmanager.NewMockPodManager(ctrl)
	var pm podmanager.PodManager = pmm

	prov := Provider{
		Pods:  pm,
		Store: &st,
	}

	ms.EXPECT().GetNodeFlag(events.NodeAddedLatency).Return(time.Duration(0), nil)
	ms.EXPECT().GetNodeFlag(events.NodePingResponse).Return(scenario.ResponseNormal, nil)

	res := prov.Ping(context.Background())
	assert.Equal(t, nil, res)
}

func TestPingError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mock_store.NewMockStore(ctrl)
	var st store.Store = ms

	pmm := mock_podmanager.NewMockPodManager(ctrl)
	var pm podmanager.PodManager = pmm

	ms.EXPECT().GetNodeFlag(events.NodeAddedLatency).Return(time.Duration(0), nil)
	ms.EXPECT().GetNodeFlag(events.NodePingResponse).Return(scenario.ResponseError, nil)

	prov := Provider{
		Pods:  pm,
		Store: &st,
	}

	res := prov.Ping(context.Background())
	assert.Error(t, res)
}

func TestConfigureNode(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mock_store.NewMockStore(ctrl)
	var st store.Store = ms

	pmm := mock_podmanager.NewMockPodManager(ctrl)
	var pm podmanager.PodManager = pmm

	pmm.EXPECT().GetAllPods().Return([]*corev1.Pod{})

	u := uuid.UUID{}
	prov := Provider{
		Pods:  pm,
		Store: &st,
		NodeInfo: &node.Info{
			NodeType:    "apate",
			Role:        "worker",
			Name:        "apate-x",
			Version:     "42",
			Namespace:   "my",
			Label:       "apate",
			MetricsPort: 123,
		},
		Cfg: &provider.InitConfig{
			ConfigPath:        "",
			NodeName:          "apate-x",
			OperatingSystem:   "not windows",
			InternalIP:        "123.123.123.123",
			DaemonPort:        4242,
			KubeClusterDomain: "",
			ResourceManager:   nil,
		},
		Resources: &scenario.NodeResources{
			UUID:             u,
			Memory:           4096,
			CPU:              1000,
			Storage:          2048,
			EphemeralStorage: 8192,
			MaxPods:          42,
			Label:            "my/apate",
		},
		Conditions: nodeConditions{
			ready:              condition.New(true, corev1.NodeReady),
			outOfDisk:          condition.New(false, corev1.NodeOutOfDisk),
			memoryPressure:     condition.New(false, corev1.NodeMemoryPressure),
			diskPressure:       condition.New(false, corev1.NodeDiskPressure),
			networkUnavailable: condition.New(false, corev1.NodeNetworkUnavailable),
			pidPressure:        condition.New(false, corev1.NodePIDPressure),
		},
		DisableTaints: false,
		Stats:         &Stats{},
	}

	newNode := &corev1.Node{}
	prov.updateStatsSummary()
	prov.ConfigureNode(context.Background(), newNode)

	assert.EqualValues(t, corev1.NodeSpec{
		Taints: []corev1.Taint{
			{
				Key:    nodeconfigv1.EmulatedLabel,
				Effect: corev1.TaintEffectNoSchedule,
			},
		},
	}, newNode.Spec)

	assert.EqualValues(t, metav1.ObjectMeta{
		Name: "apate-x",
		Labels: map[string]string{
			"type":                                       "apate",
			"kubernetes.io/role":                         "worker",
			"kubernetes.io/hostname":                     "apate-x",
			"metrics_port":                               "123",
			nodeconfigv1.EmulatedLabel:                   "yes",
			nodeconfigv1.NodeConfigurationLabel:          "apate",
			nodeconfigv1.NodeConfigurationLabelNamespace: "my",
			nodeconfigv1.NodeIDLabel:                     u.String(),
		},
	}, newNode.ObjectMeta)

	assert.EqualValues(t, corev1.ResourceList{
		corev1.ResourceCPU:              *resource.NewQuantity(1000, ""),
		corev1.ResourceMemory:           *resource.NewQuantity(4096, ""),
		corev1.ResourceStorage:          *resource.NewQuantity(2048, ""),
		corev1.ResourceEphemeralStorage: *resource.NewQuantity(8192, ""),
		corev1.ResourcePods:             *resource.NewQuantity(42, ""),
	}, newNode.Status.Capacity)

	assert.EqualValues(t, corev1.ResourceList{
		corev1.ResourceCPU:              *resource.NewQuantity(1000, ""),
		corev1.ResourceMemory:           *resource.NewQuantity(4096, ""),
		corev1.ResourceStorage:          *resource.NewQuantity(2048, ""),
		corev1.ResourceEphemeralStorage: *resource.NewQuantity(8192, ""),
		corev1.ResourcePods:             *resource.NewQuantity(42, ""),
	}, newNode.Status.Allocatable)

	assert.EqualValues(t, 2, len(newNode.Status.Addresses))

	assert.EqualValues(t, corev1.NodeDaemonEndpoints{
		KubeletEndpoint: corev1.DaemonEndpoint{
			Port: 4242,
		},
	}, newNode.Status.DaemonEndpoints)

	assert.EqualValues(t, corev1.NodeSystemInfo{
		KubeletVersion: "42",
		Architecture:   "amd64",
	}, newNode.Status.NodeInfo)
}

func TestUpdateConditionNoPressure(t *testing.T) {
	t.Parallel()

	prov, ctrl := createProviderForUpdateConditionTests(t, 500, 2048, 1024)
	defer ctrl.Finish()

	prov.updateConditions(func(node *corev1.Node) {
		assert.EqualValues(t, corev1.ConditionTrue, node.Status.Conditions[0].Status)
		assert.EqualValues(t, corev1.NodeReady, node.Status.Conditions[0].Type)

		assert.EqualValues(t, corev1.ConditionFalse, node.Status.Conditions[1].Status)
		assert.EqualValues(t, corev1.NodeOutOfDisk, node.Status.Conditions[1].Type)

		assert.EqualValues(t, corev1.ConditionFalse, node.Status.Conditions[2].Status)
		assert.EqualValues(t, corev1.NodeMemoryPressure, node.Status.Conditions[2].Type)

		assert.EqualValues(t, corev1.ConditionFalse, node.Status.Conditions[3].Status)
		assert.EqualValues(t, corev1.NodeDiskPressure, node.Status.Conditions[3].Type)

		assert.EqualValues(t, corev1.ConditionFalse, node.Status.Conditions[4].Status)
		assert.EqualValues(t, corev1.NodeNetworkUnavailable, node.Status.Conditions[4].Type)

		assert.EqualValues(t, corev1.ConditionFalse, node.Status.Conditions[5].Status)
		assert.EqualValues(t, corev1.NodePIDPressure, node.Status.Conditions[5].Type)
	})
}

func TestUpdateConditionMemoryAndDiskPressure(t *testing.T) {
	t.Parallel()

	mt := memThresh * 4096
	dt := diskThresh * 2048

	prov, ctrl := createProviderForUpdateConditionTests(t, 5000, int64(mt)+2, int64(dt)+2)
	defer ctrl.Finish()

	prov.updateConditions(func(node *corev1.Node) {
		assert.EqualValues(t, corev1.ConditionTrue, node.Status.Conditions[0].Status)
		assert.EqualValues(t, corev1.NodeReady, node.Status.Conditions[0].Type)

		assert.EqualValues(t, corev1.ConditionFalse, node.Status.Conditions[1].Status)
		assert.EqualValues(t, corev1.NodeOutOfDisk, node.Status.Conditions[1].Type)

		assert.EqualValues(t, corev1.ConditionTrue, node.Status.Conditions[2].Status)
		assert.EqualValues(t, corev1.NodeMemoryPressure, node.Status.Conditions[2].Type)

		assert.EqualValues(t, corev1.ConditionTrue, node.Status.Conditions[3].Status)
		assert.EqualValues(t, corev1.NodeDiskPressure, node.Status.Conditions[3].Type)

		assert.EqualValues(t, corev1.ConditionFalse, node.Status.Conditions[4].Status)
		assert.EqualValues(t, corev1.NodeNetworkUnavailable, node.Status.Conditions[4].Type)

		assert.EqualValues(t, corev1.ConditionFalse, node.Status.Conditions[5].Status)
		assert.EqualValues(t, corev1.NodePIDPressure, node.Status.Conditions[5].Type)
	})
}

func TestUpdateConditionDiskFull(t *testing.T) {
	t.Parallel()

	mtf := memThresh * 4096
	dtf := diskFullThresh * 2048

	prov, ctrl := createProviderForUpdateConditionTests(t, 5000, int64(mtf)+2, int64(dtf)+2)
	defer ctrl.Finish()

	prov.updateConditions(func(node *corev1.Node) {
		assert.EqualValues(t, corev1.ConditionFalse, node.Status.Conditions[0].Status)
		assert.EqualValues(t, corev1.NodeReady, node.Status.Conditions[0].Type)

		assert.EqualValues(t, corev1.ConditionTrue, node.Status.Conditions[1].Status)
		assert.EqualValues(t, corev1.NodeOutOfDisk, node.Status.Conditions[1].Type)

		assert.EqualValues(t, corev1.ConditionTrue, node.Status.Conditions[2].Status)
		assert.EqualValues(t, corev1.NodeMemoryPressure, node.Status.Conditions[2].Type)

		assert.EqualValues(t, corev1.ConditionTrue, node.Status.Conditions[3].Status)
		assert.EqualValues(t, corev1.NodeDiskPressure, node.Status.Conditions[3].Type)

		assert.EqualValues(t, corev1.ConditionFalse, node.Status.Conditions[4].Status)
		assert.EqualValues(t, corev1.NodeNetworkUnavailable, node.Status.Conditions[4].Type)

		assert.EqualValues(t, corev1.ConditionFalse, node.Status.Conditions[5].Status)
		assert.EqualValues(t, corev1.NodePIDPressure, node.Status.Conditions[5].Type)
	})
}

func createProviderForUpdateConditionTests(t *testing.T, podCPU, podMemory, podStorage int64) (Provider, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	// no defer ctrl.Finish() here because this function returns the ctrl

	ms := mock_store.NewMockStore(ctrl)
	var st store.Store = ms

	pmm := mock_podmanager.NewMockPodManager(ctrl)
	var pm podmanager.PodManager = pmm

	lbl := make(map[string]string)
	lbl[podconfigv1.PodConfigurationLabel] = "pod1"
	pod := corev1.Pod{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{Labels: lbl, Namespace: "a"},
		Spec:       corev1.PodSpec{},
		Status:     corev1.PodStatus{},
	}
	pmm.EXPECT().GetAllPods().Return([]*corev1.Pod{&pod})

	cores := uint64(podCPU)
	memory := uint64(podMemory)
	storage := uint64(podStorage)
	ms.EXPECT().GetPodFlag(&pod, events.PodResources).Return(&stats.PodStats{
		UsageNanoCores:     cores,
		UsageBytesMemory:   memory,
		UsedBytesEphemeral: storage,
	}, nil)

	ms.EXPECT().GetNodeFlag(events.NodePingResponse).Return(scenario.ResponseNormal, nil)

	u := uuid.UUID{}
	prov := Provider{
		Pods:  pm,
		Store: &st,
		Node:  &corev1.Node{},
		Resources: &scenario.NodeResources{
			UUID:             u,
			Memory:           4096,
			CPU:              1000,
			Storage:          2048,
			EphemeralStorage: 2048,
			MaxPods:          42,
			Label:            "my/apate",
		},
		NodeInfo: &node.Info{
			NodeType:    "",
			Role:        "",
			Name:        "",
			Version:     "",
			Namespace:   "",
			Label:       "",
			MetricsPort: 0,
		},
		Stats: NewStats(),
		Conditions: nodeConditions{
			ready:              condition.New(true, corev1.NodeReady),
			outOfDisk:          condition.New(false, corev1.NodeOutOfDisk),
			memoryPressure:     condition.New(false, corev1.NodeMemoryPressure),
			diskPressure:       condition.New(false, corev1.NodeDiskPressure),
			networkUnavailable: condition.New(false, corev1.NodeNetworkUnavailable),
			pidPressure:        condition.New(false, corev1.NodePIDPressure),
		},
		Cfg: &provider.InitConfig{
			DaemonPort: 100,
		},
	}

	prov.updateStatsSummary()

	return prov, ctrl
}

func TestNotifyNodeStatusNoPing(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mock_store.NewMockStore(ctrl)
	var st store.Store = ms

	pmm := mock_podmanager.NewMockPodManager(ctrl)
	var pm podmanager.PodManager = pmm

	ms.EXPECT().GetNodeFlag(events.NodePingResponse).Return(scenario.ResponseError, nil)

	prov := Provider{
		Pods:  pm,
		Store: &st,
		Node:  &corev1.Node{},
	}

	prov.updateConditions(func(node *corev1.Node) {
		assert.EqualValues(t, &corev1.Node{}, node)
	})
}
