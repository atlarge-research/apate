package main

import (
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/run"
	"io/ioutil"
	"log"
	"strconv"
)

func init() {
	// Enable line numbers in logging
	// Enables date time flags & file name + line
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	run.SetCerts()

	controlPlaneAddress := env.RetrieveFromEnvironment(env.ControlPlaneAddress, env.ControlPlaneAddressDefault)

	controlPlanePort, err := strconv.Atoi(env.RetrieveFromEnvironment(env.ControlPlanePort, env.ControlPlanePortDefault))
	if err != nil {
		log.Fatalf("Error while starting apatelet: %s", err.Error())
	}

	// Retrieve own port
	listenPort, err := strconv.Atoi(env.RetrieveFromEnvironment(env.ApateletListenPort, env.ApateletListenPortDefault))
	if err != nil {
		log.Fatalf("Error while starting apatelet: %s", err.Error())
	}

	// Retrieving connection information
	listenAddress := env.RetrieveFromEnvironment(env.ApateletListenAddress, env.ApateletListenAddressDefault)

	run.KubeConfigWriter = func(config []byte) {
		err = ioutil.WriteFile("/config", config, 0600)
		if err != nil {
			panic(err)
		}
	}

	run.StartApatelet(controlPlaneAddress, controlPlanePort, listenAddress, listenPort, 10250, 10255)
}

