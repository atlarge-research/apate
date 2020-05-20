// Package service provides an a wrapper for connection information and a small wrapper around the grpc server
package service

import (
	"fmt"
	"log"
	"net"

	"sigs.k8s.io/kind/pkg/errors"

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
		return nil, errors.Wrap(err, "failed to create GRPC listener and server")
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
		log.Fatalf("Unable to serve: %+v", err)
	}
}

func createListenerAndServer(info *ConnectionInfo) (listener net.Listener, server *grpc.Server, err error) {
	listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", info.Address, info.Port))
	if err != nil {
		err = errors.Wrapf(err, "failed to listen on %v", fmt.Sprintf("%s:%d", info.Address, info.Port))
	}

	var options []grpc.ServerOption

	// Enable TLS if needed
	if info.TLS {
		tls, err1 := getServerTLS()
		if err1 != nil {
			err = errors.Wrap(err, "failed to start TLS server")
		}

		options = []grpc.ServerOption{tls}
	}

	server = grpc.NewServer(options...)

	return
}

//TODO: Real TLS instead of test data
func getServerTLS() (grpc.ServerOption, error) {
	creds, err := credentials.NewServerTLSFromFile(testdata.Path("server1.pem"), testdata.Path("server1.key"))

	if err != nil {
		return nil, errors.Wrap(err, "failed to create TLS server")
	}

	return grpc.Creds(creds), nil
}

// CreateClientConnection creates a connection to a remote services with the given connection information
func CreateClientConnection(info *ConnectionInfo) (conn *grpc.ClientConn) {
	var options = []grpc.DialOption{grpc.WithInsecure()}

	// Enable TLS if needed
	if info.TLS {
		tls, err := getClientTLS()
		if err != nil {
			log.Fatalf("Unable to create TLS client %+v", err)
		}

		options = []grpc.DialOption{tls}
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", info.Address, info.Port), options...)

	if err != nil {
		log.Fatalf("Unable to connect to %s:%d: %+v", info.Address, info.Port, err)
	}

	return
}

func getClientTLS() (grpc.DialOption, error) {
	creds, err := credentials.NewClientTLSFromFile(testdata.Path("ca.pem"), "x.test.youtube.com")

	if err != nil {
		return nil, errors.Wrap(err, "Failed to load TLS credentials")
	}

	return grpc.WithTransportCredentials(creds), nil
}
