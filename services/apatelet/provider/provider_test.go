package provider

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/virtual-kubelet/node-cli/provider"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/normalization"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store/mock_store"
)

const podName = "pod"

func TestConfigureNode(t *testing.T) {
	resources := normalization.NodeResources{
		UUID:    uuid.New(),
		Memory:  42,
		CPU:     1337,
		MaxPods: 1001,
	}

	prov := Provider{
		pods:      NewPodManager(),
		resources: &resources,
	}

	fakeNode := corev1.Node{}

	// Run the method
	prov.ConfigureNode(context.TODO(), &fakeNode)

	assert.EqualValues(t, resources.CPU, fakeNode.Status.Capacity.Cpu().Value())
	assert.EqualValues(t, resources.Memory, fakeNode.Status.Capacity.Memory().Value())
	assert.EqualValues(t, resources.MaxPods, fakeNode.Status.Capacity.Pods().Value())
}

func TestConfigureNodeWithCreate(t *testing.T) {
	resources := normalization.NodeResources{
		UUID:    uuid.New(),
		Memory:  42,
		CPU:     1337,
		MaxPods: 1001,
	}

	st := store.NewStore()
	prov := NewProvider(&resources, provider.InitConfig{}, cluster.NodeInfo{}, &st)

	fakeNode := corev1.Node{}

	// Run the method
	prov.ConfigureNode(context.TODO(), &fakeNode)

	assert.EqualValues(t, resources.CPU, fakeNode.Status.Capacity.Cpu().Value())
	assert.EqualValues(t, resources.Memory, fakeNode.Status.Capacity.Memory().Value())
	assert.EqualValues(t, resources.MaxPods, fakeNode.Status.Capacity.Pods().Value())
}

func TestCreatePod(t *testing.T) {
	ctrl := gomock.NewController(t)
	ms := mock_store.NewMockStore(ctrl)

	/// vars
	pod := corev1.Pod{}
	pod.Name = podName
	pod.UID = types.UID(uuid.New().String())
	PCPRF := events.PodCreatePodResponse
	PCPRPF := events.PodCreatePodResponsePercentage

	// expect
	ms.EXPECT().GetNodeFlag(events.NodeAddedLatencyEnabled).Return(false, nil)
	ms.EXPECT().GetPodFlag(pod.Name, PCPRF).Return(scenario.Response_NORMAL, nil)
	ms.EXPECT().GetPodFlag(pod.Name, PCPRPF).Return(int32(100), nil)

	// sot
	var s store.Store = ms
	p := Provider{
		store: &s,
		pods:  NewPodManager(),
	}

	err := p.CreatePod(context.TODO(), &pod)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, &pod, p.pods.GetPodByUID(pod.UID))
	ctrl.Finish()
}

func TestUpdatePod(t *testing.T) {
	ctrl := gomock.NewController(t)
	ms := mock_store.NewMockStore(ctrl)

	/// vars
	pod := corev1.Pod{}
	pod.Name = podName
	pod.UID = types.UID(uuid.New().String())
	PCPRF := events.PodUpdatePodResponse
	PCPRPF := events.PodUpdatePodResponsePercentage

	// expect
	ms.EXPECT().GetNodeFlag(events.NodeAddedLatencyEnabled).Return(false, nil)

	ms.EXPECT().GetPodFlag(pod.Name, PCPRF).Return(scenario.Response_NORMAL, nil)
	ms.EXPECT().GetPodFlag(pod.Name, PCPRPF).Return(int32(100), nil)

	// sot
	var s store.Store = ms
	p := Provider{
		store: &s,
		pods:  NewPodManager(),
	}

	err := p.UpdatePod(context.TODO(), &pod)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, &pod, p.pods.GetPodByUID(pod.UID))
	ctrl.Finish()
}

func TestDeletePod(t *testing.T) {
	ctrl := gomock.NewController(t)
	ms := mock_store.NewMockStore(ctrl)

	/// vars
	pod := corev1.Pod{}
	pod.Name = podName
	pod.UID = types.UID(uuid.New().String())
	PCPRF := events.PodDeletePodResponse
	PCPRPF := events.PodDeletePodResponsePercentage

	// expect
	ms.EXPECT().GetNodeFlag(events.NodeAddedLatencyEnabled).Return(false, nil)

	ms.EXPECT().GetPodFlag(pod.Name, PCPRF).Return(scenario.Response_NORMAL, nil)
	ms.EXPECT().GetPodFlag(pod.Name, PCPRPF).Return(int32(100), nil)

	// sot
	var s store.Store = ms
	p := Provider{
		store: &s,
		pods: NewPodManager(),
	}

	err := p.DeletePod(context.TODO(), &pod)

	// assert
	assert.NoError(t, err)
	assert.NotContains(t, p.pods.uidToPod, &pod)
	ctrl.Finish()
}

