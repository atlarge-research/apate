package main

import (
	"io/ioutil"
	"log"

	"github.com/pkg/errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/run"
)

func fatal(err error) {
	log.Fatalf("an error occurred while starting the Apatelet: %+v\n", err)
}

func init() {
	// Enable line numbers in logging
	// Enables date time flags & file name + line
	log.SetFlags(log.LstdFlags | log.Llongfile)
}

// This function is only called when started separately (as Docker container for example)
// When starting from a goroutine, `StartApatelet` is called.
func main() {
	environment := env.ApateletEnv()

	if err := run.SetCerts(); err != nil {
		fatal(err)
	}

	run.KubeConfigWriter = func(config []byte) {
		kubeConfigLocation := env.ControlPlaneEnv().KubeConfigLocation
		err := ioutil.WriteFile(kubeConfigLocation, config, 0600)
		if err != nil {
			fatal(err)
		}
	}

	ch := make(chan struct{})
	if err := run.StartApatelet(environment, 10250, 10255, ch); err != nil {
		fatal(errors.Wrap(err, "error while running apatelet"))
	}
}
