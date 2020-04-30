// Package normalization provides functions to normalize and decode public scenarios.
package normalization

import (
	"github.com/google/uuid"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/normalization/events"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/apatelet"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/controlplane"
)

type normalizationContext struct {
	// The public scenario to normalize
	scenario *controlplane.PublicScenario

	// The created nodes and their resources
	nodeResources []NodeResources

	// A map from a node group to the UUIDs of corresponding nodes
	uuidsPerNodeGroup map[string][]uuid.UUID

	// A map mapping nodeType strings to node original types
	// This aids in doing lookup later on
	nodeTypeName map[string]*controlplane.Node
}

// NormalizeScenario takes a public scenario and turns it into a private scenario.
// Normalizes the structure and resolves named references.
func NormalizeScenario(scenario *controlplane.PublicScenario) (*apatelet.ApateletScenario, []NodeResources, error) {
	r := apatelet.ApateletScenario{}

	c := &normalizationContext{
		scenario:          scenario,
		nodeResources:     make([]NodeResources, 0),
		uuidsPerNodeGroup: make(map[string][]uuid.UUID),
		nodeTypeName:      make(map[string]*controlplane.Node),
	}

	// Fill the map with node cache
	for _, nodeType := range scenario.GetNodes() {
		c.nodeTypeName[nodeType.NodeType] = nodeType
	}

	// Create nodes from node group & hardware definitions
	if err := normalizeNodes(c); err != nil {
		return nil, nil, err
	}

	// Normalize the tasks
	tasks, err := normalizeTasks(c)
	if err != nil {
		return nil, nil, err
	}
	r.Task = tasks

	return &r, c.nodeResources, nil
}

// normalizeTasks translates the tasks from a public to internal scenario
func normalizeTasks(c *normalizationContext) ([]*apatelet.Task, error) {
	var tasks []*apatelet.Task

	for _, task := range c.scenario.Tasks {
		timestamp, err := events.DesugarTimestamp(task.Time)
		if err != nil {
			return nil, err
		}

		// Decode the "all" node name, also verify that all names in the nodeSet exist and
		// that there are no duplicates in the set.
		nodeGroupNames, err := desugarNodeGroups(task.NodeGroups, c.scenario.NodeGroups)
		if err != nil {
			return nil, err
		}

		var nodeSet []string
		for _, name := range nodeGroupNames {
			for _, nodeUUID := range c.uuidsPerNodeGroup[name] {
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
				return nil, err
			}
		}

		tasks = append(tasks, newTask)
	}
	return tasks, nil
}

// normalizeNodes parses the node groups in a scenario into separate nodes with a certain hardware definition
func normalizeNodes(c *normalizationContext) error {
	for _, nodeGroup := range c.scenario.NodeGroups {
		for i := 0; i < int(nodeGroup.Amount); i++ {
			id := uuid.New()

			nodeType := c.nodeTypeName[nodeGroup.NodeType]

			memory, err := events.GetInBytes(nodeType.Memory, "memory")
			if err != nil {
				return err
			}

			storage, err := events.GetInBytes(nodeType.Storage, "storage")
			if err != nil {
				return err
			}

			ephStorage, err := events.GetInBytes(nodeType.EphemeralStorage, "ephemeral storage")
			if err != nil {
				return err
			}

			c.nodeResources = append(c.nodeResources, NodeResources{
				id,
				memory,
				nodeType.Cpu,
				storage,
				ephStorage,
				nodeType.MaxPods,
			})

			c.uuidsPerNodeGroup[nodeGroup.GroupName] = append(c.uuidsPerNodeGroup[nodeGroup.GroupName], id)
		}
	}

	return nil
}
