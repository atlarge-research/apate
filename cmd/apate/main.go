package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

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
	defaultControlPlaneAddress          = "localhost"
	defaultK8sConfigurationFileLocation = "./k8s.yml"
	defaultControlPlanePort             = 8085
	defaultControlPlaneTimeout          = 30
)

func main() {
	var scenarioFileLocation string
	var k8sConfigurationFileLocation string
	var controlPlaneAddress string
	var controlPlanePort int
	var controlPlaneTimeout int
	var pullPolicy string

	ctx := context.Background()
	env := container.DefaultControlPlaneEnvironment()

	app := &cli.App{
		Name:  "apate-cli",
		Usage: "Control the Apate control plane.",
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "Runs a given scenario file on the Apate cluster",
				Action: func(c *cli.Context) error {
					return runScenario(ctx, scenarioFileLocation, controlPlaneAddress, controlPlanePort, k8sConfigurationFileLocation)
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
						Required:    false,
					},
				},
			},
			{
				Name:  "create",
				Usage: "Creates a local control plane",
				Action: func(c *cli.Context) error {
					return createControlPlane(ctx, env, controlPlaneTimeout, pullPolicy)
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "address",
						Usage:       "Listen address of control plane",
						Destination: &env.Address,
						Value:       env.Address,
						Required:    false,
					},
					&cli.StringFlag{
						Name:        "port",
						Usage:       "The port of the control plane",
						Destination: &env.Port,
						Value:       env.Port,
						Required:    false,
					},
					&cli.StringFlag{
						Name:        "config",
						Usage:       "Manager config of cluster manager",
						Destination: &env.ManagerConfig,
						Value:       env.ManagerConfig,
						Required:    false,
					},
					&cli.StringFlag{
						Name:        "external-ip",
						Usage:       "IP used by apatelets to connect to control plane",
						Destination: &env.ExternalIP,
						Value:       env.ExternalIP,
						Required:    false,
					},
					&cli.StringFlag{
						Name:        "docker-policy-cp",
						Usage:       "Docker pull policy for control plane",
						Destination: &env.DockerPolicy,
						Value:       env.DockerPolicy,
						Required:    false,
					},
					&cli.StringFlag{
						Name:        "docker-policy",
						Usage:       "Docker pull policy used for creating the control plane",
						Destination: &pullPolicy,
						Value:       container.DefaultPullPolicy,
						Required:    false,
					},
					&cli.IntFlag{
						Name:        "timeout",
						Usage:       "Time before giving up on the control plane in seconds",
						Destination: &controlPlaneTimeout,
						Value:       defaultControlPlaneTimeout,
						Required:    false,
					},
				},
			},
			{
				Name:  "kubeconfig",
				Usage: "Retrieves a kube configuration file from the control plane",
				Action: func(c *cli.Context) error {
					return getKubeConfig(ctx, controlPlaneAddress, controlPlanePort)
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "address",
						Usage:       "The address of the control plane",
						Destination: &controlPlaneAddress,
						Value:       defaultControlPlaneAddress,
						Required:    false,
					},
					&cli.IntFlag{
						Name:        "port",
						Usage:       "The port of the control plane",
						Destination: &controlPlanePort,
						Value:       defaultControlPlanePort,
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

func getKubeConfig(ctx context.Context, address string, port int) error {
	cfg, err := controlplane.GetClusterOperationClient(service.NewConnectionInfo(address, port, false)).GetKubeConfig(ctx)

	if err != nil {
		return err
	}

	fmt.Println(string(cfg))
	return nil
}

func createControlPlane(ctx context.Context, env container.ControlPlaneEnvironment, timeout int, pullPolicy string) error {
	fmt.Print("Creating control plane container ")
	port, err := strconv.Atoi(env.Port)
	if err != nil {
		return err
	}

	err = container.SpawnControlPlane(ctx, pullPolicy, env)

	if err != nil {
		return err
	}
	color.Green("DONE\n")
	fmt.Print("Waiting for control plane to be up ")

	// Polling control plane until up
	statusClient := controlplane.GetStatusClient(service.NewConnectionInfo(env.Address, port, false))
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*time.Duration(timeout)))
	defer cancel()
	err = statusClient.WaitForControlPlane(ctx)

	if err != nil {
		return err
	}

	color.Green("DONE\n")
	fmt.Printf("Apate control plane created: %v\n", env)
	return nil
}

func runScenario(ctx context.Context, scenarioFileLocation string, controlPlaneAddress string, controlPlanePort int, configFileLocation string) error {
	var scenarioDeserializer deserialize.Deserializer
	var err error

	fmt.Printf("Reading scenario file")

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

	fmt.Printf("\rReading scenario file ")
	color.Green("DONE\n")

	// Read the k8s configuration file #nosec
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

	fmt.Printf("Loading scenario ")
	// Initial call: load the scenario
	scenarioClient := controlplane.GetScenarioClient(info)

	scenario, err := scenarioDeserializer.GetScenario()
	if err != nil {
		return err
	}

	// Add k8sconfig to the scenerio
	scenario.ResourceConfig = k8sConfig

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
