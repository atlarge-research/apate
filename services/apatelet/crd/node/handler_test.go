package node

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/nodeconfiguration/v1"
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
		&v1.NodeConfigurationState{
			CustomState: &v1.NodeConfigurationCustomState{
				CreatePodResponse: v1.ResponseNormal,
			},
		})

	et2 := store.NewNodeTask(
		42*time.Millisecond,
		&v1.NodeConfigurationState{
			CustomState: &v1.NodeConfigurationCustomState{
				CreatePodResponse: v1.ResponseTimeout,
			},
		})

	ep := v1.NodeConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "TestName",
			Namespace: "TestNamespace",
		},
		Spec: v1.NodeConfigurationSpec{
			Tasks: []v1.NodeConfigurationTask{
				{
					Timestamp: 1,
					State: v1.NodeConfigurationState{
						CustomState: &v1.NodeConfigurationCustomState{
							CreatePodResponse: v1.ResponseNormal,
						},
					},
				},
				{
					Timestamp: 42,
					State: v1.NodeConfigurationState{
						CustomState: &v1.NodeConfigurationCustomState{
							CreatePodResponse: v1.ResponseTimeout,
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

	ep := v1.NodeConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "TestName",
			Namespace: "TestNamespace",
		},
		Spec: v1.NodeConfigurationSpec{
			NodeConfigurationState: v1.NodeConfigurationState{
				NetworkLatency: -1,
				NodeFailed:     true,
			},
			Tasks: []v1.NodeConfigurationTask{},
		},
	}

	gomock.InOrder(
		ms.EXPECT().SetNodeFlag(events.NodeCreatePodResponse, translateResponse(v1.ResponseTimeout)),
		ms.EXPECT().SetNodeFlag(events.NodeUpdatePodResponse, translateResponse(v1.ResponseTimeout)),
		ms.EXPECT().SetNodeFlag(events.NodeDeletePodResponse, translateResponse(v1.ResponseTimeout)),
		ms.EXPECT().SetNodeFlag(events.NodeGetPodResponse, translateResponse(v1.ResponseTimeout)),
		ms.EXPECT().SetNodeFlag(events.NodeGetPodStatusResponse, translateResponse(v1.ResponseTimeout)),
		ms.EXPECT().SetNodeFlag(events.NodeGetPodsResponse, translateResponse(v1.ResponseTimeout)),
		ms.EXPECT().SetNodeFlag(events.NodePingResponse, translateResponse(v1.ResponseTimeout)),
	)

	ms.EXPECT().SetNodeTasks(gomock.Any()).Do(func(arr []*store.Task) {
		// Test if the array is empty when no spec tasks are given
		assert.Equal(t, 0, len(arr))
	})

	err := setNodeTasks(&ep, &s)
	assert.NoError(t, err)
}
