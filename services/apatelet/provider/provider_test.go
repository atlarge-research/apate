package provider

import (
	"context"
	"testing"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario"

	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/provider/podmanager"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/virtual-kubelet/node-cli/provider"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store/mock_store"
)

const podNamespace = "podnamespace"
const podName = "pod"
const podLabel = "label"

func TestConfigureNode(t *testing.T) {
	resources := scenario.NodeResources{
		UUID:    uuid.New(),
		Memory:  42,
		CPU:     1337,
		MaxPods: 1001,
	}

	prov := Provider{
		pods:      podmanager.New(),
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
	resources := scenario.NodeResources{
		UUID:    uuid.New(),
		Memory:  42,
		CPU:     1337,
		MaxPods: 1001,
	}

	ctrl := gomock.NewController(t)
	st := store.NewStore()

	prov := NewProvider(podmanager.New(), NewStats(), &resources, provider.InitConfig{}, cluster.NodeInfo{}, &st)

	fakeNode := corev1.Node{}

	// Run the method
	prov.ConfigureNode(context.TODO(), &fakeNode)

	assert.EqualValues(t, resources.CPU, fakeNode.Status.Capacity.Cpu().Value())
	assert.EqualValues(t, resources.Memory, fakeNode.Status.Capacity.Memory().Value())
	assert.EqualValues(t, resources.MaxPods, fakeNode.Status.Capacity.Pods().Value())

	ctrl.Finish()
}

func TestCreatePod(t *testing.T) {
	ctrl := gomock.NewController(t)
	ms := mock_store.NewMockStore(ctrl)

	/// vars
	pod := corev1.Pod{}
	pod.Namespace = podNamespace
	pod.Name = podName
	pod.Labels = map[string]string{
		"apate": podLabel,
	}
	pod.UID = types.UID(uuid.New().String())
	PCPRF := events.PodCreatePodResponse

	// expect
	ms.EXPECT().GetNodeFlag(events.NodeAddedLatencyMsec).Return(int64(0), nil)
	ms.EXPECT().GetPodFlag(podNamespace+"/"+podLabel, PCPRF).Return(scenario.ResponseNormal, nil)

	// sot
	var s store.Store = ms

	p := Provider{
		store: &s,
		pods:  podmanager.New(),
	}

	err := p.CreatePod(context.TODO(), &pod)

	// assert
	assert.NoError(t, err)

	uid, ok := p.pods.GetPodByUID(pod.UID)
	assert.True(t, ok)
	assert.Equal(t, &pod, uid)
	ctrl.Finish()
}

func TestUpdatePod(t *testing.T) {
	ctrl := gomock.NewController(t)
	ms := mock_store.NewMockStore(ctrl)

	/// vars
	pod := corev1.Pod{}
	pod.Namespace = podNamespace
	pod.Name = podName
	pod.Labels = map[string]string{
		"apate": podLabel,
	}
	pod.UID = types.UID(uuid.New().String())
	PCPRF := events.PodUpdatePodResponse

	// expect
	ms.EXPECT().GetNodeFlag(events.NodeAddedLatencyMsec).Return(int64(0), nil)
	ms.EXPECT().GetPodFlag(podNamespace+"/"+podLabel, PCPRF).Return(scenario.ResponseNormal, nil)

	// sot
	var s store.Store = ms
	p := Provider{
		store: &s,
		pods:  podmanager.New(),
	}

	err := p.UpdatePod(context.TODO(), &pod)

	// assert
	assert.NoError(t, err)
	uid, ok := p.pods.GetPodByUID(pod.UID)
	assert.True(t, ok)
	assert.Equal(t, &pod, uid)
	ctrl.Finish()
}

func TestDeletePod(t *testing.T) {
	ctrl := gomock.NewController(t)
	ms := mock_store.NewMockStore(ctrl)

	/// vars
	pod := corev1.Pod{}
	pod.Namespace = podNamespace
	pod.Name = podName
	pod.Labels = map[string]string{
		"apate": podLabel,
	}
	pod.UID = types.UID(uuid.New().String())
	PCPRF := events.PodDeletePodResponse

	// expect
	ms.EXPECT().GetNodeFlag(events.NodeAddedLatencyMsec).Return(int64(0), nil)
	ms.EXPECT().GetPodFlag(podNamespace+"/"+podLabel, PCPRF).Return(scenario.ResponseNormal, nil)

	// sot
	var s store.Store = ms
	p := Provider{
		store: &s,
		pods:  podmanager.New(),
	}

	err := p.DeletePod(context.TODO(), &pod)

	// assert
	assert.NoError(t, err)
	assert.NotContains(t, p.pods.GetAllPods(), &pod)
	ctrl.Finish()
}

func TestGetPod(t *testing.T) {
	ctrl := gomock.NewController(t)
	ms := mock_store.NewMockStore(ctrl)

	/// vars
	pod := corev1.Pod{}
	pod.Namespace = podNamespace
	pod.Name = podName
	pod.Labels = map[string]string{
		"apate": podLabel,
	}
	pod.UID = types.UID(uuid.New().String())
	PCPRF := events.PodGetPodResponse

	// expect
	ms.EXPECT().GetNodeFlag(events.NodeAddedLatencyMsec).Return(int64(0), nil)
	ms.EXPECT().GetPodFlag(podNamespace+"/"+podLabel, PCPRF).Return(scenario.ResponseNormal, nil)

	// sot
	var s store.Store = ms
	prov := Provider{
		store: &s,
		pods:  podmanager.New(),
	}

	prov.pods.AddPod(pod)

	np, err := prov.GetPod(context.TODO(), podNamespace, podName)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, &pod, np)
	ctrl.Finish()
}

