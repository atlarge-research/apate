// Package run is the entry point of the actual apatelet
package run

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"

	crdPod "github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/crd/pod"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster/kubeconfig"

	cli "github.com/virtual-kubelet/node-cli"

	healthpb "github.com/atlarge-research/opendc-emulate-kubernetes/api/health"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/clients/controlplane"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/clients/health"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/normalization"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
	vkProvider "github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/provider"
	vkService "github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/services"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
)

// KubeConfigWriter is the function used for writing the kube config to the local file system
// Because we only want to write it once and not on every apatelet, this should only be set in main.go, for standalone instances
// For goroutines, this should be written to file earlier
var KubeConfigWriter func(config []byte) = nil

// StartApatelet starts the apatelet
func StartApatelet(apateletEnv env.ApateletEnvironment, kubernetesPort, metricsPort int, readyCh *chan bool) error {
	log.Println("Starting Apatelet")

	// Retrieving connection information
	connectionInfo := service.NewConnectionInfo(apateletEnv.ControlPlaneAddress, apateletEnv.ControlPlanePort, false)
	ctx := context.Background()

	// Join the apate cluster
	log.Println("Joining apate cluster")
	config, res, err := joinApateCluster(ctx, connectionInfo, apateletEnv.ListenPort)
	if err != nil {
		return errors.Wrap(err, "failed to join apate cluster")
	}

	if KubeConfigWriter != nil {
		KubeConfigWriter(config.Bytes)
	}

	// Create store
	st := store.NewStore()

	// Create virtual kubelet
	errch := make(chan error)

	crdPod.CreateCRDInformer(config, &st, &errch)

	// Setup health status
	hc, err := health.GetClient(connectionInfo, res.UUID.String())
	if err != nil {
		return errors.Wrap(err, "failed to get health client")
	}
	hc.SetStatus(healthpb.Status_UNKNOWN)
	hc.StartStreamWithRetry(ctx, 3)

	// Start the Apatelet
	nc, cancel, err := createNodeController(ctx, res, kubernetesPort, metricsPort, &st)
	if err != nil {
		return errors.Wrap(err, "failed to create node controller")
	}

	log.Println("Joining kubernetes cluster")
	go func() {
		// TODO: Notify master / proper logging
		if err = nc.Run(ctx); err != nil {
			hc.SetStatus(healthpb.Status_UNHEALTHY)
			errch <- errors.Wrap(err, "failed to run node controller")
		}
	}()

	// Start gRPC server
	server, err := createGRPC(apateletEnv.ListenPort, &st, apateletEnv.ListenAddress)
	if err != nil {
		return errors.Wrap(err, "failed to set up GRPC endpoints")
	}

	// Update status
	hc.SetStatus(healthpb.Status_HEALTHY)
	log.Printf("now accepting requests on %s:%d\n", server.Conn.Address, server.Conn.Port)
	log.Printf("now listening on :%d for kube api and :%d for metrics", kubernetesPort, metricsPort)

	// Handle signals
	signals := make(chan os.Signal, 1)
	stopped := make(chan bool, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		shutdownErr := shutdown(ctx, server, cancel, connectionInfo, res.UUID.String())
		log.Println(shutdownErr)
		stopped <- true
	}()

	// Start serving request
	go server.Serve()

	*readyCh <- true

	// Stop the server on signal or error
	select {
	case err := <-errch:
		err = errors.Wrap(err, "Apatelet stopped because of an error")
		log.Println(err)
		return err
	case <-stopped:
		log.Println("Apatelet stopped")
		return nil
	}
}

func shutdown(ctx context.Context, server *service.GRPCServer, cancel context.CancelFunc, connectionInfo *service.ConnectionInfo, uuid string) error {
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

	log.Println("Stopping provider")
	cancel()

	return nil
}

func joinApateCluster(ctx context.Context, connectionInfo *service.ConnectionInfo, listenPort int) (*kubeconfig.KubeConfig, *normalization.NodeResources, error) {
	client, err := controlplane.GetClusterOperationClient(connectionInfo)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to get cluster operation client")
	}
	defer func() {
		closeErr := client.Conn.Close()
		if closeErr != nil {
			log.Printf("could not close connection: %v\n", closeErr)
		}
	}()

	cfg, res, err := client.JoinCluster(ctx, listenPort)

	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to join cluster")
	}

	log.Printf("Joined apate cluster with resources: %v", res)

	return cfg, res, nil
}

func createNodeController(ctx context.Context, res *normalization.NodeResources, k8sPort int, metricsPort int, store *store.Store) (*cli.Command, context.CancelFunc, error) {
	ctx, cancel := context.WithCancel(ctx)
	cmd, err := vkProvider.CreateProvider(ctx, res, k8sPort, metricsPort, store)
	if err != nil {
		cancel()
		return nil, nil, errors.Wrap(err, "failed to create provider")
	}

	return cmd, cancel, nil
}

func createGRPC(listenPort int, store *store.Store, listenAddress string) (*service.GRPCServer, error) {
	// Connection settings
	connectionInfo := service.NewConnectionInfo(listenAddress, listenPort, false)

	// Create gRPC server
	server, err := service.NewGRPCServer(connectionInfo)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new GRPC server")
	}

	// Add services
	vkService.RegisterScenarioService(server, store)

	return server, nil
}
