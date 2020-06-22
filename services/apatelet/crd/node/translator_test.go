package node

import (
	"testing"
	"time"

	"github.com/atlarge-research/apate/pkg/scenario"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	nodeconfigv1 "github.com/atlarge-research/apate/pkg/apis/nodeconfiguration/v1"
	"github.com/atlarge-research/apate/pkg/scenario/events"
	"github.com/atlarge-research/apate/services/apatelet/store"
	"github.com/atlarge-research/apate/services/apatelet/store/mock_store"
)

func TestTranslateResponse(t *testing.T) {
	t.Parallel()

	assert.Equal(t, scenario.ResponseNormal, translateResponse(nodeconfigv1.ResponseNormal))
	assert.Equal(t, scenario.ResponseError, translateResponse(nodeconfigv1.ResponseError))
	assert.Equal(t, scenario.ResponseTimeout, translateResponse(nodeconfigv1.ResponseTimeout))
	assert.Equal(t, scenario.ResponseUnset, translateResponse(nodeconfigv1.ResponseUnset))
}

func TestSetNodeFlagsUnsetDirect(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mock_store.NewMockStore(ctrl)

	var s store.Store = ms

	ms.EXPECT().SetNodeFlags(store.Flags{})

	SetNodeFlags(&s, &nodeconfigv1.NodeConfigurationState{
		NetworkLatency: "unset", // default in types.go
		CustomState: &nodeconfigv1.NodeConfigurationCustomState{
			CreatePodResponse:    nodeconfigv1.ResponseUnset,
			UpdatePodResponse:    nodeconfigv1.ResponseUnset,
			DeletePodResponse:    nodeconfigv1.ResponseUnset,
			GetPodResponse:       nodeconfigv1.ResponseUnset,
			GetPodsResponse:      nodeconfigv1.ResponseUnset,
			GetPodStatusResponse: nodeconfigv1.ResponseUnset,
			NodePingResponse:     nodeconfigv1.ResponseUnset,
		},
	})
}

func TestSetNodeFlagsDirect(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mock_store.NewMockStore(ctrl)

	var s store.Store = ms

	ms.EXPECT().SetNodeFlags(store.Flags{
		events.NodeCreatePodResponse:    translateResponse(nodeconfigv1.ResponseNormal),
		events.NodeUpdatePodResponse:    translateResponse(nodeconfigv1.ResponseError),
		events.NodeDeletePodResponse:    translateResponse(nodeconfigv1.ResponseNormal),
		events.NodeGetPodResponse:       translateResponse(nodeconfigv1.ResponseNormal),
		events.NodeGetPodStatusResponse: translateResponse(nodeconfigv1.ResponseNormal),
		events.NodeGetPodsResponse:      translateResponse(nodeconfigv1.ResponseTimeout),
		events.NodePingResponse:         translateResponse(nodeconfigv1.ResponseNormal),
	})

	SetNodeFlags(&s, &nodeconfigv1.NodeConfigurationState{
		NetworkLatency: "unset", // default in types.go
		CustomState: &nodeconfigv1.NodeConfigurationCustomState{
			CreatePodResponse:    nodeconfigv1.ResponseNormal,
			UpdatePodResponse:    nodeconfigv1.ResponseError,
			DeletePodResponse:    nodeconfigv1.ResponseNormal,
			GetPodResponse:       nodeconfigv1.ResponseNormal,
			GetPodsResponse:      nodeconfigv1.ResponseTimeout,
			GetPodStatusResponse: nodeconfigv1.ResponseNormal,
			NodePingResponse:     nodeconfigv1.ResponseNormal,
		},
	})
}

func TestSetNodeFlagsHeartbeat(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mock_store.NewMockStore(ctrl)

	var s store.Store = ms

	ms.EXPECT().SetNodeFlags(store.Flags{
		events.NodePingResponse: translateResponse(nodeconfigv1.ResponseTimeout),
	})

	SetNodeFlags(&s, &nodeconfigv1.NodeConfigurationState{
		NetworkLatency:  "unset", // default in types.go
		HeartbeatFailed: true,
	})
}

func TestSetNodeFlagsLatency(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mock_store.NewMockStore(ctrl)

	var s store.Store = ms

	ms.EXPECT().SetNodeFlags(store.Flags{
		events.NodePingResponse: translateResponse(nodeconfigv1.ResponseTimeout),
		events.NodeAddedLatency: 100 * time.Millisecond,
	})

	SetNodeFlags(&s, &nodeconfigv1.NodeConfigurationState{
		HeartbeatFailed: true,
		NetworkLatency:  "100ms",
	})
}

func TestSetNodeFlagsNodeFailure(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mock_store.NewMockStore(ctrl)

	var s store.Store = ms

	ms.EXPECT().SetNodeFlags(store.Flags{
		events.NodeAddedLatency:         100 * time.Millisecond,
		events.NodeCreatePodResponse:    translateResponse(nodeconfigv1.ResponseTimeout),
		events.NodeUpdatePodResponse:    translateResponse(nodeconfigv1.ResponseTimeout),
		events.NodeDeletePodResponse:    translateResponse(nodeconfigv1.ResponseTimeout),
		events.NodeGetPodResponse:       translateResponse(nodeconfigv1.ResponseTimeout),
		events.NodeGetPodStatusResponse: translateResponse(nodeconfigv1.ResponseTimeout),
		events.NodeGetPodsResponse:      translateResponse(nodeconfigv1.ResponseTimeout),
		events.NodePingResponse:         translateResponse(nodeconfigv1.ResponseTimeout),
	})

	SetNodeFlags(&s, &nodeconfigv1.NodeConfigurationState{
		HeartbeatFailed: false,
		NetworkLatency:  "100ms",
		NodeFailed:      true,
	})
}
