// Package run is the entry point of the actual apatelet
package run

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"

	"github.com/google/uuid"

	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/scheduler"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario"

	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/crd/node"

	"github.com/pkg/errors"

	crdPod "github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/crd/pod"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes/kubeconfig"

	cli "github.com/virtual-kubelet/node-cli"

	healthpb "github.com/atlarge-research/opendc-emulate-kubernetes/api/health"
	"github.com/atlarge-research/opendc-emulate-kubernetes/internal/service"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/clients/controlplane"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/clients/health"
	vkProvider "github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/provider"
	vkService "github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/services"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
)

// KubeConfigWriter is the function used for writing the kube config to the local file system
// Because we only want to write it once and not on every apatelet, this should only be set in main.go, for standalone instances
// For goroutines, this should be written to file earlier
var KubeConfigWriter func(config []byte) = nil

// StartApatelet starts the apatelet
func StartApatelet(apateletEnv env.ApateletEnvironment, kubernetesPort, metricsPort int, readyCh chan<- struct{}) error {
	log.Println("Starting Apatelet")

	// Retrieving connection information
	connectionInfo := service.NewConnectionInfo(apateletEnv.ControlPlaneAddress, apateletEnv.ControlPlanePort, false)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create stop channel
	stop := make(chan os.Signal, 1)
	stopInformer := make(chan struct{})

	// Join the apate cluster
	log.Println("Joining apate cluster")
	config, res, startTime, err := joinApateCluster(ctx, connectionInfo, apateletEnv.ListenPort)
	if err != nil {
		return errors.Wrap(err, "failed to join apate cluster")
	}

	if KubeConfigWriter != nil {
		KubeConfigWriter(config.Bytes)
	}

	// Create store
	st := store.NewStore()

	// Create scheduler
	sch := scheduler.New(ctx, &st)
	go func() {
		ech := sch.EnableScheduler()

		for {
			select {
			case <-ctx.Done():
				return
			case err = <-ech:
				fmt.Printf("error while scheduling task occurred: %v\n", err)
			}
		}
	}()

	// Create crd informers
	err = crdPod.CreatePodInformer(config, &st, stopInformer, sch.WakeScheduler)
	if err != nil {
		return errors.Wrap(err, "failed creating crd pod informer")
	}

	err = node.CreateNodeInformer(config, &st, res.Selector, stopInformer, sch.WakeScheduler)
	if err != nil {
		return errors.Wrap(err, "failed creating crd node informer")
	}

	// Setup health status
	hc, err := startHealth(ctx, connectionInfo, res.UUID, stop)
	if err != nil {
		return errors.Wrap(err, "failed to start health client")
	}

	// Start the Apatelet
	nc, err := createNodeController(ctx, res, kubernetesPort, metricsPort, &st)
	if err != nil {
		return errors.Wrap(err, "failed to create node controller")
	}

	// Create virtual kubelet
	log.Println("Joining kubernetes cluster")
	errch := make(chan error)
	go func() {
		if err = nc.Run(ctx); err != nil {
			hc.SetStatus(healthpb.Status_UNHEALTHY)
			errch <- errors.Wrap(err, "failed to run node controller")
		}
	}()

	// Start gRPC server
	server, err := createGRPC(apateletEnv.ListenPort, &st, &sch, apateletEnv.ListenAddress, stop)
	if err != nil {
		return errors.Wrap(err, "failed to set up GRPC endpoints")
	}

	// Update status
	hc.SetStatus(healthpb.Status_HEALTHY)
	log.Printf("now accepting requests on %s:%d\n", server.Conn.Address, server.Conn.Port)
	log.Printf("now listening on :%d for kube api and :%d for metrics", kubernetesPort, metricsPort)

	// Handle stop
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Start serving request
	go func() {
		if err := server.Serve(); err != nil {
			errch <- errors.Wrap(err, "apatelet server failed")
		}
	}()

	// Start the scheduler if a scenario is already running
	if startTime >= 0 {
		sch.StartScheduler(startTime)
	}

	readyCh <- struct{}{}

	// Stop the server on signal or error
	select {
	case err := <-errch:
		err = errors.Wrap(err, "apatelet stopped because of an error")
		log.Println(err)
		return err
	case <-stop:
		stopInformer <- struct{}{}
		if err := shutdown(ctx, server, connectionInfo, res.UUID.String()); err != nil {
			log.Println(err)
		}
		return nil
	}
}

func shutdown(ctx context.Context, server *service.GRPCServer, connectionInfo *service.ConnectionInfo, uuid string) error {
	log.Println("Stopping Apatelet")

	log.Println("Stopping API")
	server.Server.Stop()

	log.Println("Leaving clusters (apate & k8s)")

	client, err := controlplane.GetClusterOperationClient(connectionInfo)
	if err != nil {
		return errors.Wrap(err, "failed to get cluster operation client")
	}
	defer func() {
		err := client.Conn.Close()
		if err != nil {
			log.Printf("could not close connection: %v\n", err)
		}
	}()

	if err := client.LeaveCluster(ctx, uuid); err != nil {
		log.Printf("An error occurred while leaving the clusters (apate & k8s): %v\n", err)
	}

	log.Println("Stopped Apatelet")

	return nil
}

func startHealth(ctx context.Context, connectionInfo *service.ConnectionInfo, uuid uuid.UUID, stop chan<- os.Signal) (*health.Client, error) {
	hc, err := health.GetClient(connectionInfo, uuid.String())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get health client")
	}
	hc.SetStatus(healthpb.Status_UNKNOWN)
	var retries int32 = 3
	hc.StartStream(ctx, func(err error) {
		if atomic.LoadInt32(&retries) < 1 {
			// Stop after retries amount of errors
			stop <- syscall.SIGTERM
			return
		}
		log.Println(err)
		atomic.AddInt32(&retries, -1)
	})
	return hc, nil
}

func joinApateCluster(ctx context.Context, connectionInfo *service.ConnectionInfo, listenPort int) (*kubeconfig.KubeConfig, *scenario.NodeResources, int64, error) {
	client, err := controlplane.GetClusterOperationClient(connectionInfo)
	if err != nil {
		return nil, nil, -1, errors.Wrap(err, "failed to get cluster operation client")
	}
	defer func() {
		closeErr := client.Conn.Close()
		if closeErr != nil {
			log.Printf("could not close connection: %v\n", closeErr)
		}
	}()

	cfg, res, startTime, err := client.JoinCluster(ctx, listenPort)

	if err != nil {
		return nil, nil, -1, errors.Wrap(err, "failed to join cluster")
	}

	log.Printf("Joined apate cluster with resources: %v", res)

	return cfg, res, startTime, nil
}

func createNodeController(ctx context.Context, res *scenario.NodeResources, k8sPort int, metricsPort int, store *store.Store) (*cli.Command, error) {
	cmd, err := vkProvider.CreateProvider(ctx, res, k8sPort, metricsPort, store)
	return cmd, errors.Wrap(err, "failed to create provider")
}

func createGRPC(listenPort int, store *store.Store, sch *scheduler.Scheduler, listenAddress string, stopChannel chan<- os.Signal) (*service.GRPCServer, error) {
	// Connection settings
	connectionInfo := service.NewConnectionInfo(listenAddress, listenPort, false)

	// Create gRPC server
	server, err := service.NewGRPCServer(connectionInfo)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new GRPC server")
	}

	// Add services
	vkService.RegisterScenarioService(server, store, sch)
	vkService.RegisterApateletService(server, stopChannel)

	return server, nil
}