func TestGetPod(t *testing.T) {
	ctrl := gomock.NewController(t)
	ms := mock_store.NewMockStore(ctrl)

	/// vars
	p := corev1.Pod{}
	p.Name = podName
	p.UID = types.UID(uuid.New().String())
	PCPRF := events.PodGetPodResponse
	PCPRPF := events.PodGetPodResponsePercentage

	// expect
	ms.EXPECT().GetNodeFlag(events.NodeAddedLatencyEnabled).Return(false, nil)

	ms.EXPECT().GetPodFlag(p.Name, PCPRF).Return(scenario.Response_NORMAL, nil)
	ms.EXPECT().GetPodFlag(p.Name, PCPRPF).Return(int32(100), nil)

	// sot
	var s store.Store = ms
	prov := Provider{
		store: &s,
		pods:  NewPodManager(),
	}

	prov.pods.AddPod(p)

	np, err := prov.GetPod(context.TODO(), "", p.Name)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, &p, np)
	ctrl.Finish()
}

func TestGetPods(t *testing.T) {
	ctrl := gomock.NewController(t)
	ms := mock_store.NewMockStore(ctrl)

	/// vars
	p := corev1.Pod{}
	p.Name = podName
	p.UID = types.UID(uuid.New().String())
	PCPRF := events.NodeGetPodsResponse
	PCPRPF := events.NodeGetPodsResponsePercentage

	// expect
	ms.EXPECT().GetNodeFlag(events.NodeAddedLatencyEnabled).Return(false, nil)

	ms.EXPECT().GetNodeFlag(PCPRF).Return(scenario.Response_NORMAL, nil)
	ms.EXPECT().GetNodeFlag(PCPRPF).Return(int32(100), nil)

	// sot
	var s store.Store = ms
	prov := Provider{
		store: &s,
		pods: NewPodManager(),
	}
	prov.pods.AddPod(p)

	ps, err := prov.GetPods(context.TODO())

	// assert
	assert.NoError(t, err)
	assert.Contains(t, ps, prov.pods.GetPodByUID(p.UID))
	ctrl.Finish()
}

func TestGetPodStatus100(t *testing.T) {
	ctrl := gomock.NewController(t)
	ms := mock_store.NewMockStore(ctrl)

	// vars
	p := corev1.Pod{}
	p.Name = podName
	p.UID = types.UID(uuid.New().String())
	PCPRF := events.PodGetPodStatusResponse
	PCPRPF := events.PodGetPodStatusResponsePercentage

	// expect
	ms.EXPECT().GetNodeFlag(events.NodeAddedLatencyEnabled).Return(false, nil)

	ms.EXPECT().GetPodFlag(p.Name, PCPRF).Return(scenario.Response_NORMAL, nil)
	ms.EXPECT().GetPodFlag(p.Name, PCPRPF).Return(int32(100), nil)

	ms.EXPECT().GetPodFlag(p.Name, events.PodUpdatePodStatus).Return(scenario.PodStatus_POD_SUCCEEDED, nil)
	ms.EXPECT().GetPodFlag(p.Name, events.PodUpdatePodStatusPercentage).Return(int32(100), nil)

	// sot
	var s store.Store = ms
	prov := Provider{
		store: &s,
		pods:  NewPodManager(),
	}
	prov.pods.AddPod(p)

	ps, err := prov.GetPodStatus(context.TODO(), "", p.Name)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, ps.Phase, corev1.PodSucceeded)
	ctrl.Finish()
}

func TestGetPodStatus0(t *testing.T) {
	ctrl := gomock.NewController(t)
	ms := mock_store.NewMockStore(ctrl)

	// vars
	p := corev1.Pod{}
	p.Name = podName
	p.UID = types.UID(uuid.New().String())
	PCPRF := events.PodGetPodStatusResponse
	PCPRPF := events.PodGetPodStatusResponsePercentage

	// expect
	ms.EXPECT().GetNodeFlag(events.NodeAddedLatencyEnabled).Return(false, nil)

	ms.EXPECT().GetPodFlag(p.Name, PCPRF).Return(scenario.Response_NORMAL, nil)
	ms.EXPECT().GetPodFlag(p.Name, PCPRPF).Return(int32(100), nil)

	ms.EXPECT().GetPodFlag(p.Name, events.PodUpdatePodStatus).Return(scenario.PodStatus_POD_SUCCEEDED, nil)
	ms.EXPECT().GetPodFlag(p.Name, events.PodUpdatePodStatusPercentage).Return(int32(0), nil)

	// sot
	var s store.Store = ms
	prov := Provider{
		store: &s,
		pods: NewPodManager(),
	}
	prov.pods.AddPod(p)

	ps, err := prov.GetPodStatus(context.TODO(), "", p.Name)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, ps.Phase, corev1.PodRunning)
	ctrl.Finish()
}
