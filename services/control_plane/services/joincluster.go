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
func (s *joinClusterService) JoinCluster(ctx context.Context, _ *empty.Empty) (*join_cluster.JoinInformation, error) {
	//TODO: TLS bool from somewhere?

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

	return &join_cluster.JoinInformation{
		KubeConfig:  config,
		KubeContext: "kind-Apate",
		NodeUUID:    node.UUID.String(),
	}, nil
}

func getKubeConfigData(path string) []byte {
	data, err := ioutil.ReadFile(filepath.Join(os.TempDir(), filepath.Clean(path)))

	// TODO: Better error handling
	if err != nil {
		log.Fatalf("Could not read kube config: %v", err)
	}

	return data
}
