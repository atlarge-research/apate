package main

import (
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/cmd"
	"log"
	"os"
)

func main() {
	if err := cmd.Run(os.Args[1:]); err != nil {
		log.Fatalf("An error occured: %s", err)
	}
}
