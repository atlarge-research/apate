package main

import (
	"io/ioutil"
	"log"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/run"
)

func init() {
	// Enable line numbers in logging
	// Enables date time flags & file name + line
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	environment, err := env.ApateletEnvironmentFromEnv()
	if err != nil {
		log.Fatalf("Error while starting apatelet: %s", err.Error())
	}

	if err = run.SetCerts(); err != nil {
		log.Fatal(err)
	}

	run.KubeConfigWriter = func(config []byte) {
		err = ioutil.WriteFile("/config", config, 0600)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = run.StartApatelet(environment, 10250, 10255)
	if err != nil {
		log.Fatalf("Error while running apatelet: %s", err.Error())
	}
}
