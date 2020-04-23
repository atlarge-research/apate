package main

import (
	"os"

	"github.com/atlarge-research/opendc-emulate-kubernetes/cli/apate/app"
)

func main() {
	app.Main(os.Args[1:])
}
