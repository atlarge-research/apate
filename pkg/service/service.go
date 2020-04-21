package service

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"
	"log"
	"net"
)

func CreateListenerAndServer(port int, tls bool) (listener net.Listener, server *grpc.Server) {
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	var options []grpc.ServerOption

	// Enable TLS if needed
	if tls {
		options = []grpc.ServerOption{GetServerTLS()}
	}

	server = grpc.NewServer(options...)

	if err != nil {
		log.Fatalf("Failed to listen on port 8080: %v", err)
	}

	return
}

func CreateConnection(address string, port int, tls bool) (conn *grpc.ClientConn) {
	var options = []grpc.DialOption{grpc.WithInsecure()}

	// Enable TLS if needed
	if tls {
		options = []grpc.DialOption{GetClientTLS()}
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", address, port), options...)

	if err != nil {
		log.Fatalf("Unable to connect to localhost:%d: %v", port, err)
	}

	return
}

//TODO: Real TLS instead of test data
func GetServerTLS() grpc.ServerOption {
	creds, err := credentials.NewServerTLSFromFile(testdata.Path("server1.pem"), testdata.Path("server1.key"))

	if err != nil {
		log.Fatalf("Failed to create TLS credentials: %v", err)
	}

	return grpc.Creds(creds)
}

func GetClientTLS() grpc.DialOption {
	creds, err := credentials.NewClientTLSFromFile(testdata.Path("ca.pem"), "x.test.youtube.com")

	if err != nil {
		log.Fatalf("Failed to load TLS credentials: %v", err)
	}

	return grpc.WithTransportCredentials(creds)
}