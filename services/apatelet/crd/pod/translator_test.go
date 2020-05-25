package pod

import (
	"testing"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	podconfigv1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/podconfiguration/v1"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store/mock_store"
)

func TestTranslateResponse(t *testing.T) {
	assert.Equal(t, scenario.ResponseNormal, translateResponse(podconfigv1.ResponseNormal))
	assert.Equal(t, scenario.ResponseError, translateResponse(podconfigv1.ResponseError))
	assert.Equal(t, scenario.ResponseTimeout, translateResponse(podconfigv1.ResponseTimeout))
	assert.Equal(t, scenario.ResponseUnset, translateResponse(podconfigv1.ResponseUnset))
	assert.Equal(t, scenario.ResponseUnset, translateResponse(podconfigv1.PodResponse(20)))
}

func TestTranslatePodStatus(t *testing.T) {
	assert.Equal(t, scenario.PodStatusPending, translatePodStatus(podconfigv1.PodStatusPending))
	assert.Equal(t, scenario.PodStatusRunning, translatePodStatus(podconfigv1.PodStatusRunning))
	assert.Equal(t, scenario.PodStatusSucceeded, translatePodStatus(podconfigv1.PodStatusSucceeded))
	assert.Equal(t, scenario.PodStatusFailed, translatePodStatus(podconfigv1.PodStatusFailed))
	assert.Equal(t, scenario.PodStatusUnknown, translatePodStatus(podconfigv1.PodStatusUnknown))
	assert.Equal(t, scenario.PodStatusUnset, translatePodStatus(podconfigv1.PodStatusUnset))
	assert.Equal(t, scenario.PodStatusUnset, translatePodStatus(podconfigv1.PodStatus(20)))
}

func TestTranslatePodResources(t *testing.T) {
	r, err := translatePodResources(&podconfigv1.PodResources{
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
	_, err := translatePodResources(&podconfigv1.PodResources{
		Memory:           "-1B",
		CPU:              50,
		Storage:          "1B",
		EphemeralStorage: "1B",
	})
	assert.Error(t, err)
}

func TestTranslatePodResourcesErrorStorage(t *testing.T) {
	_, err := translatePodResources(&podconfigv1.PodResources{
		Memory:           "1B",
		CPU:              50,
		Storage:          "-1B",
		EphemeralStorage: "1B",
	})
	assert.Error(t, err)
}

func TestTranslatePodResourcesErrorEphemeralStorage(t *testing.T) {
	_, err := translatePodResources(&podconfigv1.PodResources{
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

	err := SetPodFlags(&s, "test", &podconfigv1.PodConfigurationState{
		CreatePodResponse:    podconfigv1.ResponseUnset,
		UpdatePodResponse:    podconfigv1.ResponseUnset,
		DeletePodResponse:    podconfigv1.ResponseUnset,
		GetPodResponse:       podconfigv1.ResponseUnset,
		GetPodStatusResponse: podconfigv1.ResponseUnset,
		PodResources:         nil,
		PodStatus:            podconfigv1.PodStatusUnset,
	})

	assert.NoError(t, err)
}

func TestSetPodFlags(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := mock_store.NewMockStore(ctrl)

	var s store.Store = ms

	gomock.InOrder(
		ms.EXPECT().SetPodFlag("test", events.PodCreatePodResponse, translateResponse(podconfigv1.ResponseNormal)),
		ms.EXPECT().SetPodFlag("test", events.PodUpdatePodResponse, translateResponse(podconfigv1.ResponseNormal)),
		ms.EXPECT().SetPodFlag("test", events.PodDeletePodResponse, translateResponse(podconfigv1.ResponseNormal)),
		ms.EXPECT().SetPodFlag("test", events.PodGetPodResponse, translateResponse(podconfigv1.ResponseNormal)),
		ms.EXPECT().SetPodFlag("test", events.PodGetPodStatusResponse, translateResponse(podconfigv1.ResponseNormal)),
		// Any here because the resource usage has times which can't be compared. This is tested better in TestTranslatePodResources
		ms.EXPECT().SetPodFlag("test", events.PodResources, gomock.Any()),
		ms.EXPECT().SetPodFlag("test", events.PodStatus, translatePodStatus(podconfigv1.PodStatusRunning)),
	)

	err := SetPodFlags(&s, "test", &podconfigv1.PodConfigurationState{
		CreatePodResponse:    podconfigv1.ResponseNormal,
		UpdatePodResponse:    podconfigv1.ResponseNormal,
		DeletePodResponse:    podconfigv1.ResponseNormal,
		GetPodResponse:       podconfigv1.ResponseNormal,
		GetPodStatusResponse: podconfigv1.ResponseNormal,
		PodResources: &podconfigv1.PodResources{
			Memory:           "1B",
			CPU:              1,
			Storage:          "1B",
			EphemeralStorage: "1B",
		},
		PodStatus: podconfigv1.PodStatusRunning,
	})

	assert.NoError(t, err)
}

func TestSetPodFlagsErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := mock_store.NewMockStore(ctrl)

	var s store.Store = ms

	ms.EXPECT().SetPodFlag(gomock.Any(), gomock.Any(), gomock.Any()).MinTimes(0)

	err := SetPodFlags(&s, "test", &podconfigv1.PodConfigurationState{
		CreatePodResponse:    podconfigv1.ResponseNormal,
		UpdatePodResponse:    podconfigv1.ResponseNormal,
		DeletePodResponse:    podconfigv1.ResponseNormal,
		GetPodResponse:       podconfigv1.ResponseNormal,
		GetPodStatusResponse: podconfigv1.ResponseNormal,
		PodResources: &podconfigv1.PodResources{
			Memory:           "-1B",
			CPU:              1,
			Storage:          "1B",
			EphemeralStorage: "1B",
		},
		PodStatus: podconfigv1.PodStatusRunning,
	})

	assert.Error(t, err)
}
