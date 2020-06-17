package provider

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"testing"
	"time"

	"github.com/virtual-kubelet/virtual-kubelet/node/api"

	"github.com/finitum/node-cli/stats"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes/node"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/google/uuid"
	"k8s.io/apimachinery/pkg/types"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario"

	podconfigv1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/podconfiguration/v1"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store/mock_store"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"

	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/provider/podmanager"
)

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

	assert.Error(t, p.UpdatePod(ctx, &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "Test", Namespace: "Test"}}))
	assert.Error(t, p.CreatePod(ctx, &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "Test", Namespace: "Test"}}))
	assert.Error(t, p.DeletePod(ctx, &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "Test", Namespace: "Test"}}))
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

	assert.Error(t, p.UpdatePod(ctx, &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "Test", Namespace: "Test"}}))
	assert.Error(t, p.CreatePod(ctx, &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "Test", Namespace: "Test"}}))
	assert.Error(t, p.DeletePod(ctx, &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "Test", Namespace: "Test"}}))
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
	assert.Error(t, p.UpdatePod(ctx, &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "Test", Namespace: "Test"}}))

	ctx, cancel = context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	assert.Error(t, p.CreatePod(ctx, &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "Test", Namespace: "Test"}}))

	ctx, cancel = context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	assert.Error(t, p.DeletePod(ctx, &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "Test", Namespace: "Test"}}))

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

	// expect
	ms.EXPECT().GetNodeFlag(events.NodeAddedLatency).Return(time.Duration(0), nil)
	ms.EXPECT().GetPodFlag(&pod, events.PodCreatePodResponse).Return(scenario.ResponseNormal, nil)
	ms.EXPECT().GetPodFlag(&pod, events.PodResources).Return(&stats.PodStats{}, nil)
	ms.EXPECT().GetNodeFlag(events.NodeCreatePodResponse).Return(scenario.ResponseUnset, nil)

	// sot
	var s store.Store = ms

	p := Provider{
		Store:     &s,
		Pods:      podmanager.New(),
		NodeInfo:  &node.Info{},
		Resources: &scenario.NodeResources{},
		Stats:     NewStats(),
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

	// expect
	ms.EXPECT().GetNodeFlag(events.NodeAddedLatency).Return(time.Duration(0), nil)
	ms.EXPECT().GetPodFlag(&pod, events.PodUpdatePodResponse).Return(scenario.ResponseNormal, nil)
	ms.EXPECT().GetPodFlag(&pod, events.PodResources).Return(&stats.PodStats{}, nil)
	ms.EXPECT().GetNodeFlag(events.NodeUpdatePodResponse).Return(scenario.ResponseUnset, nil)

	// sot
	var s store.Store = ms
	p := Provider{
		Store:     &s,
		Pods:      podmanager.New(),
		NodeInfo:  &node.Info{},
		Resources: &scenario.NodeResources{},
		Stats:     NewStats(),
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

	// expect
	ms.EXPECT().GetNodeFlag(events.NodeAddedLatency).Return(time.Duration(0), nil)
	ms.EXPECT().GetPodFlag(&pod, events.PodDeletePodResponse).Return(scenario.ResponseNormal, nil)
	ms.EXPECT().GetNodeFlag(events.NodeDeletePodResponse).Return(scenario.ResponseUnset, nil)

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

	// expect
	ms.EXPECT().GetNodeFlag(events.NodeAddedLatency).Return(time.Duration(0), nil)
	ms.EXPECT().GetPodFlag(&pod, events.PodGetPodResponse).Return(scenario.ResponseNormal, nil)
	ms.EXPECT().GetNodeFlag(events.NodeGetPodResponse).Return(scenario.ResponseUnset, nil)

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

func TestGetContainerLogs(t *testing.T) {
	prov := Provider{}
	logs, err := prov.GetContainerLogs(context.Background(), "ns", "name", "", api.ContainerLogOpts{})
	assert.NoError(t, err)

	all, err := ioutil.ReadAll(logs)
	assert.NoError(t, err)

	assert.Equal(t, "This container is emulated by Apate\n", string(all))
}

type fakeAttachIO struct{}

func (f fakeAttachIO) Stdin() io.Reader {
	return nil
}

func (f fakeAttachIO) Stdout() io.WriteCloser {
	return nil
}

func (f fakeAttachIO) Stderr() io.WriteCloser {
	return nil
}

func (f fakeAttachIO) TTY() bool {
	return false
}

func (f fakeAttachIO) Resize() <-chan api.TermSize {
	return make(chan api.TermSize)
}

func TestRunInContainerNoError(t *testing.T) {
	prov := Provider{}
	err := prov.RunInContainer(context.Background(), "ns", "name", "", []string{}, fakeAttachIO{})
	assert.NoError(t, err)
}
