package services

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
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

	// TODO: Skipped error wrapping in normalize scenario because it will be removed
	normalizedScenario, resources, err := normalization.NormalizeScenario(scenario)
	if err != nil {
		err = errors.Wrap(err, "failed to normalize scenario")
		log.Printf("%+v", err)
		return nil, err
	}

	log.Printf("Adding %v to the queue", len(resources))
	if err = (*s.store).AddResourcesToQueue(resources); err != nil {
		err = errors.Wrap(err, "failed to add resources to queue")
		log.Printf("%+v", err)

		return nil, err
	}

	if err = (*s.store).SetApateletScenario(normalizedScenario); err != nil {
		err = errors.Wrap(err, "failed to set scenario on Apatelet")
		log.Printf("%+v", err)
		return nil, err
	}

	// Retrieve pull policy
	pullPolicy := env.RetrieveFromEnvironment(env.ControlPlaneDockerPolicy, env.ControlPlaneDockerPolicyDefault)
	fmt.Printf("Using pull policy %s to spawn apatelets\n", pullPolicy)

	// Create environment for apatelets
	environment, err := env.DefaultApateletEnvironment()
	if err != nil {
		err = errors.Wrap(err, "failed to create Apatelet environment")
		log.Printf("%+v", err)
		return nil, err
	}

	environment.AddConnectionInfo(s.info.Address, s.info.Port)

	// Start the apatelets
	if err = run.StartApatelets(ctx, len(resources), environment); err != nil {
		err = errors.Wrap(err, "failed to start Apatelets")
		log.Printf("%+v", err)
		return nil, err
	}

	return new(empty.Empty), nil
}

func (s *scenarioService) StartScenario(ctx context.Context, config *controlplane.StartScenarioConfig) (*empty.Empty, error) {
	nodes, err := (*s.store).GetNodes()
	if err != nil {
		err = errors.Wrap(err, "failed to get nodes")

		scenario.Failed(err)
		log.Println(err)
		return nil, err
	}

	apateletScenario, err := (*s.store).GetApateletScenario()
	if err != nil {
		err = errors.Wrap(err, "failed to get Apatelet scenario")

		scenario.Failed(err)
		log.Println(err)
		return nil, err
	}

	startTime := time.Now().Add(time.Second * amountOfSecondsToWait).UnixNano()
	apateletScenario.StartTime = startTime

	err = startOnNodes(ctx, nodes, apateletScenario)
	if err != nil {
		err = errors.Wrap(err, "failed to get start scenario on nodes")

		scenario.Failed(err)
		log.Println(err)
		return nil, err
	}

	cfg, err := (*s.store).GetKubeConfig()
	if err != nil {
		err = errors.Wrap(err, "failed to get get Kubeconfig")

		scenario.Failed(err)
		log.Println(err)
		return nil, err
	}

	// TODO: This is probably very flaky
	err = kubectl.Create(config.ResourceConfig, &cfg)
	if err != nil {
		err = errors.Wrap(err, "failed to create resource config")

		scenario.Failed(err)
		log.Println(err)
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
				return errors.Wrap(err, "failed to start scenario on client")
			}

			return scenarioClient.Conn.Close()
		})
	}

	return errors.Wrap(errs.Wait(), "failed to start scenario on nodes")
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
