package controlplane

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/controlplane"
	"github.com/atlarge-research/opendc-emulate-kubernetes/internal/service"
)

// StatusClient is the client for the StatusService containing the connection and gRPC client
type StatusClient struct {
	Conn   *grpc.ClientConn
	Client controlplane.StatusClient
}

// GetStatusClient returns client for the StatusService
func GetStatusClient(info *service.ConnectionInfo) (*StatusClient, error) {
	conn, err := service.CreateClientConnection(info)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create client connection")
	}

	return &StatusClient{
		Conn:   conn,
		Client: controlplane.NewStatusClient(conn),
	}, nil
}

// WaitForControlPlane waits for the control plane to be up and running
func (c *StatusClient) WaitForControlPlane(ctx context.Context, timeout time.Duration) error {
	ready := make(chan struct{})

	// Checks if the cp is up
	go func() {
		for {
			_, err := c.Client.Status(ctx, new(empty.Empty))
			if err == nil {
				ready <- struct{}{}
				return
			}

			select {
			case <-ctx.Done():
				// cancel
				return
			case <-time.After(time.Millisecond * 500):
			}
		}
	}()

	select {
	case <-ctx.Done():
		return errors.Wrap(ctx.Err(), "context canceled while waiting for control plane status")
	case <-ready:
		return nil
	case <-time.After(timeout):
		return errors.New("timeout reached but control plane still offline")
	}
}

// WaitForTrigger polls the server every second to retrieve the latest amount of healthy nodes and calls the
// given update function after every poll until the trigger channel is written to
func (c *StatusClient) WaitForTrigger(ctx context.Context, trigger <-chan struct{}, update func(int)) error {
	for {
		select {
		case <-trigger:
			return nil
		default:
		}

		res, err := c.Client.Status(ctx, new(empty.Empty))
		if err != nil {
			return errors.Wrap(err, "failed to get cluster status")
		}

		healthy := int(res.HealthyNodes)

		update(healthy)

		// Sleep for a second before trying again
		time.Sleep(time.Second)
	}
}
