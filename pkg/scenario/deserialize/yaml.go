package deserialize

import (
	"io/ioutil"
	"sigs.k8s.io/yaml"
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
	json, err := yaml.JSONToYAML(data)
	if err != nil {
		return YamlScenario{}, err
	}

	return JsonScenario{}.FromBytes(json)
}