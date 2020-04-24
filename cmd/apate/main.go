package main

import (
	"os"

	"github.com/atlarge-research/opendc-emulate-kubernetes/cmd/apate/app"
)

func main() {
	app.Main(os.Args[1:])
}
