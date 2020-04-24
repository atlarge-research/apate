package deserialize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeserializeJson(t *testing.T) {
	scenario, err := JSONScenario{}.FromBytes([]byte(`
{
	"nodes" : [
		{
			"node_type": "testnode",
			"ram": "2G",
			"cpu_percent": 42,
			"max_pods": 42
		}
	],
	"node_groups" : [
		{
			"group_name": "testgroup",
			"node_type": "testnode",
			"amount": 42
		}
	],
	"tasks" : [
		{
			"name": "testtask",
			"time": "10s",
			"node_groups": [
				"testgroup"
			]
		}
	]
}
	`))

	assert.NoError(t, err)

	assert.NotNil(t, scenario.GetScenario().NodeGroups)
	assert.NotNil(t, scenario.GetScenario().Nodes)
	assert.NotNil(t, scenario.GetScenario().Tasks)
	assert.Equal(t, 1, len(scenario.GetScenario().Tasks))

	assert.Equal(t, "testtask", scenario.GetScenario().Tasks[0].Name)
	assert.Equal(t, "10s", scenario.GetScenario().Tasks[0].Time)
	assert.Equal(t, 1, len(scenario.GetScenario().Tasks[0].NodeGroups))
	assert.Equal(t, "testgroup", scenario.GetScenario().Tasks[0].NodeGroups[0])

	assert.Equal(t, 1, len(scenario.GetScenario().NodeGroups))
	assert.Equal(t, "testgroup", scenario.GetScenario().NodeGroups[0].GroupName)
	assert.Equal(t, "testnode", scenario.GetScenario().NodeGroups[0].NodeType)
	assert.Equal(t, int32(42), scenario.GetScenario().NodeGroups[0].Amount)

	assert.Equal(t, 1, len(scenario.GetScenario().Nodes))

	assert.Equal(t, "testnode", scenario.GetScenario().Nodes[0].NodeType)
	assert.Equal(t, "2G", scenario.GetScenario().Nodes[0].Ram)
	assert.Equal(t, int32(42), scenario.GetScenario().Nodes[0].CpuPercent)
	assert.Equal(t, int32(42), scenario.GetScenario().Nodes[0].MaxPods)
}
