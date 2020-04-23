package services

import (
	"context"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"

	"google.golang.org/grpc/peer"

	"github.com/golang/protobuf/ptypes/empty"

	cluster_ops "github.com/atlarge-research/opendc-emulate-kubernetes/api/cluster_operations"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/control_plane/cluster"
)

type clusterOperationService struct {
	cluster *cluster.ApateCluster
}

// RegisterJoinClusterService registers a new clusterOperationService with the given gRPC server
func RegisterJoinClusterService(server *service.GRPCServer, cluster *cluster.ApateCluster) {
	cluster_ops.RegisterClusterOperationsServer(server.Server, &clusterOperationService{
		cluster: cluster,
	})
}

// JoinCluster accepts an incoming request from a virtual kubelet to join the cluster
func (s *clusterOperationService) JoinCluster(ctx context.Context, _ *empty.Empty) (*cluster_ops.JoinInformation, error) {
	//TODO: TLS bool from somewhere?

	// TODO: Check if node is already joined based on remote address
	// Get connection information
	p, _ := peer.FromContext(ctx)
	addr := p.Addr.(*net.TCPAddr)
	connectionInfo := *service.NewConnectionInfo(addr.IP.String(), addr.Port, false)
	log.Printf("Received request to join apate cluster from %v\n", connectionInfo)

	//TODO: Retrieve path from somewhere else
	// Retrieving kube config
	config := getKubeConfigData("/apate/config")

	// Get connection information and create node
	node := cluster.NewNode(connectionInfo)

	// Add to apate cluster
	err := (*s.cluster).AddNode(node)

	if err != nil {
		return nil, err
	}

	log.Printf("Added node to apate cluster: %v\n", node)

	return &cluster_ops.JoinInformation{
		KubeConfig:  config,
		KubeContext: "kind-Apate",
		NodeUUID:    node.UUID.String(),
	}, nil
}

// LeaveCluster removes the node from the cluster
// This will maybe also remove it from k8s itself, TBD
func (s *clusterOperationService) LeaveCluster(_ context.Context, leaveInformation *cluster_ops.LeaveInformation) (*empty.Empty, error) {
	// TODO: Maybe check if the remote address is still the same? idk

	// TODO: Remove node from cluster and maybe from k8s too?
	log.Printf("Received request to leave apate cluster from node %s\n", leaveInformation.NodeUUID)
	return &empty.Empty{}, nil
}

func getKubeConfigData(path string) []byte {
	data, err := ioutil.ReadFile(filepath.Join(os.TempDir(), filepath.Clean(path)))

	// TODO: Better error handling
	if err != nil {
		log.Fatalf("Could not read kube config: %v", err)
	}

	return data
}
