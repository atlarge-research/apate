// Package clients contains all the GRPC clients that can be used by the control plane
package clients

import (
	"google.golang.org/grpc"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/kubelet"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
)

// ScenarioClient is the client for the ScenarioHandler containing the connection and gRPC client
type ScenarioClient struct {
	Conn   *grpc.ClientConn
	Client kubelet.ScenarioClient
}

// GetScenarioClient returns client for the ScenarioHandler
func GetScenarioClient(info *service.ConnectionInfo) *ScenarioClient {
	conn := service.CreateClientConnection(info)
	return &ScenarioClient{
		Conn:   conn,
		Client: kubelet.NewScenarioClient(conn),
	}
}
