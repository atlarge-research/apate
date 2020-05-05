package main

import (
	"context"
	"io/ioutil"
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
	defaultK8sConfigurationFileLocation = "./k8s.yml"
	defaultControlPlanePort    = 8083
)

func main() {
	var scenarioFileLocation string
	var k8sConfigurationFileLocation string
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
					return runScenario(scenarioFileLocation, controlPlaneAddress, controlPlanePort, k8sConfigurationFileLocation)
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
					&cli.StringFlag{
						Name:        "k8s-config",
						Usage:       "The location of the kubernetes configuration for the resources to be created",
						EnvVars:     []string{"K8S_CONFIG_LOCATION"},
						Required:    true,
						Destination: &k8sConfigurationFileLocation,
						Value:       defaultK8sConfigurationFileLocation,
						DefaultText: defaultControlPlaneAddress,
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

func runScenario(scenarioFileLocation string, controlPlaneAddress string, controlPlanePort int, configFileLocation string) error {
	var scenarioDeserializer deserialize.Deserializer
	var err error

	if scenarioFileLocation == "-" {
		// Read the file given by stdin
		var bytes []byte

		bytes, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			return err
		}

		scenarioDeserializer, err = deserialize.YamlScenario{}.FromBytes(bytes)
	} else {
		// Read the file given by the argument
		scenarioDeserializer, err = deserialize.YamlScenario{}.FromFile(scenarioFileLocation)
	}

	if err != nil {
		return err
	}

	// Read the k8s configuration file
	k8sConfig, err := ioutil.ReadFile(configFileLocation)
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

	scenario, err := scenarioDeserializer.GetScenario()
	if err != nil {
		return err
	}

	// Add k8sconfig to the scenerio
	scenario.K8SConfiguration = k8sConfig

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
