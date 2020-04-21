package main

import (
	"control_plane/cluster"
	"log"
)

func main() {
	clusterbuilder := cluster.New()
	c, err := clusterbuilder.WithName("Apate").Create()
	if err != nil {
		log.Fatalf("An error occured: %s", err.Error())
	}

	numberOfPods, err := c.GetNumberOfPods()
	if err != nil {
		c.Delete()
		log.Fatalf("An error occured: %s", err.Error())
	}

	log.Printf("There are %d pods in the cluster", numberOfPods)

	c.Delete()
}
