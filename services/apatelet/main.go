package main

import (
	"context"
	"log"
	"os"

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

	// This should be only set in ci to fix DinD problems
	k8sAddr := environment.CIKubernetesAddress
	if k8sAddr != "" {
		f, err := os.OpenFile("/etc/hosts", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			log.Printf("Error opening hosts file: %v\n", err)
		}
		if _, err = f.WriteString("\n" + k8sAddr + " docker"); err != nil {
			log.Printf("Error appending hosts file: %v\n", err)
		}
		if err := f.Close(); err != nil {
			log.Printf("Error closing hosts file: %v\n", err)
		}
	}

	// Set the certificates to communicate with the kubelet API
	if err := run.SetCerts(); err != nil {
		panicf(errors.Wrap(err, "error while setting certs"))
	}

	ch := make(chan struct{}, 1)
	if err := run.StartApatelet(context.Background(), environment, ch); err != nil {
		panicf(errors.Wrap(err, "error while running apatelet"))
	} else {
		<-ch
	}
}
