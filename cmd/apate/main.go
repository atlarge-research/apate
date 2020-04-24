package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/control_plane"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/deserialize"
)

const (
	defaultControlPlaneAddress = "localhost:8083"
)

func main() {
	// TODO: Do arg parsing here with a proper library.

	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatalf("Please give a scenario filename as first argument")
	}

	log.Printf("Parsing yaml file")

	yaml, err := deserialize.YamlScenario{}.FromFile(args[0])
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Dialling server")

	var controlPlaneAddress string
	if len(args) == 1 {
		controlPlaneAddress = defaultControlPlaneAddress
	} else {
		controlPlaneAddress = args[1]
	}

	ctx, cancel, conn := createClient(controlPlaneAddress)
	defer func() {
		if err = conn.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	defer cancel()

	// Initial call: load the scenario
	scenarioClient := control_plane.NewScenarioClient(conn)
	_, err = scenarioClient.LoadScenario(ctx, yaml.GetScenario())
	if err != nil {
		log.Fatal(err)
	}

	statusClient := control_plane.NewStatusClient(conn)
	_, err = statusClient.Status(ctx, new(empty.Empty)) // TODO poll for status
	if err != nil {
		log.Fatal(err)
	}

	if _, err := scenarioClient.StartScenario(ctx, new(empty.Empty)); err != nil {
		log.Fatal(err)
	}
}

func createClient(controlPlaneAddress string) (context.Context, context.CancelFunc, *grpc.ClientConn) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	conn, err := grpc.DialContext(ctx, controlPlaneAddress, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	log.Printf("Registering client")

	return ctx, cancel, conn
}
