package scenario

import (
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/deserialize"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/normalize"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScenarioIntegration(t *testing.T) {

	// CLI

	scenario, err := deserialize.YamlScenario{}.FromBytes([]byte(`
nodes:
    - nodetype: testnode
      ram: 2G
      cpu_percent: 42
      max_pods: 42

nodegroups:
    - groupname: testgroup1
      nodetype: testnode
      amount: 42
    - groupname: testgroup2
      nodetype: testnode
      amount: 10

tasks:
    - name: testtask1
      time: 10s
      nodegroups: 
          - testgroup1

    - name: testtask2
      time: 10s
      nodegroups: 
          - all

    - name: testtask2
      time: 20s
      revert: true
`))
	assert.NoError(t, err)

	// Control plane

	nodecounter := 0

	var uuids []uuid.UUID

	normalize.IterNodes(scenario.GetScenario(), func(_ int) {
		nodecounter += 1

		// Nodes would be spawned here
		uuids = append(uuids, uuid.New())
	})

	assert.Equal(t, nodecounter, 52)

	ps, err := normalize.GetPrivateScenario(scenario.GetScenario(), uuids)
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