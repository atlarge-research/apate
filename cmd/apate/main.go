package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/container"

	"github.com/golang/protobuf/ptypes/empty"

	api "github.com/atlarge-research/opendc-emulate-kubernetes/api/controlplane"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/deserialize"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/clients/controlplane"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
)

const (
	defaultControlPlaneAddress = "localhost"
	defaultControlPlanePort    = 8085
)

func main() {
	var scenarioFileLocation string
	var controlPlaneAddress string
	var controlPlanePort int

	env := container.DefaultControlPlaneEnvironment()

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
			{
				Name:  "create",
				Usage: "Creates a local control plane",
				Action: func(c *cli.Context) error {
					return createControlPlane(env)
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "address",
						Usage:       "Listen address of control plane",
						Destination: &env.Address,
						Value:       env.Address,
						DefaultText: container.ControlPlaneListenAddressDefault,
						Required:    false,
					},
					&cli.StringFlag{
						Name:        "port",
						Usage:       "The port of the control plane",
						Destination: &env.Port,
						Value:       env.Port,
						DefaultText: container.ControlPlaneListenPortDefault,
						Required:    false,
					},
					&cli.StringFlag{
						Name:        "config",
						Usage:       "Manager config of cluster manager",
						Destination: &env.ManagerConfig,
						Value:       env.ManagerConfig,
						DefaultText: container.ManagedClusterConfigDefault,
						Required:    false,
					},
					&cli.StringFlag{
						Name:        "external-ip",
						Usage:       "IP used by apatelets to connect to control plane",
						Destination: &env.ExternalIP,
						Value:       env.ExternalIP,
						DefaultText: container.ControlPlaneExternalIPDefault,
						Required:    false,
					},
					&cli.StringFlag{
						Name:        "docker-policy",
						Usage:       "Docker pull policy",
						Destination: &env.DockerPolicy,
						Value:       env.DockerPolicy,
						DefaultText: container.DefaultPullPolicy,
						Required:    false,
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		_, _ = color.New(color.FgRed).Printf("FAILED\nERROR: ")
		fmt.Printf("%s\n", err.Error())
	}
}

func createControlPlane(env container.ControlPlaneEnvironment) error {
	fmt.Printf("%v\n", env)
	fmt.Println("Creating control plane container")

	err := container.SpawnControlPlane(context.Background(), container.PullIfNotLocal, env)

	color.Green("DONE\n")
	return err
}

func runScenario(scenarioFileLocation string, controlPlaneAddress string, controlPlanePort int) error {
	var deserializer deserialize.Deserializer
	var err error

	fmt.Printf("Reading scenario file")

	if scenarioFileLocation == "-" {
		// Read the file given by stdin
		var bytes []byte

		bytes, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			return err
		}

		deserializer, err = deserialize.YamlScenario{}.FromBytes(bytes)
	} else {
		// Read the file given by the argument
		deserializer, err = deserialize.YamlScenario{}.FromFile(scenarioFileLocation)
	}

	if err != nil {
		return err
	}

	fmt.Printf("\rReading scenario file ")
	color.Green("DONE\n")

	// The connectionInfo that will be used to connect to the control plane
	info := &service.ConnectionInfo{
		Address: controlPlaneAddress,
		Port:    controlPlanePort,
		TLS:     false,
	}

	ctx := context.Background()

	fmt.Printf("Loading scenario ")
	// Initial call: load the scenario
	scenarioClient := controlplane.GetScenarioClient(info)

	scenario, err := deserializer.GetScenario()
	if err != nil {
		return err
	}

	_, err = scenarioClient.Client.LoadScenario(ctx, scenario)
	if err != nil {
		return err
	}
	color.Green("DONE\n")

	// Next: keep polling until the control plane is happy
	expectedApatelets := getAmountOfApatelets(scenario)
	statusClient := controlplane.GetStatusClient(info)
	err = statusClient.WaitForHealthy(ctx, expectedApatelets, func(healthy int) {
		fmt.Printf("\rWaiting for healthy apatelets (%d/%d) ", healthy, expectedApatelets)
	})

	if err != nil {
		return err
	}

	color.Green("DONE\n")
	fmt.Printf("Starting scenario ")

	//Finally: actually start the scenario
	if _, err := scenarioClient.Client.StartScenario(ctx, new(empty.Empty)); err != nil {
		return err
	}

	color.Green("DONE\n")

	return nil
}

func getAmountOfApatelets(scenario *api.PublicScenario) int {
	var cnt int32

	for _, j := range scenario.NodeGroups {
		cnt += j.Amount
	}

	return int(cnt)
}
