// Package apatelet contains all the GRPC clients that can be used to interact with the apatelet
package apatelet

import (
	"google.golang.org/grpc"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/apatelet"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
)

// ScenarioClient is the client for the ScenarioHandler containing the connection and gRPC client
type ScenarioClient struct {
	Conn   *grpc.ClientConn
	Client apatelet.ScenarioClient
}

// GetScenarioClient returns client for the ScenarioHandler
func GetScenarioClient(info *service.ConnectionInfo) *ScenarioClient {
	conn := service.CreateClientConnection(info)
	return &ScenarioClient{
		Conn:   conn,
		Client: apatelet.NewScenarioClient(conn),
	}
}
