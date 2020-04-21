package service

import (
	"context"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/heartbeat"
	"log"
)

type Service struct{}

func RegisterService(server *GRPCServer) {
	heartbeat.RegisterHeartbeatServer(server.Server, &Service{})
}

func (s *Service) Ping(_ context.Context, in *heartbeat.HeartbeatMessage) (*heartbeat.HeartbeatMessage, error) {
	log.Println("Received heartbeat")
	return &heartbeat.HeartbeatMessage{Message: in.Message}, nil
}
