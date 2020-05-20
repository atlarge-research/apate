package services

import (
	"context"
	"github.com/pkg/errors"
	"log"
	"sync/atomic"
	"time"

	"github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/scenario"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/health"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
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
	octx := server.Context()

	var id uuid.UUID

	// Keeps track of consecutive network errors
	var cnt int32 = 0

	// Sends a heartbeat to the client
	go h.sendHeartbeat(server, &cnt)

	for {
		if atomic.LoadInt32(&cnt) >= maxNetworkErrors {
			break
		}

		ctx, cancel := context.WithTimeout(octx, recvTimeout)
		c := make(chan bool)
		// Exit if context is done
		go func() {
			select {
			case <-ctx.Done():
				atomic.AddInt32(&cnt, 1)
				_ = (*h.store).SetNodeStatus(id, health.Status_UNKNOWN)
			case <-c:
				cancel()
			}
		}()

		// receive data
		req, err := server.Recv()
		c <- true

		if err != nil {
			log.Println("Receive error")
			atomic.AddInt32(&cnt, 1)
			continue
		}

		id, err = uuid.Parse(req.NodeUuid)
		if err != nil {
			log.Println("stopping a stream due to invalid uuid")
			break
		}

		if err = (*h.store).SetNodeStatus(id, req.Status); err != nil {
			log.Printf("%+v\n", err)
		}

		// TODO: Improves
		// atomic.StoreInt32(&cnt, 0)
	}

	// If the loop is broken -> node status unknown
	if err := (*h.store).SetNodeStatus(id, health.Status_UNKNOWN); err != nil {
		scenario.Failed(errors.Wrap(err, "failed to set node status"))
		return nil
	}

	scenario.Failed(errors.New("node healthcheck disconnected"))
	return nil
}

func (h healthService) sendHeartbeat(server health.Health_HealthStreamServer, cnt *int32) {
	for {
		if atomic.LoadInt32(cnt) >= maxNetworkErrors {
			break
		}

		if err := server.Send(&empty.Empty{}); err != nil {
			log.Println("send error")
			atomic.AddInt32(cnt, 1)
		}
		time.Sleep(sendInterval)
	}
}

// RegisterHealthService registers the HealthService on a GRPCServer
func RegisterHealthService(server *service.GRPCServer, store *store.Store) {
	health.RegisterHealthServer(server.Server, &healthService{store: store})
}
