package main

import (
	"strconv"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/container"

	"github.com/virtual-kubelet/virtual-kubelet/node"
	"k8s.io/client-go/kubernetes"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/clients/controlplane"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/normalization"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
	vkProvider "github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/provider"
	vkService "github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/services"

	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	k8sVersion = "v1.15.2" // This should follow the version of k8s.io/kubernetes we are importing
)

func init() {
	// Enable line numbers in logging
	// Enables date time flags & file name + line
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	log.Println("Starting Apatelet")

	// Retrieving connection information
	controlPlaneAddress := container.RetrieveFromEnvironment(container.ControlPlaneAddress, container.ControlPlaneAddressDefault)
	controlPlanePort, err := strconv.Atoi(container.RetrieveFromEnvironment(container.ControlPlanePort, container.ControlPlanePortDefault))

	if err != nil {
		log.Fatalf("Error while starting apatelet: %s", err.Error())
	}

	connectionInfo := service.NewConnectionInfo(controlPlaneAddress, controlPlanePort, false)
	ctx := context.Background()

	// Retrieve own port
	listenPort, err := strconv.Atoi(container.RetrieveFromEnvironment(container.ApateletListenPort, container.ApateletListenPortDefault))

	if err != nil {
		log.Fatalf("Error while starting apatelet: %s", err.Error())
	}

	// Join the apate cluster
	log.Println("Joining apate cluster")
	kubeConfig, res := joinApateCluster(ctx, connectionInfo, listenPort)

	// Setup health status
	//hc := health.GetClient(connectionInfo, res.UUID.String())
	//hc.SetStatus(healthpb.Status_UNKNOWN)
	//hc.StartStreamWithRetry(ctx, 3)

	// Start the Apatelet
	ctx, nc, cancel := createNodeController(ctx, kubeConfig, res)

	log.Println("Joining kubernetes cluster")
	go func() {
		// TODO: Notify master / proper logging
		if err = nc.Run(ctx); err != nil {
			//hc.SetStatus(healthpb.Status_UNHEALTHY)
			log.Fatalf("Unable to start apatelet: %v", err)
		}
	}()

	// Start gRPC server
	server := createGRPC(listenPort)

	// Update status
	//hc.SetStatus(healthpb.Status_HEALTHY)
	log.Printf("Now accepting requests on %s:%d\n", server.Conn.Address, server.Conn.Port)

	// Handle signals
	signals := make(chan os.Signal, 1)
	stopped := make(chan bool, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		shutdown(ctx, server, cancel, connectionInfo, res.UUID.String())
		stopped <- true
	}()

	// Start serving request
	server.Serve()

	// Stop the server on signal
	<-stopped
	log.Println("Apatelet stopped")
}

func shutdown(ctx context.Context, server *service.GRPCServer, cancel context.CancelFunc, connectionInfo *service.ConnectionInfo, uuid string) {
	log.Println("Stopping Apatelet")

	log.Println("Stopping API")
	server.Server.Stop()

	log.Println("Leaving clusters (apate & k8s)")

	client := controlplane.GetClusterOperationClient(connectionInfo)
	defer func() {
		_ = client.Conn.Close()
	}()

	if err := client.LeaveCluster(ctx, uuid); err != nil {
		log.Printf("An error occurred while leaving the clusters (apate & k8s): %s", err.Error())
	}

	log.Println("Stopping provider")
	cancel()
}

func joinApateCluster(ctx context.Context, connectionInfo *service.ConnectionInfo, listenPort int) (cluster.KubeConfig, *normalization.NodeResources) {
	client := controlplane.GetClusterOperationClient(connectionInfo)
	defer func() {
		_ = client.Conn.Close()
	}()

	kubeconfig, res, err := client.JoinCluster(ctx, listenPort)

	// TODO: Better error handling
	if err != nil {
		log.Fatalf("Unable to join cluster: %v", err)
	}

	log.Printf("Joined apate cluster with resources: %v", res)

	return kubeconfig, res
}

func createNodeController(ctx context.Context, kubeConfig cluster.KubeConfig, res *normalization.NodeResources) (context.Context, *node.NodeController, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	restconfig, err := kubeConfig.GetConfig()
	if err != nil {
		log.Fatal("Could not parse config.")
	}

	client := kubernetes.NewForConfigOrDie(restconfig)
	n := cluster.NewNode("apatelet", "agent", "apatelet-"+res.UUID.String(), k8sVersion)
	nc, _ := node.NewNodeController(node.NaiveNodeProvider{},
		cluster.CreateKubernetesNode(ctx, *n, vkProvider.CreateProvider(res)),
		client.CoreV1().Nodes())

	return ctx, nc, cancel
}

func createGRPC(listenPort int) *service.GRPCServer {
	// Retrieving connection information
	listenAddress := container.RetrieveFromEnvironment(container.ApateletListenAddress, container.ApateletListenAddressDefault)

	// Connection settings
	connectionInfo := service.NewConnectionInfo(listenAddress, listenPort, false)

	// Create gRPC server
	server := service.NewGRPCServer(connectionInfo)

	// Add services
	vkService.RegisterScenarioService(server)

	return server
}
