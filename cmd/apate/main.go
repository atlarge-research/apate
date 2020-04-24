package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/clients/controlplane"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/deserialize"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
)

const (
	defaultControlPlaneAddress = "localhost"
	defaultControlPlanePort    = 8083
)

func main() {
	var err error

	// TODO: Do arg parsing here with a proper library.
	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatalf("Please give a scenario filename as first argument")
	}

	log.Printf("Parsing yaml file")

	yaml, err := deserialize.YamlScenario{}.FromFile(args[0])
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Dialling server")

	// TODO: Do arg parsing here with a proper library.
	var controlPlaneAddress string
	if len(args) == 1 {
		controlPlaneAddress = defaultControlPlaneAddress
	} else {
		controlPlaneAddress = args[1]
	}

	var controlPlanePort int
	if len(args) == 2 {
		controlPlanePort = defaultControlPlanePort
	} else {
		controlPlanePort, err = strconv.Atoi(args[2])
		if err != nil {
			log.Fatal(err)
		}
	}

	info := &service.ConnectionInfo{
		Address: controlPlaneAddress,
		Port:    controlPlanePort,
		TLS:     false,
	}

	ctx := context.Background()

	// Initial call: load the scenario
	scenarioClient := controlplane.GetScenarioClient(info)

	scenario, err := yaml.GetScenario()
	if err != nil {
		log.Fatal(err)
	}

	_, err = scenarioClient.Client.LoadScenario(ctx, scenario)
	if err != nil {
		log.Fatal(err)
	}

	statusClient := controlplane.GetStatusClient(info)
	_, err = statusClient.Client.Status(ctx, new(empty.Empty)) // TODO poll for status
	if err != nil {
		log.Fatal(err)
	}

	if _, err := scenarioClient.Client.StartScenario(ctx, new(empty.Empty)); err != nil {
		log.Fatal(err)
	}
}
