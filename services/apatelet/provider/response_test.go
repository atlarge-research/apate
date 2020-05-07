package provider

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario"
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
	PCPRPF := events.PodCreatePodResponsePercentage

	// Expectations
	ms.EXPECT().GetPodFlag(podName, PCPRF).Return(scenario.Response_NORMAL, nil)
	ms.EXPECT().GetPodFlag(podName, PCPRPF).Return(int32(100), nil)

	// SOT
	var s store.Store = ms

	out, err := podAndNodeResponse(podNodeResponse{
		responseArgs: responseArgs{
			ctx: context.TODO(),
			p:   &VKProvider{store: &s},
			action: func() (i interface{}, err error) {
				return tStr, nil
			},
		},
		podResponseArgs: podResponseArgs{
			name:              podName,
			podResponseFlag:   PCPRF,
			podPercentageFlag: PCPRPF,
		},
	})

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
	PCPRPF := events.PodCreatePodResponsePercentage

	NCPRF := events.PodCreatePodResponse
	NCPRPF := events.PodCreatePodResponsePercentage

	// Expectations
	ms.EXPECT().GetPodFlag(podName, PCPRF).Return(scenario.Response_NORMAL, nil)
	ms.EXPECT().GetPodFlag(podName, PCPRPF).Return(int32(0), nil)

	ms.EXPECT().GetNodeFlag(NCPRF).Return(scenario.Response_NORMAL, nil)
	ms.EXPECT().GetNodeFlag(NCPRPF).Return(int32(100), nil)

	// SOT
	var s store.Store = ms

	out, err := podAndNodeResponse(podNodeResponse{
		responseArgs: responseArgs{
			ctx: context.TODO(),
			p:   &VKProvider{store: &s},
			action: func() (i interface{}, err error) {
				return tStr, nil
			},
		},
		podResponseArgs: podResponseArgs{
			name:              podName,
			podResponseFlag:   PCPRF,
			podPercentageFlag: PCPRPF,
		},
		nodeResponseArgs: nodeResponseArgs{
			nodeResponseFlag:   NCPRF,
			nodePercentageFlag: NCPRPF,
		},
	})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, tStr, out)
}
