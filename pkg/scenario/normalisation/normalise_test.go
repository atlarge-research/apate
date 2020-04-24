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
	assert.Equal(t, int32(0), ps.StartTime)
	assert.Equal(t, 3, len(ps.Task))
	assert.Equal(t, "testtask1", ps.Task[0].Name)
	assert.Equal(t, false, ps.Task[0].RevertTask)
	assert.Equal(t, 42, len(ps.Task[0].NodeSet))

	assert.Equal(t, "testtask2", ps.Task[1].Name)
	assert.Equal(t, false, ps.Task[1].RevertTask)
	assert.Equal(t, 52, len(ps.Task[1].NodeSet))

	assert.Equal(t, "testtask2", ps.Task[2].Name)
	assert.Equal(t, true, ps.Task[2].RevertTask)

	assert.Equal(t, 52, len(nodes))

	alreadySeenUUID := make(map[uuid.UUID]bool)
	alreadySeenType1 := 0
	alreadySeenType2 := 0

	for _, node := range nodes {
		assert.False(t, alreadySeenUUID[node.UUID])
		alreadySeenUUID[node.UUID] = true

		switch node.RAM {
		case 2 * units.GiB:
			assert.Equal(t, 42, node.CPUPercent)
			assert.Equal(t, 42, node.MaxPods)
			alreadySeenType1++
		case 42 * units.GiB:
			assert.Equal(t, 24, node.CPUPercent)
			assert.Equal(t, 24, node.MaxPods)
			alreadySeenType2++
		default:
			assert.Fail(t, "This unit doesn't exist")
		}
	}

	assert.Equal(t, 42, alreadySeenType1)
	assert.Equal(t, 10, alreadySeenType2)
}
