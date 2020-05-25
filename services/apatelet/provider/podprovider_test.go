package provider

import (
	"context"
	"errors"
	"testing"
	"time"

	podconfigv1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/podconfiguration/v1"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store/mock_store"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/provider/podmanager"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/provider/podmanager/mock_podmanager"
)

func TestGetPodLabelByPod(t *testing.T) {
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
	pod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "TestNamespace",
			Labels: map[string]string{
				podconfigv1.PodConfigurationLabel + "xxx": "TestLabel",
			},
		},
	}

	assert.Equal(t, "TestNamespace/", getPodLabelByPod(&pod))
}

func TestGetPodLabelByNameOk(t *testing.T) {
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

func TestPodStatusToPhase(t *testing.T) {
	assert.Equal(t, corev1.PodPending, podStatusToPhase(scenario.PodStatusPending))
	assert.Equal(t, corev1.PodRunning, podStatusToPhase(scenario.PodStatusRunning))
	assert.Equal(t, corev1.PodSucceeded, podStatusToPhase(scenario.PodStatusSucceeded))
	assert.Equal(t, corev1.PodFailed, podStatusToPhase(scenario.PodStatusFailed))
	assert.Equal(t, corev1.PodUnknown, podStatusToPhase(scenario.PodStatusUnknown))
	assert.Equal(t, corev1.PodUnknown, podStatusToPhase(scenario.PodStatus(20)))
}

func TestRunLatencyError(t *testing.T) {
	ctx := context.TODO()

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
	ctx, cancel := context.WithCancel(context.TODO())
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := mock_store.NewMockStore(ctrl)

	var s store.Store = ms

	p := Provider{
		Store: &s,
	}

	ms.EXPECT().GetNodeFlag(events.NodeAddedLatency).Return(100000*time.Millisecond, nil).Times(6)

	ctx, cancel := context.WithTimeout(context.TODO(), 500*time.Millisecond)
	defer cancel()
	assert.Error(t, p.UpdatePod(ctx, nil))

	ctx, cancel = context.WithTimeout(context.TODO(), 500*time.Millisecond)
	defer cancel()
	assert.Error(t, p.CreatePod(ctx, nil))

	ctx, cancel = context.WithTimeout(context.TODO(), 500*time.Millisecond)
	defer cancel()
	assert.Error(t, p.DeletePod(ctx, nil))

	ctx, cancel = context.WithTimeout(context.TODO(), 500*time.Millisecond)
	defer cancel()
	_, err := p.GetPod(ctx, "", "")
	assert.Error(t, err)

	ctx, cancel = context.WithTimeout(context.TODO(), 500*time.Millisecond)
	defer cancel()
	_, err = p.GetPodStatus(ctx, "", "")
	assert.Error(t, err)

	ctx, cancel = context.WithTimeout(context.TODO(), 500*time.Millisecond)
	defer cancel()
	_, err = p.GetPods(ctx)
	assert.Error(t, err)
}
