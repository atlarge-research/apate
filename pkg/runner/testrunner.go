package runner

import (
	"context"
	"github.com/pkg/errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"
	apateRun "github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/run"
)

// RoutineRunner spawns apatelets using go routines
type TestRoutineRunner struct{}

// SpawnApatelets spawns apatelets using go routines
func (d *TestRoutineRunner) SpawnApatelets(_ context.Context, _ int, _ env.ApateletEnvironment) error {
	if err := apateRun.SetCerts(); err != nil {
		return errors.Wrap(err, "failed to set certificates")
	}


	return nil
}
