package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
)

func TestGetCorrespondingNodeEventFlag(t *testing.T) {
	nodeFlag, err := getCorrespondingNodeEventFlag(events.PodCreatePodResponse)
	assert.Equal(t, events.NodeCreatePodResponse, nodeFlag)
	assert.NoError(t, err)

	nodeFlag, err = getCorrespondingNodeEventFlag(events.PodUpdatePodResponse)
	assert.Equal(t, events.NodeUpdatePodResponse, nodeFlag)
	assert.NoError(t, err)

	nodeFlag, err = getCorrespondingNodeEventFlag(events.PodDeletePodResponse)
	assert.Equal(t, events.NodeDeletePodResponse, nodeFlag)
	assert.NoError(t, err)

	nodeFlag, err = getCorrespondingNodeEventFlag(events.PodGetPodResponse)
	assert.Equal(t, events.NodeGetPodResponse, nodeFlag)
	assert.NoError(t, err)

	nodeFlag, err = getCorrespondingNodeEventFlag(events.PodGetPodStatusResponse)
	assert.Equal(t, events.NodeGetPodStatusResponse, nodeFlag)
	assert.NoError(t, err)

	unsetFlag, err := getCorrespondingNodeEventFlag(events.PodResources)
	assert.Error(t, err)
	assert.Equal(t, int32(-1), unsetFlag)
}
