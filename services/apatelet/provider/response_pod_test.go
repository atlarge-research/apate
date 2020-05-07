package provider

import (
	"context"
	"errors"
	"math/rand"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store/mock_store"
)

func TestMagicPodNormal100(t *testing.T) {
	ctrl := gomock.NewController(t)

	ms := mock_store.NewMockStore(ctrl)

	// vars
	PCPRF := events.PodCreatePodResponse
	PCPRPF := events.PodCreatePodResponsePercentage

	// Expectations
	ms.EXPECT().GetPodFlag(podName, PCPRF).Return(scenario.Response_NORMAL, nil)
	ms.EXPECT().GetPodFlag(podName, PCPRPF).Return(int32(100), nil)

	var s store.Store = ms

	// Run code under test
	out, err := podResponse(responseArgs{
		ctx: context.TODO(),
		p:   &Provider{store: &s},
		action: func() (i interface{}, err error) {
			return tStr, nil
		},
	},
		podResponseArgs{
			name:              podName,
			podResponseFlag:   PCPRF,
			podPercentageFlag: PCPRPF,
		})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, tStr, out)

	ctrl.Finish()
}

func TestMagicPodNormal0(t *testing.T) {
	ctrl := gomock.NewController(t)

	ms := mock_store.NewMockStore(ctrl)

	// vars
	PCPRF := events.PodCreatePodResponse
	PCPRPF := events.PodCreatePodResponsePercentage

	// Expectations
	ms.EXPECT().GetPodFlag(podName, PCPRF).Return(scenario.Response_NORMAL, nil)
	ms.EXPECT().GetPodFlag(podName, PCPRPF).Return(int32(0), nil)

	var s store.Store = ms

	// Run code under test
	out, err := podResponse(responseArgs{
		ctx: context.TODO(),
		p:   &Provider{store: &s},
		action: func() (i interface{}, err error) {
			return tStr, nil
		},
	},
		podResponseArgs{
			name:              podName,
			podResponseFlag:   PCPRF,
			podPercentageFlag: PCPRPF,
		})

	// Assert
	assert.NotNil(t, err)
	assert.EqualError(t, flagNotSetError, err.Error())
	assert.Nil(t, out)

	ctrl.Finish()
}

func TestMagicPodNormal50A(t *testing.T) {
	ctrl := gomock.NewController(t)

	ms := mock_store.NewMockStore(ctrl)

	// vars
	PCPRF := events.PodCreatePodResponse
	PCPRPF := events.PodCreatePodResponsePercentage

	rand.Seed(69)

	// Expectations
	ms.EXPECT().GetPodFlag(podName, PCPRF).Return(scenario.Response_ERROR, nil)
	ms.EXPECT().GetPodFlag(podName, PCPRPF).Return(int32(50), nil)

	var s store.Store = ms

	// Run code under test
	out, err := podResponse(responseArgs{
		ctx: context.TODO(),
		p:   &Provider{store: &s},
		action: func() (i interface{}, err error) {
			return tStr, nil
		},
	},
		podResponseArgs{
			name:              podName,
			podResponseFlag:   PCPRF,
			podPercentageFlag: PCPRPF,
		})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, tStr, out)

	ctrl.Finish()
}

func TestMagicPodNormal50B(t *testing.T) {
	ctrl := gomock.NewController(t)

	ms := mock_store.NewMockStore(ctrl)

	// vars
	PCPRF := events.PodCreatePodResponse
	PCPRPF := events.PodCreatePodResponsePercentage

	rand.Seed(42)

	// Expectations
	ms.EXPECT().GetPodFlag(podName, PCPRF).Return(scenario.Response_ERROR, nil)
	ms.EXPECT().GetPodFlag(podName, PCPRPF).Return(int32(50), nil)

	var s store.Store = ms

	// Run code under test
	out, err := podResponse(responseArgs{
		ctx: context.TODO(),
		p:   &Provider{store: &s},
		action: func() (i interface{}, err error) {
			return tStr, nil
		},
	},
		podResponseArgs{
			name:              podName,
			podResponseFlag:   PCPRF,
			podPercentageFlag: PCPRPF,
		})

	// Assert
	assert.NotNil(t, err)
	assert.EqualError(t, expectedError, err.Error())
	assert.Nil(t, out)

	ctrl.Finish()
}

func TestMagicPodStoreError1(t *testing.T) {
	ctrl := gomock.NewController(t)
	ms := mock_store.NewMockStore(ctrl)

	// vars
	PCPRF := events.PodCreatePodResponse
	PCPRPF := events.PodCreatePodResponsePercentage
	genericError := errors.New("some error")

	// Expectations
	ms.EXPECT().GetPodFlag(podName, PCPRF).Return(nil, genericError)

	var s store.Store = ms

	// Run code under test
	out, err := podResponse(responseArgs{
		ctx: context.TODO(),
		p:   &Provider{store: &s},
		action: func() (i interface{}, err error) {
			return tStr, nil
		},
	},
		podResponseArgs{
			name:              podName,
			podResponseFlag:   PCPRF,
			podPercentageFlag: PCPRPF,
		})

	// Assert
	assert.NotNil(t, err)
	assert.EqualError(t, genericError, err.Error())
	assert.Nil(t, out)

	ctrl.Finish()
}

