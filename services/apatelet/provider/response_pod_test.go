package provider

import (
	"context"
	"testing"
	"time"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario"

	"github.com/pkg/errors"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store/mock_store"
)

func setup(t *testing.T) (*mock_store.MockStore, *gomock.Controller, func(podFlag events.PodEventFlag) (interface{}, error)) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	// no defer ctrl.Finish() here because this function returns the ctrl

	ms := mock_store.NewMockStore(ctrl)
	var s store.Store = ms

	return ms, ctrl, func(podFlag events.PodEventFlag) (interface{}, error) {
		// Run code under test
		return podResponse(responseArgs{
			ctx:      context.Background(),
			provider: &Provider{Store: &s},
			action: func() (i interface{}, err error) {
				return tStr, nil
			}},
			podName,
			podFlag,
		)
	}
}

func TestPodNormal(t *testing.T) {
	ms, ctrl, executor := setup(t)
	defer ctrl.Finish()

	// Expectations
	ms.EXPECT().GetPodFlag(podName, events.PodCreatePodResponse).Return(scenario.ResponseNormal, nil)
	ms.EXPECT().GetNodeFlag(events.NodeCreatePodResponse).Return(scenario.ResponseUnset, nil)

	// Execute
	out, err := executor(events.PodCreatePodResponse)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, tStr, out)
}

func TestPodStoreError1(t *testing.T) {
	ms, ctrl, executor := setup(t)
	defer ctrl.Finish()

	// vars
	genericError := errors.New("some error")

	// Expectations
	ms.EXPECT().GetPodFlag(podName, events.PodCreatePodResponse).Return(nil, genericError)

	// Run code under test
	out, err := executor(events.PodCreatePodResponse)

	// Assert
	assert.Error(t, err)
	assert.False(t, IsExpected(err))
	assert.Nil(t, out)
}

func TestPodStoreError2(t *testing.T) {
	ms, ctrl, executor := setup(t)
	defer ctrl.Finish()

	// Expectations
	ms.EXPECT().GetPodFlag(podName, events.PodCreatePodResponse).Return(scenario.ResponseError, nil)
	ms.EXPECT().GetNodeFlag(events.NodeCreatePodResponse).Return(scenario.ResponseUnset, nil)

	// Run code under test
	out, err := executor(events.PodCreatePodResponse)

	// Assert
	assert.Error(t, err)
	assert.True(t, IsExpected(err))
	assert.Nil(t, out)
}

func TestPodUnset(t *testing.T) {
	ms, ctrl, executor := setup(t)
	defer ctrl.Finish()

	// Expectations
	ms.EXPECT().GetPodFlag(podName, events.PodCreatePodResponse).Return(scenario.ResponseUnset, nil)
	ms.EXPECT().GetNodeFlag(events.NodeCreatePodResponse).Return(scenario.ResponseUnset, nil)

	// Run code under test
	out, err := executor(events.PodCreatePodResponse)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, tStr, out)
}

func TestPodInvalidResponseType(t *testing.T) {
	ms, ctrl, executor := setup(t)
	defer ctrl.Finish()

	// Expectations
	ms.EXPECT().GetPodFlag(podName, events.PodCreatePodResponse).Return(42, nil)

	// Run code under test
	out, err := executor(events.PodCreatePodResponse)

	// Assert
	assert.Error(t, err)
	assert.False(t, IsExpected(err))
	assert.Nil(t, out)
}

func TestPodInvalidResponse(t *testing.T) {
	ms, ctrl, executor := setup(t)
	defer ctrl.Finish()

	// Expectations
	ms.EXPECT().GetPodFlag(podName, events.PodCreatePodResponse).Return(scenario.Response(42), nil)
	ms.EXPECT().GetNodeFlag(events.NodeCreatePodResponse).Return(scenario.ResponseUnset, nil)

	// Run code under test
	out, err := executor(events.PodCreatePodResponse)

	// Assert
	assert.NoError(t, err)
	assert.False(t, IsExpected(err))
	assert.Equal(t, tStr, out)
}

func TestPodTimeOut(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	ctrl, _ := gomock.WithContext(ctx, t)
	defer ctrl.Finish()

	ms := mock_store.NewMockStore(ctrl)
	var s store.Store = ms

	// Expectations
	ms.EXPECT().GetPodFlag(podName, events.PodCreatePodResponse).Return(scenario.ResponseTimeout, nil)
	ms.EXPECT().GetNodeFlag(events.NodeCreatePodResponse).Return(scenario.ResponseUnset, nil)

	assert.NoError(t, ctx.Err())

	// Run code under test
	out, err := podResponse(responseArgs{
		ctx:      ctx,
		provider: &Provider{Store: &s},
		action: func() (i interface{}, err error) {
			return tStr, nil
		},
	},
		podName,
		events.PodCreatePodResponse,
	)

	assert.Error(t, ctx.Err())

	// Assert
	assert.Nil(t, err)
	assert.Nil(t, out)
}

func TestTimeoutMostImportant(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	ctrl, ctx := gomock.WithContext(ctx, t)
	defer ctrl.Finish()

	ms := mock_store.NewMockStore(ctrl)
	var s store.Store = ms

	// vars
	podFlag := events.PodCreatePodResponse

	// Expectations
	ms.EXPECT().GetPodFlag(podName, podFlag).Return(scenario.ResponseError, nil)
	ms.EXPECT().GetNodeFlag(events.NodeCreatePodResponse).Return(scenario.ResponseTimeout, nil)

	assert.NoError(t, ctx.Err())

	// Run code under test
	out, err := podResponse(responseArgs{
		ctx:      ctx,
		provider: &Provider{Store: &s},
		action: func() (i interface{}, err error) {
			return tStr, nil
		},
	},
		podName,
		podFlag,
	)

	assert.Error(t, ctx.Err())

	// Assert
	assert.Nil(t, err)
	assert.Nil(t, out)
}

func TestErrorVsNormal(t *testing.T) {
	ms, ctrl, executor := setup(t)
	defer ctrl.Finish()

	// Expectations
	ms.EXPECT().GetPodFlag(podName, events.PodCreatePodResponse).Return(scenario.ResponseNormal, nil)
	ms.EXPECT().GetNodeFlag(events.NodeCreatePodResponse).Return(scenario.ResponseError, nil)

	// Run code under test
	out, err := executor(events.PodCreatePodResponse)

	// Assert
	assert.Error(t, err)
	assert.True(t, IsExpected(err))
	assert.Nil(t, out)
}
