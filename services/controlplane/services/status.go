package services

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/controlplane"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/store"
)

type statusService struct {
	store *store.Store
}

// RegisterStatusService registers a new statusService with the given gRPC server
func RegisterStatusService(server *service.GRPCServer, store *store.Store) {
	controlplane.RegisterStatusServer(server.Server, &statusService{store: store})
}

func (s statusService) Status(_ context.Context, _ *empty.Empty) (*controlplane.ClusterStatus, error) {
	return &controlplane.ClusterStatus{HealthyNodes: 5}, nil
}
