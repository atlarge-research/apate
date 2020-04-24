package scenario

import (
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/deserialize"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/normalise"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/control_plane/cluster"
	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"

	"testing"
)

func TestScenario(t *testing.T) {
	// CLI

	scenario, err := deserialize.YamlScenario{}.FromBytes([]byte(`
nodes:
    - node_type: testnode
      ram: 2G
      cpu_percent: 42
      max_pods: 42

node_groups:
    - group_name: testgroup1
      node_type: testnode
      amount: 42
    - group_name: testgroup2
      node_type: testnode
      amount: 10

tasks:
    - name: testtask1
      time: 10s
      node_groups: 
          - testgroup1

    - name: testtask2
      time: 10s
      node_groups: 
          - all

    - name: testtask2
      time: 20s
      revert: true
`))
	assert.NoError(t, err)

	// Control plane



	nodecount := normalise.NumNodes(scenario.GetScenario())

	assert.Equal(t, nodecount, 52)

	nodes := make([]cluster.Node, 0, nodecount)
	for i := 0; i < nodecount; i++ {
		nodes = append(nodes, cluster.Node{UUID: uuid.New()})
	}

	ps, _, err := normalise.NormaliseScenario(scenario.GetScenario(), nodes)
	assert.NoError(t, err)

	// Should be 0 because this is set when the scenario is started.
	assert.Equal(t, ps.StartTime, int32(0))
	assert.Equal(t, len(ps.Task), 3)
	assert.Equal(t, ps.Task[0].Name, "testtask1")
	assert.Equal(t, ps.Task[0].RevertTask, false)
	assert.Equal(t, len(ps.Task[0].NodeSet), 42)

	assert.Equal(t, ps.Task[1].Name, "testtask2")
	assert.Equal(t, ps.Task[1].RevertTask, false)
	assert.Equal(t, len(ps.Task[1].NodeSet), 52)

	assert.Equal(t, ps.Task[2].Name, "testtask2")
	assert.Equal(t, ps.Task[2].RevertTask, true)
}
