package main

import (
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario/public"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
	apatecluster "github.com/atlarge-research/opendc-emulate-kubernetes/services/control_plane/cluster"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/control_plane/services"
)

func init() {
	// Enable line numbers in logging
	// Enables date time flags & file name + line
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	log.Println("Starting Apate control plane")

	// Create kubernetes cluster
	log.Println("Starting kubernetes control plane")
	kubernetesCluster := createCluster()

	// Create apate cluster state
	apateCluster := apatecluster.NewApateCluster()

	// Start gRPC server
	log.Println("Now accepting requests")
	server := createGRPC(&apateCluster)

	// Handle signals
	signals := make(chan os.Signal, 1)
	stopped := make(chan bool, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		shutdown(&apateCluster, &kubernetesCluster, server)
		stopped <- true
	}()

	// Register listeners
	public.RegisterScenarioSenderServer(server.Server, &services.SendScenarioServer{})

	// Start serving request
	server.Serve()

	// Stop the server on signal
	<-stopped
	log.Printf("Apate control plane stopped")
}

func shutdown(cluster *apatecluster.Store, kubernetesCluster *cluster.KubernetesCluster, server *service.GRPCServer) {
	log.Println("Stopping Apate control plane")

	log.Println("Stopping API")
	server.Server.Stop()

	// TODO: Actual cleanup for other nodes, for now just wipe state
	if err := (*cluster).ClearNodes(); err != nil {
		log.Printf("An error occurred while cleaning the apate cluster: %s", err.Error())
	}

	// TODO: Cleanup /tmp/ dir we used for kube config etc

	log.Println("Stopping kubernetes control plane")
	if err := kubernetesCluster.Delete(); err != nil {
		log.Printf("An error occurred while deleting the kubernetes cluster: %s", err.Error())
	}
}

func createGRPC(apateCluster *apatecluster.Store) *service.GRPCServer {
	// TODO: Get grpc settings from env
	// Connection settings
	connectionInfo := service.NewConnectionInfo("localhost", 8083, true)

	// Create gRPC server
	server := service.NewGRPCServer(connectionInfo)

	// Add services
	services.RegisterClusterOperationService(server, apateCluster)

	return server
}

func createCluster() cluster.KubernetesCluster {
	cb := cluster.Default()
	c, err := cb.WithName("Apate").ForceCreate()
	if err != nil {
		log.Fatalf("An error occurred: %s", err.Error())
	}

	numberOfPods, err := c.GetNumberOfPods("kube-system")
	if err != nil {
		log.Fatalf("An error occurred: %s", err.Error())
	}

	log.Printf("There are %d pods in the cluster", numberOfPods)

	return c
}

//func scheduleOnNodes(sc *private.Scenario, c *cluster.KubernetesCluster) {
//	for _, port := range c.GetNodePorts() {
//		// Connection settings
//		connectionInfo := service.NewConnectionInfo("localhost", port, true)
//
//		// Client
//		scenarioClient := services.GetScenarioClient(connectionInfo)
//
//		_, err := scenarioClient.Client.StartScenario(context.Background(), sc)
//
//		if err != nil {
//			log.Fatalf("Could not complete call: %v", err)
//		}
//
//		if err := scenarioClient.Conn.Close(); err != nil {
//			log.Fatal("Failed to close connection")
//		}
//	}
//}
