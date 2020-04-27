package health

import (
	"context"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/health"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
	"google.golang.org/grpc"
	"log"
	"sync"
)

type Client struct {
	Conn   *grpc.ClientConn
	Client health.HealthClient

	uuid       string
	status     health.Status
	statusLock sync.RWMutex
}

func GetClient(info *service.ConnectionInfo, UUID string) *Client {
	conn := service.CreateClientConnection(info)

	return &Client{
		Conn:   conn,
		Client: health.NewHealthClient(conn),
		uuid:   UUID,
		status: health.Status_UNKNOWN,
	}
}

func (c *Client) StartStream(errCallback func(error)) {
	ctx := context.Background()
	stream, err := c.Client.HealthStream(ctx)
	if err != nil {
		log.Print("can't connect")
		errCallback(err)
	}

	// Send health status
	go func() {
		for {
			c.statusLock.RLock()
			err = stream.Send(&health.NodeStatus{
				NodeUUID: c.uuid,
				Status:   c.status,
			})

			if err != nil {
				errCallback(err)
			}

			c.statusLock.RUnlock()
		}
	}()

	// Receive heartbeat from server
	go func() {
		for {
			_, err := stream.Recv()
			if err != nil {
				log.Print("error in recv")
				errCallback(err)
			}
		}
	}()

	return
}

func (c *Client) SetStatus(s health.Status) {
	c.statusLock.Lock()
	defer c.statusLock.Unlock()
	c.status = s
}
