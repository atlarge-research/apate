package services

import (
	"context"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/health"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/store"
)

type healthService struct {
	store *store.Store
}

func (h healthService) HealthStream(server health.Health_HealthStreamServer) error {
	log.Println("Starting new health stream")

	octx := server.Context()

	var id uuid.UUID

	cnt := 0

	for {
		if cnt > 2 {
			break
		}

		ctx, cancel := context.WithTimeout(octx, time.Second*15)
		c := make(chan bool)
		// Exit if context is done
		go func() {
			select {
			case <-ctx.Done():
				cnt++
				_ = (*h.store).SetNodeStatus(id, health.Status_DISCONNECTED)
			case <-c:
				cancel()
			}
		}()

		// receive data
		req, err := server.Recv()
		c <- true

		if err != nil {
			log.Println("Receive error")
			cnt++
			continue
		}

		id, err = uuid.Parse(req.NodeUuid)
		if err != nil {
			log.Println("stopping a stream due to invalid uuid")
			break
		}

		if err = (*h.store).SetNodeStatus(id, req.Status); err != nil {
			log.Println(err)
			continue
		}

		// Send heartbeat back
		if err := server.Send(&empty.Empty{}); err != nil {
			log.Println("send error")
			cnt++
			continue
		}
	}

	// If the loop is broken -> node unhealthy
	if err := (*h.store).SetNodeStatus(id, health.Status_DISCONNECTED); err != nil {
		log.Print(err)
	}

	// TODO: We should stop the scenario here (here as in here in time not place)
	log.Println("Node healthcheck disconnected...")

	return nil
}

// RegisterHealthService registers the HealthService on a GRPCServer
func RegisterHealthService(server *service.GRPCServer, store *store.Store) {
	health.RegisterHealthServer(server.Server, &healthService{store: store})
}
