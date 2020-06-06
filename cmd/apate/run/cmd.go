// Package run is the main package for the apate cli
package run

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/atlarge-research/opendc-emulate-kubernetes/internal/service"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/clients/controlplane"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/container"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"
)

type commandLineArgs struct {
	kubeConfigFileLocation string

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
	defaultControlPlaneTimeout = 300
)

func panicf(err error) {
	log.Panicf("An error occurred while running the CLI: %+v\n", err)
}

// StartCmd is the cmd entrypoint
func StartCmd(cmdArgs []string) {
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
						Name:        "manager-location",
						Usage:       "Manager config of cluster manager",
						TakesFile:   true,
						Destination: &cpEnv.ManagerConfigLocation,
						Value:       cpEnv.ManagerConfigLocation,
						Required:    false,
					},
					&cli.StringFlag{
						Name:        "kubeconfig-location",
						Usage:       "Location of the kubeconfig. If set, the managed cluster will be disabled",
						TakesFile:   true,
						Value:       args.kubeConfigFileLocation,
						Destination: &args.kubeConfigFileLocation,
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

	err := app.Run(cmdArgs)
	if err != nil {
		_, _ = color.New(color.FgRed).Printf("FAILED\nERROR: ")
		fmt.Printf("%+v\n", err)
	}
}

func printKubeConfig(ctx context.Context, args *commandLineArgs) error {
	client, err := controlplane.GetClusterOperationClient(service.NewConnectionInfo(args.controlPlaneAddress, args.controlPlanePort))
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

	if len(args.kubeConfigFileLocation) != 0 {
		bytes, err := ioutil.ReadFile(filepath.Clean(args.kubeConfigFileLocation))
		if err != nil {
			return errors.Wrapf(err, "failed to read kubeconfig from file at %v", args.kubeConfigFileLocation)
		}
		cpEnv.KubeConfig = string(bytes)
	}

	err := container.SpawnControlPlaneContainer(ctx, pp, cpEnv)
	if err != nil {
		return errors.Wrap(err, "couldn't spawn control plane container")
	}

	color.Green("DONE\n")
	fmt.Print("Waiting for control plane to be up ")

	// Polling control plane until up
	statusClient, _ := controlplane.GetStatusClient(service.NewConnectionInfo(cpEnv.ListenAddress, cpEnv.ListenPort))
	err = statusClient.WaitForControlPlane(ctx, time.Duration(args.controlPlaneTimeout)*time.Second)
	if err != nil {
		return errors.Wrap(err, "waiting for control plane on the client failed")
	}

	color.Green("DONE\n")
	fmt.Printf("Apate control plane created: %v\n", cpEnv)
	return statusClient.Conn.Close()
}

func runScenario(ctx context.Context, args *commandLineArgs) error {
	// The connectionInfo that will be used to connect to the control plane
	info := &service.ConnectionInfo{
		Address: args.controlPlaneAddress,
		Port:    args.controlPlanePort,
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
	if _, err = scenarioClient.Client.StartScenario(ctx, &empty.Empty{}); err != nil {
		return errors.Wrap(err, "couldn't start scenario")
	}
	err = scenarioClient.Conn.Close()
	if err != nil {
		return errors.Wrap(err, "couldn't close connection to scenario client")
	}

	color.Green("DONE\n")
	return statusClient.Conn.Close()
}
