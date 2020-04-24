package normalise

import (
	scenario2 "github.com/atlarge-research/opendc-emulate-kubernetes/api/control_plane"
	"testing"

	"gotest.tools/assert"
)

func TestIterNodes(t *testing.T) {
	node1 := scenario2.Node{
		NodeType:   "test1",
		Ram:        "",
		CpuPercent: 0,
		MaxPods:    0,
	}

	node2 := scenario2.Node{
		NodeType:   "test2",
		Ram:        "",
		CpuPercent: 0,
		MaxPods:    0,
	}

	nodegroup1 := scenario2.NodeGroup{
		GroupName: "testgroup1",
		NodeType:  "test1",
		Amount:    27,
	}

	nodegroup2 := scenario2.NodeGroup{
		GroupName: "testgroup2",
		NodeType:  "test1",
		Amount:    42,
	}

	nodegroup3 := scenario2.NodeGroup{
		GroupName: "testgroup3",
		NodeType:  "test2",
		Amount:    42,
	}

	scenario := scenario2.Scenario{
		Nodes: []*scenario2.Node{
			&node1,
			&node2,
		},
		NodeGroups: []*scenario2.NodeGroup{
			&nodegroup1,
			&nodegroup2,
			&nodegroup3,
		},
		Tasks: nil,
	}
}
