package services

import (
	"context"
	"log"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/health"
	"github.com/atlarge-research/opendc-emulate-kubernetes/internal/service"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/store"
)

type healthService struct {
	store *store.Store
}

const (
	sendInterval     = 1 * time.Second
	recvTimeout      = 5 * time.Second
	maxNetworkErrors = 3
)

func (h healthService) HealthStream(server health.Health_HealthStreamServer) error {
	log.Println("Starting new health stream")

	// Outer/Original context
	ctx := server.Context()

	var id uuid.UUID

	// Keeps track of consecutive network errors
	var cnt int32 = 0

	// Sends a heartbeat to the client
	go h.sendHeartbeat(ctx, server, &cnt)

	for {
		if atomic.LoadInt32(&cnt) >= maxNetworkErrors {
			break
		}

		c := make(chan struct{})
		go func() {
			select {
			case <-ctx.Done():
				return
			case <-time.After(recvTimeout):
				atomic.AddInt32(&cnt, 1)
				_ = (*h.store).SetNodeStatus(id, health.Status_UNKNOWN)
			case <-c:
			}
		}()

		// receive data
		req, err := server.Recv()
		c <- struct{}{}

		if err != nil {
			log.Printf("Receive error: %v\n", err)
			atomic.AddInt32(&cnt, 1)
			continue
		}

		id, err = uuid.Parse(req.NodeUuid)
		if err != nil {
			log.Println("stopping a stream due to invalid uuid")
			break
		}

		if err = (*h.store).SetNodeStatus(id, req.Status); err != nil {
			log.Printf("%v\n", err)
		}
	}

	// If the loop is broken -> node status unknown
	if err := (*h.store).SetNodeStatus(id, health.Status_UNKNOWN); err != nil {
		log.Println(errors.Wrap(err, "failed to set node status"))
		return nil
	}

	log.Println("node healthcheck disconnected")
	return nil
}

func (h healthService) sendHeartbeat(ctx context.Context, server health.Health_HealthStreamServer, cnt *int32) {
	for {
		if atomic.LoadInt32(cnt) >= maxNetworkErrors {
			break
		}

		if err := server.Send(&empty.Empty{}); err != nil {
			log.Println("send error")
			atomic.AddInt32(cnt, 1)
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(sendInterval):
		}
	}
}

// RegisterHealthService registers the HealthService on a GRPCServer
func RegisterHealthService(server *service.GRPCServer, store *store.Store) {
	health.RegisterHealthServer(server.Server, &healthService{store: store})
}
