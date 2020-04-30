package services

import (
	"context"
	"log"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/clients/apatelet"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/controlplane"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/normalization"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/store"

	"github.com/golang/protobuf/ptypes/empty"
)

type scenarioService struct {
	store *store.Store
}

// RegisterScenarioService registers a new scenarioService with the given gRPC server
func RegisterScenarioService(server *service.GRPCServer, store *store.Store) {
	controlplane.RegisterScenarioServer(server.Server, &scenarioService{store: store})
}

func (s *scenarioService) LoadScenario(ctx context.Context, scenario *controlplane.PublicScenario) (*empty.Empty, error) {
	log.Printf("Loading new scenario")

	normalizedScenario, resources, err := normalization.NormalizeScenario(scenario)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	log.Printf("Adding %v to the queue", len(resources))
	if err := (*s.store).AddResourcesToQueue(resources); err != nil {
		log.Print(err)
		return nil, err
	}

	if err := (*s.store).SetApateletScenario(normalizedScenario); err != nil {
		log.Print(err)
		return nil, err
	}

	if err := cluster.SpawnNodes(ctx, len(resources)); err != nil {
		log.Print(err)
		return nil, err
	}

	return new(empty.Empty), nil
}

func (s *scenarioService) StartScenario(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	nodes, err := (*s.store).GetNodes()
	if err != nil {
		log.Print(err)
		return nil, err
	}

	apateletScenario, err := (*s.store).GetApateletScenario()
	if err != nil {
		log.Print(err)
		return nil, err
	}

	// TODO set the start time on the apatelet scenario (apateletScenario.startTime)
	// TODO make the task times absolute

	// TODO make async
	for _, node := range nodes {
		scenarioClient := apatelet.GetScenarioClient(&node.ConnectionInfo)
		_, err := scenarioClient.Client.StartScenario(ctx, apateletScenario)

		if err != nil {
			log.Fatalf("Could not complete call: %v", err)
		}

		if err := scenarioClient.Conn.Close(); err != nil {
			log.Fatal("Failed to close connection")
		}
	}

	return new(empty.Empty), nil
}
