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

const tStr = "test"

func TestNodeNormal100(t *testing.T) {
	ctrl := gomock.NewController(t)

	ms := mock_store.NewMockStore(ctrl)

	// vars
	PCPRF := events.NodeCreatePodResponse
	PCPRPF := events.NodeCreatePodResponsePercentage

	// Expectations
	ms.EXPECT().GetNodeFlag(PCPRF).Return(scenario.Response_NORMAL, nil)
	ms.EXPECT().GetNodeFlag(PCPRPF).Return(int32(100), nil)

	var s store.Store = ms

	// Run code under test
	out, err := nodeResponse(responseArgs{
		ctx:      context.TODO(),
		provider: &Provider{store: &s},
		action: func() (i interface{}, err error) {
			return tStr, nil
		},
	},
		nodeResponseArgs{
			nodeResponseFlag:   PCPRF,
			nodePercentageFlag: PCPRPF,
		})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, tStr, out)

	ctrl.Finish()
}

func TestNodeNormal0(t *testing.T) {
	ctrl := gomock.NewController(t)

	ms := mock_store.NewMockStore(ctrl)

	// vars
	PCPRF := events.NodeCreatePodResponse
	PCPRPF := events.NodeCreatePodResponsePercentage

	// Expectations
	ms.EXPECT().GetNodeFlag(PCPRF).Return(scenario.Response_NORMAL, nil)
	ms.EXPECT().GetNodeFlag(PCPRPF).Return(int32(0), nil)

	var s store.Store = ms

	// Run code under test
	out, err := nodeResponse(responseArgs{
		ctx:      context.TODO(),
		provider: &Provider{store: &s},
		action: func() (i interface{}, err error) {
			return tStr, nil
		},
	},
		nodeResponseArgs{
			nodeResponseFlag:   PCPRF,
			nodePercentageFlag: PCPRPF,
		})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, tStr, out)

	ctrl.Finish()
}

func TestNodeNormal50A(t *testing.T) {
	ctrl := gomock.NewController(t)

	ms := mock_store.NewMockStore(ctrl)

	// vars
	PCPRF := events.NodeCreatePodResponse
	PCPRPF := events.NodeCreatePodResponsePercentage

	rand.Seed(69)

	// Expectations
	ms.EXPECT().GetNodeFlag(PCPRF).Return(scenario.Response_ERROR, nil)
	ms.EXPECT().GetNodeFlag(PCPRPF).Return(int32(50), nil)

	var s store.Store = ms

	// Run code under test
	out, err := nodeResponse(responseArgs{
		ctx:      context.TODO(),
		provider: &Provider{store: &s},
		action: func() (i interface{}, err error) {
			return tStr, nil
		},
	},
		nodeResponseArgs{
			nodeResponseFlag:   PCPRF,
			nodePercentageFlag: PCPRPF,
		})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, tStr, out)

	ctrl.Finish()
}

func TestNodeNormal50B(t *testing.T) {
	ctrl := gomock.NewController(t)

	ms := mock_store.NewMockStore(ctrl)

	// vars
	PCPRF := events.NodeCreatePodResponse
	PCPRPF := events.NodeCreatePodResponsePercentage

	rand.Seed(42)

	// Expectations
	ms.EXPECT().GetNodeFlag(PCPRF).Return(scenario.Response_ERROR, nil)
	ms.EXPECT().GetNodeFlag(PCPRPF).Return(int32(50), nil)

	var s store.Store = ms

	// Run code under test
	out, err := nodeResponse(responseArgs{
		ctx:      context.TODO(),
		provider: &Provider{store: &s},
		action: func() (i interface{}, err error) {
			return tStr, nil
		},
	},
		nodeResponseArgs{
			nodeResponseFlag:   PCPRF,
			nodePercentageFlag: PCPRPF,
		})

	// Assert
	assert.NotNil(t, err)
	assert.EqualError(t, expectedError, err.Error())
	assert.Nil(t, out)

	ctrl.Finish()
}

func TestNodeStoreError1(t *testing.T) {
	ctrl := gomock.NewController(t)
	ms := mock_store.NewMockStore(ctrl)

	// vars
	PCPRF := events.NodeCreatePodResponse
	PCPRPF := events.NodeCreatePodResponsePercentage
	genericError := errors.New("some error")

	// Expectations
	ms.EXPECT().GetNodeFlag(PCPRF).Return(nil, genericError)

	var s store.Store = ms

	// Run code under test
	out, err := nodeResponse(responseArgs{
		ctx:      context.TODO(),
		provider: &Provider{store: &s},
		action: func() (i interface{}, err error) {
			return tStr, nil
		},
	},
		nodeResponseArgs{
			nodeResponseFlag:   PCPRF,
			nodePercentageFlag: PCPRPF,
		})

	// Assert
	assert.NotNil(t, err)
	assert.EqualError(t, genericError, err.Error())
	assert.Nil(t, out)

	ctrl.Finish()
}

