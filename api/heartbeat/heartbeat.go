package heartbeat

import (
	"emulating-k8s/pkg/service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
)

type Server struct {}

type Client struct {
	Conn *grpc.ClientConn
	Client HeartbeatClient
}

func RegisterService(server *service.GRPCServer) {
	RegisterHeartbeatServer(server.Server, &Server{})
}

func GetClient(info *service.ConnectionInfo) *Client {
	conn := service.CreateClientConnection(info)
	return &Client{
		Conn: conn,
		Client: NewHeartbeatClient(conn),
	}
}

func (s *Server) Ping(ctx context.Context, in *HeartbeatMessage) (*HeartbeatMessage, error) {
	log.Println("Received heartbeat")
	return &HeartbeatMessage{Message: in.Message}, nil
}