func TestMagicPodStoreError2(t *testing.T) {
	ctrl := gomock.NewController(t)
	ms := mock_store.NewMockStore(ctrl)

	// vars
	PCPRF := events.PodCreatePodResponse
	PCPRPF := events.PodCreatePodResponsePercentage
	genericError := errors.New("some error")

	// Expectations
	ms.EXPECT().GetPodFlag(podName, PCPRF).Return(scenario.Response_ERROR, nil)
	ms.EXPECT().GetPodFlag(podName, PCPRPF).Return(nil, genericError)

	var s store.Store = ms

	// Run code under test
	out, err := podResponse(responseArgs{
		ctx: context.TODO(),
		p:   &Provider{store: &s},
		action: func() (i interface{}, err error) {
			return tStr, nil
		},
	},
		podResponseArgs{
			name:              podName,
			podResponseFlag:   PCPRF,
			podPercentageFlag: PCPRPF,
		})

	// Assert
	assert.NotNil(t, err)
	assert.EqualError(t, genericError, err.Error())
	assert.Nil(t, out)

	ctrl.Finish()
}

func TestMagicPodInvalidPercentage(t *testing.T) {
	ctrl := gomock.NewController(t)
	ms := mock_store.NewMockStore(ctrl)

	// vars
	PCPRF := events.PodCreatePodResponse
	PCPRPF := events.PodCreatePodResponsePercentage

	// Expectations
	ms.EXPECT().GetPodFlag(podName, PCPRF).Return(scenario.Response_ERROR, nil)
	ms.EXPECT().GetPodFlag(podName, PCPRPF).Return(nil, nil)

	var s store.Store = ms

	// Run code under test
	out, err := podResponse(responseArgs{
		ctx: context.TODO(),
		p:   &Provider{store: &s},
		action: func() (i interface{}, err error) {
			return tStr, nil
		},
	},
		podResponseArgs{
			name:              podName,
			podResponseFlag:   PCPRF,
			podPercentageFlag: PCPRPF,
		})

	// Assert
	assert.NotNil(t, err)
	assert.EqualError(t, invalidPercentage, err.Error())
	assert.Nil(t, out)

	ctrl.Finish()
}

func TestMagicPodInvalidResponseType(t *testing.T) {
	ctrl := gomock.NewController(t)
	ms := mock_store.NewMockStore(ctrl)

	// vars
	PCPRF := events.PodCreatePodResponse
	PCPRPF := events.PodCreatePodResponsePercentage

	// Expectations
	ms.EXPECT().GetPodFlag(podName, PCPRF).Return(42, nil)

	var s store.Store = ms

	// Run code under test
	out, err := podResponse(responseArgs{
		ctx: context.TODO(),
		p:   &Provider{store: &s},
		action: func() (i interface{}, err error) {
			return tStr, nil
		},
	},
		podResponseArgs{
			name:              podName,
			podResponseFlag:   PCPRF,
			podPercentageFlag: PCPRPF,
		})

	// Assert
	assert.NotNil(t, err)
	assert.EqualError(t, invalidFlag, err.Error())
	assert.Nil(t, out)

	ctrl.Finish()
}

func TestMagicPodInvalidResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	ms := mock_store.NewMockStore(ctrl)

	// vars
	PCPRF := events.PodCreatePodResponse
	PCPRPF := events.PodCreatePodResponsePercentage

	rand.Seed(42)

	// Expectations
	ms.EXPECT().GetPodFlag(podName, PCPRF).Return(scenario.Response(42), nil)
	ms.EXPECT().GetPodFlag(podName, PCPRPF).Return(int32(100), nil)

	var s store.Store = ms

	// Run code under test
	out, err := podResponse(responseArgs{
		ctx: context.TODO(),
		p:   &Provider{store: &s},
		action: func() (i interface{}, err error) {
			return tStr, nil
		},
	},
		podResponseArgs{
			name:              podName,
			podResponseFlag:   PCPRF,
			podPercentageFlag: PCPRPF,
		})

	// Assert
	assert.NotNil(t, err)
	assert.EqualError(t, invalidResponse, err.Error())
	assert.Nil(t, out)

	ctrl.Finish()
}

func TestMagicPodTimeOut(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3)
	defer cancel()
	ctrl, ctx := gomock.WithContext(ctx, t)
	ms := mock_store.NewMockStore(ctrl)

	// vars
	PCPRF := events.PodCreatePodResponse
	PCPRPF := events.PodCreatePodResponsePercentage

	rand.Seed(42)

	// Expectations
	ms.EXPECT().GetPodFlag(podName, PCPRF).Return(scenario.Response_TIMEOUT, nil)
	ms.EXPECT().GetPodFlag(podName, PCPRPF).Return(int32(100), nil)

	var s store.Store = ms

	// Run code under test
	out, err := podResponse(responseArgs{
		ctx: ctx,
		p:   &Provider{store: &s},
		action: func() (i interface{}, err error) {
			return tStr, nil
		},
	},
		podResponseArgs{
			name:              podName,
			podResponseFlag:   PCPRF,
			podPercentageFlag: PCPRPF,
		})

	// Assert
	assert.Nil(t, err)
	assert.Nil(t, out)

	ctrl.Finish()
}
