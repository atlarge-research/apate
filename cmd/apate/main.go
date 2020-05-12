package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/container"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"

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
	defaultControlPlaneTimeout = 45
)

func main() {
	var scenarioFileLocation string
	var k8sConfigurationFileLocation string
	var controlPlaneAddress string
	var controlPlanePort int
	var controlPlaneTimeout int
	var pullPolicy env.PullPolicy

	ctx := context.Background()
	cpEnv := env.DefaultControlPlaneEnvironment()

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
						Required:    false,
						Value:       "",
						Destination: &k8sConfigurationFileLocation,
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
					return createControlPlane(ctx, cpEnv, controlPlaneTimeout, pullPolicy)
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "address",
						Usage:       "Listen address of control plane",
						Destination: &cpEnv.Address,
						Value:       cpEnv.Address,
						Required:    false,
					},
					&cli.StringFlag{
						Name:        "port",
						Usage:       "The port of the control plane",
						Destination: &cpEnv.Port,
						Value:       cpEnv.Port,
						Required:    false,
					},
					&cli.StringFlag{
						Name:        "config",
						Usage:       "Manager config of cluster manager",
						Destination: &cpEnv.ManagerConfig,
						Value:       cpEnv.ManagerConfig,
						Required:    false,
					},
					&cli.StringFlag{
						Name:        "external-ip",
						Usage:       "IP used by apatelets to connect to control plane",
						Destination: &cpEnv.ExternalIP,
						Value:       cpEnv.ExternalIP,
						Required:    false,
					},
					&cli.StringFlag{
						Name:        "docker-policy-cp",
						Usage:       "Docker pull policy for control plane",
						Destination: &cpEnv.DockerPolicy,
						Value:       cpEnv.DockerPolicy,
						Required:    false,
					},
					&cli.StringFlag{
						Name:        "docker-policy",
						Usage:       "Docker pull policy used for creating the control plane",
						Destination: &pullPolicy,
						Value:       env.DefaultPullPolicy,
						Required:    false,
					},
					&cli.IntFlag{
						Name:        "timeout",
						Usage:       "Time before giving up on the control plane in seconds",
						Destination: &controlPlaneTimeout,
						Value:       defaultControlPlaneTimeout,
						Required:    false,
					},
					&cli.StringFlag{
						Name:        "runtype",
						Usage:       "How the control plane runs new apatelets. Can be DOCKER or ROUTINE.",
						Destination: &cpEnv.ApateletRunType,
						Value:       cpEnv.ApateletRunType,
						Required:    false,
					},
					&cli.StringFlag{
						Name:        "prometheus-enabled",
						Usage:       "If the control plane start a Prometheus stack. Can be TRUE or FALSE.",
						Destination: &cpEnv.PrometheusStackEnabled,
						Value:       cpEnv.PrometheusStackEnabled,
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

func createControlPlane(ctx context.Context, cpEnv env.ControlPlaneEnvironment, timeout int, pullPolicy env.PullPolicy) error {
	fmt.Print("Creating control plane container ")
	port, err := strconv.Atoi(cpEnv.Port)
	if err != nil {
		return err
	}

	err = container.SpawnControlPlaneContainer(ctx, pullPolicy, cpEnv)

	if err != nil {
		return err
	}
	color.Green("DONE\n")
	fmt.Print("Waiting for control plane to be up ")

	// Polling control plane until up
	statusClient := controlplane.GetStatusClient(service.NewConnectionInfo(cpEnv.Address, port, false))
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*time.Duration(timeout)))
	defer cancel()
	err = statusClient.WaitForControlPlane(ctx)

	if err != nil {
		return err
	}

	color.Green("DONE\n")
	fmt.Printf("Apate control plane created: %v\n", cpEnv)
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

	var k8sConfig []byte
	if len(configFileLocation) > 0 {
		// Read the k8s configuration file #nosec
		k8sConfig, err = ioutil.ReadFile(configFileLocation)
		if err != nil {
			return err
		}
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
	if _, err := scenarioClient.Client.StartScenario(ctx, &api.StartScenarioConfig{ResourceConfig: k8sConfig}); err != nil {
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
