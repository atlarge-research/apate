package node

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/atlarge-research/apate/internal/crd/node"
	nodeconfigv1 "github.com/atlarge-research/apate/pkg/apis/nodeconfiguration/v1"
	"github.com/atlarge-research/apate/pkg/scenario"
)

func createResources(needed int, base scenario.NodeResources) []scenario.NodeResources {
	var resources []scenario.NodeResources

	for i := 0; i < needed; i++ {
		res := base
		res.UUID = uuid.New()

		resources = append(resources, res)
	}

	return resources
}

func getNodeResources(nodeCfg *nodeconfigv1.NodeConfiguration) (scenario.NodeResources, error) {
	res := nodeCfg.Spec.Resources
	mem, err := scenario.GetInBytes(res.Memory, "memory")
	if err != nil {
		return scenario.NodeResources{}, errors.Wrap(err, "couldn't convert memory to bytes")
	}

	storage, err := scenario.GetInBytes(res.Storage, "storage")
	if err != nil {
		return scenario.NodeResources{}, errors.Wrap(err, "couldn't convert storage to bytes")
	}

	ephemeralStorage, err := scenario.GetInBytes(res.EphemeralStorage, "ephemeral storage")
	if err != nil {
		return scenario.NodeResources{}, errors.Wrap(err, "couldn't convert ephemeral storage to bytes")
	}

	return scenario.NodeResources{
		Memory:           mem,
		CPU:              res.CPU,
		Storage:          storage,
		EphemeralStorage: ephemeralStorage,
		MaxPods:          res.MaxPods,
		Label:            node.GetCrdLabel(nodeCfg),
	}, nil
}
