// Package controlplane contains all the GRPC clients that can be used to interact with the control plane
package controlplane

import (
	"context"
	"io/ioutil"
	"os"
	"path"

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

// JoinCluster joins the apate cluster, saves the received kube config and returns the kube context and node resources
func (c *ClusterOperationClient) JoinCluster(location string) (string, *normalization.NodeResources, error) {
	res, err := c.Client.JoinCluster(context.Background(), &empty.Empty{})

	// Check for any grpc error
	if err != nil {
		return "", nil, err
	}

	// Ensure all dirs exist
	if _, err = os.Stat(location); os.IsNotExist(err) {
		if err = os.MkdirAll(path.Dir(location), os.ModePerm); err != nil {
			return "", nil, err
		}
	}

	// Save kube config
	err = ioutil.WriteFile(location, res.KubeConfig, 0600)

	// Check for write error
	if err != nil {
		return "", nil, err
	}

	// Parse the uuid and check for errors
	id, err := uuid.Parse(res.NodeUuid)

	if err != nil {
		return "", nil, err
	}

	// Return final join information
	return res.KubeContext, &normalization.NodeResources{
		UUID:    id,
		Memory:  res.Hardware.Memory,
		CPU:     res.Hardware.Cpu,
		MaxPods: res.Hardware.MaxPods,
	}, nil
}

// LeaveCluster signals to the apate control panel that this node is leaving the cluster
func (c *ClusterOperationClient) LeaveCluster(uuid string) error {
	_, err := c.Client.LeaveCluster(context.Background(), &controlplane.LeaveInformation{NodeUuid: uuid})
	return err
}
