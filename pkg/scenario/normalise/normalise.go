// Package normalise provides functions to normalise and decode public scenarios.
package normalise

import (
	"github.com/google/uuid"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario/private"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario/public"
)

// IterNodes returns only the part of a scenario relevant to creating nodes.
// This is necessary because it's impossible to entirely normalise a
// scenario without knowing the UUIDs of each spawned node. They need to
// be spawned first. They can be spawned based on this function.
//
// This function returns a iterable channel. This is useful because storing all
// (potentially very many) nodes in for example an array would be extremely inefficient.
// Especially because most of the nodes are the same. Using a channel it's possible to
// return a reference to the same node multiple times
func IterNodes(scenario *public.Scenario, callback func(i int)) {
	// Iterate over every nodegroup
	for _, nodegroup := range scenario.GetNodeGroups() {
		// Yield every type of node as many times as the `amount` field
		// in the nodegroup says.
		for i := 0; i < int(nodegroup.Amount); i++ {
			callback(i)
		}
	}
}

// GetPrivateScenario takes a public scenario and turns it into a private scenario. Normalises the structure and resolves named references.
func GetPrivateScenario(scenario *public.Scenario, uuids []uuid.UUID) (*private.Scenario, error) {
	r := private.Scenario{}

	// This function does not need to set this field. This is set by the control plane
	// Whenever the scenario is started.
	r.StartTime = 0

	groups := make(map[string][]string)

	// A variable holding which uuid was used last.
	// With this, every node can get a new uuid.
	uuidindex := 0

	for _, nodegroup := range scenario.NodeGroups {
		for i := 0; i < int(nodegroup.Amount); i++ {
			id := uuids[uuidindex]
			uuidindex++

			groups[nodegroup.GroupName] = append(groups[nodegroup.GroupName], id.String())
		}
	}

	var tasks []*private.Task

	for _, task := range scenario.Tasks {
		time, err := desugarTimestamp(task.Time)
		if err != nil {
			return nil, err
		}

		nodegroupnames, err := desugarNodeSet(task.NodeGroups, scenario.NodeGroups)
		if err != nil {
			return nil, err
		}

		var nodeset []string

		for _, name := range nodegroupnames {
			nodeset = append(nodeset, groups[name]...)
		}

		tasks = append(tasks, &private.Task{
			Name:       task.Name,
			RevertTask: task.Revert,
			Timestamp:  int32(time),
			NodeSet:    nodeset,
			Event:      nil,
		})
	}

	r.Task = tasks

	return &r, nil
}
