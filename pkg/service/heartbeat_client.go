package service

import (
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/heartbeat"
	"google.golang.org/grpc"
)

type Client struct {
	Conn   *grpc.ClientConn
	Client heartbeat.HeartbeatClient
}

func GetClient(info *ConnectionInfo) *Client {
	conn := CreateClientConnection(info)
	return &Client{
		Conn:   conn,
		Client: heartbeat.NewHeartbeatClient(conn),
	}
}
