package main

import (
	"context"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/control_plane"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/deserialize"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"log"
	"os"
	"time"
)

const (
	address = "localhost:8083"
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

	log.Printf("Dialing server")

	ctx, cancel, conn := createClient(err)
	defer conn.Close()

	ctx, cancel = context.WithTimeout(ctx, time.Second)
	defer cancel()

	// Initial call: load the scenario
	scenarioClient := control_plane.NewScenarioClient(conn)
	r, err := scenarioClient.LoadScenario(ctx, yaml.GetScenario())
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

	log.Print(r)
}

func createClient(err error) (context.Context, context.CancelFunc, *grpc.ClientConn) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	log.Printf("Registering client")

	return ctx, cancel, conn
}
