package node

import (
	"testing"
	"time"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	v1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/nodeconfiguration/v1"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store/mock_store"
)

func TestTranslateResponse(t *testing.T) {
	assert.Equal(t, scenario.ResponseNormal, translateResponse(v1.ResponseNormal))
	assert.Equal(t, scenario.ResponseError, translateResponse(v1.ResponseError))
	assert.Equal(t, scenario.ResponseTimeout, translateResponse(v1.ResponseTimeout))
	assert.Equal(t, scenario.ResponseUnset, translateResponse(v1.ResponseUnset))
}

func TestSetNodeFlagsUnsetDirect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := mock_store.NewMockStore(ctrl)

	var s store.Store = ms

	SetNodeFlags(&s, &v1.NodeConfigurationState{
		NetworkLatency: "unset", // default in types.go
		CustomState: &v1.NodeConfigurationCustomState{
			CreatePodResponse:    v1.ResponseUnset,
			UpdatePodResponse:    v1.ResponseUnset,
			DeletePodResponse:    v1.ResponseUnset,
			GetPodResponse:       v1.ResponseUnset,
			GetPodsResponse:      v1.ResponseUnset,
			GetPodStatusResponse: v1.ResponseUnset,
			NodePingResponse:     v1.ResponseUnset,
		},
	})
}

func TestSetNodeFlagsDirect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := mock_store.NewMockStore(ctrl)

	var s store.Store = ms

	gomock.InOrder(
		ms.EXPECT().SetNodeFlag(events.NodeCreatePodResponse, translateResponse(v1.ResponseNormal)),
		ms.EXPECT().SetNodeFlag(events.NodeUpdatePodResponse, translateResponse(v1.ResponseError)),
		ms.EXPECT().SetNodeFlag(events.NodeDeletePodResponse, translateResponse(v1.ResponseNormal)),
		ms.EXPECT().SetNodeFlag(events.NodeGetPodResponse, translateResponse(v1.ResponseNormal)),
		ms.EXPECT().SetNodeFlag(events.NodeGetPodStatusResponse, translateResponse(v1.ResponseNormal)),
		ms.EXPECT().SetNodeFlag(events.NodeGetPodsResponse, translateResponse(v1.ResponseTimeout)),
		ms.EXPECT().SetNodeFlag(events.NodePingResponse, translateResponse(v1.ResponseNormal)),
	)

	SetNodeFlags(&s, &v1.NodeConfigurationState{
		NetworkLatency: "unset", // default in types.go
		CustomState: &v1.NodeConfigurationCustomState{
			CreatePodResponse:    v1.ResponseNormal,
			UpdatePodResponse:    v1.ResponseError,
			DeletePodResponse:    v1.ResponseNormal,
			GetPodResponse:       v1.ResponseNormal,
			GetPodsResponse:      v1.ResponseTimeout,
			GetPodStatusResponse: v1.ResponseNormal,
			NodePingResponse:     v1.ResponseNormal,
		},
	})
}

func TestSetNodeFlagsHeartbeat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := mock_store.NewMockStore(ctrl)

	var s store.Store = ms

	gomock.InOrder(
		ms.EXPECT().SetNodeFlag(events.NodePingResponse, translateResponse(v1.ResponseTimeout)),
	)

	SetNodeFlags(&s, &v1.NodeConfigurationState{
		NetworkLatency:  "unset", // default in types.go
		HeartbeatFailed: true,
	})
}

func TestSetNodeFlagsLatency(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := mock_store.NewMockStore(ctrl)

	var s store.Store = ms

	gomock.InOrder(
		ms.EXPECT().SetNodeFlag(events.NodePingResponse, translateResponse(v1.ResponseTimeout)),
		ms.EXPECT().SetNodeFlag(events.NodeAddedLatency, 100*time.Millisecond),
	)

	SetNodeFlags(&s, &v1.NodeConfigurationState{
		HeartbeatFailed: true,
		NetworkLatency:  "100ms",
	})
}

func TestSetNodeFlagsNodeFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := mock_store.NewMockStore(ctrl)

	var s store.Store = ms

	gomock.InOrder(
		ms.EXPECT().SetNodeFlag(events.NodeAddedLatency, 100*time.Millisecond),
		ms.EXPECT().SetNodeFlag(events.NodeCreatePodResponse, translateResponse(v1.ResponseTimeout)),
		ms.EXPECT().SetNodeFlag(events.NodeUpdatePodResponse, translateResponse(v1.ResponseTimeout)),
		ms.EXPECT().SetNodeFlag(events.NodeDeletePodResponse, translateResponse(v1.ResponseTimeout)),
		ms.EXPECT().SetNodeFlag(events.NodeGetPodResponse, translateResponse(v1.ResponseTimeout)),
		ms.EXPECT().SetNodeFlag(events.NodeGetPodStatusResponse, translateResponse(v1.ResponseTimeout)),
		ms.EXPECT().SetNodeFlag(events.NodeGetPodsResponse, translateResponse(v1.ResponseTimeout)),
		ms.EXPECT().SetNodeFlag(events.NodePingResponse, translateResponse(v1.ResponseTimeout)),
	)

	SetNodeFlags(&s, &v1.NodeConfigurationState{
		HeartbeatFailed: false,
		NetworkLatency:  "100ms",
		NodeFailed:      true,
	})
}
