package apatelet

import (
	"google.golang.org/grpc"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/apatelet"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
)

//ApateClient is the client used for the Apatelet service on the apatelets
type ApateClient struct {
	Conn   *grpc.ClientConn
	Client apatelet.ApateletClient
}

// GetApateletClient returns a new apatelet handler
func GetApateletClient(info *service.ConnectionInfo) *ApateClient {
	conn := service.CreateClientConnection(info)
	return &ApateClient{
		Conn:   conn,
		Client: apatelet.NewApateletClient(conn),
	}
}
