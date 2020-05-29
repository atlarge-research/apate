package main

import (
	"context"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/runner"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/run"
)

func main() {
	ctx := context.Background()
	run.StartControlPlane(ctx, runner.New())
}
