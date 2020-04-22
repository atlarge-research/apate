package deserialize

import (
	"context"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario/public"
)

type Deserializer interface {
	FromFile(filename string) (Deserializer, error)
	FromBytes(data []byte) (Deserializer, error)
	Send(client public.ScenarioSenderClient, ctx context.Context) (*public.SendScenarioResponse, error)
}
