package services

import (
	"context"
	"log"
	"time"

	"github.com/atlarge-research/apate/pkg/channel"

	"github.com/pkg/errors"

	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/sync/errgroup"

	apiApatelet "github.com/atlarge-research/apate/api/apatelet"
	"github.com/atlarge-research/apate/api/controlplane"
	"github.com/atlarge-research/apate/internal/service"
	"github.com/atlarge-research/apate/pkg/clients/apatelet"
	"github.com/atlarge-research/apate/services/controlplane/store"
)

// The amount of seconds to wait with starting the scenario
// Should be large enough so that all apatelets have received the scenario and are ready
const amountOfSecondsToWait = 5

type scenarioService struct {
	store          *store.Store
	info           *service.ConnectionInfo
	stopInformerCh *channel.StopChannel
}

// RegisterScenarioService registers a new scenarioService with the given gRPC server
func RegisterScenarioService(server *service.GRPCServer, store *store.Store, info *service.ConnectionInfo, stopInformerCh *channel.StopChannel) {
	controlplane.RegisterScenarioServer(server.Server, &scenarioService{
		store:          store,
		info:           info,
		stopInformerCh: stopInformerCh,
	})
}

func (s *scenarioService) StartScenario(ctx context.Context, startScenario *controlplane.StartScenario) (*empty.Empty, error) {
	nodes, err := (*s.store).GetNodes()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	apateletScenario := &apiApatelet.ApateletScenario{
		StartTime:       time.Now().Add(time.Second * amountOfSecondsToWait).UnixNano(),
		DisableWatchers: startScenario.DisableWatchers,
	}

	if err = (*s.store).SetApateletScenario(apateletScenario); err != nil {
		err = errors.Wrap(err, "failed to get Apatelet scenario")
		log.Println(err)
		return nil, err
	}

	log.Println("Starting scenario on nodes")
	if err = startOnNodes(ctx, nodes, apateletScenario); err != nil {
		err = errors.Wrap(err, "failed to get start scenario on nodes")
		log.Println(err)
		return nil, err
	}

	if startScenario.DisableWatchers {
		s.stopInformerCh.Close()
	}

	return new(empty.Empty), nil
}

func startOnNodes(ctx context.Context, nodes []store.Node, apateletScenario *apiApatelet.ApateletScenario) error {
	errs, ctx := errgroup.WithContext(ctx)

	for i := range nodes {
		node := nodes[i]
		errs.Go(func() error {
			scenarioClient, err := apatelet.GetScenarioClient(&node.ConnectionInfo)
			if err != nil {
				return errors.Wrap(err, "failed to get scenario client")
			}

			_, err = scenarioClient.Client.StartScenario(ctx, apateletScenario)

			if err != nil {
				return errors.Wrapf(err, "failed to start scenario on Apatelet with uuid %v", node.UUID.String())
			}

			return scenarioClient.Conn.Close()
		})
	}

	return errors.Wrap(errs.Wait(), "failed to start scenario on nodes")
}
