package services

import (
	"context"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/health"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/store"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/store/mock_store"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

//go:generate sh -c "cd ../../../ && make mockgen"

func TestStatusSimple(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl, ctx := gomock.WithContext(ctx, t)

	// Created mocked store
	ms := mock_store.NewMockStore(ctrl)

	// Create expectations
	ms.EXPECT().GetNodes().Return([]store.Node{
		{
			ConnectionInfo: service.ConnectionInfo{},
			UUID:           uuid.UUID{},
			Status:         health.Status_HEALTHY,
		},
	}, nil)

	var s store.Store = ms
	ss := statusService{&s}

	ret, err := ss.Status(nil, nil)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, ret.HealthyNodes)

}

func TestStatusExtensive(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl, ctx := gomock.WithContext(ctx, t)

	// Created mocked store
	ms := mock_store.NewMockStore(ctrl)

	// Create expectations
	ms.EXPECT().GetNodes().Return([]store.Node{
		{
			ConnectionInfo: service.ConnectionInfo{},
			UUID:           uuid.UUID{},
			Status:         health.Status_HEALTHY,
		},
		{
			ConnectionInfo: service.ConnectionInfo{},
			UUID:           uuid.UUID{},
			Status:         health.Status_UNKNOWN,
		},
		{
			ConnectionInfo: service.ConnectionInfo{},
			UUID:           uuid.UUID{},
			Status:         health.Status_HEALTHY,
		},
		{
			ConnectionInfo: service.ConnectionInfo{},
			UUID:           uuid.UUID{},
			Status:         health.Status_UNHEALTHY,
		},
		{
			ConnectionInfo: service.ConnectionInfo{},
			UUID:           uuid.UUID{},
			Status:         health.Status_HEALTHY,
		},
	}, nil)

	var s store.Store = ms
	ss := statusService{&s}

	ret, err := ss.Status(nil, nil)
	assert.NoError(t, err)
	assert.EqualValues(t, 3, ret.HealthyNodes)
}
