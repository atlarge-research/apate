package services

import (
	"github.com/atlarge-research/apate/api/health"
	"github.com/atlarge-research/apate/api/health/mock_health"
	"github.com/atlarge-research/apate/services/controlplane/store"
	"github.com/atlarge-research/apate/services/controlplane/store/mock_store"

	"github.com/pkg/errors"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"context"
	"testing"
	"time"
)

//go:generate sh -c "cd ../../../ && make mock_gen"

func TestHealthStream(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("Skipping health stream test")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl, ctx := gomock.WithContext(ctx, t)

	server := mock_health.NewMockHealth_HealthStreamServer(ctrl)

	ms := mock_store.NewMockStore(ctrl)

	msg := health.NodeStatus{
		NodeUuid: uuid.New().String(),
		Status:   health.Status_HEALTHY,
	}

	// Return our context when asked for
	server.EXPECT().Context().Return(ctx)

	// We return an error on send so both methods should be called thrice
	server.EXPECT().Recv().Return(&msg, nil).MinTimes(1)
	server.EXPECT().Send(gomock.Any()).Return(errors.New("some error")).Times(3)

	// Just accept all store calls
	ms.EXPECT().SetNodeStatus(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	var s store.Store = ms
	hs := healthService{&s}

	err := hs.HealthStream(server)
	assert.NoError(t, err)

	time.Sleep(time.Second)
}
