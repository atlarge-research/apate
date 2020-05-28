// Package runner exposes an interface to run new apatelets, and a registry to track all runners
package runner

import (
	"context"
	"sync"

	"github.com/pkg/errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"
)

// ApateletRunner is something that can be used to create new apatelets
type ApateletRunner interface {
	// SpawnApatelets spawns n apatelets
	SpawnApatelets(context.Context, int, env.ApateletEnvironment) error
}

type registration struct {
	runner *ApateletRunner
}

// RunnerRegistry is the registry that holds all runners
type Registry struct {
	sync.RWMutex
	registrations map[env.RunType]registration
}

// New returns a new registry
func New() *Registry {
	return &Registry{
		registrations: make(map[env.RunType]registration),
	}
}

// RegisterRunner registers a new runner
func (r *Registry) RegisterRunner(name env.RunType, runner *ApateletRunner) {
	r.Lock()
	defer r.Unlock()

	r.registrations[name] = registration{
		runner: runner,
	}
}

// Run call the appropriate ApateletRunner based on the current environment
func (r *Registry) Run(ctx context.Context, amount int, environment env.ApateletEnvironment) error {
	r.RLock()
	defer r.RUnlock()

	runType := env.ControlPlaneEnv().ApateletRunType
	if runner, ok := r.registrations[runType]; ok {
		return errors.Wrapf((*runner.runner).SpawnApatelets(ctx, amount, environment), "failed to spawn apatelet using runner %v", runner)
	}

	return errors.Errorf("unable to find runner type %v, have you registered it?", runType)
}
