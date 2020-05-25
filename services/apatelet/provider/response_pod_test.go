package provider

import (
	"context"
	"math/rand"
	"testing"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario"

	"github.com/pkg/errors"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store/mock_store"
)

func TestPodNormal(t *testing.T) {
	ctrl := gomock.NewController(t)

	ms := mock_store.NewMockStore(ctrl)

	// vars
	PCPRF := events.PodCreatePodResponse

	// Expectations
	ms.EXPECT().GetPodFlag(podName, PCPRF).Return(scenario.ResponseNormal, nil)

	var s store.Store = ms

	// Run code under test
	out, changed, err := podResponse(responseArgs{
		ctx:      context.Background(),
		provider: &Provider{Store: &s},
		action: func() (i interface{}, err error) {
			return tStr, nil
		}},
		podName,
		PCPRF,
	)

	// Assert
	assert.NoError(t, err)
	assert.True(t, changed)
	assert.Equal(t, tStr, out)

	ctrl.Finish()
}

func TestPodStoreError1(t *testing.T) {
	ctrl := gomock.NewController(t)
	ms := mock_store.NewMockStore(ctrl)

	// vars
	PCPRF := events.PodCreatePodResponse
	genericError := errors.New("some error")

	// Expectations
	ms.EXPECT().GetPodFlag(podName, PCPRF).Return(nil, genericError)

	var s store.Store = ms

	// Run code under test
	out, changed, err := podResponse(responseArgs{
		ctx:      context.Background(),
		provider: &Provider{Store: &s},
		action: func() (i interface{}, err error) {
			return tStr, nil
		}},
		podName,
		PCPRF,
	)

	// Assert
	assert.Error(t, err)
	assert.False(t, changed)
	assert.Nil(t, out)

	ctrl.Finish()
}

func TestPodStoreError2(t *testing.T) {
	ctrl := gomock.NewController(t)
	ms := mock_store.NewMockStore(ctrl)

	// vars
	PCPRF := events.PodCreatePodResponse

	// Expectations
	ms.EXPECT().GetPodFlag(podName, PCPRF).Return(scenario.ResponseError, nil)

	var s store.Store = ms

	// Run code under test
	out, changed, err := podResponse(responseArgs{
		ctx:      context.Background(),
		provider: &Provider{Store: &s},
		action: func() (i interface{}, err error) {
			return tStr, nil
		}},
		podName,
		PCPRF,
	)

	// Assert
	assert.Error(t, err)
	assert.True(t, changed)
	assert.Nil(t, out)

	ctrl.Finish()
}

func TestPodUnset(t *testing.T) {
	ctrl := gomock.NewController(t)
	ms := mock_store.NewMockStore(ctrl)

	// vars
	PCPRF := events.PodCreatePodResponse

	// Expectations
	ms.EXPECT().GetPodFlag(podName, PCPRF).Return(scenario.ResponseUnset, nil)

	var s store.Store = ms

	// Run code under test
	out, changed, err := podResponse(responseArgs{
		ctx:      context.Background(),
		provider: &Provider{Store: &s},
		action: func() (i interface{}, err error) {
			return tStr, nil
		}},
		podName,
		PCPRF,
	)

	// Assert
	assert.NoError(t, err)
	assert.False(t, changed)
	assert.Nil(t, out)

	ctrl.Finish()
}

func TestPodInvalidResponseType(t *testing.T) {
	ctrl := gomock.NewController(t)
	ms := mock_store.NewMockStore(ctrl)

	// vars
	PCPRF := events.PodCreatePodResponse

	// Expectations
	ms.EXPECT().GetPodFlag(podName, PCPRF).Return(42, nil)

	var s store.Store = ms

	// Run code under test
	out, changed, err := podResponse(responseArgs{
		ctx:      context.Background(),
		provider: &Provider{Store: &s},
		action: func() (i interface{}, err error) {
			return tStr, nil
		}},
		podName,
		PCPRF,
	)

	// Assert
	assert.Error(t, err)
	assert.False(t, changed)
	assert.Nil(t, out)

	ctrl.Finish()
}

func TestPodInvalidResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	ms := mock_store.NewMockStore(ctrl)

	// vars
	PCPRF := events.PodCreatePodResponse

	rand.Seed(42)

	// Expectations
	ms.EXPECT().GetPodFlag(podName, PCPRF).Return(scenario.Response(42), nil)

	var s store.Store = ms

	// Run code under test
	out, changed, err := podResponse(responseArgs{
		ctx:      context.Background(),
		provider: &Provider{Store: &s},
		action: func() (i interface{}, err error) {
			return tStr, nil
		}},
		podName,
		PCPRF,
	)

	// Assert
	assert.False(t, changed)
	assert.Error(t, err)
	assert.Nil(t, out)

	ctrl.Finish()
}

func TestPodTimeOut(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3)
	defer cancel()
	ctrl, ctx := gomock.WithContext(ctx, t)
	ms := mock_store.NewMockStore(ctrl)

	// vars
	PCPRF := events.PodCreatePodResponse

	// Expectations
	ms.EXPECT().GetPodFlag(podName, PCPRF).Return(scenario.ResponseTimeout, nil)

	var s store.Store = ms

	// Run code under test
	out, changed, err := podResponse(responseArgs{
		ctx:      ctx,
		provider: &Provider{Store: &s},
		action: func() (i interface{}, err error) {
			return tStr, nil
		},
	},
		podName,
		PCPRF,
	)

	// Assert
	assert.True(t, changed)
	assert.Nil(t, err)
	assert.Nil(t, out)

	ctrl.Finish()
}
