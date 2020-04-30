package events

import (
	"bytes"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/apatelet"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/controlplane"
	"github.com/golang/protobuf/jsonpb"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ghodss/yaml"
)

func TestScenario(t *testing.T) {
	task := translateYaml(t, []byte(`
node_failure: {}
`))
	newTask := &apatelet.Task{}
	err := NewEventTranslator(task, newTask).TranslateEvent()
	assert.NoError(t, err)
	//assert.Equal(t, controlplane.Task_NodeFailure{NodeFailure: &events.NodeFailure{}}, newTask.Event)
}

func translateYaml(t *testing.T, data []byte) *controlplane.Task {
	json, err := yaml.YAMLToJSON(data)
	assert.NoError(t, err)

	var task controlplane.Task
	err = jsonpb.Unmarshal(bytes.NewReader(json), &task)
	assert.NoError(t, err)

	return &task
}
