package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/urfave/cli/v2"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/clients/controlplane"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/deserialize"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
)

const (
	defaultControlPlaneAddress = "localhost"
	defaultControlPlanePort    = 8083
)

func main() {
	var scenarioFileLocation string
	var controlPlaneAddress string
	var controlPlanePort int

	app := &cli.App{
		Name:  "apate-cli",
		Usage: "Control the Apate control plane.",
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "Runs a given scenario file on the Apate cluster",
				Action: func(c *cli.Context) error {
					return runScenario(scenarioFileLocation, controlPlaneAddress, controlPlanePort)
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "scenario",
						Usage:       "Load the scenario from `FILE`",
						Destination: &scenarioFileLocation,
						Required:    true,
					},
					&cli.StringFlag{
						Name:        "address",
						Usage:       "The address of the control plane",
						Destination: &controlPlaneAddress,
						Value:       defaultControlPlaneAddress,
						DefaultText: defaultControlPlaneAddress,
						Required:    false,
					},
					&cli.IntFlag{
						Name:        "port",
						Usage:       "The port of the control plane",
						Destination: &controlPlanePort,
						Value:       defaultControlPlanePort,
						DefaultText: strconv.Itoa(defaultControlPlanePort),
						Required:    false,
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func runScenario(scenarioFileLocation string, controlPlaneAddress string, controlPlanePort int) error {
	// Read the file given by the argument
	yaml, err := deserialize.YamlScenario{}.FromFile(scenarioFileLocation)
	if err != nil {
		return err
	}

	// The connectionInfo that will be used to connect to the control plane
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
		return err
	}

	_, err = scenarioClient.Client.LoadScenario(ctx, scenario)
	if err != nil {
		return err
	}

	// Next: keep polling until the control plane is happy
	statusClient := controlplane.GetStatusClient(info)
	_, err = statusClient.Client.Status(ctx, new(empty.Empty)) // TODO poll for status
	if err != nil {
		return err
	}

	// Finally: actually start the scenario
	if _, err := scenarioClient.Client.StartScenario(ctx, new(empty.Empty)); err != nil {
		return err
	}

	return nil
}
