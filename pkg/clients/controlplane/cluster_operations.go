// Package controlplane contains all the GRPC clients that can be used to interact with the control plane
package controlplane

import (
	"context"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster/kubeconfig"

	"github.com/google/uuid"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/normalization"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/controlplane"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
)

// ClusterOperationClient is the client for the ClusterOperationService containing the connection and gRPC client
type ClusterOperationClient struct {
	Conn   *grpc.ClientConn
	Client controlplane.ClusterOperationsClient
}

// GetClusterOperationClient returns client for the JoinClusterService
func GetClusterOperationClient(info *service.ConnectionInfo) *ClusterOperationClient {
	conn := service.CreateClientConnection(info)
	return &ClusterOperationClient{
		Conn:   conn,
		Client: controlplane.NewClusterOperationsClient(conn),
	}
}

// JoinCluster joins the apate cluster, saves the received kube config and returns the node resources
func (c *ClusterOperationClient) JoinCluster(ctx context.Context) (*kubeconfig.KubeConfig, *normalization.NodeResources, error) {
	res, err := c.Client.JoinCluster(ctx, &empty.Empty{})

	// Check for any grpc error
	if err != nil {
		return nil, nil, err
	}

	// Parse the uuid and check for errors
	id, err := uuid.Parse(res.NodeUuid)

	if err != nil {
		return nil, nil, err
	}

	cfg, err := kubeconfig.FromBytes(res.KubeConfig)
	if err != nil {
		return nil, nil, err
	}

	// Return final join information
	return cfg, &normalization.NodeResources{
		UUID:             id,
		Memory:           res.Hardware.Memory,
		CPU:              res.Hardware.Cpu,
		Storage:          res.Hardware.Storage,
		EphemeralStorage: res.Hardware.EphemeralStorage,
		MaxPods:          res.Hardware.MaxPods,
	}, nil
}

// LeaveCluster signals to the apate control panel that this node is leaving the cluster
func (c *ClusterOperationClient) LeaveCluster(ctx context.Context, uuid string) error {
	_, err := c.Client.LeaveCluster(ctx, &controlplane.LeaveInformation{NodeUuid: uuid})
	return err
}
