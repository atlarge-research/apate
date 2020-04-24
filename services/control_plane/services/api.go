package services

import (
	"context"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario/public"
	"log"
)

type SendScenarioServer struct {
	public.UnimplementedScenarioSenderServer
}

func (s *SendScenarioServer) SendScenario(ctx context.Context, scenario *public.Scenario) (*public.SendScenarioResponse, error) {
	log.Printf("Received: %v", scenario.GetNodes())

	return nil, nil
}
