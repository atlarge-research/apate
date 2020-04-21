package main

import (
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster"
	"log"
)

func main() {
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
