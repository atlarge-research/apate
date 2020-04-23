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
			"nodetype": "testnode",
			"ram": "2G",
			"cpu_percent": 42,
			"max_pods": 42
		}
	],
	"nodegroups" : [
		{
			"groupname": "testgroup",
			"nodetype": "testnode",
			"amount": 42
		}
	],
	"tasks" : [
		{
			"name": "testtask",
			"time": "10s",
			"nodegroups": [
				"testgroup"
			]
		}
	]
}
	`))

	assert.NoError(t, err)

	assert.NotNil(t, scenario.GetScenario().Nodegroups)
	assert.NotNil(t, scenario.GetScenario().Nodes)
	assert.NotNil(t, scenario.GetScenario().Tasks)
	assert.Equal(t, len(scenario.GetScenario().Tasks), 1)

	assert.Equal(t, scenario.GetScenario().Tasks[0].Name, "testtask")
	assert.Equal(t, scenario.GetScenario().Tasks[0].Time, "10s")
	assert.Equal(t, len(scenario.GetScenario().Tasks[0].Nodegroups), 1)
	assert.Equal(t, scenario.GetScenario().Tasks[0].Nodegroups[0], "testgroup")

	assert.Equal(t, len(scenario.GetScenario().Nodegroups), 1)
	assert.Equal(t, scenario.GetScenario().Nodegroups[0].Groupname, "testgroup")
	assert.Equal(t, scenario.GetScenario().Nodegroups[0].Nodetype, "testnode")
	assert.Equal(t, scenario.GetScenario().Nodegroups[0].Amount, int32(42))

	assert.Equal(t, len(scenario.GetScenario().Nodes), 1)

	assert.Equal(t, scenario.GetScenario().Nodes[0].Nodetype, "testnode")
	assert.Equal(t, scenario.GetScenario().Nodes[0].Ram, "2G")
	assert.Equal(t, scenario.GetScenario().Nodes[0].CpuPercent, int32(42))
	assert.Equal(t, scenario.GetScenario().Nodes[0].MaxPods, int32(42))
}
