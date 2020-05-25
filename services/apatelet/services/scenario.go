// Package services contains all the clients and servers for the services
package services

import (
	"context"
	"log"

	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/scheduler"

	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/apatelet"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/atlarge-research/opendc-emulate-kubernetes/internal/service"
)

// scenarioHandlerService will contain the implementation for the scenarioService
type scenarioHandlerService struct {
	store *store.Store
	sch   *scheduler.Scheduler
}

// RegisterScenarioService registers the scenarioHandlerService to the given GRPCServer
func RegisterScenarioService(server *service.GRPCServer, store *store.Store, sch *scheduler.Scheduler) {
	apatelet.RegisterScenarioServer(server.Server, &scenarioHandlerService{
		store: store,
		sch:   sch,
	})
}

// StartScenario starts a given scenario on the current Apatelet
func (s *scenarioHandlerService) StartScenario(_ context.Context, scenario *apatelet.ApateletScenario) (*empty.Empty, error) {
	log.Printf("Scenario starting at %v\n", scenario.StartTime)

	s.sch.StartScheduler(scenario.StartTime)
	return new(empty.Empty), nil
}
