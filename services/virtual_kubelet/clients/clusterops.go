package clients

import (
	"context"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/control_plane"
	"io/ioutil"
	"os"
	"path"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
)

// ClusterOperationClient is the client for the ClusterOperationService containing the connection and gRPC client
type ClusterOperationClient struct {
	Conn   *grpc.ClientConn
	Client control_plane.ClusterOperationsClient
}

// GetClusterOperationClient returns client for the JoinClusterService
func GetClusterOperationClient(info *service.ConnectionInfo) *ClusterOperationClient {
	conn := service.CreateClientConnection(info)
	return &ClusterOperationClient{
		Conn:   conn,
		Client: control_plane.NewClusterOperationsClient(conn),
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
	_, err := c.Client.LeaveCluster(context.Background(), &control_plane.LeaveInformation{NodeUUID: uuid})
	return err
}
