package services

import (
	"context"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/health"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/store"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"log"
	"time"
)

type healthService struct {
	store *store.Store
}

func (h healthService) HealthStream(server health.Health_HealthStreamServer) error {
	log.Println("Starting new health stream")

	ctx := server.Context()
	ctx, _ = context.WithDeadline(ctx, time.Now().Add(time.Minute*10)) // Disconnect after 10 minutes

	var id uuid.UUID

	cnt := 0

	for {
		if cnt > 2 {
			break
		}

		ctx, cancel := context.WithTimeout(ctx, time.Second*15)
		c := make(chan bool)
		// Exit if context is done
		go func() {
			select {
			case <-ctx.Done():
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

		id, err = uuid.Parse(req.NodeUUID)
		if err != nil {
			log.Println("stopping a stream due to invalid uuid")
			break
		}

		if err = (*h.store).SetNodeStatus(id, req.Status); err != nil {
			log.Println(err)
			continue
		}

		if err := server.Send(&empty.Empty{}); err != nil {
			log.Println("send error")
			cnt++
			continue
		}
	}

	// If the loop is broken -> node unhealthy
	if err := (*h.store).SetNodeStatus(id, health.Status_DISCONNECTED); err != nil {
		log.Println(err)
	}

	return nil
}

func RegisterHealthService(server *service.GRPCServer, store *store.Store) {
	health.RegisterHealthServer(server.Server, &healthService{store: store})
}
