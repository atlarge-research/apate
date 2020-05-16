package pod

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario"
	v1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/emulatedpod/v1"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store/mock_store"
)

func TestTranslateResponse(t *testing.T) {
	assert.Equal(t, scenario.Response_RESPONSE_NORMAL, translateResponse(v1.ResponseNormal))
	assert.Equal(t, scenario.Response_RESPONSE_ERROR, translateResponse(v1.ResponseError))
	assert.Equal(t, scenario.Response_RESPONSE_TIMEOUT, translateResponse(v1.ResponseTimeout))
	assert.Equal(t, scenario.Response_RESPONSE_UNSET, translateResponse(v1.ResponseUnset))
	assert.Equal(t, scenario.Response_RESPONSE_UNSET, translateResponse(v1.EmulatedPodResponse(20)))
}

func TestTranslatePodStatus(t *testing.T) {
	assert.Equal(t, scenario.PodStatus_POD_STATUS_PENDING, translatePodStatus(v1.PodStatusPending))
	assert.Equal(t, scenario.PodStatus_POD_STATUS_RUNNING, translatePodStatus(v1.PodStatusRunning))
	assert.Equal(t, scenario.PodStatus_POD_STATUS_SUCCEEDED, translatePodStatus(v1.PodStatusSucceeded))
	assert.Equal(t, scenario.PodStatus_POD_STATUS_FAILED, translatePodStatus(v1.PodStatusFailed))
	assert.Equal(t, scenario.PodStatus_POD_STATUS_UNKNOWN, translatePodStatus(v1.PodStatusUnknown))
	assert.Equal(t, scenario.PodStatus_POD_STATUS_UNSET, translatePodStatus(v1.PodStatusUnset))
	assert.Equal(t, scenario.PodStatus_POD_STATUS_UNSET, translatePodStatus(v1.EmulatedPodStatus(20)))
}

func TestTranslatePodResources(t *testing.T) {
	r, err := translatePodResources(&v1.EmulatedPodResourceUsage{
		Memory:           "1B",
		CPU:              50,
		Storage:          "1B",
		EphemeralStorage: "1B",
	})

	assert.NoError(t, err)
	// All these asserts are necessary because time fields contain time.Now() which are impossible to compare.
	assert.Equal(t, uint64(50), *r.CPU.UsageNanoCores)
	assert.Equal(t, uint64(1), *r.Memory.UsageBytes)
	assert.Equal(t, 1, len(r.VolumeStats))
	assert.Equal(t, uint64(1), *r.VolumeStats[0].UsedBytes)
	assert.Equal(t, uint64(1), *r.EphemeralStorage.UsedBytes)
}

func TestTranslatePodResourcesErrorMemory(t *testing.T) {
	_, err := translatePodResources(&v1.EmulatedPodResourceUsage{
		Memory:           "-1B",
		CPU:              50,
		Storage:          "1B",
		EphemeralStorage: "1B",
	})
	assert.Error(t, err)
}

func TestTranslatePodResourcesErrorStorage(t *testing.T) {
	_, err := translatePodResources(&v1.EmulatedPodResourceUsage{
		Memory:           "1B",
		CPU:              50,
		Storage:          "-1B",
		EphemeralStorage: "1B",
	})
	assert.Error(t, err)
}

func TestTranslatePodResourcesErrorEphemeralStorage(t *testing.T) {
	_, err := translatePodResources(&v1.EmulatedPodResourceUsage{
		Memory:           "1B",
		CPU:              50,
		Storage:          "1B",
		EphemeralStorage: "-1B",
	})
	assert.Error(t, err)
}

func TestSetPodFlagsUnset(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := mock_store.NewMockStore(ctrl)

	var s store.Store = ms

	err := SetPodFlags(&s, "test", &v1.EmulatedPodState{
		CreatePodResponse:    v1.ResponseUnset,
		UpdatePodResponse:    v1.ResponseUnset,
		DeletePodResponse:    v1.ResponseUnset,
		GetPodResponse:       v1.ResponseUnset,
		GetPodStatusResponse: v1.ResponseUnset,
		PodResources:         nil,
		PodStatus:            v1.PodStatusUnset,
	})

	assert.NoError(t, err)
}

func TestSetPodFlags(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := mock_store.NewMockStore(ctrl)

	var s store.Store = ms

	gomock.InOrder(
		ms.EXPECT().SetPodFlag("test", events.PodCreatePodResponse, translateResponse(v1.ResponseNormal)),
		ms.EXPECT().SetPodFlag("test", events.PodUpdatePodResponse, translateResponse(v1.ResponseNormal)),
		ms.EXPECT().SetPodFlag("test", events.PodDeletePodResponse, translateResponse(v1.ResponseNormal)),
		ms.EXPECT().SetPodFlag("test", events.PodGetPodResponse, translateResponse(v1.ResponseNormal)),
		ms.EXPECT().SetPodFlag("test", events.PodGetPodStatusResponse, translateResponse(v1.ResponseNormal)),
		// Any here because the resource usage has times which can't be compared. This is tested better in TestTranslatePodResources
		ms.EXPECT().SetPodFlag("test", events.PodResources, gomock.Any()),
		ms.EXPECT().SetPodFlag("test", events.PodStatus, translatePodStatus(v1.PodStatusRunning)),
	)

	err := SetPodFlags(&s, "test", &v1.EmulatedPodState{
		CreatePodResponse:    v1.ResponseNormal,
		UpdatePodResponse:    v1.ResponseNormal,
		DeletePodResponse:    v1.ResponseNormal,
		GetPodResponse:       v1.ResponseNormal,
		GetPodStatusResponse: v1.ResponseNormal,
		PodResources: &v1.EmulatedPodResourceUsage{
			Memory:           "1B",
			CPU:              1,
			Storage:          "1B",
			EphemeralStorage: "1B",
		},
		PodStatus: v1.PodStatusRunning,
	})

	assert.NoError(t, err)
}

func TestSetPodFlagsErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := mock_store.NewMockStore(ctrl)

	var s store.Store = ms

	ms.EXPECT().SetPodFlag(gomock.Any(), gomock.Any(), gomock.Any()).MinTimes(0)

	err := SetPodFlags(&s, "test", &v1.EmulatedPodState{
		CreatePodResponse:    v1.ResponseNormal,
		UpdatePodResponse:    v1.ResponseNormal,
		DeletePodResponse:    v1.ResponseNormal,
		GetPodResponse:       v1.ResponseNormal,
		GetPodStatusResponse: v1.ResponseNormal,
		PodResources: &v1.EmulatedPodResourceUsage{
			Memory:           "-1B",
			CPU:              1,
			Storage:          "1B",
			EphemeralStorage: "1B",
		},
		PodStatus: v1.PodStatusRunning,
	})

	assert.Error(t, err)
}
