package main

import (
	"context"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/runner"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/run"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	ctx := context.Background()
	run.StartControlPlane(ctx, runner.New())
}
