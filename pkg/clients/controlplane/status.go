package controlplane

import (
	"context"
	"github.com/pkg/errors"
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

// WaitForControlPlane waits for the control plane to be up and running
func (c *StatusClient) WaitForControlPlane(ctx context.Context) error {
	for {
		cancellable, cancel := context.WithCancel(ctx)
		_, err := c.Client.Status(cancellable, new(empty.Empty))

		deadline, _ := ctx.Deadline()

		if deadline.Before(time.Now()) {
			cancel()
			return errors.New("waiting too long on control plane, giving up")
		}

		if err == nil {
			cancel()
			return nil
		}

		time.Sleep(time.Millisecond * 500)
		cancel()
	}
}

// WaitForHealthy polls the server every second to retrieve the latest amount of healthy nodes and calls the
// given update function after every poll
func (c *StatusClient) WaitForHealthy(ctx context.Context, expectedApatelets int, update func(int)) error {
	for {
		res, err := c.Client.Status(ctx, new(empty.Empty))

		if err != nil {
			return errors.Wrap(err, "failed to get client status")
		}

		healthy := int(res.HealthyNodes)

		update(healthy)

		if healthy >= expectedApatelets {
			return nil
		}

		// Sleep for a second before trying again
		time.Sleep(time.Second)
	}
}
