package main

import (
	"context"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
	provider2 "github.com/atlarge-research/opendc-emulate-kubernetes/services/virtual_kubelet/provider"
	vkService "github.com/atlarge-research/opendc-emulate-kubernetes/services/virtual_kubelet/services"
	"github.com/virtual-kubelet/virtual-kubelet/node"
	"k8s.io/client-go/kubernetes"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	k8sVersion = "v1.15.2" // This should follow the version of k8s.io/kubernetes we are importing
)

func main() {
	log.Println("Starting Apate virtual kubelet")

	// TODO: Get these from envvars
	connectionInfo := service.NewConnectionInfo("localhost", 8080, true)
	location := os.TempDir() + "/apate/vk/config"

	// Join the apate cluster and start the kubelet
	log.Println("Joining apate cluster")
	kubeContext, _ := joinApateCluster(location, connectionInfo)
	ctx, nc, cancel := getVirtualKubelet(location, kubeContext)

	log.Println("Joining kubernetes cluster")
	go nc.Run(ctx)

	// Start gRPC server
	log.Println("Now accepting requests")
	server := createGRPC()

	// Handle signals
	signals := make(chan os.Signal, 1)
	stopped := make(chan bool, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		shutdown(server, cancel)
		stopped <- true
	}()

	// Start serving request
	server.Serve()

	// Stop the server on signal
	<-stopped
	log.Println("Apate virtual kubelet stopped")
}

func shutdown(server *service.GRPCServer, cancel context.CancelFunc) {
	log.Println("Stopping Apate virtual kubelet")

	log.Println("Stopping API")
	server.Server.Stop()

	log.Println("Stopping provider")
	cancel()

	// TODO: Send message to control plane which deletes us from cluster (both apate and k8s)
}

func joinApateCluster(location string, connectionInfo *service.ConnectionInfo) (string, string) {
	c := vkService.GetJoinClusterClient(connectionInfo)
	defer func() {
		_ = c.Conn.Close()
	}()

	ctx, uuid, err := c.JoinCluster(location)

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
	nc, _ := node.NewNodeController(node.NaiveNodeProvider{},
		cluster.CreateKubernetesNode(ctx,
			"virtual-kubelet",
			"agent",
			"apatelet",
			provider2.CreateProvider(),
			k8sVersion),
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
