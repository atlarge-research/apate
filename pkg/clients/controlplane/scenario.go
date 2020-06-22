package controlplane

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/atlarge-research/apate/api/controlplane"
	"github.com/atlarge-research/apate/internal/service"
)

// ScenarioClient is the client for the ScenarioService containing the connection and gRPC client
type ScenarioClient struct {
	Conn   *grpc.ClientConn
	Client controlplane.ScenarioClient
}

// GetScenarioClient returns client for the ScenarioService
func GetScenarioClient(info *service.ConnectionInfo) (*ScenarioClient, error) {
	conn, err := service.CreateClientConnection(info)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create client connection")
	}

	return &ScenarioClient{
		Conn:   conn,
		Client: controlplane.NewScenarioClient(conn),
	}, nil
}
