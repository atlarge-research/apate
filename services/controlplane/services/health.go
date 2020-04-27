package services

import (
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/health"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/store"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"io"
	"log"
)

type healthService struct {
	store *store.Store
}

func (h healthService) HealthStream(server health.Health_HealthStreamServer) error {
	log.Println("Starting new health stream")

	ctx := server.Context()

	var id uuid.UUID

	for {
		// Exit if context is done
		select {
		case <-ctx.Done():
			_ = (*h.store).SetNodeStatus(id, health.Status_UNHEALTHY)
			return ctx.Err()
		default:
		}

		// receive data
		req, err := server.Recv()

		if err == io.EOF {
			log.Println("stopping a stream due to EOF")
			break
		}
		if err != nil {
			log.Println("Receive error")
			continue
		}

		id, err = uuid.Parse(req.NodeUUID)
		if err != nil {
			log.Println("stopping a stream due to invalid uuid")
			break
		}

		if err := server.Send(&empty.Empty{}); err != nil {
			log.Println("send error")
		}
	}

	// If the loop is broken -> node unhealthy
	if err := (*h.store).SetNodeStatus(id, health.Status_UNHEALTHY); err != nil {
		log.Println(err)
	}

	return nil
}

func RegisterHealthService(server *service.GRPCServer, store *store.Store) {
	health.RegisterHealthServer(server.Server, &healthService{store: store})
}
