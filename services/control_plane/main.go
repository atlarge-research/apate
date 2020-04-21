package main

import (
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/heartbeat"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
	"golang.org/x/net/context"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/control_plane/cluster"
	"log"
)

func main()  {
	// Create and delete cluster for now
	createAndDeleteCluster()

	log.Println("Starting Apate control plane")

	// Start GRPC server/client
	startGRPC()
}

func startGRPC() {
	// Connection settings
	connectionInfo := service.NewConnectionInfo("localhost", 8080, true)

	// Server
	server := service.NewGRPServer(connectionInfo)
	heartbeat.RegisterService(server)
	server.Serve()

	// Client
	c := heartbeat.GetClient(connectionInfo)
	defer c.Conn.Close()

	res, err := c.Client.Ping(context.Background(), &heartbeat.HeartbeatMessage{Message: "ping"})

	if err != nil {
		log.Fatalf("Could not complete call: %v", err)
	}

	log.Printf("Got back from server: %v\n", res)
}

func createAndDeleteCluster() {
	cb := cluster.Default()
	c, err := cb.WithName("Apate").Create()
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

	if err := c.Delete(); err != nil {
		log.Fatalf("An error occured: %s", err.Error())
	}
}