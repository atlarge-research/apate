package provider

import (
	"context"
	"testing"

	"github.com/atlarge-research/apate/pkg/scenario"

	"github.com/pkg/errors"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/atlarge-research/apate/pkg/scenario/events"
	"github.com/atlarge-research/apate/services/apatelet/store"
	"github.com/atlarge-research/apate/services/apatelet/store/mock_store"
)

const tStr = "test"

func TestNodeNormal(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mock_store.NewMockStore(ctrl)

	// Expectations
	ms.EXPECT().GetNodeFlag(events.NodeCreatePodResponse).Return(scenario.ResponseNormal, nil)

	var s store.Store = ms

	// Run code under test
	out, err := nodeResponse(responseArgs{
		ctx:      context.Background(),
		provider: &Provider{Store: &s},
		action: func() (i interface{}, err error) {
			return tStr, nil
		}},
		events.NodeCreatePodResponse,
	)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, tStr, out)
}

func TestNodeStoreError1(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mock_store.NewMockStore(ctrl)

	// vars
	genericError := errors.New("some error")

	// Expectations
	ms.EXPECT().GetNodeFlag(events.NodeCreatePodResponse).Return(nil, genericError)

	var s store.Store = ms

	// Run code under test
	out, err := nodeResponse(responseArgs{
		ctx:      context.Background(),
		provider: &Provider{Store: &s},
		action: func() (i interface{}, err error) {
			return tStr, nil
		}},
		events.NodeCreatePodResponse,
	)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, out)
}

func TestNodeErrorAction(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mock_store.NewMockStore(ctrl)

	// Expectations
	ms.EXPECT().GetNodeFlag(events.NodeCreatePodResponse).Return(scenario.ResponseNormal, nil)

	var s store.Store = ms

	// Run code under test
	out, err := nodeResponse(responseArgs{
		ctx:      context.Background(),
		provider: &Provider{Store: &s},
		action: func() (i interface{}, err error) {
			return nil, errors.New("some error")
		}},
		events.NodeCreatePodResponse,
	)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, out)
}

func TestNodeInvalidResponseType(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mock_store.NewMockStore(ctrl)

	// Expectations
	ms.EXPECT().GetNodeFlag(events.NodeCreatePodResponse).Return(42, nil)

	var s store.Store = ms

	// Run code under test
	out, err := nodeResponse(responseArgs{
		ctx:      context.Background(),
		provider: &Provider{Store: &s},
		action: func() (i interface{}, err error) {
			return tStr, nil
		}},
		events.NodeCreatePodResponse,
	)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, out)
}

func TestNodeInvalidResponse(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mock_store.NewMockStore(ctrl)

	// Expectations
	ms.EXPECT().GetNodeFlag(events.NodeCreatePodResponse).Return(scenario.Response(42), nil)

	var s store.Store = ms

	// Run code under test
	out, err := nodeResponse(responseArgs{
		ctx:      context.Background(),
		provider: &Provider{Store: &s},
		action: func() (i interface{}, err error) {
			return tStr, nil
		}},
		events.NodeCreatePodResponse,
	)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, tStr, out)
}

func TestNodeTimeOut(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 3)
	defer cancel()

	ctrl, ctx := gomock.WithContext(ctx, t)
	defer ctrl.Finish()

	ms := mock_store.NewMockStore(ctrl)

	// Expectations
	ms.EXPECT().GetNodeFlag(events.NodeCreatePodResponse).Return(scenario.ResponseTimeout, nil)

	var s store.Store = ms

	// Run code under test
	out, err := nodeResponse(responseArgs{
		ctx:      ctx,
		provider: &Provider{Store: &s},
		action: func() (i interface{}, err error) {
			return tStr, nil
		}},
		events.NodeCreatePodResponse,
	)

	// Assert
	assert.Nil(t, err)
	assert.Nil(t, out)
}
