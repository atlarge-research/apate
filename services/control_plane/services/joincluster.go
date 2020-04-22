package services

import (
	"context"
	"log"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/join_cluster"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/control_plane/cluster"
)

type joinClusterService struct {
	cluster *cluster.ApateCluster
}

// RegisterJoinClusterService registers a new joinClusterService with the given gRPC server
func RegisterJoinClusterService(server *service.GRPCServer, cluster *cluster.ApateCluster) {
	join_cluster.RegisterJoinClusterServer(server.Server, &joinClusterService{
		cluster: cluster,
	})
}

// JoinCluster accepts an incoming request from a virtual kubelet to join the cluster
func (s *joinClusterService) JoinCluster(context.Context, *empty.Empty) (*join_cluster.JoinInformation, error) {
	//TODO: TLS bool from somewhere?

	// Get connection information
	connectionInfo := *service.NewConnectionInfo("localhost", 8080, false)
	log.Printf("Received request to join apate cluster from %v\n", connectionInfo)

	// Get connection information and create node
	node := cluster.NewNode(connectionInfo)

	// Add to apate cluster
	s.cluster.AddNode(node)
	log.Printf("Added node to apate cluster: %v\n", node)

	return &join_cluster.JoinInformation{
		KubernetesJoinToken: "TODO",
		NodeUUID:            node.UUID.String(),
	}, nil
}
