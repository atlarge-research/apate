package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/clients"

	"github.com/virtual-kubelet/virtual-kubelet/node"
	"k8s.io/client-go/kubernetes"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
	vkProvider "github.com/atlarge-research/opendc-emulate-kubernetes/services/virtual_kubelet/provider"
	vkService "github.com/atlarge-research/opendc-emulate-kubernetes/services/virtual_kubelet/services"
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
	log.Println("Starting Apate virtual kubelet")

	// TODO: Get these from envvars
	connectionInfo := service.NewConnectionInfo("localhost", 8083, false)
	location := os.TempDir() + "/apate/vk/config"

	// Join the apate cluster and start the kubelet
	log.Println("Joining apate cluster")
	kubeContext, uuid := joinApateCluster(location, connectionInfo)
	ctx, nc, cancel := getVirtualKubelet(location, kubeContext)

	log.Println("Joining kubernetes cluster")
	go func() {
		// TODO: Notify master / proper logging
		if err := nc.Run(ctx); err != nil {
			log.Fatalf("Unable to start apatelet: %v", err)
		}
	}()

	// Start gRPC server
	log.Println("Now accepting requests")
	server := createGRPC()

	// Handle signals
	signals := make(chan os.Signal, 1)
	stopped := make(chan bool, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		shutdown(server, cancel, connectionInfo, uuid)
		stopped <- true
	}()

	// Start serving request
	server.Serve()

	// Stop the server on signal
	<-stopped
	log.Println("Apate virtual kubelet stopped")
}

func shutdown(server *service.GRPCServer, cancel context.CancelFunc, connectionInfo *service.ConnectionInfo, uuid string) {
	log.Println("Stopping Apate virtual kubelet")

	log.Println("Stopping API")
	server.Server.Stop()

	log.Println("Leaving clusters (apate & k8s)")

	// TODO: Maybe leave k8s? Or will control plane do that?
	client := clients.GetClusterOperationClient(connectionInfo)
	defer func() {
		_ = client.Conn.Close()
	}()

	if err := client.LeaveCluster(uuid); err != nil {
		log.Printf("An error occurred while leaving the clusters (apate & k8s): %s", err.Error())
	}

	log.Println("Stopping provider")
	cancel()
}

func joinApateCluster(location string, connectionInfo *service.ConnectionInfo) (string, string) {
	client := clients.GetClusterOperationClient(connectionInfo)
	defer func() {
		_ = client.Conn.Close()
	}()

	ctx, uuid, err := client.JoinCluster(location)

	// TODO: Better error handling
	if err != nil {
		log.Fatalf("Unable to join cluster: %v", err)
	}

	log.Printf("Joined apate cluster with uuid %s", uuid)

	return ctx, uuid
}

func getVirtualKubelet(location string, kubeContext string) (context.Context, *node.NodeController, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	config, _ := cluster.GetConfigForContext(kubeContext, location)
	client := kubernetes.NewForConfigOrDie(config)
	n := cluster.NewNode("virtual-kubelet", "agent", "apatelet", k8sVersion)
	nc, _ := node.NewNodeController(node.NaiveNodeProvider{},
		cluster.CreateKubernetesNode(ctx, *n, vkProvider.CreateProvider()),
		client.CoreV1().Nodes())

	return ctx, nc, cancel
}

func createGRPC() *service.GRPCServer {
	// TODO: Get grpc settings from env
	// Connection settings
	connectionInfo := service.NewConnectionInfo("localhost", 8081, true)

	// Create gRPC server
	server := service.NewGRPCServer(connectionInfo)

	// Add services
	vkService.RegisterScenarioService(server)

	return server
}
