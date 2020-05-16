package provider

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/kind/pkg/errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario"
	v1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/emulatedpod/v1"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/crd/pod"
	mock_cache "github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/provider/mock_cache_store"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/provider/podmanager"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/provider/podmanager/mock_podmanager"
)

func TestGetPodLabelByPod(t *testing.T) {
	pod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "TestNamespace",
			Labels: map[string]string{
				"apate": "TestLabel",
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
				"apate1": "TestLabel",
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
		pods: pm,
	}
	name := "Apate"
	namespace := "TestNamespace"

	pod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "TestNamespace",
			Labels: map[string]string{
				"apate": "TestLabel",
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
		pods: pm,
	}
	name := "Apate"
	namespace := "TestNamespace"

	pmm.EXPECT().GetPodByName(namespace, name).Return(nil, false)

	res := prov.getPodLabelByName(namespace, name)

	assert.Equal(t, "", res)
}

func TestFindCrdErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cs := mock_cache.NewMockStore(ctrl)
	var c cache.Store = cs

	prov := Provider{
		crdInformer: pod.NewInformer(&c),
	}

	cs.EXPECT().GetByKey("TestNamespace/TestLabel").Return(nil, false, errors.New("TestError"))

	pod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "TestNamespace",
			Labels: map[string]string{
				"apate": "TestLabel",
			},
		},
	}

	_, err := prov.findCRD(&pod)
	assert.Error(t, err)
}

func TestFindCrd(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cs := mock_cache.NewMockStore(ctrl)
	var c cache.Store = cs

	prov := Provider{
		crdInformer: pod.NewInformer(&c),
	}

	p := v1.EmulatedPod{}

	cs.EXPECT().GetByKey("TestNamespace/TestLabel").
		Return(&p, true, nil)

	pod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "TestNamespace",
			Labels: map[string]string{
				"apate": "TestLabel",
			},
		},
	}

	res, err := prov.findCRD(&pod)
	assert.NoError(t, err)
	assert.Equal(t, res, &p)
}

func TestPodStatusToPhase(t *testing.T) {
	assert.Equal(t, corev1.PodPending, podStatusToPhase(scenario.PodStatus_POD_STATUS_PENDING))
	assert.Equal(t, corev1.PodRunning, podStatusToPhase(scenario.PodStatus_POD_STATUS_RUNNING))
	assert.Equal(t, corev1.PodRunning, podStatusToPhase(scenario.PodStatus_POD_STATUS_UNSET))
	assert.Equal(t, corev1.PodSucceeded, podStatusToPhase(scenario.PodStatus_POD_STATUS_SUCCEEDED))
	assert.Equal(t, corev1.PodFailed, podStatusToPhase(scenario.PodStatus_POD_STATUS_FAILED))
	assert.Equal(t, corev1.PodUnknown, podStatusToPhase(scenario.PodStatus_POD_STATUS_UNKNOWN))
	assert.Equal(t, corev1.PodUnknown, podStatusToPhase(scenario.PodStatus(20)))
}
