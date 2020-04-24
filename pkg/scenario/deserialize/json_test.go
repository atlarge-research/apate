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
			"cpu": 42,
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

	obtainedScenario, err := scenario.GetScenario()

	assert.NoError(t, err)

	assert.NotNil(t, obtainedScenario.NodeGroups)
	assert.NotNil(t, obtainedScenario.Nodes)
	assert.NotNil(t, obtainedScenario.Tasks)
	assert.Equal(t, 1, len(obtainedScenario.Tasks))

	assert.Equal(t, "testtask", obtainedScenario.Tasks[0].Name)
	assert.Equal(t, "10s", obtainedScenario.Tasks[0].Time)
	assert.Equal(t, 1, len(obtainedScenario.Tasks[0].NodeGroups))
	assert.Equal(t, "testgroup", obtainedScenario.Tasks[0].NodeGroups[0])

	assert.Equal(t, 1, len(obtainedScenario.NodeGroups))
	assert.Equal(t, "testgroup", obtainedScenario.NodeGroups[0].GroupName)
	assert.Equal(t, "testnode", obtainedScenario.NodeGroups[0].NodeType)
	assert.Equal(t, int32(42), obtainedScenario.NodeGroups[0].Amount)

	assert.Equal(t, 1, len(obtainedScenario.Nodes))

	assert.Equal(t, "testnode", obtainedScenario.Nodes[0].NodeType)
	assert.Equal(t, "2G", obtainedScenario.Nodes[0].RAM)
	assert.Equal(t, int32(42), obtainedScenario.Nodes[0].CPU)
	assert.Equal(t, int32(42), obtainedScenario.Nodes[0].MaxPods)
}
