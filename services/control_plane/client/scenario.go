package client

import (
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario/private"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
	"google.golang.org/grpc"
)

type ScenarioClient struct {
	Conn   *grpc.ClientConn
	Client private.ScenarioHandlerClient
}

func GetScenarioClient(info *service.ConnectionInfo) *ScenarioClient {
	conn := service.CreateClientConnection(info)
	return &ScenarioClient{
		Conn:   conn,
		Client: private.NewScenarioHandlerClient(conn),
	}
}
