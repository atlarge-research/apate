package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/sync/errgroup"

	apiApatelet "github.com/atlarge-research/opendc-emulate-kubernetes/api/apatelet"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/controlplane"
	"github.com/atlarge-research/opendc-emulate-kubernetes/internal/run"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/clients/apatelet"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubectl"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/normalization"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/scenario"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/store"
)

// The amount of seconds to wait with starting the scenario
// Should be large enough so that all apatelets have received the scenario and are ready
const amountOfSecondsToWait = 15

type scenarioService struct {
	store *store.Store
	info  *service.ConnectionInfo
}

// RegisterScenarioService registers a new scenarioService with the given gRPC server
func RegisterScenarioService(server *service.GRPCServer, store *store.Store, info *service.ConnectionInfo) {
	controlplane.RegisterScenarioServer(server.Server, &scenarioService{
		store: store,
		info:  info,
	})
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

	// Retrieve pull policy
	pullPolicy := env.RetrieveFromEnvironment(env.ControlPlaneDockerPolicy, env.ControlPlaneDockerPolicyDefault)
	fmt.Printf("Using pull policy %s to spawn apatelets\n", pullPolicy)

	// Create environment for apatelets
	environment := env.DefaultApateletEnvironment()
	environment.AddConnectionInfo(s.info.Address, s.info.Port)

	// Start the apatelets
	if err := run.StartApatelets(ctx, len(resources), environment); err != nil {
		log.Print(err)
		return nil, err
	}

	return new(empty.Empty), nil
}

func (s *scenarioService) StartScenario(ctx context.Context, config *controlplane.StartScenarioConfig) (*empty.Empty, error) {
	nodes, err := (*s.store).GetNodes()
	if err != nil {
		scenario.Failed(err)
		return nil, err
	}

	apateletScenario, err := (*s.store).GetApateletScenario()
	if err != nil {
		scenario.Failed(err)
		return nil, err
	}

	startTime := time.Now().Add(time.Second * amountOfSecondsToWait).UnixNano()
	convertToAbsoluteTimestamp(apateletScenario, startTime)

	err = startOnNodes(ctx, nodes, apateletScenario)
	if err != nil {
		scenario.Failed(err)
		return nil, err
	}

	cfg, err := (*s.store).GetKubeConfig()
	if err != nil {
		scenario.Failed(err)
		return nil, err
	}

	// TODO: This is probably very flaky
	err = kubectl.Create(config.ResourceConfig, cfg)
	if err != nil {
		scenario.Failed(err)
		return nil, err
	}

	return new(empty.Empty), nil
}

func startOnNodes(ctx context.Context, nodes []store.Node, apateletScenario *apiApatelet.ApateletScenario) error {
	errs, ctx := errgroup.WithContext(ctx)

	for i := range nodes {
		node := nodes[i]
		errs.Go(func() error {
			ft := filterTasksForNode(apateletScenario.Task, node.UUID.String())
			nodeSc := &apiApatelet.ApateletScenario{Task: ft}

			scenarioClient := apatelet.GetScenarioClient(&node.ConnectionInfo)
			_, err := scenarioClient.Client.StartScenario(ctx, nodeSc)

			if err != nil {
				return err
			}

			return scenarioClient.Conn.Close()
		})
	}

	return errs.Wait()
}

func convertToAbsoluteTimestamp(as *apiApatelet.ApateletScenario, unixNanoStartTime int64) {
	for _, t := range as.Task {
		newStartTimeNano := t.RelativeTimestamp + unixNanoStartTime
		t.AbsoluteTimestamp = newStartTimeNano
		t.RelativeTimestamp = 0 // Reset
	}
}

func filterTasksForNode(inputTasks []*apiApatelet.Task, uuid string) []*apiApatelet.Task {
	filteredTasks := make([]*apiApatelet.Task, 0)
	for _, t := range inputTasks {
		if t.NodeSet[uuid] {
			filteredTasks = append(filteredTasks, t)
		}
	}
	return filteredTasks
}
