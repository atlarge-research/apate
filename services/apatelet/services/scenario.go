// Package services contains all the clients and servers for the services
package services

import (
	"context"
	"log"

	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/apatelet"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
)

// scenarioHandlerService will contain the implementation for the scenarioService
type scenarioHandlerService struct {
	store *store.Store
}

// RegisterScenarioService registers the scenarioHandlerService to the given GRPCServer
func RegisterScenarioService(server *service.GRPCServer, store *store.Store) {
	apatelet.RegisterScenarioServer(server.Server, &scenarioHandlerService{
		store: store,
	})
}

// StartScenario starts a given scenario on the current Apatelet
func (s *scenarioHandlerService) StartScenario(_ context.Context, scenario *apatelet.ApateletScenario) (*empty.Empty, error) {
	log.Print("Received scenario with ", len(scenario.Task), " tasks")

	// TODO remove when moving node to CRD
	var tasks []*store.Task
	for _, task := range scenario.Task {
		tasks = append(tasks, &store.Task{
			RelativeTimestamp: task.RelativeTimestamp,
			NodeTask:          task,
		})
	}
	(*s.store).EnqueueTasks(tasks)

	(*s.store).SetStartTime(scenario.StartTime)

	return new(empty.Empty), nil
}
