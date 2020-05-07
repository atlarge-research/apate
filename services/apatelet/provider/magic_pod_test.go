package provider

import (
	"context"
	"errors"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store/mock_store"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestMagicPodNormal100(t *testing.T) {
	ctrl := gomock.NewController(t)

	ms := mock_store.NewMockStore(ctrl)

	// vars
	tStr := "test"
	podName := "madjik"
	PCPRF := events.PodCreatePodResponse
	PCPRPF := events.PodCreatePodResponsePercentage

	// Expectations
	ms.EXPECT().GetPodFlag(podName, PCPRF).Return(scenario.Response_NORMAL, nil)
	ms.EXPECT().GetPodFlag(podName, PCPRPF).Return(int32(100), nil)
	
	var s store.Store = ms

	// Run code under test
	out, err := magicPod(magicArgs{
		ctx:    context.TODO(),
		p: &VKProvider{store: &s},
		action: func() (i interface{}, err error) {
			return tStr, nil
		},
	},
	magicPodArgs{
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
	tStr := "test"
	podName := "madjik"
	PCPRF := events.PodCreatePodResponse
	PCPRPF := events.PodCreatePodResponsePercentage

	// Expectations
	ms.EXPECT().GetPodFlag(podName, PCPRF).Return(scenario.Response_NORMAL, nil)
	ms.EXPECT().GetPodFlag(podName, PCPRPF).Return(int32(0), nil)

	var s store.Store = ms

	// Run code under test
	out, err := magicPod(magicArgs{
		ctx:    context.TODO(),
		p: &VKProvider{store: &s},
		action: func() (i interface{}, err error) {
			return tStr, nil
		},
	},
		magicPodArgs{
			name:              podName,
			podResponseFlag:   PCPRF,
			podPercentageFlag: PCPRPF,
		})

	// Assert
	assert.NotNil(t, err)
	assert.EqualError(t, FlagNotSetError, err.Error())
	assert.Nil(t, out)

	ctrl.Finish()
}


func TestMagicPodNormal50A(t *testing.T) {
	ctrl := gomock.NewController(t)

	ms := mock_store.NewMockStore(ctrl)

	// vars
	tStr := "test"
	podName := "madjik"
	PCPRF := events.PodCreatePodResponse
	PCPRPF := events.PodCreatePodResponsePercentage

	rand.Seed(69)

	// Expectations
	ms.EXPECT().GetPodFlag(podName, PCPRF).Return(scenario.Response_ERROR, nil)
	ms.EXPECT().GetPodFlag(podName, PCPRPF).Return(int32(50), nil)

	var s store.Store = ms

	// Run code under test
	out, err := magicPod(magicArgs{
		ctx:    context.TODO(),
		p: &VKProvider{store: &s},
		action: func() (i interface{}, err error) {
			return tStr, nil
		},
	},
		magicPodArgs{
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
	tStr := "test"
	podName := "madjik"
	PCPRF := events.PodCreatePodResponse
	PCPRPF := events.PodCreatePodResponsePercentage

	rand.Seed(42)

	// Expectations
	ms.EXPECT().GetPodFlag(podName, PCPRF).Return(scenario.Response_ERROR, nil)
	ms.EXPECT().GetPodFlag(podName, PCPRPF).Return(int32(50), nil)

	var s store.Store = ms

	// Run code under test
	out, err := magicPod(magicArgs{
		ctx:    context.TODO(),
		p: &VKProvider{store: &s},
		action: func() (i interface{}, err error) {
			return tStr, nil
		},
	},
		magicPodArgs{
			name:              podName,
			podResponseFlag:   PCPRF,
			podPercentageFlag: PCPRPF,
		})

	// Assert
	assert.NotNil(t, err)
	assert.EqualError(t, ExpectedError, err.Error())
	assert.Nil(t, out)

	ctrl.Finish()
}

func TestMagicPodStoreError1(t *testing.T) {
	ctrl := gomock.NewController(t)
	ms := mock_store.NewMockStore(ctrl)

	// vars
	tStr := "test"
	podName := "madjik"
	PCPRF := events.PodCreatePodResponse
	PCPRPF := events.PodCreatePodResponsePercentage
	genericError := errors.New("some error")

	// Expectations
	ms.EXPECT().GetPodFlag(podName, PCPRF).Return(nil, genericError)

	var s store.Store = ms

	// Run code under test
	out, err := magicPod(magicArgs{
		ctx:    context.TODO(),
		p: &VKProvider{store: &s},
		action: func() (i interface{}, err error) {
			return tStr, nil
		},
	},
		magicPodArgs{
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
	tStr := "test"
	podName := "madjik"
	PCPRF := events.PodCreatePodResponse
	PCPRPF := events.PodCreatePodResponsePercentage
	genericError := errors.New("some error")

	// Expectations
	ms.EXPECT().GetPodFlag(podName, PCPRF).Return(scenario.Response_ERROR, nil)
	ms.EXPECT().GetPodFlag(podName, PCPRPF).Return(nil, genericError)

	var s store.Store = ms

	// Run code under test
	out, err := magicPod(magicArgs{
		ctx:    context.TODO(),
		p: &VKProvider{store: &s},
		action: func() (i interface{}, err error) {
			return tStr, nil
		},
	},
		magicPodArgs{
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
	tStr := "test"
	podName := "madjik"
	PCPRF := events.PodCreatePodResponse
	PCPRPF := events.PodCreatePodResponsePercentage

	// Expectations
	ms.EXPECT().GetPodFlag(podName, PCPRF).Return(scenario.Response_ERROR, nil)
	ms.EXPECT().GetPodFlag(podName, PCPRPF).Return(nil, nil)

	var s store.Store = ms

	// Run code under test
	out, err := magicPod(magicArgs{
		ctx:    context.TODO(),
		p: &VKProvider{store: &s},
		action: func() (i interface{}, err error) {
			return tStr, nil
		},
	},
		magicPodArgs{
			name:              podName,
			podResponseFlag:   PCPRF,
			podPercentageFlag: PCPRPF,
		})

	// Assert
	assert.NotNil(t, err)
	assert.EqualError(t, InvalidPercentage, err.Error())
	assert.Nil(t, out)

	ctrl.Finish()
}

func TestMagicPodInvalidResponseType(t *testing.T) {
	ctrl := gomock.NewController(t)
	ms := mock_store.NewMockStore(ctrl)

	// vars
	tStr := "test"
	podName := "madjik"
	PCPRF := events.PodCreatePodResponse
	PCPRPF := events.PodCreatePodResponsePercentage

	// Expectations
	ms.EXPECT().GetPodFlag(podName, PCPRF).Return(42, nil)

	var s store.Store = ms

	// Run code under test
	out, err := magicPod(magicArgs{
		ctx:    context.TODO(),
		p: &VKProvider{store: &s},
		action: func() (i interface{}, err error) {
			return tStr, nil
		},
	},
		magicPodArgs{
			name:              podName,
			podResponseFlag:   PCPRF,
			podPercentageFlag: PCPRPF,
		})

	// Assert
	assert.NotNil(t, err)
	assert.EqualError(t, InvalidFlag, err.Error())
	assert.Nil(t, out)

	ctrl.Finish()
}



func TestMagicPodInvalidResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	ms := mock_store.NewMockStore(ctrl)

	// vars
	tStr := "test"
	podName := "madjik"
	PCPRF := events.PodCreatePodResponse
	PCPRPF := events.PodCreatePodResponsePercentage

	rand.Seed(42)

	// Expectations
	ms.EXPECT().GetPodFlag(podName, PCPRF).Return(scenario.Response(42), nil)
	ms.EXPECT().GetPodFlag(podName, PCPRPF).Return(int32(100), nil)

	var s store.Store = ms

	// Run code under test
	out, err := magicPod(magicArgs{
		ctx:    context.TODO(),
		p: &VKProvider{store: &s},
		action: func() (i interface{}, err error) {
			return tStr, nil
		},
	},
		magicPodArgs{
			name:              podName,
			podResponseFlag:   PCPRF,
			podPercentageFlag: PCPRPF,
		})

	// Assert
	assert.NotNil(t, err)
	assert.EqualError(t, InvalidResponse, err.Error())
	assert.Nil(t, out)

	ctrl.Finish()
}

func TestMagicPodTimeOut(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 3)
	ctrl, ctx := gomock.WithContext(ctx, t)
	ms := mock_store.NewMockStore(ctrl)

	// vars
	tStr := "test"
	podName := "madjik"
	PCPRF := events.PodCreatePodResponse
	PCPRPF := events.PodCreatePodResponsePercentage

	rand.Seed(42)

	// Expectations
	ms.EXPECT().GetPodFlag(podName, PCPRF).Return(scenario.Response_TIMEOUT, nil)
	ms.EXPECT().GetPodFlag(podName, PCPRPF).Return(int32(100), nil)

	var s store.Store = ms

	// Run code under test
	out, err :=  magicPod(magicArgs{
		ctx:    ctx,
		p: &VKProvider{store: &s},
		action: func() (i interface{}, err error) {
			return tStr, nil
		},
	},
		magicPodArgs{
			name:              podName,
			podResponseFlag:   PCPRF,
			podPercentageFlag: PCPRPF,
		})

	// Assert
	assert.Nil(t, err)
	assert.Nil(t, out)

	ctrl.Finish()
}
