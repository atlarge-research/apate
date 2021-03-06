package pod

import (
	"testing"

	"github.com/atlarge-research/apate/pkg/scenario/events"

	"github.com/atlarge-research/apate/pkg/scenario"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	podconfigv1 "github.com/atlarge-research/apate/pkg/apis/podconfiguration/v1"
	"github.com/atlarge-research/apate/services/apatelet/store"
	"github.com/atlarge-research/apate/services/apatelet/store/mock_store"
)

func TestTranslateResponse(t *testing.T) {
	t.Parallel()

	assert.Equal(t, scenario.ResponseNormal, translateResponse(podconfigv1.ResponseNormal))
	assert.Equal(t, scenario.ResponseError, translateResponse(podconfigv1.ResponseError))
	assert.Equal(t, scenario.ResponseTimeout, translateResponse(podconfigv1.ResponseTimeout))
	assert.Equal(t, scenario.ResponseUnset, translateResponse(podconfigv1.ResponseUnset))
	assert.Equal(t, scenario.ResponseUnset, translateResponse(podconfigv1.PodResponse(20)))
}

func TestTranslatePodStatus(t *testing.T) {
	t.Parallel()

	assert.Equal(t, scenario.PodStatusPending, translatePodStatus(podconfigv1.PodStatusPending))
	assert.Equal(t, scenario.PodStatusRunning, translatePodStatus(podconfigv1.PodStatusRunning))
	assert.Equal(t, scenario.PodStatusSucceeded, translatePodStatus(podconfigv1.PodStatusSucceeded))
	assert.Equal(t, scenario.PodStatusFailed, translatePodStatus(podconfigv1.PodStatusFailed))
	assert.Equal(t, scenario.PodStatusUnknown, translatePodStatus(podconfigv1.PodStatusUnknown))
	assert.Equal(t, scenario.PodStatusUnset, translatePodStatus(podconfigv1.PodStatusUnset))
	assert.Equal(t, scenario.PodStatusUnset, translatePodStatus(podconfigv1.PodStatus(20)))
}

func TestTranslatePodResources(t *testing.T) {
	t.Parallel()

	r, err := translatePodResources(&podconfigv1.PodResources{
		Memory:           "1B",
		CPU:              50,
		Storage:          "1B",
		EphemeralStorage: "1B",
	})

	assert.NoError(t, err)
	// All these asserts are necessary because time fields contain time.Now() which are impossible to compare.
	assert.Equal(t, uint64(50), r.UsageNanoCores)
	assert.Equal(t, uint64(1), r.UsageBytesMemory)
	assert.Equal(t, uint64(1), r.UsedBytesStorage)
	assert.Equal(t, uint64(1), r.UsedBytesEphemeral)
}

func TestTranslatePodResourcesErrorMemory(t *testing.T) {
	t.Parallel()

	_, err := translatePodResources(&podconfigv1.PodResources{
		Memory:           "-1B",
		CPU:              50,
		Storage:          "1B",
		EphemeralStorage: "1B",
	})
	assert.Error(t, err)
}

func TestTranslatePodResourcesErrorStorage(t *testing.T) {
	t.Parallel()

	_, err := translatePodResources(&podconfigv1.PodResources{
		Memory:           "1B",
		CPU:              50,
		Storage:          "-1B",
		EphemeralStorage: "1B",
	})
	assert.Error(t, err)
}

func TestTranslatePodResourcesErrorEphemeralStorage(t *testing.T) {
	t.Parallel()

	_, err := translatePodResources(&podconfigv1.PodResources{
		Memory:           "1B",
		CPU:              50,
		Storage:          "1B",
		EphemeralStorage: "-1B",
	})
	assert.Error(t, err)
}

func TestSetPodFlagsUnset(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mock_store.NewMockStore(ctrl)

	var s store.Store = ms

	ms.EXPECT().SetPodFlags("test", store.Flags{})

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
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mock_store.NewMockStore(ctrl)

	var s store.Store = ms

	ms.EXPECT().SetPodFlags("test", gomock.Any()).Do(func(_ string, flags store.Flags) {
		assert.Equal(t, translateResponse(podconfigv1.ResponseNormal), flags[events.PodCreatePodResponse])
		assert.Equal(t, translateResponse(podconfigv1.ResponseNormal), flags[events.PodUpdatePodResponse])
		assert.Equal(t, translateResponse(podconfigv1.ResponseNormal), flags[events.PodDeletePodResponse])
		assert.Equal(t, translateResponse(podconfigv1.ResponseNormal), flags[events.PodGetPodResponse])
		assert.Equal(t, translateResponse(podconfigv1.ResponseNormal), flags[events.PodGetPodStatusResponse])
		assert.NotNil(t, flags[events.PodResources])
		assert.Equal(t, translatePodStatus(podconfigv1.PodStatusRunning), flags[events.PodStatus])
	})

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
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mock_store.NewMockStore(ctrl)

	var s store.Store = ms

	ms.EXPECT().SetPodFlags(gomock.Any(), gomock.Any()).MinTimes(0)

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
