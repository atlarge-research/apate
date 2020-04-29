package controlplane

import (
	"google.golang.org/grpc"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/controlplane"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
)

// ScenarioClient is the client for the ScenarioService containing the connection and gRPC client
type ScenarioClient struct {
	Conn   *grpc.ClientConn
	Client controlplane.ScenarioClient
}

// GetScenarioClient returns client for the ScenarioService
func GetScenarioClient(info *service.ConnectionInfo) *ScenarioClient {
	conn := service.CreateClientConnection(info)
	return &ScenarioClient{
		Conn:   conn,
		Client: controlplane.NewScenarioClient(conn),
	}
}
