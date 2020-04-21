package main

import (
	"emulating-k8s/api/heartbeat"
	"emulating-k8s/pkg/service"
	"golang.org/x/net/context"
	"log"
)

func main()  {
	log.Println("Starting Apate control plane")

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
