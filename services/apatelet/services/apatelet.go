package services

import (
	"context"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/apatelet"
	"github.com/atlarge-research/opendc-emulate-kubernetes/internal/service"
)

type apateletService struct {
	stopChannel chan<- struct{}
}

// RegisterApateletService registers the apateletService to the given GRPCServer
func RegisterApateletService(server *service.GRPCServer, stopChannel chan<- struct{}) {
	apatelet.RegisterApateletServer(server.Server, &apateletService{stopChannel: stopChannel})
}

// StopApatelet stops the apatelet
func (s *apateletService) StopApatelet(context.Context, *empty.Empty) (*empty.Empty, error) {
	log.Printf("received request to stop")

	go func() {
		time.Sleep(time.Second) // Wait a bit to properly answer the control plane

		s.stopChannel <- struct{}{}
	}()

	return new(empty.Empty), nil
}
