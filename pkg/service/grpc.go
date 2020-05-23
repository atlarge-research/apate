// Package service provides an a wrapper for connection information and a small wrapper around the grpc server
package service

import (
	"fmt"
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
func (s *GRPCServer) Serve() error {
	return errors.Wrap(s.Server.Serve(s.listener), "error while serving GRPC")
}

func createListenerAndServer(info *ConnectionInfo) (listener net.Listener, server *grpc.Server, err error) {
	listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", info.Address, info.Port))
	if err != nil {
		err = errors.Wrapf(err, "failed to listen on %s:%d", info.Address, info.Port)
		return
	}

	var options []grpc.ServerOption

	// Enable TLS if needed
	if info.TLS {
		tls, tlsErr := getServerTLS()
		if tlsErr != nil {
			err = errors.Wrap(err, "failed to get server TLS")
			return
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
func CreateClientConnection(info *ConnectionInfo) (*grpc.ClientConn, error) {
	var options = []grpc.DialOption{grpc.WithInsecure()}

	// Enable TLS if needed
	if info.TLS {
		tls, err := getClientTLS()
		if err != nil {
			return nil, errors.Wrap(err, "Unable to create TLS client")
		}

		options = []grpc.DialOption{tls}
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", info.Address, info.Port), options...)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to connect to %s:%d", info.Address, info.Port)
	}

	return conn, nil
}

func getClientTLS() (grpc.DialOption, error) {
	creds, err := credentials.NewClientTLSFromFile(testdata.Path("ca.pem"), "x.test.youtube.com")

	if err != nil {
		return nil, errors.Wrap(err, "Failed to create TLS client")
	}

	return grpc.WithTransportCredentials(creds), nil
}
