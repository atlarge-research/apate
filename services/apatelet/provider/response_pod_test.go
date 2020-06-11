package provider

import (
	"context"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	podconfigv1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/podconfiguration/v1"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario"

	"github.com/pkg/errors"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store/mock_store"
)

func setup(t *testing.T) (*mock_store.MockStore, *gomock.Controller, *corev1.Pod, func(podFlag events.PodEventFlag) (interface{}, error)) {
	ctrl := gomock.NewController(t)
	ms := mock_store.NewMockStore(ctrl)
	var s store.Store = ms

	pod := createPodWithLabel(podNamespace, podLabel)

	return ms, ctrl, pod, func(podFlag events.PodEventFlag) (interface{}, error) {
		// Run code under test
		return podResponse(responseArgs{
			ctx:      context.Background(),
			provider: &Provider{Store: &s},
			action: func() (i interface{}, err error) {
				return tStr, nil
			}},
			pod,
			podFlag,
		)
	}
}

func createPodWithLabel(ns string, label string) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				podconfigv1.PodConfigurationLabel: label,
			},
			Namespace: ns,
		},
	}
}

func TestPodNormal(t *testing.T) {
	t.Parallel()

	ms, ctrl, pod, executor := setup(t)
	defer ctrl.Finish()

	// Expectations
	ms.EXPECT().GetPodFlag(pod, events.PodCreatePodResponse).Return(scenario.ResponseNormal, nil)
	ms.EXPECT().GetNodeFlag(events.NodeCreatePodResponse).Return(scenario.ResponseUnset, nil)

	// Execute
	out, err := executor(events.PodCreatePodResponse)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, tStr, out)
}

func TestPodStoreError1(t *testing.T) {
	t.Parallel()

	ms, ctrl, pod, executor := setup(t)
	defer ctrl.Finish()

	// vars
	genericError := errors.New("some error")

	// Expectations
	ms.EXPECT().GetPodFlag(pod, events.PodCreatePodResponse).Return(nil, genericError)

	// Run code under test
	out, err := executor(events.PodCreatePodResponse)

	// Assert
	assert.Error(t, err)
	assert.False(t, IsExpected(err))
	assert.Nil(t, out)
}

func TestPodStoreError2(t *testing.T) {
	t.Parallel()

	ms, ctrl, pod, executor := setup(t)
	defer ctrl.Finish()

	// Expectations
	ms.EXPECT().GetPodFlag(pod, events.PodCreatePodResponse).Return(scenario.ResponseError, nil)
	ms.EXPECT().GetNodeFlag(events.NodeCreatePodResponse).Return(scenario.ResponseUnset, nil)

	// Run code under test
	out, err := executor(events.PodCreatePodResponse)

	// Assert
	assert.Error(t, err)
	assert.True(t, IsExpected(err))
	assert.Nil(t, out)
}

func TestPodUnset(t *testing.T) {
	t.Parallel()

	ms, ctrl, pod, executor := setup(t)
	defer ctrl.Finish()

	// Expectations
	ms.EXPECT().GetPodFlag(pod, events.PodCreatePodResponse).Return(scenario.ResponseUnset, nil)
	ms.EXPECT().GetNodeFlag(events.NodeCreatePodResponse).Return(scenario.ResponseUnset, nil)

	// Run code under test
	out, err := executor(events.PodCreatePodResponse)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, tStr, out)
}

func TestPodInvalidResponseType(t *testing.T) {
	t.Parallel()

	ms, ctrl, pod, executor := setup(t)
	defer ctrl.Finish()

	// Expectations
	ms.EXPECT().GetPodFlag(pod, events.PodCreatePodResponse).Return(42, nil)

	// Run code under test
	out, err := executor(events.PodCreatePodResponse)

	// Assert
	assert.Error(t, err)
	assert.False(t, IsExpected(err))
	assert.Nil(t, out)
}

func TestPodInvalidResponse(t *testing.T) {
	t.Parallel()

	ms, ctrl, pod, executor := setup(t)
	defer ctrl.Finish()

	// Expectations
	ms.EXPECT().GetPodFlag(pod, events.PodCreatePodResponse).Return(scenario.Response(42), nil)
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

	pod := createPodWithLabel(podNamespace, podLabel)

	// Expectations
	ms.EXPECT().GetPodFlag(pod, events.PodCreatePodResponse).Return(scenario.ResponseTimeout, nil)
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
		pod,
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

	pod := createPodWithLabel(podNamespace, podLabel)

	// vars
	podFlag := events.PodCreatePodResponse

	// Expectations
	ms.EXPECT().GetPodFlag(pod, podFlag).Return(scenario.ResponseError, nil)
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
		pod,
		podFlag,
	)

	assert.Error(t, ctx.Err())

	// Assert
	assert.Nil(t, err)
	assert.Nil(t, out)
}

func TestErrorVsNormal(t *testing.T) {
	t.Parallel()

	ms, ctrl, pod, executor := setup(t)
	defer ctrl.Finish()

	// Expectations
	ms.EXPECT().GetPodFlag(pod, events.PodCreatePodResponse).Return(scenario.ResponseNormal, nil)
	ms.EXPECT().GetNodeFlag(events.NodeCreatePodResponse).Return(scenario.ResponseError, nil)

	// Run code under test
	out, err := executor(events.PodCreatePodResponse)

	// Assert
	assert.Error(t, err)
	assert.True(t, IsExpected(err))
	assert.Nil(t, out)
}
