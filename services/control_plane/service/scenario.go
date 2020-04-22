// Package service contains all the clients and servers for the services
package service

import (
	"google.golang.org/grpc"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario/private"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
)

// ScenarioClient is the client for the ScenarioHandler containing the connection and gRPC client
type ScenarioClient struct {
	Conn   *grpc.ClientConn
	Client private.ScenarioHandlerClient
}

// GetScenarioClient returns client for the ScenarioHandler
func GetScenarioClient(info *service.ConnectionInfo) *ScenarioClient {
	conn := service.CreateClientConnection(info)
	return &ScenarioClient{
		Conn:   conn,
		Client: private.NewScenarioHandlerClient(conn),
	}
}
