package deserialize

import (
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario/public"
)

// A Deserializer is any struct that has the ability to take either
// a string or file encoded in some format, and convert that to a
// public.Scenario struct which can be sent over gRPC.
type Deserializer interface {
	// Loads from file. Usually opens and reads a file and then immediately calls FromBytes
	FromFile(filename string) (Deserializer, error)

	// Loads from a byte array. Useful for tests.
	FromBytes(data []byte) (Deserializer, error)

	// Gets the internal public.Scenario to for example be sent over gRPC.
	GetScenario() (public.Scenario)
}
