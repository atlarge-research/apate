package service

import (
	"context"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario/private"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
	"github.com/golang/protobuf/ptypes/empty"
	"log"
)

type ScenarioService struct{}

func RegisterScenarioService(server *service.GRPCServer) {
	private.RegisterScenarioHandlerServer(server.Server, &ScenarioService{})
}

func (s ScenarioService) StartScenario(ctx context.Context, scenario *private.Scenario) (*empty.Empty, error) {
	log.Print("Received scenario with start time ", scenario.StartTime, " and ", len(scenario.Task), " tasks")
	return new(empty.Empty), nil
}
