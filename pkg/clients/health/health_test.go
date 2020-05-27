package health

import (
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/health"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/health/mock_health"

	"github.com/pkg/errors"

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/empty"

	"context"
	"testing"
	"time"
)

//go:generate sh -c "cd ../../../ && make mockgen"

func TestHealthClient(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	ctrl, ctx := gomock.WithContext(ctx, t)
	defer ctrl.Finish()

	tuuid := "test"
	tstatus := health.Status_HEALTHY

	mClient := mock_health.NewMockHealthClient(ctrl)
	mStream := mock_health.NewMockHealth_HealthStreamClient(ctrl)

	err := errors.New("stopping test")

	// Set expectations
	mClient.EXPECT().HealthStream(gomock.Any()).Return(mStream, nil)

	mStream.EXPECT().Send(gomock.Eq(&health.NodeStatus{
		NodeUuid: tuuid,
		Status:   tstatus,
	})).Return(nil).MinTimes(1)

	// And now just fail to clean up
	mStream.EXPECT().Recv().Return(&empty.Empty{}, err).MinTimes(1)

	// Create client
	c := Client{
		Client: mClient,
		status: tstatus,
		uuid:   tuuid,
	}

	c.StartStream(ctx, func(r error) {
		cancel()
	})
	time.Sleep(time.Second)
}
