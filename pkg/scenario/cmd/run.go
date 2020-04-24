package cmd

import (
	"context"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario/public"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/deserialize"
	"google.golang.org/grpc"
	"log"
	"time"
)

const (
	address     = "localhost:8083"
)

func Run(args []string) error {
	// TODO: Do arg parsing here with a proper library.

	if len(args) < 1 {
		log.Fatalf("Please give a scenario filename as first argument")
	}

	log.Printf("Parsing yaml file")

	yaml, err := deserialize.YamlScenario{}.FromFile(args[0])
	if err != nil {
		return err
	}

	log.Printf("Dialing server")


	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(1 * time.Second))
	defer cancel()
	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure())


	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	log.Printf("Registering client")

	c := public.NewScenarioSenderClient(conn)

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SendScenario(ctx, yaml.GetScenario())

	if err != nil {
		return err
	}

	println(r)

	return nil
}
