// Package apatelet contains all the GRPC clients that can be used to interact with the apatelet
package apatelet

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/atlarge-research/apate/api/apatelet"

	"github.com/atlarge-research/apate/internal/service"
)

// ScenarioClient is the client for the ScenarioHandler containing the connection and gRPC client
type ScenarioClient struct {
	Conn   *grpc.ClientConn
	Client apatelet.ScenarioClient
}

// GetScenarioClient returns client for the ScenarioHandler
func GetScenarioClient(info *service.ConnectionInfo) (*ScenarioClient, error) {
	conn, err := service.CreateClientConnection(info)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create scenario client connection")
	}

	return &ScenarioClient{
		Conn:   conn,
		Client: apatelet.NewScenarioClient(conn),
	}, nil
}
