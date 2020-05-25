package main

import (
	"log"

	"github.com/pkg/errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/run"
)

func panicf(err error) {
	log.Panicf("an error occurred while starting the Apatelet: %+v\n", err)
}

func init() {
	// Enable line numbers in logging
	// Enables date time flags & file name + line
	log.SetFlags(log.LstdFlags | log.Llongfile)
}

// This function is only called when started separately (as Docker container for example)
// When starting from a goroutine, `StartApatelet` is called.
func main() {
	// Create Apatelet environment
	environment, err := env.ApateletEnv()
	if err != nil {
		panicf(errors.Wrap(err, "error while creating apatelet environment"))
	}

	// Set the certificates to communicate with the kubelet API
	if err := run.SetCerts(); err != nil {
		panicf(errors.Wrap(err, "error while setting certs"))
	}

	ch := make(chan struct{})
	if err := run.StartApatelet(environment, 10250, 10255, ch); err != nil {
		panicf(errors.Wrap(err, "error while running apatelet"))
	}
}
