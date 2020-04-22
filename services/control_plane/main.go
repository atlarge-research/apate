package main

import (
	"context"
	privateScenario "github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario/private"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
	cpService "github.com/atlarge-research/opendc-emulate-kubernetes/services/control_plane/service"
	"log"
)

func main() {
	// Create and delete cluster for now
	c := createCluster()

	log.Println("Starting Apate control plane")

	sc := &privateScenario.Scenario{
		Task:      nil,
		StartTime: 0,
	}
	scheduleOnPods(sc, &c)
}

func createCluster() cluster.KubernetesCluster {
	cb := cluster.Default()
	c, err := cb.WithName("Apate").ForceCreate()
	if err != nil {
		log.Fatalf("An error occured: %s", err.Error())
	}

	numberOfPods, err := c.GetNumberOfPods()
	if err != nil {
		if err := c.Delete(); err != nil {
			log.Printf("An error occured: %s", err.Error())
		}
		log.Fatalf("An error occured: %s", err.Error())
	}

	log.Printf("There are %d pods in the cluster", numberOfPods)

	return c
}

func scheduleOnPods(sc *privateScenario.Scenario, c *cluster.KubernetesCluster) {

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
			log.Fatalf("Failed to close connection")
		}
	}
}
