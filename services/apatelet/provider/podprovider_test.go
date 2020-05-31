package provider

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/virtual-kubelet/node-cli/provider"
	"k8s.io/apimachinery/pkg/types"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario"

	podconfigv1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/podconfiguration/v1"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store/mock_store"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/provider/podmanager"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/provider/podmanager/mock_podmanager"
)

func TestGetPodLabelByPod(t *testing.T) {
	t.Parallel()

	pod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "TestNamespace",
			Labels: map[string]string{
				podconfigv1.PodConfigurationLabel: "TestLabel",
			},
		},
	}

	assert.Equal(t, "TestNamespace/TestLabel", getPodLabelByPod(&pod))
}

func TestGetPodLabelByPodApateNotFound(t *testing.T) {
	t.Parallel()

	pod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "TestNamespace",
			Labels: map[string]string{
				podconfigv1.PodConfigurationLabel + "xxx": "TestLabel",
			},
		},
	}

	assert.Equal(t, "", getPodLabelByPod(&pod))
}

func TestGetPodLabelByNameOk(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pmm := mock_podmanager.NewMockPodManager(ctrl)
	var pm podmanager.PodManager = pmm

	prov := Provider{
		Pods: pm,
	}
	name := "Apate"
	namespace := "TestNamespace"

	pod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "TestNamespace",
			Labels: map[string]string{
				podconfigv1.PodConfigurationLabel: "TestLabel",
			},
		},
	}

	pmm.EXPECT().GetPodByName(namespace, name).Return(&pod, true)

	res := prov.getPodLabelByName(namespace, name)

	assert.Equal(t, "TestNamespace/TestLabel", res)
}

func TestGetPodLabelByNameFail(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pmm := mock_podmanager.NewMockPodManager(ctrl)
	var pm podmanager.PodManager = pmm

	prov := Provider{
		Pods: pm,
	}
	name := "Apate"
	namespace := "TestNamespace"

	pmm.EXPECT().GetPodByName(namespace, name).Return(nil, false)

	res := prov.getPodLabelByName(namespace, name)

	assert.Equal(t, "", res)
}

func TestRunLatencyError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := mock_store.NewMockStore(ctrl)

	var s store.Store = ms

	p := Provider{
		Store: &s,
	}

	ms.EXPECT().GetNodeFlag(events.NodeAddedLatency).Return(time.Duration(0), errors.New("test error")).Times(6)

	assert.Error(t, p.UpdatePod(ctx, nil))
	assert.Error(t, p.CreatePod(ctx, nil))
	assert.Error(t, p.DeletePod(ctx, nil))
	_, err := p.GetPod(ctx, "", "")
	assert.Error(t, err)
	_, err = p.GetPodStatus(ctx, "", "")
	assert.Error(t, err)
	_, err = p.GetPods(ctx)
	assert.Error(t, err)
}

func TestCancelContextEarlyReturn(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := mock_store.NewMockStore(ctrl)

	var s store.Store = ms

	p := Provider{
		Store: &s,
	}

	assert.Error(t, p.UpdatePod(ctx, nil))
	assert.Error(t, p.CreatePod(ctx, nil))
	assert.Error(t, p.DeletePod(ctx, nil))
	_, err := p.GetPod(ctx, "", "")
	assert.Error(t, err)
	_, err = p.GetPodStatus(ctx, "", "")
	assert.Error(t, err)
	_, err = p.GetPods(ctx)
	assert.Error(t, err)
}

func TestCancelContextWhileRunningLatency(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := mock_store.NewMockStore(ctrl)

	var s store.Store = ms

	p := Provider{
		Store: &s,
	}

	ms.EXPECT().GetNodeFlag(events.NodeAddedLatency).Return(100000*time.Millisecond, nil).Times(6)

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	assert.Error(t, p.UpdatePod(ctx, nil))

	ctx, cancel = context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	assert.Error(t, p.CreatePod(ctx, nil))

	ctx, cancel = context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	assert.Error(t, p.DeletePod(ctx, nil))

	ctx, cancel = context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	_, err := p.GetPod(ctx, "", "")
	assert.Error(t, err)

	ctx, cancel = context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	_, err = p.GetPodStatus(ctx, "", "")
	assert.Error(t, err)

	ctx, cancel = context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	_, err = p.GetPods(ctx)
	assert.Error(t, err)
}

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
