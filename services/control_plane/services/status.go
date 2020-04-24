package services

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/control_plane"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/control_plane/store"
)

type statusService struct {
	store *store.Store
}

// RegisterStatusService registers a new statusService with the given gRPC server
func RegisterStatusService(server *service.GRPCServer, store *store.Store) {
	control_plane.RegisterStatusServer(server.Server, &statusService{store: store})
}

func (s statusService) Status(_ context.Context, _ *empty.Empty) (*control_plane.ClusterStatus, error) {
	return &control_plane.ClusterStatus{HealthyNodes: 5}, nil
}
