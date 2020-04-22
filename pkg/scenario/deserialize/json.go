package deserialize

import (
	"context"
	"encoding/json"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario/public"
	"io/ioutil"
)

type JsonScenario struct {
	scenario public.Scenario
}

func (s JsonScenario) FromFile(filename string) (Deserializer, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return JsonScenario{}, err
	}

	return s.FromBytes(data)
}

func (JsonScenario) FromBytes(data []byte) (Deserializer, error) {
	var scenario public.Scenario
	if err := json.Unmarshal(data, &scenario); err != nil {
		return JsonScenario{}, err
	}
	return JsonScenario{scenario}, nil
}

func (s JsonScenario) Send(client public.ScenarioSenderClient,  ctx context.Context) (*public.SendScenarioResponse, error) {
	return client.SendScenario(ctx, &s.scenario)
}