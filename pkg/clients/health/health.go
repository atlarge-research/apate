// Package health provides the client side gRPC API for the healthcheck service
package health

import (
	"context"
	"log"
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

func (c *Client) close() error {
	if c.Conn != nil {
		return c.Conn.Close()
	}

	return nil
}

const (
	sendInterval = 20 * time.Second
	recvTimeout  = 30 * time.Second
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

	go c.sendLoop(ctx, errCallback, stream)
	go c.recvLoop(ctx, errCallback, stream)
}

func (c *Client) sendLoop(ctx context.Context, errCallback func(error) bool, stream health.Health_HealthStreamClient) {
	// Send health status
	timeoutDelay := time.NewTimer(sendInterval)
	defer timeoutDelay.Stop()

	for {
		c.statusLock.RLock()
		err := stream.Send(&health.NodeStatus{
			NodeUuid: c.uuid,
			Status:   c.status,
		})
		c.statusLock.RUnlock()

		if err != nil {
			if c.cancelSend(errCallback, err, stream, "failed to send health status message over stream") {
				return
			}
		}

		timeoutDelay.Reset(sendInterval)
		select {
		case <-ctx.Done():
			if c.cancelSend(errCallback, err, stream, "unable to send heartbeat to server") {
				return
			}
		case <-timeoutDelay.C:
		}
	}
}

func (c *Client) cancelSend(errCallback func(error) bool, err error, stream health.Health_HealthStreamClient, msg string) bool {
	if errCallback(errors.Wrap(err, msg)) {
		if err := stream.CloseSend(); err != nil {
			log.Printf("%v\n", errors.Wrap(err, "error while closing health stream client (send)"))
		}
		return true
	}
	return false
}

func (c *Client) recvLoop(ctx context.Context, errCallback func(error) bool, stream health.Health_HealthStreamClient) {
	// Receive heartbeat from server
	timeoutDelay := time.NewTimer(recvTimeout)
	defer timeoutDelay.Stop()

	for {
		r := make(chan struct{}, 1)
		go c.recv(ctx, stream, errCallback, r)

		timeoutDelay.Reset(recvTimeout)
		select {
		case <-ctx.Done():
			// On context cancel stop
			if c.cancelRecv(ctx, errCallback, "unable to receive heartbeat from server") {
				return
			}
		case <-timeoutDelay.C:
			if c.cancelRecv(ctx, errCallback, "health stream died") {
				return
			}
		case <-r:
		}
	}
}

func (c *Client) cancelRecv(ctx context.Context, errCallback func(error) bool, msg string) bool {
	if errCallback(errors.Wrap(ctx.Err(), msg)) {
		if err := c.close(); err != nil {
			log.Printf("%v\n", errors.Wrap(err, "error while closing health stream client (recv)"))
		}
		return true
	}
	return false
}

func (c *Client) recv(ctx context.Context, stream health.Health_HealthStreamClient, errCallback func(error) bool, r chan struct{}) {
	_, err := stream.Recv()

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
}

// SetStatus sets the internal health status which is reported back to the control plane
func (c *Client) SetStatus(s health.Status) {
	c.statusLock.Lock()
	defer c.statusLock.Unlock()
	c.status = s
}
