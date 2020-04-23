package main

import (
	"context"
	"log"

	privateScenario "github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario/private"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
	cpService "github.com/atlarge-research/opendc-emulate-kubernetes/services/control_plane/service"
)

func main() {
	// Create and delete cluster for now
	c := createCluster()

	log.Println("Starting Apate control plane")

	if err := cluster.SpawnNodes(); err != nil {
		log.Fatal(err)
	}

	sc := &privateScenario.Scenario{
		Task:      nil,
		StartTime: 0,
	}
	scheduleOnNodes(sc, &c)
}

func createCluster() cluster.KubernetesCluster {
	cb := cluster.Default()
	c, err := cb.WithName("Apate").ForceCreate()
	if err != nil {
		log.Fatalf("An error occurred: %s", err.Error())
	}

	return c
}

func scheduleOnNodes(sc *privateScenario.Scenario, c *cluster.KubernetesCluster) {
	for _, port := range c.GetNodePorts() {
		// Connection settings
		connectionInfo := service.NewConnectionInfo("localhost", port, true)

		// Client
		scenarioClient := cpService.GetScenarioClient(connectionInfo)

		_, err := scenarioClient.Client.StartScenario(context.Background(), sc)

		if err != nil {
			log.Fatalf("Could not complete call: %v", err)
		}

		if err := scenarioClient.Conn.Close(); err != nil {
			log.Fatal("Failed to close connection")
		}
	}
}
