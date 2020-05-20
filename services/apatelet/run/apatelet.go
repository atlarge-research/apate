// Package run is the entry point of the actual apatelet
package run

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario"

	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/crd/node"

	crdPod "github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/crd/pod"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster/kubeconfig"

	cli "github.com/virtual-kubelet/node-cli"

	healthpb "github.com/atlarge-research/opendc-emulate-kubernetes/api/health"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/clients/controlplane"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/clients/health"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create stop channel
	stop := make(chan os.Signal)
	stopInformer := make(chan struct{})

	// Join the apate cluster
	log.Println("Joining apate cluster")
	config, res, err := joinApateCluster(ctx, connectionInfo, apateletEnv.ListenPort)
	if err != nil {
		return err
	}

	if KubeConfigWriter != nil {
		KubeConfigWriter(config.Bytes)
	}

	// Create store
	st := store.NewStore()

	// Create crd informers
	err = crdPod.CreatePodInformer(config, &st, stopInformer)
	if err != nil {
		return err
	}
	err = node.CreateNodeInformer(config, &st, res.Selector, stopInformer)
	if err != nil {
		return err
	}

	// Setup health status
	hc := health.GetClient(connectionInfo, res.UUID.String())
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

	// Start the Apatelet
	nc, err := createNodeController(ctx, res, kubernetesPort, metricsPort, &st)
	if err != nil {
		return err
	}

	// Create virtual kubelet
	log.Println("Joining kubernetes cluster")
	errch := make(chan error)
	go func() {
		if err = nc.Run(ctx); err != nil {
			hc.SetStatus(healthpb.Status_UNHEALTHY)
			errch <- err
		}
	}()

	// Start gRPC server
	server, err := createGRPC(apateletEnv.ListenPort, &st, apateletEnv.ListenAddress, &stop)
	if err != nil {
		return err
	}

	// Update status
	hc.SetStatus(healthpb.Status_HEALTHY)
	log.Printf("now accepting requests on %s:%d\n", server.Conn.Address, server.Conn.Port)
	log.Printf("now listening on :%d for kube api and :%d for metrics", kubernetesPort, metricsPort)

	// Handle stop
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Start serving request
	go server.Serve()

	*readyCh <- true

	// Stop the server on signal or error
	select {
	case err := <-errch:
		log.Printf("Apatelet stopped because of an error %v\n", err)
		return err
	case <-stop:
		stopInformer <- struct{}{}
		shutdown(ctx, server, connectionInfo, res.UUID.String())
		return nil
	}
}

func shutdown(ctx context.Context, server *service.GRPCServer, connectionInfo *service.ConnectionInfo, uuid string) {
	log.Println("Stopping Apatelet")

	log.Println("Stopping API")
	server.Server.Stop()

	log.Println("Leaving clusters (apate & k8s)")

	client := controlplane.GetClusterOperationClient(connectionInfo)
	defer func() {
		err := client.Conn.Close()
		if err != nil {
			log.Printf("could not close connection: %v\n", err)
		}
	}()

	if err := client.LeaveCluster(ctx, uuid); err != nil {
		log.Printf("An error occurred while leaving the clusters (apate & k8s): %s", err.Error())
	}

	log.Println("Stopped Apatelet")
}

func joinApateCluster(ctx context.Context, connectionInfo *service.ConnectionInfo, listenPort int) (*kubeconfig.KubeConfig, *scenario.NodeResources, error) {
	client := controlplane.GetClusterOperationClient(connectionInfo)
	defer func() {
		err := client.Conn.Close()
		if err != nil {
			log.Printf("could not close connection: %v\n", err)
		}
	}()

	cfg, res, err := client.JoinCluster(ctx, listenPort)

	if err != nil {
		return nil, nil, err
	}

	log.Printf("Joined apate cluster with resources: %v", res)

	return cfg, res, nil
}

func createNodeController(ctx context.Context, res *scenario.NodeResources, k8sPort int, metricsPort int, store *store.Store) (*cli.Command, error) {
	cmd, err := vkProvider.CreateProvider(ctx, res, k8sPort, metricsPort, store)
	return cmd, err
}

func createGRPC(listenPort int, store *store.Store, listenAddress string, stopChannel *chan os.Signal) (*service.GRPCServer, error) {
	// Connection settings
	connectionInfo := service.NewConnectionInfo(listenAddress, listenPort, false)

	// Create gRPC server
	server, err := service.NewGRPCServer(connectionInfo)
	if err != nil {
		return nil, err
	}

	// Add services
	vkService.RegisterScenarioService(server, store)
	vkService.RegisterApateletService(server, stopChannel)

	return server, nil
}
