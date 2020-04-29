// Package health provides the client side gRPC API for the healthcheck service
package health

import (
	"context"
	"log"
	"sync"
	"sync/atomic"
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

// StartStreamWithRetry calls StartStream but will retry n times to re-establish a connection
func (c *Client) StartStreamWithRetry(n int32) {
	ctx := context.Background()
	c.StartStream(ctx, func(err error) {
		if atomic.LoadInt32(&n) < 1 {
			log.Fatal(err)
			return
		}
		log.Println(err)
		atomic.AddInt32(&n, -1)
	})
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

			time.Sleep(sendInterval)
		}
	}()

	// Receive heartbeat from server
	go func() {
		for {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, recvTimeout)
			c := make(chan bool)
			go func() {
				select {
				case <-ctx.Done():
					// timeout reached
					errCallback(ctx.Err())
				case <-c:
					cancel()
				}
			}()
			_, err := stream.Recv()
			c <- true

			// Stream dead
			if err != nil {
				errCallback(err)
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