func TestGetPods(t *testing.T) {
	ctrl := gomock.NewController(t)
	ms := mock_store.NewMockStore(ctrl)

	/// vars
	pod := corev1.Pod{}
	pod.Namespace = podNamespace
	pod.Name = podName
	pod.Labels = map[string]string{
		"apate": podLabel,
	}
	pod.UID = types.UID(uuid.New().String())
	PCPRF := events.NodeGetPodsResponse

	// expect
	ms.EXPECT().GetNodeFlag(PCPRF).Return(scenario.ResponseNormal, nil)
	ms.EXPECT().GetNodeFlag(events.NodeAddedLatencyMsec).Return(int64(0), nil)

	// sot
	var s store.Store = ms
	prov := Provider{
		store: &s,
		pods:  podmanager.New(),
	}
	prov.pods.AddPod(pod)

	ps, err := prov.GetPods(context.TODO())

	// assert
	assert.NoError(t, err)
	uid, ok := prov.pods.GetPodByUID(pod.UID)
	assert.True(t, ok)
	assert.Contains(t, ps, uid)
	ctrl.Finish()
}

func TestGetPodStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	ms := mock_store.NewMockStore(ctrl)

	// vars
	pod := corev1.Pod{}
	pod.Namespace = podNamespace
	pod.Name = podName
	pod.Labels = map[string]string{
		"apate": podLabel,
	}
	pod.UID = types.UID(uuid.New().String())
	PCPRF := events.PodGetPodStatusResponse

	// expect
	ms.EXPECT().GetNodeFlag(events.NodeAddedLatencyMsec).Return(int64(0), nil)
	ms.EXPECT().GetPodFlag(podNamespace+"/"+podLabel, PCPRF).Return(scenario.ResponseNormal, nil)
	ms.EXPECT().GetPodFlag(podNamespace+"/"+podLabel, events.PodStatus).Return(scenario.PodStatusSucceeded, nil)

	// sot
	var s store.Store = ms
	prov := Provider{
		store: &s,
		pods:  podmanager.New(),
	}
	prov.pods.AddPod(pod)

	ps, err := prov.GetPodStatus(context.TODO(), podNamespace, podName)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, ps.Phase, corev1.PodSucceeded)
	ctrl.Finish()
}
