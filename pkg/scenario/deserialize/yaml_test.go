package deserialize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeserializeYaml(t *testing.T) {
	scenario, err := YamlScenario{}.FromBytes([]byte(`
nodes:
    - node_type: testnode
      ram: 2G
      cpu_percent: 42
      max_pods: 42

node_groups:
    - group_name: testgroup
      node_type: testnode
      amount: 42

tasks:
    - name: testtask
      time: 10s
      node_groups: 
          - testgroup
`))

	assert.NoError(t, err)

	assert.NotNil(t, scenario.GetScenario().NodeGroups)
	assert.NotNil(t, scenario.GetScenario().Nodes)
	assert.NotNil(t, scenario.GetScenario().Tasks)
	assert.Equal(t, len(scenario.GetScenario().Tasks), 1)

	assert.Equal(t, scenario.GetScenario().Tasks[0].Name, "testtask")
	assert.Equal(t, scenario.GetScenario().Tasks[0].Time, "10s")
	assert.Equal(t, len(scenario.GetScenario().Tasks[0].NodeGroups), 1)
	assert.Equal(t, scenario.GetScenario().Tasks[0].NodeGroups[0], "testgroup")

	assert.Equal(t, len(scenario.GetScenario().NodeGroups), 1)
	assert.Equal(t, scenario.GetScenario().NodeGroups[0].GroupName, "testgroup")
	assert.Equal(t, scenario.GetScenario().NodeGroups[0].NodeType, "testnode")
	assert.Equal(t, scenario.GetScenario().NodeGroups[0].Amount, int32(42))

	assert.Equal(t, len(scenario.GetScenario().Nodes), 1)

	assert.Equal(t, scenario.GetScenario().Nodes[0].NodeType, "testnode")
	assert.Equal(t, scenario.GetScenario().Nodes[0].Ram, "2G")
	assert.Equal(t, scenario.GetScenario().Nodes[0].CpuPercent, int32(42))
	assert.Equal(t, scenario.GetScenario().Nodes[0].MaxPods, int32(42))
}
