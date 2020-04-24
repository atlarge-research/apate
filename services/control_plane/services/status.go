package services

import (
	"context"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/control_plane"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/control_plane/store"
	"github.com/golang/protobuf/ptypes/empty"
)

type statusService struct {
	store *store.Store
}

func RegisterStatusService(server *service.GRPCServer, store *store.Store) {
	control_plane.RegisterStatusServer(server.Server, &statusService{store: store})
}

func (s statusService) Status(ctx context.Context, e *empty.Empty) (*control_plane.ClusterStatus, error) {
	panic("implement me")
}
