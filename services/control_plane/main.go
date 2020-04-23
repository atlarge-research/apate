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

	numberOfPods, err := c.GetNumberOfPods("kube-system")
	if err != nil {
		if err1 := c.Delete(); err1 != nil {
			err = err1
		}
		log.Fatalf("An error occurred: %s", err.Error())
	}

	log.Printf("There are %d pods in the cluster", numberOfPods)

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
