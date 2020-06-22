package apatelet

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/atlarge-research/apate/api/apatelet"
	"github.com/atlarge-research/apate/internal/service"
)

//ApateClient is the client used for the Apatelet service on the apatelets
type ApateClient struct {
	Conn   *grpc.ClientConn
	Client apatelet.ApateletClient
}

// GetApateletClient returns a new apatelet handler
func GetApateletClient(info *service.ConnectionInfo) (*ApateClient, error) {
	conn, err := service.CreateClientConnection(info)
	if err != nil {
		return nil, errors.Wrap(err, "creating client connection info failed")
	}

	return &ApateClient{
		Conn:   conn,
		Client: apatelet.NewApateletClient(conn),
	}, nil
}
