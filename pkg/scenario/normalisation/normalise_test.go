package normalisation

import (
	"testing"

	"github.com/docker/go-units"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/deserialize"
)

func TestScenario(t *testing.T) {
	// CLI

	scenario, err := deserialize.YamlScenario{}.FromBytes([]byte(`
nodes:
    - node_type: testnode
      ram: 2G
      cpu_percent: 42
      max_pods: 42
    - node_type: testnode2
      ram: 42G
      cpu_percent: 24
      max_pods: 24

node_groups:
    - group_name: testgroup1
      node_type: testnode
      amount: 42
    - group_name: testgroup2
      node_type: testnode2
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

	ps, nodes, err := NormaliseScenario(scenario.GetScenario())
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

	assert.Equal(t, 52, len(nodes))

	alreadySeenUUID := make(map[uuid.UUID]bool)
	alreadySeenType1 := 0
	alreadySeenType2 := 0

	for _, node := range nodes {
		assert.False(t, alreadySeenUUID[node.UUID])
		alreadySeenUUID[node.UUID] = true

		switch node.RAM {
		case 2 * units.GiB:
			assert.Equal(t, node.CPUPercent, 42)
			assert.Equal(t, node.MaxPods, 42)
			alreadySeenType1++
		case 42 * units.GiB:
			assert.Equal(t, node.CPUPercent, 24)
			assert.Equal(t, node.MaxPods, 24)
			alreadySeenType2++
		default:
			assert.Fail(t, "This unit doesn't exist")
		}
	}

	assert.Equal(t, 42, alreadySeenType1)
	assert.Equal(t, 10, alreadySeenType2)
}
