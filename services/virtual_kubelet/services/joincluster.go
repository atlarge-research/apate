package services

import (
	"context"
	"io/ioutil"
	"os"
	"path"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"

	cluster_ops "github.com/atlarge-research/opendc-emulate-kubernetes/api/cluster_operations"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
)

// ClusterOperationClient is the client for the ClusterOperationService containing the connection and gRPC client
type ClusterOperationClient struct {
	Conn   *grpc.ClientConn
	Client cluster_ops.ClusterOperationsClient
}

// GetJoinClusterClient returns client for the JoinClusterService
func GetJoinClusterClient(info *service.ConnectionInfo) *ClusterOperationClient {
	conn := service.CreateClientConnection(info)
	return &ClusterOperationClient{
		Conn:   conn,
		Client: cluster_ops.NewClusterOperationsClient(conn),
	}
}

// JoinCluster joins the apate cluster, saves the received kube config and returns the kube context and uuid
func (c *ClusterOperationClient) JoinCluster(location string) (string, string, error) {
	res, err := c.Client.JoinCluster(context.Background(), &empty.Empty{})

	if err != nil {
		return "", "", err
	}

	if _, err = os.Stat(location); os.IsNotExist(err) {
		if err = os.MkdirAll(path.Dir(location), os.ModePerm); err != nil {
			return "", "", err
		}
	}

	// Save kube config
	err = ioutil.WriteFile(location, res.KubeConfig, 0644)

	if err != nil {
		return "", "", err
	}

	return res.KubeContext, res.NodeUUID, nil
}

// LeaveCluster signals to the apate control panel that this node is leaving the cluster
func (c *ClusterOperationClient) LeaveCluster(uuid string) error {
	_, err := c.Client.LeaveCluster(context.Background(), &cluster_ops.LeaveInformation{NodeUUID: uuid})
	return err
}
