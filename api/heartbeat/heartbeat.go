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

func (s *Server) Ping(ctx context.Context, in *HeartbeatMessage) (*HeartbeatMessage, error) {
	log.Println("Received heartbeat")
	return &HeartbeatMessage{Message: in.Message}, nil
}

func StartServer() {
	listener, server := service.CreateListenerAndServer(8080, true)
	s := Server{}

	RegisterHeartbeatServer(server, &s)

	if err := server.Serve(listener); err != nil {
		log.Fatalf("Unable to serve: %v", err)
	}
}

func GetClient() *Client {
	conn := service.CreateConnection("localhost", 8080, true)
	return &Client{
		Conn: conn,
		Client: NewHeartbeatClient(conn),
	}
}