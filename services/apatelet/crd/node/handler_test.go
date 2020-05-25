package node

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	nodeconfigv1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/nodeconfiguration/v1"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store/mock_store"
)

func TestEnqueueNodeTasks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := mock_store.NewMockStore(ctrl)

	var s store.Store = ms

	et1 := store.NewNodeTask(
		1*time.Millisecond,
		&nodeconfigv1.NodeConfigurationState{
			CustomState: &nodeconfigv1.NodeConfigurationCustomState{
				CreatePodResponse: nodeconfigv1.ResponseNormal,
			},
		})

	et2 := store.NewNodeTask(
		42*time.Millisecond,
		&nodeconfigv1.NodeConfigurationState{
			CustomState: &nodeconfigv1.NodeConfigurationCustomState{
				CreatePodResponse: nodeconfigv1.ResponseTimeout,
			},
		})

	ep := nodeconfigv1.NodeConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "TestName",
			Namespace: "TestNamespace",
		},
		Spec: nodeconfigv1.NodeConfigurationSpec{
			Tasks: []nodeconfigv1.NodeConfigurationTask{
				{
					Timestamp: 1,
					State: nodeconfigv1.NodeConfigurationState{
						CustomState: &nodeconfigv1.NodeConfigurationCustomState{
							CreatePodResponse: nodeconfigv1.ResponseNormal,
						},
					},
				},
				{
					Timestamp: 42,
					State: nodeconfigv1.NodeConfigurationState{
						CustomState: &nodeconfigv1.NodeConfigurationCustomState{
							CreatePodResponse: nodeconfigv1.ResponseTimeout,
						},
					},
				},
			},
		},
	}

	ms.EXPECT().SetNodeTasks(gomock.Any()).Do(func(arr []*store.Task) {
		assert.Equal(t, 2, len(arr))
		assert.EqualValues(t, et1, arr[0])
		assert.EqualValues(t, et2, arr[1])
	})

	err := setNodeTasks(&ep, &s)
	assert.NoError(t, err)
}

func TestEnqueueCRDDirect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := mock_store.NewMockStore(ctrl)

	var s store.Store = ms

	ep := nodeconfigv1.NodeConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "TestName",
			Namespace: "TestNamespace",
		},
		Spec: nodeconfigv1.NodeConfigurationSpec{
			NodeConfigurationState: nodeconfigv1.NodeConfigurationState{
				NetworkLatency: "unset",
				NodeFailed:     true,
			},
			Tasks: []nodeconfigv1.NodeConfigurationTask{},
		},
	}

	gomock.InOrder(
		ms.EXPECT().SetNodeFlag(events.NodeCreatePodResponse, translateResponse(nodeconfigv1.ResponseTimeout)),
		ms.EXPECT().SetNodeFlag(events.NodeUpdatePodResponse, translateResponse(nodeconfigv1.ResponseTimeout)),
		ms.EXPECT().SetNodeFlag(events.NodeDeletePodResponse, translateResponse(nodeconfigv1.ResponseTimeout)),
		ms.EXPECT().SetNodeFlag(events.NodeGetPodResponse, translateResponse(nodeconfigv1.ResponseTimeout)),
		ms.EXPECT().SetNodeFlag(events.NodeGetPodStatusResponse, translateResponse(nodeconfigv1.ResponseTimeout)),
		ms.EXPECT().SetNodeFlag(events.NodeGetPodsResponse, translateResponse(nodeconfigv1.ResponseTimeout)),
		ms.EXPECT().SetNodeFlag(events.NodePingResponse, translateResponse(nodeconfigv1.ResponseTimeout)),
	)

	ms.EXPECT().SetNodeTasks(gomock.Any()).Do(func(arr []*store.Task) {
		// Test if the array is empty when no spec tasks are given
		assert.Equal(t, 0, len(arr))
	})

	err := setNodeTasks(&ep, &s)
	assert.NoError(t, err)
}
