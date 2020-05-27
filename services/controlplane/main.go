package main

import (
	"context"

	"github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/app"
)

func main() {
	ctx := context.Background()
	app.Main(ctx)
}
