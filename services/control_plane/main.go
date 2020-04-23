package main

import (
	"context"
	"log"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario/private"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
	apatecluster "github.com/atlarge-research/opendc-emulate-kubernetes/services/control_plane/cluster"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/control_plane/services"
)

func main() {
	// Create and delete cluster for now
	c := createCluster()

	log.Println("Starting Apate control plane")

	// Create apate cluster state
	//apateCluster := apatecluster.NewApateCluster()

	// Start gRPC server/client
	//startGRPC(&apateCluster)

	//sc := &private.Scenario{
	//	Task:      nil,
	//	StartTime: 0,
	//}
	//scheduleOnNodes(sc, &c)

	if deleteErr := c.Delete(); deleteErr != nil {
		log.Printf("An error occurred: %s", deleteErr.Error())
	}
}

func startGRPC(apateCluster *apatecluster.ApateCluster) {
	// Connection settings
	connectionInfo := service.NewConnectionInfo("localhost", 8080, true)

	// Create gRPC server
	server := service.NewGRPCServer(connectionInfo)

	// Add services
	services.RegisterJoinClusterService(server, apateCluster)

	// Start serving request
	server.Serve()
}

func createCluster() cluster.KubernetesCluster {
	cb := cluster.Default()
	c, err := cb.WithName("Apate").ForceCreate()
	if err != nil {
		log.Fatalf("An error occurred: %s", err.Error())
	}

	numberOfPods, err := c.GetNumberOfPods()
	if err != nil {
		log.Fatalf("An error occurred: %s", err.Error())
	}

	log.Printf("There are %d pods in the cluster", numberOfPods)

	return c
}

func scheduleOnNodes(sc *private.Scenario, c *cluster.KubernetesCluster) {
	for _, port := range c.GetNodePorts() {
		// Connection settings
		connectionInfo := service.NewConnectionInfo("localhost", port, true)

		// Client
		scenarioClient := services.GetScenarioClient(connectionInfo)

		_, err := scenarioClient.Client.StartScenario(context.Background(), sc)

		if err != nil {
			log.Fatalf("Could not complete call: %v", err)
		}

		if err := scenarioClient.Conn.Close(); err != nil {
			log.Fatal("Failed to close connection")
		}
	}
}
