// Package service provides a GRPC service. TODO: Tim?
package service

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"
)

// GRPCServer is a GRPC server. TODO: Tim
type GRPCServer struct {
	listener net.Listener
	Server   *grpc.Server
}

// ConnectionInfo gives connection info TODO: Tim
type ConnectionInfo struct {
	address string
	port    int
	tls     bool
}

// NewGRPCServer creates a new GRPC server. TODO: Tim (maybe put the server in it's own file and just call this New so you get server.New())
func NewGRPCServer(info *ConnectionInfo) *GRPCServer {
	lis, server := createListenerAndServer(info)
	return &GRPCServer{
		listener: lis,
		Server:   server,
	}
}

// NewConnectionInfo gets information about new connections. TODO: Tim
func NewConnectionInfo(address string, port int, tls bool) *ConnectionInfo {
	return &ConnectionInfo{
		address: address,
		port:    port,
		tls:     tls,
	}
}

// Serve starts serving the GRPC server TODO: Tim
func (s *GRPCServer) Serve() {
	if err := s.Server.Serve(s.listener); err != nil {
		log.Fatalf("Unable to serve: %v", err)
	}
}

func createListenerAndServer(info *ConnectionInfo) (listener net.Listener, server *grpc.Server) {
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", info.port))
	var options []grpc.ServerOption

	// Enable TLS if needed
	if info.tls {
		options = []grpc.ServerOption{getServerTLS()}
	}

	server = grpc.NewServer(options...)

	if err != nil {
		log.Fatalf("Failed to listen on port %d: %v", info.port, err)
	}

	return
}

// CreateClientConnection creates a new client connection. TODO: Tim
func CreateClientConnection(info *ConnectionInfo) (conn *grpc.ClientConn) {
	var options = []grpc.DialOption{grpc.WithInsecure()}

	// Enable TLS if needed
	if info.tls {
		options = []grpc.DialOption{getClientTLS()}
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", info.address, info.port), options...)

	if err != nil {
		log.Fatalf("Unable to connect to %s:%d: %v", info.address, info.port, err)
	}

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

func getClientTLS() grpc.DialOption {
	creds, err := credentials.NewClientTLSFromFile(testdata.Path("ca.pem"), "x.test.youtube.com")

	if err != nil {
		log.Fatalf("Failed to load TLS credentials: %v", err)
	}

	return grpc.WithTransportCredentials(creds)
}
