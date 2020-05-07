// Package service provides an a wrapper for connection information and a small wrapper around the grpc server
package service

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"
)

// GRPCServer represents the gRPC server and listener
type GRPCServer struct {
	listener net.Listener
	Server   *grpc.Server
	Conn     *ConnectionInfo
}

// NewGRPCServer creates new gGRP server based on connection information
func NewGRPCServer(info *ConnectionInfo) (*GRPCServer, error) {
	lis, server, err := createListenerAndServer(info)
	if err != nil {
		return nil, err
	}

	return &GRPCServer{
		listener: lis,
		Server:   server,
		Conn:     info,
	}, nil
}

// Serve starts listening for incoming requests
func (s *GRPCServer) Serve() {
	if err := s.Server.Serve(s.listener); err != nil {
		log.Fatalf("Unable to serve: %v", err)
	}
}

func createListenerAndServer(info *ConnectionInfo) (listener net.Listener, server *grpc.Server, err error) {
	listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", info.Address, info.Port))
	var options []grpc.ServerOption

	// Enable TLS if needed
	if info.TLS {
		options = []grpc.ServerOption{getServerTLS()}
	}

	server = grpc.NewServer(options...)

	return
}

//TODO: Real TLS instead of test data
func getServerTLS() grpc.ServerOption {
	creds, err := credentials.NewServerTLSFromFile(testdata.Path("server1.pem"), testdata.Path("server1.key"))

	if err != nil {
		log.Fatalf("Failed to create TLS credentials: %v", err)
	}

	return grpc.Creds(creds)
}

// CreateClientConnection creates a connection to a remote services with the given connection information
func CreateClientConnection(info *ConnectionInfo) (conn *grpc.ClientConn) {
	var options = []grpc.DialOption{grpc.WithInsecure()}

	// Enable TLS if needed
	if info.TLS {
		options = []grpc.DialOption{getClientTLS()}
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", info.Address, info.Port), options...)

	if err != nil {
		log.Fatalf("Unable to connect to %s:%d: %v", info.Address, info.Port, err)
	}

	return
}

func getClientTLS() grpc.DialOption {
	creds, err := credentials.NewClientTLSFromFile(testdata.Path("ca.pem"), "x.test.youtube.com")

	if err != nil {
		log.Fatalf("Failed to load TLS credentials: %v", err)
	}

	return grpc.WithTransportCredentials(creds)
}
