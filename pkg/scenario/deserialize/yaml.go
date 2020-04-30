package deserialize

import (
	"bytes"
	"github.com/golang/protobuf/jsonpb"
	"io/ioutil"
	"path/filepath"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/controlplane"

	"github.com/ghodss/yaml"
)

// YamlScenario is a struct with methods to deserialize Yaml configurations.
type YamlScenario struct {
	JSONScenario
}

// FromFile creates a new YamlScenario from a file in yaml format.
func (s YamlScenario) FromFile(filename string) (Deserializer, error) {
	data, err := ioutil.ReadFile(filepath.Clean(filename))
	if err != nil {
		return YamlScenario{}, err
	}

	return s.FromBytes(data)
}

// FromBytes creates a new YamlScenario from a byte array of data.
func (s YamlScenario) FromBytes(data []byte) (Deserializer, error) {

	json, err := yaml.YAMLToJSON(data)
	if err != nil {
		return nil, err
	}

	var scenario controlplane.PublicScenario
	if err := jsonpb.Unmarshal(bytes.NewReader(json), &scenario); err != nil {
		return nil, err
	}

	return YamlScenario{JSONScenario{&scenario}}, nil
}
