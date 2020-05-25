package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/container"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"

	api "github.com/atlarge-research/opendc-emulate-kubernetes/api/controlplane"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"

	"github.com/atlarge-research/opendc-emulate-kubernetes/internal/service"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/clients/controlplane"
)

type commandLineArgs struct {
	k8sConfigurationFileLocation string

	controlPlaneAddress string
	controlPlanePort    int
	controlPlaneTimeout int

	apateletRunType        string
	pullPolicyControlPlane string
	pullPolicyCreate       string
}

const (
	defaultControlPlaneAddress = "localhost"
	defaultControlPlanePort    = 8085
	defaultControlPlaneTimeout = 45
)

func panicf(err error) {
	log.Panicf("An error occurred while running the CLI: %+v\n", err)
}

func main() {
	args := &commandLineArgs{}

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
					return errors.Wrap(runScenario(ctx, args), "failed to run scenario")
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "address",
						Usage:       "The address of the control plane",
						Destination: &args.controlPlaneAddress,
						Value:       defaultControlPlaneAddress,
						Required:    false,
					},
					&cli.StringFlag{
						Name:        "k8s-config",
						Usage:       "The location of the kubernetes configuration for the resources to be created",
						EnvVars:     []string{"K8S_CONFIG_LOCATION"},
						Required:    false,
						Value:       "",
						Destination: &args.k8sConfigurationFileLocation,
					},
					&cli.IntFlag{
						Name:        "port",
						Usage:       "The port of the control plane",
						Destination: &args.controlPlanePort,
						Value:       defaultControlPlanePort,
						Required:    false,
					},
				},
			},
			{
				Name:  "create",
				Usage: "Creates a local control plane",
				Action: func(c *cli.Context) error {
					return errors.Wrap(createControlPlane(ctx, cpEnv, args), "failed to create control plane")
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "address",
						Usage:       "Listen address of control plane",
						Destination: &cpEnv.ListenAddress,
						Value:       cpEnv.ListenAddress,
						Required:    false,
					},
					&cli.IntFlag{
						Name:        "port",
						Usage:       "The port of the control plane",
						Destination: &cpEnv.ListenPort,
						Value:       cpEnv.ListenPort,
						Required:    false,
					},
					&cli.StringFlag{
						Name:        "config",
						Usage:       "Manager config of cluster manager",
						Destination: &cpEnv.ManagerConfigLocation,
						Value:       cpEnv.ManagerConfigLocation,
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
						Destination: &args.pullPolicyControlPlane,
						Value:       string(cpEnv.DockerPolicy),
						Required:    false,
					},
					&cli.StringFlag{
						Name:        "docker-policy",
						Usage:       "Docker pull policy used for creating the control plane",
						Destination: &args.pullPolicyCreate,
						Value:       string(env.DefaultPullPolicy),
						Required:    false,
					},
					&cli.IntFlag{
						Name:        "timeout",
						Usage:       "Time before giving up on the control plane in seconds",
						Destination: &args.controlPlaneTimeout,
						Value:       defaultControlPlaneTimeout,
						Required:    false,
					},
					&cli.StringFlag{
						Name:        "runtype",
						Usage:       "How the control plane runs new apatelets. Can be DOCKER or ROUTINE.",
						Destination: &args.apateletRunType,
						Value:       string(cpEnv.ApateletRunType),
						Required:    false,
					},
					&cli.BoolFlag{
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
					return errors.Wrap(printKubeConfig(ctx, args), "failed to get Kubeconfig")
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "address",
						Usage:       "The address of the control plane",
						Destination: &args.controlPlaneAddress,
						Value:       defaultControlPlaneAddress,
						Required:    false,
					},
					&cli.IntFlag{
						Name:        "port",
						Usage:       "The port of the control plane",
						Destination: &args.controlPlanePort,
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
		fmt.Printf("%+v\n", err)
	}
}

