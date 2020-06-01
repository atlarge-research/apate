// Package service provides an a wrapper for connection information and a small wrapper around the grpc server
package service

import (
	"fmt"
	"net"

	"sigs.k8s.io/kind/pkg/errors"

	"google.golang.org/grpc"
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
	return errors.Wrap(s.Server.Serve(s.listener), "error while serving gRPC")
}

func createListenerAndServer(info *ConnectionInfo) (listener net.Listener, server *grpc.Server, err error) {
	listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", info.Address, info.Port))
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to listen on %s:%d", info.Address, info.Port)
	}

	info.Port = listener.Addr().(*net.TCPAddr).Port

	server = grpc.NewServer()
	return
}

// CreateClientConnection creates a connection to a remote services with the given connection information
func CreateClientConnection(info *ConnectionInfo) (*grpc.ClientConn, error) {
	var options = []grpc.DialOption{grpc.WithInsecure()}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", info.Address, info.Port), options...)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to connect to %s:%d", info.Address, info.Port)
	}

	return conn, nil
}
