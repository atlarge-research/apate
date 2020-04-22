package main

import (
	"github.com/atlarge-research/opendc-emulate-kubernetes/cli/apate/app"
	"os"
)

func main() {
	app.Main(os.Args[1:])
}