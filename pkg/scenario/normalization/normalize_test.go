package normalization

import (
	"testing"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/any"

	"github.com/docker/go-units"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/normalization/translate"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/deserialize"
)

func TestScenario(t *testing.T) {
	sc, err := deserialize.YamlScenario{}.FromBytes([]byte(`
nodes:
   -
       node_type: testnode
       memory: 2G
       cpu: 42
       storage: 2G
       ephemeral_storage: 2M
       max_pods: 42
   -
       node_type: testnode2
       memory: 42G
       storage: 22G
       ephemeral_storage: 21K
       cpu: 24
       max_pods: 24
node_groups:
   -
       group_name: testgroup1
       node_type: testnode
       amount: 42
   -
       group_name: testgroup2
       node_type: testnode2
       amount: 10
tasks:
   -
       name: testtask1
       time: 10s
       node_groups:
           - testgroup1
       node_response_state:
           type: GET_POD
           response: ERROR
           percentage: 14
   -
       name: testtask2
       time: 10s
       node_groups:
           - all
       node_response_state:
           type: DELETE_POD
           response: TIMEOUT
           percentage: 42
   -
       name: testtask2
       time: 20s
       revert: true
`))
	assert.NoError(t, err)

	getScenario, err := sc.GetScenario()
	assert.NoError(t, err)

	ps, nodes, err := NormalizeScenario(getScenario)
	assert.NoError(t, err)

	assert.Equal(t, 3, len(ps.Task))
	assert.Equal(t, false, ps.Task[0].RevertTask)
	assert.Equal(t, 42, len(ps.Task[0].NodeSet))
	assert.EqualValues(t, translate.EventFlags{
		events.NodeGetPodResponse:           any.MarshalOrDie(scenario.Response_ERROR),
		events.NodeGetPodResponsePercentage: any.MarshalOrDie(14),
	}, ps.Task[0].NodeEventFlags) // Is tested more in translator_test

	assert.Equal(t, false, ps.Task[1].RevertTask)
	assert.Equal(t, 52, len(ps.Task[1].NodeSet))
	assert.EqualValues(t, translate.EventFlags{
		events.NodeDeletePodResponse:           any.MarshalOrDie(scenario.Response_TIMEOUT),
		events.NodeDeletePodResponsePercentage: any.MarshalOrDie(42),
	}, ps.Task[1].NodeEventFlags)

	assert.Equal(t, true, ps.Task[2].RevertTask)
	assert.Equal(t, 52, len(ps.Task[2].NodeSet))
	assert.EqualValues(t, translate.EventFlags{
		events.NodeDeletePodResponse:           any.MarshalOrDie(scenario.Response_TIMEOUT),
		events.NodeDeletePodResponsePercentage: any.MarshalOrDie(42),
	}, ps.Task[2].NodeEventFlags)

	assert.Equal(t, 52, len(nodes))

	alreadySeenUUID := make(map[uuid.UUID]bool)
	alreadySeenType1 := 0
	alreadySeenType2 := 0

	for _, node := range nodes {
		assert.False(t, alreadySeenUUID[node.UUID])
		alreadySeenUUID[node.UUID] = true

		switch node.Memory {
		case 2 * units.GiB:
			assert.Equal(t, int64(42), node.CPU)
			assert.Equal(t, int64(42), node.MaxPods)
			assert.Equal(t, int64(2*units.GiB), node.Storage)
			assert.Equal(t, int64(2*units.MiB), node.EphemeralStorage)
			alreadySeenType1++
		case 42 * units.GiB:
			assert.Equal(t, int64(24), node.CPU)
			assert.Equal(t, int64(24), node.MaxPods)
			assert.Equal(t, int64(22*units.GiB), node.Storage)
			assert.Equal(t, int64(21*units.KiB), node.EphemeralStorage)
			alreadySeenType2++
		default:
			assert.Fail(t, "This unit doesn't exist")
		}
	}

	assert.Equal(t, 42, alreadySeenType1)
	assert.Equal(t, 10, alreadySeenType2)
}

func TestScenarioRevertUnknown(t *testing.T) {
	sc, err := deserialize.YamlScenario{}.FromBytes([]byte(`
nodes:
   -
       node_type: testnode
       memory: 2G
       cpu: 42
       storage: 2G
       ephemeral_storage: 2M
       max_pods: 42
node_groups:
   -
       group_name: testgroup1
       node_type: testnode
       amount: 42
tasks:
   -
       name: a
       time: 10s
       node_groups:
           - testgroup1
       node_failure: {}
   -
       name: b
       time: 10s
       revert: true
`))
	assert.NoError(t, err)

	getScenario, err := sc.GetScenario()
	assert.NoError(t, err)

	_, _, err = NormalizeScenario(getScenario)
	assert.Error(t, err, "you can't revert task with name 'b' as you have never used it before")
}

func TestScenarioSameNameTwice(t *testing.T) {
	sc, err := deserialize.YamlScenario{}.FromBytes([]byte(`
nodes:
   -
       node_type: testnode
       memory: 2G
       cpu: 42
       storage: 2G
       ephemeral_storage: 2M
       max_pods: 42
node_groups:
   -
       group_name: testgroup1
       node_type: testnode
       amount: 42
tasks:
   -
       name: a
       time: 10s
       node_groups:
           - testgroup1
       node_failure: {}
   -
       name: a
       time: 10s
       node_groups:
           - testgroup1
       node_failure: {}
`))
	assert.NoError(t, err)

	getScenario, err := sc.GetScenario()
	assert.NoError(t, err)

	_, _, err = NormalizeScenario(getScenario)
	assert.Error(t, err, "you can't use the task with name 'a' twice")
}

func TestScenarioRevertNameless(t *testing.T) {
	sc, err := deserialize.YamlScenario{}.FromBytes([]byte(`
nodes:
   -
       node_type: testnode
       memory: 2G
       cpu: 42
       storage: 2G
       ephemeral_storage: 2M
       max_pods: 42
node_groups:
   -
       group_name: testgroup1
       node_type: testnode
       amount: 42
tasks:
   -
       time: 10s
       revert: true
`))
	assert.NoError(t, err)

	getScenario, err := sc.GetScenario()
	assert.NoError(t, err)

	_, _, err = NormalizeScenario(getScenario)
	assert.Error(t, err, "you can't revert a task with an empty task name")
}

func TestScenarioNameless(t *testing.T) {
	sc, err := deserialize.YamlScenario{}.FromBytes([]byte(`
nodes:
   -
       node_type: testnode
       memory: 2G
       cpu: 42
       storage: 2G
       ephemeral_storage: 2M
       max_pods: 42
node_groups:
   -
       group_name: testgroup1
       node_type: testnode
       amount: 42
tasks:
   -
       time: 10s
       node_groups:
           - testgroup1
       node_failure: {}
   -
       time: 10s
       node_groups:
           - testgroup1
       node_failure: {}
`))
	assert.NoError(t, err)

	getScenario, err := sc.GetScenario()
	assert.NoError(t, err)

	_, _, err = NormalizeScenario(getScenario)
	assert.NoError(t, err)
}
