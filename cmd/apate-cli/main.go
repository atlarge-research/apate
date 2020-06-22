package main

import (
	"os"

	"github.com/atlarge-research/apate/cmd/apate-cli/run"
)

func main() {
	run.StartCmd(os.Args)
}
