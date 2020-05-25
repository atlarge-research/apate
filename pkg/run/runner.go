// Package run exposes an interface to run new apatelets, and a registry to track all runners
package run

import (
	"context"
	"sync"

	"github.com/pkg/errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"
)

// ApateletRunner is something that can be used to create new apatelets
type ApateletRunner interface {
	// SpawnApatelets spawns n apatelets
	SpawnApatelets(context.Context, int, env.ApateletEnvironment, ...interface{}) error
}

type registration struct {
	args   []interface{}
	runner *ApateletRunner
}

// RunnerRegistry is the registry that holds all runners
type RunnerRegistry struct {
	sync.RWMutex
	registrations map[env.RunType]registration
}

// New returns a new registry
func New() *RunnerRegistry {
	return &RunnerRegistry{
		registrations: make(map[env.RunType]registration),
	}
}

// RegisterRunner registers a new runner
func (rr *RunnerRegistry) RegisterRunner(name env.RunType, runner *ApateletRunner, args ...interface{}) {
	rr.Lock()
	defer rr.Unlock()

	rr.registrations[name] = registration{
		args:   args,
		runner: runner,
	}
}

// Run call the appropriate ApateletRunner based on the current environment
func (rr *RunnerRegistry) Run(ctx context.Context, amount int, environment env.ApateletEnvironment) error {
	rr.RLock()
	defer rr.RUnlock()

	runType := env.ControlPlaneEnv().ApateletRunType
	if runner, ok := rr.registrations[runType]; ok {
		return errors.Wrapf((*runner.runner).SpawnApatelets(ctx, amount, environment, runner.args...), "failed to spawn apatelet using runner %v", runner)
	}

	return errors.Errorf("unable to find runner type %v, have you registered it?", runType)
}
