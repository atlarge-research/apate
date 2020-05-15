package crd

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/emulatedpod/v1"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store/mock_store"
)

func TestGetCRDAndLabel(t *testing.T) {
	ep := v1.EmulatedPod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "TestName",
			Namespace: "TestNamespace",
		},
	}

	rep, lbl := getCRDAndLabel(&ep)

	assert.Equal(t, lbl, "TestNamespace/TestName")
	assert.Equal(t, &ep, rep)
}

func TestEnqueueCRD(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := mock_store.NewMockStore(ctrl)

	var s store.Store = ms

	et := store.Task{
		RelativeTimestamp: 1,
		IsPod:             true,
		PodTask: &store.PodTask{
			Label: "TestNamespace/TestName",
			State: &v1.EmulatedPodState{},
		},
		NodeTask: nil,
	}

	ep := v1.EmulatedPod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "TestName",
			Namespace: "TestNamespace",
		},
		Spec: v1.EmulatedPodSpec{
			Tasks: []v1.EmulatedPodTask{
				{
					Timestamp: 1,
					State:     v1.EmulatedPodState{},
				},
			},
		},
	}

	ms.EXPECT().EnqueueCRDTasks(
		"TestNamespace/TestName",
		gomock.Any(),
	).Do(func(_ string, arr []*store.Task) {
		assert.Equal(t, 1, len(arr))
		assert.EqualValues(t, arr[0], &et)
	})

	err := enqueueCRD(&ep, &s)
	assert.NoError(t, err)
}

func TestEnqueueCRDDirect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := mock_store.NewMockStore(ctrl)

	var s store.Store = ms

	ep := v1.EmulatedPod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "TestName",
			Namespace: "TestNamespace",
		},
		Spec: v1.EmulatedPodSpec{
			DirectState: v1.EmulatedPodState{
				CreatePodResponse:    v1.ResponseNormal,
				UpdatePodResponse:    v1.ResponseNormal,
				DeletePodResponse:    v1.ResponseNormal,
				GetPodResponse:       v1.ResponseNormal,
				GetPodStatusResponse: v1.ResponseNormal,
				PodResources:         nil,
				PodStatus:            v1.PodStatusRunning,
			},
			Tasks: []v1.EmulatedPodTask{},
		},
	}

	gomock.InOrder(
		ms.EXPECT().SetPodFlag("TestNamespace/TestName", events.PodCreatePodResponse, translateResponse(v1.ResponseNormal)),
		ms.EXPECT().SetPodFlag("TestNamespace/TestName", events.PodUpdatePodResponse, translateResponse(v1.ResponseNormal)),
		ms.EXPECT().SetPodFlag("TestNamespace/TestName", events.PodDeletePodResponse, translateResponse(v1.ResponseNormal)),
		ms.EXPECT().SetPodFlag("TestNamespace/TestName", events.PodGetPodResponse, translateResponse(v1.ResponseNormal)),
		ms.EXPECT().SetPodFlag("TestNamespace/TestName", events.PodGetPodStatusResponse, translateResponse(v1.ResponseNormal)),
		ms.EXPECT().SetPodFlag("TestNamespace/TestName", events.PodStatus, translatePodStatus(v1.PodStatusRunning)),
	)

	ms.EXPECT().EnqueueCRDTasks(
		"TestNamespace/TestName",
		gomock.Any(),
	).Do(func(_ string, arr []*store.Task) {
		// Test if the array is empty when no spec tasks are given
		assert.Equal(t, 0, len(arr))
	})

	err := enqueueCRD(&ep, &s)
	assert.NoError(t, err)
}
