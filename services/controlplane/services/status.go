package services

import (
	"context"

	"github.com/pkg/errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/health"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/controlplane"
	"github.com/atlarge-research/opendc-emulate-kubernetes/internal/service"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/store"
)

type statusService struct {
	store *store.Store
}

// RegisterStatusService registers a new statusService with the given gRPC server
func RegisterStatusService(server *service.GRPCServer, store *store.Store) {
	controlplane.RegisterStatusServer(server.Server, &statusService{store: store})
}

func (s *statusService) Status(_ context.Context, _ *empty.Empty) (*controlplane.ClusterStatus, error) {
	nodes, err := (*s.store).GetNodes()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get nodes list")
	}

	// OPTIMISE
	var cnt int32
	for _, node := range nodes {
		if node.Status == health.Status_HEALTHY {
			cnt++
		}
	}

	return &controlplane.ClusterStatus{HealthyNodes: cnt}, nil
}
