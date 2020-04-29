// Package services contains all the clients and servers for the services
package services

import (
	"context"
	"log"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/apatelet"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
)

// scenarioHandlerService will contain the implementation for the scenarioService
type scenarioHandlerService struct{}

// RegisterScenarioService registers the scenarioHandlerService to the given GRPCServer
func RegisterScenarioService(server *service.GRPCServer) {
	apatelet.RegisterScenarioServer(server.Server, &scenarioHandlerService{})
}

// TODO implement for real
// StartScenario starts a given scenario on the current Apatelet
func (s *scenarioHandlerService) StartScenario(_ context.Context, scenario *apatelet.ApateletScenario) (*empty.Empty, error) {
	log.Print("Received scenario with start time ", scenario.StartTime, " and ", len(scenario.Task), " tasks")
	return new(empty.Empty), nil
}
