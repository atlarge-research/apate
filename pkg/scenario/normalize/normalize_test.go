package normalize

import (
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario/public"
	"gotest.tools/assert"
	"testing"
)

func TestIterNodes(t *testing.T) {

	node1 := public.Node{
		Nodetype:   "test1",
		Ram:        "",
		CpuPercent: 0,
		MaxPods:    0,
	}

	node2 := public.Node{
		Nodetype:   "test2",
		Ram:        "",
		CpuPercent: 0,
		MaxPods:    0,
	}

	nodegroup1 := public.NodeGroup{
		Groupname: "testgroup1",
		Nodetype:  "test1",
		Amount:    27,
	}

	nodegroup2 := public.NodeGroup{
		Groupname: "testgroup2",
		Nodetype:  "test1",
		Amount:    42,
	}

	nodegroup3 := public.NodeGroup{
		Groupname: "testgroup3",
		Nodetype:  "test2",
		Amount:    42,
	}

	scenario := public.Scenario {
		Nodes: []*public.Node {
			&node1,
			&node2,
		},
		Nodegroups: []*public.NodeGroup {
			&nodegroup1,
			&nodegroup2,
			&nodegroup3,
		},
		Tasks:      nil,
	}

	nodecounter := 0

	IterNodes(scenario, func(_ int) {
		nodecounter += 1
	})

	assert.Equal(t, nodecounter,111)
}