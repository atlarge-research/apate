package main

import (
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
	apatecluster "github.com/atlarge-research/opendc-emulate-kubernetes/services/control_plane/cluster"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/control_plane/services"
	"log"
)

func main() {
	// Create and delete cluster for now
	createAndDeleteCluster()

	log.Println("Starting Apate control plane")

	// Create apate cluster state
	apateCluster := apatecluster.NewApateCluster()

	// Start gRPC server/client
	startGRPC(apateCluster)
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

//func startGRPC() {
//	// Connection settings
//	connectionInfo := service.NewConnectionInfo("localhost", 8080, true)
//
//	// HeartbeatService
//	server := service.NewGRPCServer(connectionInfo)
//	service.RegisterHeartbeatService(server)
//	server.Serve()
//
//	// HeartbeatClient
//	c := service.GetHeartbeatClient(connectionInfo)
//	defer func() {
//		if err := c.Conn.Close(); err != nil {
//			log.Fatalf("Failed to close connection")
//		}
//	}()
//
//	res, err := c.Client.Ping(context.Background(), &heartbeat.HeartbeatMessage{Message: "ping"})
//
//	if err != nil {
//		log.Fatalf("Could not complete call: %v", err)
//	}
//
//	log.Printf("Got back from server: %v\n", res)
//}

func createAndDeleteCluster() {
	cb := cluster.Default()
	c, err := cb.WithName("Apate").Create()
	if err != nil {
		log.Fatalf("An error occurred: %s", err.Error())
	}

	numberOfPods, err := c.GetNumberOfPods()
	if err != nil {
		if err := c.Delete(); err != nil {
			log.Printf("An error occurred: %s", err.Error())
		}
		log.Fatalf("An error occurred: %s", err.Error())
	}

	log.Printf("There are %d pods in the cluster", numberOfPods)

	if err := c.Delete(); err != nil {
		log.Fatalf("An error occurred: %s", err.Error())
	}
}
