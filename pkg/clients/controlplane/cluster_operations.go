// Package controlplane contains all the GRPC clients that can be used to interact with the control plane
package controlplane

import (
	"context"

	"github.com/pkg/errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster/kubeconfig"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/google/uuid"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/normalization"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/controlplane"

	"google.golang.org/grpc"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
)

// ClusterOperationClient is the client for the ClusterOperationService containing the connection and gRPC client
type ClusterOperationClient struct {
	Conn   *grpc.ClientConn
	Client controlplane.ClusterOperationsClient
}

// GetClusterOperationClient returns client for the JoinClusterService
func GetClusterOperationClient(info *service.ConnectionInfo) (*ClusterOperationClient, error) {
	conn, err := service.CreateClientConnection(info)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create client connection")
	}
	return &ClusterOperationClient{
		Conn:   conn,
		Client: controlplane.NewClusterOperationsClient(conn),
	}, nil
}

// JoinCluster joins the apate cluster, saves the received kube config and returns the node resources
func (c *ClusterOperationClient) JoinCluster(ctx context.Context, listenPort int) (*kubeconfig.KubeConfig, *normalization.NodeResources, error) {
	res, err := c.Client.JoinCluster(ctx, &controlplane.ApateletInformation{Port: int32(listenPort)})

	// Check for any grpc error
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to join cluster")
	}

	// Parse the uuid and check for errors
	id, err := uuid.Parse(res.NodeUuid)

	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to parse node uuid (%v)", res.NodeUuid)
	}

	cfg, err := kubeconfig.FromBytes(res.KubeConfig)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to load kubeconfig")
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
	return errors.Wrap(err, "failed to leave cluster")
}

// GetKubeConfig returns the kubeconfig file
func (c *ClusterOperationClient) GetKubeConfig(ctx context.Context) ([]byte, error) {
	cfg, err := c.Client.GetKubeConfig(ctx, new(empty.Empty))

	if err != nil {
		return nil, errors.Wrap(err, "failed to get kubeconfig")
	}

	return cfg.Config, nil
}
