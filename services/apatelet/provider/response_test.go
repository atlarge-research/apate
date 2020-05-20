package provider

import (
	"context"
	"testing"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store/mock_store"
)

func TestMagicPodAndNodePod(t *testing.T) {
	ctrl := gomock.NewController(t)

	ms := mock_store.NewMockStore(ctrl)

	// vars
	tStr := "test"
	podName := "madjik"
	PCPRF := events.PodCreatePodResponse

	// Expectations
	ms.EXPECT().GetPodFlag(podName, PCPRF).Return(scenario.ResponseNormal, nil)

	// SOT
	var s store.Store = ms

	out, err := podAndNodeResponse(responseArgs{
		ctx:      context.TODO(),
		provider: &Provider{store: &s},
		action: func() (i interface{}, err error) {
			return tStr, nil
		}},
		podName,
		PCPRF,
		events.NodeCreatePodResponse,
	)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, tStr, out)
}

func TestMagicPodAndNodeNode(t *testing.T) {
	ctrl := gomock.NewController(t)

	ms := mock_store.NewMockStore(ctrl)

	// vars
	tStr := "test"
	podName := "madjik"
	PCPRF := events.PodCreatePodResponse

	NCPRF := events.PodCreatePodResponse

	// Expectations
	ms.EXPECT().GetPodFlag(podName, PCPRF).Return(scenario.ResponseNormal, nil)
	ms.EXPECT().GetNodeFlag(NCPRF).Return(scenario.ResponseNormal, nil)

	// SOT
	var s store.Store = ms

	out, err := podAndNodeResponse(
		responseArgs{
			ctx:      context.TODO(),
			provider: &Provider{store: &s},
			action: func() (i interface{}, err error) {
				return tStr, nil
			},
		},
		podName,
		PCPRF,
		NCPRF,
	)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, tStr, out)
}
