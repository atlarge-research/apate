package services

import (
	"context"
	"log"
	"time"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/clients/apatelet"

	apiApatelet "github.com/atlarge-research/opendc-emulate-kubernetes/api/apatelet"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/controlplane"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/normalization"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/store"

	"github.com/golang/protobuf/ptypes/empty"
)

// The amount of seconds to wait with starting the scenario
// Should be large enough so that all apatelets have received the scenario and are ready
const amountOfSecondsToWait = 15

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

	startTime := time.Now().Local().Add(time.Second * time.Duration(amountOfSecondsToWait)).UnixNano()
	convertToAbsoluteTimestamp(apateletScenario, startTime)

	for i := range nodes {
		go func(node *store.Node) {
			ft := prepareTasksForNode(apateletScenario.Task, node.UUID.String())
			nodeSc := &apiApatelet.ApateletScenario{Task: ft}

			scenarioClient := apatelet.GetScenarioClient(&node.ConnectionInfo)
			_, err := scenarioClient.Client.StartScenario(ctx, nodeSc)

			if err != nil {
				log.Fatalf("Could not complete call: %v", err)
			}

			if err := scenarioClient.Conn.Close(); err != nil {
				log.Fatal("Failed to close connection")
			}
		}(&nodes[i])
	}

	return new(empty.Empty), nil
}

func convertToAbsoluteTimestamp(as *apiApatelet.ApateletScenario, unixNanoStartTime int64) {
	for _, t := range as.Task {
		newStartTimeNano := t.RelativeTimestamp + unixNanoStartTime
		t.AbsoluteTimestamp = newStartTimeNano
	}
}

func prepareTasksForNode(inputTasks []*apiApatelet.Task, uuid string) []*apiApatelet.Task {
	filteredTasks := make([]*apiApatelet.Task, 0)
	for _, t := range inputTasks {
		if t.NodeSet[uuid] {
			filteredTasks = append(filteredTasks, t)
		}
	}
	return filteredTasks
}