func printKubeConfig(ctx context.Context, args *commandLineArgs) error {
	client, err := controlplane.GetClusterOperationClient(service.NewConnectionInfo(args.controlPlaneAddress, args.controlPlanePort, false))
	if err != nil {
		return errors.Wrap(err, "couldn't get cluster operation client for kube config")
	}

	cfg, err := client.GetKubeConfig(ctx)
	if err != nil {
		return errors.Wrap(err, "couldn't get kube config from control plane")
	}

	fmt.Println(string(cfg))

	if err := client.Conn.Close(); err != nil {
		return errors.Wrap(err, "error closing connection to cluster operation client")
	}

	return nil
}

func createControlPlane(ctx context.Context, cpEnv env.ControlPlaneEnvironment, args *commandLineArgs) error {
	fmt.Print("Creating control plane container ")

	pp := env.PullPolicy(args.pullPolicyCreate)
	if !pp.Valid() {
		return errors.Errorf("invalid pull policy %v", cpEnv.DockerPolicy)
	}

	cpEnv.DockerPolicy = env.PullPolicy(args.pullPolicyControlPlane)
	if !cpEnv.DockerPolicy.Valid() {
		return errors.Errorf("invalid pull policy for control plane %v", cpEnv.DockerPolicy)
	}

	cpEnv.ApateletRunType = env.RunType(args.apateletRunType)

	err := container.SpawnControlPlaneContainer(ctx, pp, cpEnv)
	if err != nil {
		return errors.Wrap(err, "couldn't spawn control plane container")
	}

	color.Green("DONE\n")
	fmt.Print("Waiting for control plane to be up ")

	// Polling control plane until up
	statusClient, _ := controlplane.GetStatusClient(service.NewConnectionInfo(cpEnv.ListenAddress, cpEnv.ListenPort, false))
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*time.Duration(args.controlPlaneTimeout)))
	defer cancel()
	err = statusClient.WaitForControlPlane(ctx)
	if err != nil {
		return errors.Wrap(err, "waiting for control plane on the client failed")
	}

	color.Green("DONE\n")
	fmt.Printf("Apate control plane created: %v\n", cpEnv)
	return nil
}

func runScenario(ctx context.Context, args *commandLineArgs) error {
	k8sConfig, err := func() ([]byte, error) {
		if len(args.k8sConfigurationFileLocation) > 0 {
			// #nosec
			k8sConfig, err := ioutil.ReadFile(args.k8sConfigurationFileLocation)
			if err != nil {
				return nil, errors.Wrap(err, "reading k8sconfig failed")
			}
			return k8sConfig, nil
		}
		return []byte{}, nil
	}()

	// The connectionInfo that will be used to connect to the control plane
	info := &service.ConnectionInfo{
		Address: args.controlPlaneAddress,
		Port:    args.controlPlanePort,
		TLS:     false,
	}

	// Initial call: load the scenario
	scenarioClient, err := controlplane.GetScenarioClient(info)
	if err != nil {
		return errors.Wrap(err, "failed to get scenario client")
	}

	// Next: poll amount of healthy nodes
	trigger := make(chan struct{})

	go func() {
		_, err = bufio.NewReader(os.Stdin).ReadBytes('\n')
		if err != nil {
			panicf(err)
		}
		trigger <- struct{}{}
	}()

	statusClient, err := controlplane.GetStatusClient(info)
	if err != nil {
		return errors.Wrap(err, "getting status client for runScenario failed")
	}
	err = statusClient.WaitForTrigger(ctx, trigger, func(healthy int) {
		fmt.Printf("\rGot %d healthy apatelets - Press enter to start scenario...", healthy)
	})
	if err != nil {
		return errors.Wrap(err, "waiting for healthy Apatelets failed")
	}

	fmt.Printf("Starting scenario ")

	//Finally: actually start the scenario
	if _, err = scenarioClient.Client.StartScenario(ctx, &api.StartScenarioConfig{ResourceConfig: k8sConfig}); err != nil {
		return errors.Wrap(err, "couldn't start scenario")
	}
	err = scenarioClient.Conn.Close()
	if err != nil {
		return errors.Wrap(err, "couldn't close connection to scenario client")
	}

	color.Green("DONE\n")
	return nil
}
