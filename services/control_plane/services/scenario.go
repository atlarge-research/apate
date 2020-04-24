package services

import (
	"context"
	"log"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/control_plane"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/clients"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/normalisation"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/control_plane/store"

	"github.com/golang/protobuf/ptypes/empty"
)

type scenarioService struct {
	store *store.Store
}

// RegisterScenarioService registers a new scenarioService with the given gRPC server
func RegisterScenarioService(server *service.GRPCServer, store *store.Store) {
	control_plane.RegisterScenarioServer(server.Server, &scenarioService{store: store})
}

func (s *scenarioService) LoadScenario(_ context.Context, scenario *control_plane.PublicScenario) (*empty.Empty, error) {
	log.Printf("Loading new scenario")

	normalisedScenario, resources, err := normalisation.NormaliseScenario(scenario)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	log.Printf("Adding %v to the queue", len(resources))
	if err := (*s.store).AddResourcesToQueue(resources); err != nil {
		log.Print(err)
		return nil, err
	}

	if err := (*s.store).AddKubeletScenario(normalisedScenario); err != nil {
		log.Print(err)
		return nil, err
	}

	if err := cluster.SpawnNodes(len(resources)); err != nil {
		log.Print(err)
		return nil, err
	}

	return new(empty.Empty), nil
}

func (s *scenarioService) StartScenario(_ context.Context, _ *empty.Empty) (*empty.Empty, error) {
	nodes, err := (*s.store).GetNodes()
	if err != nil {
		log.Print(err)
		return nil, err
	}

	kubeletScenario, err := (*s.store).GetKubeletScenario()
	if err != nil {
		log.Print(err)
		return nil, err
	}

	for _, node := range nodes {
		scenarioClient := clients.GetScenarioClient(&node.ConnectionInfo)
		_, err := scenarioClient.Client.StartScenario(context.Background(), kubeletScenario)

		if err != nil {
			log.Fatalf("Could not complete call: %v", err)
		}

		if err := scenarioClient.Conn.Close(); err != nil {
			log.Fatal("Failed to close connection")
		}
	}

	return new(empty.Empty), nil
}
