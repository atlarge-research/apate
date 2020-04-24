package normalise

import (
	"testing"

	"gotest.tools/assert"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario/public"
)

func TestIterNodes(t *testing.T) {
	node1 := public.Node{
		NodeType:   "test1",
		Ram:        "",
		CpuPercent: 0,
		MaxPods:    0,
	}

	node2 := public.Node{
		NodeType:   "test2",
		Ram:        "",
		CpuPercent: 0,
		MaxPods:    0,
	}

	nodegroup1 := public.NodeGroup{
		GroupName: "testgroup1",
		NodeType:  "test1",
		Amount:    27,
	}

	nodegroup2 := public.NodeGroup{
		GroupName: "testgroup2",
		NodeType:  "test1",
		Amount:    42,
	}

	nodegroup3 := public.NodeGroup{
		GroupName: "testgroup3",
		NodeType:  "test2",
		Amount:    42,
	}

	scenario := public.Scenario{
		Nodes: []*public.Node{
			&node1,
			&node2,
		},
		NodeGroups: []*public.NodeGroup{
			&nodegroup1,
			&nodegroup2,
			&nodegroup3,
		},
		Tasks: nil,
	}

	nodecounter := 0

	IterNodes(&scenario, func(_ int) {
		nodecounter++
	})

	assert.Equal(t, nodecounter, 111)
}
