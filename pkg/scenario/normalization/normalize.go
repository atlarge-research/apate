// Package normalization provides functions to normalize and decode public scenarios.
package normalization

import (
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/normalization/events"
	"time"

	"github.com/google/uuid"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/apatelet"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/controlplane"
	"github.com/docker/go-units"
)

// NormalizeScenario takes a public scenario and turns it into a private scenario.
// Normalizes the structure and resolves named references.
func NormalizeScenario(scenario *controlplane.PublicScenario) (*apatelet.ApateletScenario, []NodeResources, error) {
	r := apatelet.ApateletScenario{}

	nodeResources := make([]NodeResources, 0)
	uuidsPerNodeGroup := make(map[string][]uuid.UUID)

	// First make a lookup mapping nodeType strings to node types.
	// This makes later lookup O(1)
	nodeTypeName := make(map[string]*controlplane.Node)
	for _, nodeType := range scenario.GetNodes() {
		nodeTypeName[nodeType.NodeType] = nodeType
	}

	for _, nodeGroup := range scenario.NodeGroups {
		for i := 0; i < int(nodeGroup.Amount); i++ {
			id := uuid.New()

			nodeType := nodeTypeName[nodeGroup.NodeType]

			memory, err := units.RAMInBytes(nodeType.Memory)
			if err != nil {
				return nil, nil, err
			}

			storage, err := units.RAMInBytes(nodeType.Storage)
			if err != nil {
				return nil, nil, err
			}

			ephStorage, err := units.RAMInBytes(nodeType.EphemeralStorage)
			if err != nil {
				return nil, nil, err
			}

			nodeResources = append(nodeResources, NodeResources{
				id,
				memory,
				nodeType.Cpu,
				storage,
				ephStorage,
				nodeType.MaxPods,
			})

			uuidsPerNodeGroup[nodeGroup.GroupName] = append(uuidsPerNodeGroup[nodeGroup.GroupName], id)
		}
	}

	var tasks []*apatelet.Task

	for _, task := range scenario.Tasks {
		timestamp, err := desugarTimestamp(task.Time)
		if err != nil {
			return nil, nil, err
		}

		// Decode the "all" node name, also verify that all names in the nodeSet exist and
		// that there are no duplicates in the set.
		nodeGroupNames, err := desugarNodeGroups(task.NodeGroups, scenario.NodeGroups)
		if err != nil {
			return nil, nil, err
		}

		var nodeSet []string

		for _, name := range nodeGroupNames {
			for _, nodeUUID := range uuidsPerNodeGroup[name] {
				nodeSet = append(nodeSet, nodeUUID.String())
			}
		}

		newTask := &apatelet.Task{
			Name:       task.Name,
			RevertTask: task.Revert,
			Timestamp:  int32(timestamp),
			NodeSet:    nodeSet,
		}

		if !task.Revert {
			if err := events.NewEventTranslator(task, newTask).TranslateEvent(); err != nil {
				return nil, nil, err
			}
		}

		tasks = append(tasks, newTask)
	}

	r.Task = tasks

	return &r, nodeResources, nil
}

func desugarTimestamp(t string) (int, error) {
	duration, err := time.ParseDuration(t)
	if err != nil {
		return 0, err
	}

	return int(duration.Milliseconds()), nil
}
