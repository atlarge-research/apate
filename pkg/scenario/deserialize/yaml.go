package deserialize

import (
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario/public"
	"github.com/ghodss/yaml"
	"io/ioutil"
)

type YamlScenario struct {
	JsonScenario
}

func (s YamlScenario) FromFile(filename string) (Deserializer, error) {

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return YamlScenario{}, err
	}

	return s.FromBytes(data)
}

func (s YamlScenario) FromBytes(data []byte) (Deserializer, error) {
	var scenario public.Scenario
	if err := yaml.Unmarshal(data, &scenario); err != nil {
		return JsonScenario{}, err
	}
	return YamlScenario{JsonScenario{scenario}}, nil
}