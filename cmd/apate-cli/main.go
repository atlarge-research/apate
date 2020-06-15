package main

import (
	"os"

	"github.com/atlarge-research/opendc-emulate-kubernetes/cmd/apate-cli/run"
)

func main() {
	run.StartCmd(os.Args)
}
