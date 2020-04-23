package services

import (
	"context"
	"io/ioutil"
	"os"
	"path"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/join_cluster"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
)

// JoinClusterClient is the client for the JoinClusterService containing the connection and gRPC client
type JoinClusterClient struct {
	Conn   *grpc.ClientConn
	Client join_cluster.JoinClusterClient
}

// GetJoinClusterClient returns client for the JoinClusterService
func GetJoinClusterClient(info *service.ConnectionInfo) *JoinClusterClient {
	conn := service.CreateClientConnection(info)
	return &JoinClusterClient{
		Conn:   conn,
		Client: join_cluster.NewJoinClusterClient(conn),
	}
}

// JoinCluster joins the apate cluster, saves the received kube config and returns the kube context and uuid
func (c *JoinClusterClient) JoinCluster(location string) (string, string, error) {
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
