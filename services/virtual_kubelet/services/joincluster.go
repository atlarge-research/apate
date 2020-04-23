package services

import (
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
