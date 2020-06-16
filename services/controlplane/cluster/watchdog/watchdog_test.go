package watchdog

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/health"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes/mock_kubernetes"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/store"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/store/mock_store"
)

func TestWatchdogNoApateError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mock_store.NewMockStore(ctrl)
	var st store.Store = ms

	mapi := mock_kubernetes.NewMockClusterAPI(ctrl)
	var api kubernetes.ClusterAPI = mapi

	unhealthyUUID := uuid.New()
	ms.EXPECT().GetNodes().Return([]store.Node{
		{
			Status: health.Status_HEALTHY,
		},
		{
			UUID:   unhealthyUUID,
			Status: health.Status_UNHEALTHY,
			Resources: &scenario.NodeResources{
				CPU: 1000,
			},
		},
	}, nil)
	ms.EXPECT().GetNodes().Return([]store.Node{
		{
			Status: health.Status_HEALTHY,
		},
	}, nil)
	ms.EXPECT().RemoveNode(unhealthyUUID).Return(nil)
	ms.EXPECT().AddResourcesToQueue([]scenario.NodeResources{
		{
			CPU: 1000,
		},
	}).Return(nil)
	mapi.EXPECT().RemoveNodeFromCluster("apatelet-" + unhealthyUUID.String()).Return(nil)

	StartWatchDog(ctx, 1*time.Second, &st, &api)

	time.Sleep(3 * time.Second)

	cancel()
}

func TestWatchdogApateError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := mock_store.NewMockStore(ctrl)
	var st store.Store = ms

	mapi := mock_kubernetes.NewMockClusterAPI(ctrl)
	var api kubernetes.ClusterAPI = mapi

	unhealthyUUID := uuid.New()
	ms.EXPECT().GetNodes().Return([]store.Node{
		{
			Status: health.Status_HEALTHY,
		},
		{
			UUID:   unhealthyUUID,
			Status: health.Status_UNHEALTHY,
			Resources: &scenario.NodeResources{
				CPU: 1000,
			},
		},
	}, nil).Times(2)
	ms.EXPECT().RemoveNode(unhealthyUUID).Return(errors.New("f")).Times(2)
	mapi.EXPECT().RemoveNodeFromCluster("apatelet-" + unhealthyUUID.String()).Return(nil).Times(2)

	StartWatchDog(ctx, 1*time.Second, &st, &api)

	time.Sleep(3 * time.Second)

	cancel()
}