func TestNodeStoreError2(t *testing.T) {
	ctrl := gomock.NewController(t)
	ms := mock_store.NewMockStore(ctrl)

	// vars
	PCPRF := events.NodeCreatePodResponse
	PCPRPF := events.NodeCreatePodResponsePercentage
	genericError := errors.New("some error")

	// Expectations
	ms.EXPECT().GetNodeFlag(PCPRF).Return(scenario.Response_ERROR, nil)
	ms.EXPECT().GetNodeFlag(PCPRPF).Return(nil, genericError)

	var s store.Store = ms

	// Run code under test
	out, err := nodeResponse(responseArgs{
		ctx:      context.TODO(),
		provider: &Provider{store: &s},
		action: func() (i interface{}, err error) {
			return tStr, nil
		},
	},
		nodeResponseArgs{
			nodeResponseFlag:   PCPRF,
			nodePercentageFlag: PCPRPF,
		})

	// Assert
	assert.NotNil(t, err)
	assert.EqualError(t, genericError, err.Error())
	assert.Nil(t, out)

	ctrl.Finish()
}

func TestNodeInvalidPercentage(t *testing.T) {
	ctrl := gomock.NewController(t)
	ms := mock_store.NewMockStore(ctrl)

	// vars
	PCPRF := events.NodeCreatePodResponse
	PCPRPF := events.NodeCreatePodResponsePercentage

	// Expectations
	ms.EXPECT().GetNodeFlag(PCPRF).Return(scenario.Response_ERROR, nil)
	ms.EXPECT().GetNodeFlag(PCPRPF).Return(nil, nil)

	var s store.Store = ms

	// Run code under test
	out, err := nodeResponse(responseArgs{
		ctx:      context.TODO(),
		provider: &Provider{store: &s},
		action: func() (i interface{}, err error) {
			return tStr, nil
		},
	},
		nodeResponseArgs{
			nodeResponseFlag:   PCPRF,
			nodePercentageFlag: PCPRPF,
		})

	// Assert
	assert.NotNil(t, err)
	assert.EqualError(t, invalidPercentage, err.Error())
	assert.Nil(t, out)

	ctrl.Finish()
}

func TestNodeInvalidResponseType(t *testing.T) {
	ctrl := gomock.NewController(t)
	ms := mock_store.NewMockStore(ctrl)

	// vars
	PCPRF := events.NodeCreatePodResponse
	PCPRPF := events.NodeCreatePodResponsePercentage

	// Expectations
	ms.EXPECT().GetNodeFlag(PCPRF).Return(42, nil)

	var s store.Store = ms

	// Run code under test
	out, err := nodeResponse(responseArgs{
		ctx:      context.TODO(),
		provider: &Provider{store: &s},
		action: func() (i interface{}, err error) {
			return tStr, nil
		},
	},
		nodeResponseArgs{
			nodeResponseFlag:   PCPRF,
			nodePercentageFlag: PCPRPF,
		})

	// Assert
	assert.NotNil(t, err)
	assert.EqualError(t, invalidFlag, err.Error())
	assert.Nil(t, out)

	ctrl.Finish()
}

func TestNodeInvalidResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	ms := mock_store.NewMockStore(ctrl)

	// vars
	PCPRF := events.NodeCreatePodResponse
	PCPRPF := events.NodeCreatePodResponsePercentage

	rand.Seed(42)

	// Expectations
	ms.EXPECT().GetNodeFlag(PCPRF).Return(scenario.Response(42), nil)
	ms.EXPECT().GetNodeFlag(PCPRPF).Return(int32(100), nil)

	var s store.Store = ms

	// Run code under test
	out, err := nodeResponse(responseArgs{
		ctx:      context.TODO(),
		provider: &Provider{store: &s},
		action: func() (i interface{}, err error) {
			return tStr, nil
		},
	},
		nodeResponseArgs{
			nodeResponseFlag:   PCPRF,
			nodePercentageFlag: PCPRPF,
		})

	// Assert
	assert.NotNil(t, err)
	assert.EqualError(t, invalidResponse, err.Error())
	assert.Nil(t, out)

	ctrl.Finish()
}

func TestNodeTimeOut(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3)
	defer cancel()

	ctrl, ctx := gomock.WithContext(ctx, t)
	ms := mock_store.NewMockStore(ctrl)

	// vars
	PCPRF := events.NodeCreatePodResponse
	PCPRPF := events.NodeCreatePodResponsePercentage

	rand.Seed(42)

	// Expectations
	ms.EXPECT().GetNodeFlag(PCPRF).Return(scenario.Response_TIMEOUT, nil)
	ms.EXPECT().GetNodeFlag(PCPRPF).Return(int32(100), nil)

	var s store.Store = ms

	// Run code under test
	out, err := nodeResponse(responseArgs{
		ctx:      ctx,
		provider: &Provider{store: &s},
		action: func() (i interface{}, err error) {
			return tStr, nil
		},
	},
		nodeResponseArgs{
			nodeResponseFlag:   PCPRF,
			nodePercentageFlag: PCPRPF,
		})

	// Assert
	assert.Nil(t, err)
	assert.Nil(t, out)

	ctrl.Finish()
}
