package controlplane

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
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

// WaitForHealthy polls the server every second to retrieve the latest amount of healthy nodes and calls the
// given update function after every poll
func (c *StatusClient) WaitForHealthy(ctx context.Context, expectedApatelets int, update func(int)) error {
	for {
		res, err := c.Client.Status(ctx, new(empty.Empty))

		if err != nil {
			return err
		}

		healthy := int(res.HealthyNodes)

		update(healthy)

		if healthy >= expectedApatelets {
			return nil
		}

		// Sleep for a second before trying again
		time.Sleep(1000000000)
	}
}
