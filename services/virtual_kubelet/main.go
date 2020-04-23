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
	"strconv"
)

var (
	k8sVersion = "v1.15.2" // This should follow the version of k8s.io/kubernetes we are importing
)

func main() {
	connectionInfo := service.NewConnectionInfo("localhost", 8080, true)

	location := os.TempDir() + "/apate/vk/config"

	// Join the apate cluster and start the kubelet
	kubeCtx, _ := joinApateCluster(location, connectionInfo)
	startVirtualKubelet(location, kubeCtx)

	// Start serving gRPC request
	//startGRPC()
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

func startVirtualKubelet(location string, kubeCtx string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config, _ := cluster.GetConfigForContext(kubeCtx, location)
	client := kubernetes.NewForConfigOrDie(config)
	nc, _ := node.NewNodeController(node.NaiveNodeProvider{},
		cluster.CreateKubernetesNode(ctx,
			"virtual-kubelet",
			"agent",
			"apatelet",
			provider2.CreateProvider(),
			k8sVersion),
		client.CoreV1().Nodes())

	nc.Run(ctx)
}

func startGRPC() {
	// Connection settings
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal("Port not found in env")
	}

	connectionInfo := service.NewConnectionInfo("localhost", port, true)

	// Service
	server := service.NewGRPCServer(connectionInfo)
	vkService.RegisterScenarioService(server)
	server.Serve()
}
