package controlplane

import (
	"google.golang.org/grpc"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/controlplane"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
)

// StatusClient is the client for the StatusService containing the connection and gRPC client
type StatusClient struct {
	Conn   *grpc.ClientConn
	Client controlplane.StatusClient
}

// GetStatusClient returns client for the StatusService
func GetStatusClient(info *service.ConnectionInfo) *StatusClient {
	conn := service.CreateClientConnection(info)
	return &StatusClient{
		Conn:   conn,
		Client: controlplane.NewStatusClient(conn),
	}
}
