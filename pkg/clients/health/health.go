// Package health provides the client side gRPC API for the healthcheck service
package health

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/pkg/errors"

	"google.golang.org/grpc"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/health"
	"github.com/atlarge-research/opendc-emulate-kubernetes/internal/service"
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
func GetClient(info *service.ConnectionInfo, uuid string) (*Client, error) {
	conn, err := service.CreateClientConnection(info)

	if err != nil {
		return nil, errors.Wrap(err, "failed to create GRPC health client")
	}

	return &Client{
		Conn:   conn,
		Client: health.NewHealthClient(conn),
		uuid:   uuid,
		status: health.Status_UNKNOWN,
	}, nil
}

// StartStream starts the bidirectional health stream, errCallback is called upon any error
func (c *Client) StartStream(ctx context.Context, errCallback func(error) bool) {
	stream, err := c.Client.HealthStream(ctx)
	if err != nil {
		errCallback(errors.Wrap(err, "failed to set up health stream"))
	}

	// Send health status
	go func() {
		timeoutDuration := sendInterval
		timeoutDelay := time.NewTimer(timeoutDuration)
		defer timeoutDelay.Stop()

		for {
			c.statusLock.RLock()
			err = stream.Send(&health.NodeStatus{
				NodeUuid: c.uuid,
				Status:   c.status,
			})
			c.statusLock.RUnlock()

			if err != nil {
				if errCallback(errors.Wrap(err, "failed to send health status message over stream")) {
					_ = stream.CloseSend()
					return
				}
			}

			timeoutDelay.Reset(timeoutDuration)
			select {
			case <-ctx.Done():
				if errCallback(errors.Wrap(ctx.Err(), "context cancelled")) {
					_ = stream.CloseSend()
					return
				}
			case <-timeoutDelay.C:
			}
		}
	}()

	// Receive heartbeat from server
	go func() {
		timeoutDuration := time.Second * recvTimeout
		timeoutDelay := time.NewTimer(timeoutDuration)
		defer timeoutDelay.Stop()

		for {
			r := make(chan struct{}, 1)

			go func() {
				_, err := stream.Recv()

				if err == io.EOF {
					return
				}

				if err != nil {
					if errCallback(errors.Wrap(err, "health stream timed out")) {
						return
					}
				}

				select {
				case <-ctx.Done():
					if errCallback(errors.Wrap(ctx.Err(), "unable to receive heartbeat from server")) {
						return
					}
				case r <- struct{}{}:
					//
				}
			}()

			timeoutDelay.Reset(timeoutDuration)
			select {
			case <-ctx.Done():
				// On context cancel stop
				if errCallback(errors.Wrap(ctx.Err(), "unable to receive heartbeat from server")) {
					_ = c.Conn.Close()
					return
				}
			case <-timeoutDelay.C:
				if errCallback(errors.Errorf("health stream died")) {
					_ = c.Conn.Close()
					return
				}
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
