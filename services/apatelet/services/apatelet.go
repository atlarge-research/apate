package services

import (
	"context"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/apatelet"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
)

type apateletService struct {
	stopChannel chan<- os.Signal
}

// RegisterApateletService registers the apateletService to the given GRPCServer
func RegisterApateletService(server *service.GRPCServer, stopChannel chan<- os.Signal) {
	apatelet.RegisterApateletServer(server.Server, &apateletService{stopChannel: stopChannel})
}

// StopApatelet stops the apatelet
func (s *apateletService) StopApatelet(context.Context, *empty.Empty) (*empty.Empty, error) {
	log.Printf("received request to stop")

	go func() {
		time.Sleep(time.Second) // Wait a bit to properly answer the control plane

		s.stopChannel <- syscall.SIGTERM
	}()

	return new(empty.Empty), nil
}
