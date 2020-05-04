package deserialize

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path/filepath"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/controlplane"
)

// JSONScenario is a struct with methods to deserialize JSON configurations.
type JSONScenario struct {
	scenario *controlplane.PublicScenario
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

	var scenario controlplane.PublicScenario
	if err := json.Unmarshal(data, &scenario); err != nil {
		return nil, err
	}

	cfp := customFlagParser{
		scenario: &scenario,
	}

	cfp.Get(data)

	return JSONScenario{&scenario}, nil
}

// GetScenario returns the inner stored public scenario.
func (s JSONScenario) GetScenario() (*controlplane.PublicScenario, error) {
	if s.scenario == nil {
		return nil, errors.New("scenario is not set yet")
	}
	return s.scenario, nil
}
