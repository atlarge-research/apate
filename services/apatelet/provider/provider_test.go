package provider

import (
	"context"
	"testing"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/normalization"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

// TODO: Expand when provider has more functionality

func TestConfigureNode(t *testing.T) {
	resources := normalization.NodeResources{
		UUID:    uuid.New(),
		Memory:  42,
		CPU:     1337,
		MaxPods: 1001,
	}

	prov := VKProvider{
		Pods:      nil,
		resources: &resources,
	}

	fakeNode := corev1.Node{}

	// Run the method
	prov.ConfigureNode(context.TODO(), &fakeNode)

	assert.EqualValues(t, resources.CPU, fakeNode.Status.Capacity.Cpu().Value())
	assert.EqualValues(t, resources.Memory, fakeNode.Status.Capacity.Memory().Value())
	assert.EqualValues(t, resources.MaxPods, fakeNode.Status.Capacity.Pods().Value())
}

func TestConfigureNodeWithCreate(t *testing.T) {
	resources := normalization.NodeResources{
		UUID:    uuid.New(),
		Memory:  42,
		CPU:     1337,
		MaxPods: 1001,
	}

	prov := CreateProvider(&resources)

	fakeNode := corev1.Node{}

	// Run the method
	prov.ConfigureNode(context.TODO(), &fakeNode)

	assert.EqualValues(t, resources.CPU, fakeNode.Status.Capacity.Cpu().Value())
	assert.EqualValues(t, resources.Memory, fakeNode.Status.Capacity.Memory().Value())
	assert.EqualValues(t, resources.MaxPods, fakeNode.Status.Capacity.Pods().Value())
}
