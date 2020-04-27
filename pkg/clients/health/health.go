package health

import (
	"context"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/health"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
	"google.golang.org/grpc"
	"io"
	"log"
	"sync"
	"time"
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
			c.statusLock.RUnlock()

			if err != nil {
				errCallback(err)
			}
		}
	}()

	// Receive heartbeat from server
	go func() {
		cnt := 0
		for {
			ctx, cancel := context.WithTimeout(ctx, time.Second*15)
			c := make(chan bool)
			go func() {
				select {
				case <-ctx.Done():
					// oof
					errCallback(ctx.Err())
				case <-c:
					// anti-oof
					cancel()
				}
			}()
			_, err := stream.Recv()
			c <- true

			// Stream dead
			if err == io.EOF {
				errCallback(err)
			}

			// Other error?
			if err != nil {
				cnt++
				log.Print("error in recv")

				if cnt > 2 {
					errCallback(err)
				}
				continue
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
