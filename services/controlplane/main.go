package main

import (
	"context"
	"log"
	"net/http"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"

	// #nosec exposing debug statistics is not a problem for this application
	_ "net/http/pprof"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/runner"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/run"
)

func main() {
	cpEnv := env.ControlPlaneEnv()

	// Start debug server if debug is enabled
	if cpEnv.DebugEnabled {
		go log.Println(http.ListenAndServe("localhost:6060", nil))
	}

	ctx := context.Background()
	run.StartControlPlane(ctx, runner.New())
}
