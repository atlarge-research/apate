package main

import (
	"io/ioutil"
	"log"
	"os"

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
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// This function is only called when started separately (as Docker container for example)
// When starting from a goroutine, `StartApatelet` is called.
func main() {
	environment, err := env.ApateletEnvironmentFromEnv()
	if err != nil {
		fatal(errors.Wrap(err, "error while starting apatelet"))
	}

	if err = run.SetCerts(); err != nil {
		fatal(err)
	}

	run.KubeConfigWriter = func(config []byte) {
		err = ioutil.WriteFile(os.TempDir()+"/apate/config", config, 0600)
		if err != nil {
			fatal(err)
		}
	}

	ch := make(chan bool)
	err = run.StartApatelet(environment, 10250, 10255, &ch)
	if err != nil {
		fatal(errors.Wrap(err, "error while running apatelet"))
	}
}
