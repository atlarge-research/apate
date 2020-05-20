// Package health provides the client side gRPC API for the healthcheck service
package health

import (
	"context"
	"sync"
	"time"

	"google.golang.org/grpc"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/health"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
)

// Client holds all the information used to communicate with the server
type Client struct {
	Conn   *grpc.ClientConn
	Client health.HealthClient

	uuid       string
	status     health.Status
	statusLock sync.RWMutex
}

const (
	sendInterval = 1 * time.Second
	recvTimeout  = 5 * time.Second
)

// GetClient creates a new health client
func GetClient(info *service.ConnectionInfo, uuid string) *Client {
	conn := service.CreateClientConnection(info)

	return &Client{
		Conn:   conn,
		Client: health.NewHealthClient(conn),
		uuid:   uuid,
		status: health.Status_UNKNOWN,
	}
}

// StartStream starts the bidirectional health stream, errCallback is called upon any error
func (c *Client) StartStream(ctx context.Context, errCallback func(error)) {
	stream, err := c.Client.HealthStream(ctx)
	if err != nil {
		errCallback(err)
	}

	// Send health status
	go func() {
		for {
			c.statusLock.RLock()
			err = stream.Send(&health.NodeStatus{
				NodeUuid: c.uuid,
				Status:   c.status,
			})
			c.statusLock.RUnlock()

			if err != nil {
				errCallback(err)
			}

			select {
			case <-ctx.Done():
				return
			case <-time.After(sendInterval):
			}
		}
	}()

	// Receive heartbeat from server
	go func() {
		for {
			r := make(chan bool)

			go func() {
				_, err := stream.Recv()
				if err != nil {
					errCallback(err)
				}
				r <- true
			}()

			select {
			case <-ctx.Done():
				// On context cancel stop
				return
			case <-time.After(time.Second * recvTimeout):
				errCallback(ctx.Err())
			case <-r:
			}
		}
	}()
}

// SetStatus sets the internal health status which is reported back to the control plane
func (c *Client) SetStatus(s health.Status) {
	c.statusLock.Lock()
	defer c.statusLock.Unlock()
	c.status = s
}
