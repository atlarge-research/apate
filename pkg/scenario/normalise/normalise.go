// Package normalise provides functions to normalise and decode public scenarios.
package normalise

import (
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario/private"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario/public"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/control_plane/cluster"
	"github.com/docker/go-units"
	"time"
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

// NormaliseScenario takes a public scenario and turns it into a private scenario. Normalises the structure and resolves named references.
func NormaliseScenario(scenario *public.Scenario, nodes []cluster.Node) (*private.Scenario, map[cluster.Node]NodeResources, error) {
	r := private.Scenario{}

	// This function does not need to set this field. This is set by the control plane
	// Whenever the scenario is started.
	r.StartTime = 0

	groups := make(map[string][]string)
	resources := make(map[cluster.Node]NodeResources)

	nodetypes := make(map[string]*public.Node)
	// First make a lookup mapping nodetype strings to node types.
	// This makes later lookup O(1)
	for _, nodetype := range scenario.GetNodes() {
		nodetypes[nodetype.NodeType] = nodetype
	}

	// A variable holding which  ode was used last.
	// With this, every node can get a new node.
	index := 0

	for _, nodegroup := range scenario.NodeGroups {
		for i := 0; i < int(nodegroup.Amount); i++ {
			node := nodes[index]
			index++

			groups[nodegroup.GroupName] = append(groups[nodegroup.GroupName], node.UUID.String())

			nodetype := nodetypes[nodegroup.NodeType]
			memory, err := units.RAMInBytes(nodetype.Ram)
			if err != nil {
				return nil, nil, err
			}

			resources[node] = NodeResources{
				memory,
				int(nodetype.CpuPercent),
				int(nodetype.MaxPods),
			}
		}
	}

	var tasks []*private.Task

	for _, task := range scenario.Tasks {
		// Desugar the timestamp postfix.
		time, err := desugarTimestamp(task.Time)
		if err != nil {
			return nil, nil, err
		}

		// Decode the "all" node name, also verify that all names in the nodeset exist and
		// that there are no duplicates in the set.
		nodegroupnames, err := desugarNodeSet(task.NodeGroups, scenario.NodeGroups)
		if err != nil {
			return nil, nil, err
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

	return &r, nil, nil
}

func desugarTimestamp(t string) (int, error) {
	duration, err := time.ParseDuration(t)
	if err != nil {
		return 0, err
	}

	return int(duration.Milliseconds()), nil
}
