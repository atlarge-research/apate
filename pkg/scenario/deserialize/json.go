package deserialize

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario/public"
)

// JSONScenario is a struct with methods to deserialize JSON configurations.
type JSONScenario struct {
	scenario *public.Scenario
}

// FromFile creates a new JSONScenario from a file in yaml format.
func (s JSONScenario) FromFile(filename string) (Deserializer, error) {
	data, err := ioutil.ReadFile(filepath.Clean(filename))
	if err != nil {
		return JSONScenario{}, err
	}

	return s.FromBytes(data)
}

// FromBytes creates a new JSONScenario from a byte array of data.
func (JSONScenario) FromBytes(data []byte) (Deserializer, error) {
	var scenario public.Scenario
	if err := json.Unmarshal(data, &scenario); err != nil {
		return JSONScenario{}, err
	}
	return JSONScenario{&scenario}, nil
}

// GetScenario returns the inner stored public scenario.
func (s JSONScenario) GetScenario() *public.Scenario {
	return s.scenario
}
