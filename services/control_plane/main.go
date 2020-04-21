package main

import (
	"emulating-k8s/api/heartbeat"
	"golang.org/x/net/context"
	"log"
)

func main()  {
	log.Println("Starting Apate control plane")

	// Server
	heartbeat.StartServer()

	// Client
	c := heartbeat.GetClient()
	defer c.Conn.Close()

	res, err := c.Client.Ping(context.Background(), &heartbeat.HeartbeatMessage{Message: "ping"})

	if err != nil {
		log.Fatalf("Could not complete call: %v", err)
	}

	log.Printf("Got back from server: %v\n", res)
}
