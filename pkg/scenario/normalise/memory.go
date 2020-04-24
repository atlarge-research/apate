package normalise

import (
	"github.com/docker/go-units"
)

func desugarMemory(memory string) (int64, error) {
	return units.RAMInBytes(memory)
}
