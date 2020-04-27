// Package services contains the implementation of the GRPC services of the control plane
package services

import (
	"context"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/controlplane"

	"google.golang.org/grpc/peer"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/store"
)

type clusterOperationService struct {
	store *store.Store
}

// RegisterClusterOperationService registers a new clusterOperationService with the given gRPC server
func RegisterClusterOperationService(server *service.GRPCServer, cluster *store.Store) {
	controlplane.RegisterClusterOperationsServer(server.Server, &clusterOperationService{
		store: cluster,
	})
}

// JoinCluster accepts an incoming request from an Apatelet to join the store
func (s *clusterOperationService) JoinCluster(ctx context.Context, _ *empty.Empty) (*controlplane.JoinInformation, error) {
	//TODO: TLS bool from somewhere?

	// TODO: Check if node is already joined based on remote address
	// Get connection information
	p, _ := peer.FromContext(ctx)
	addr := p.Addr.(*net.TCPAddr)
	connectionInfo := *service.NewConnectionInfo(addr.IP.String(), addr.Port, false)
	log.Printf("Received request to join apate store from %v\n", connectionInfo)

	// Retrieve node resources
	st := *s.store
	nodeResources, err := st.GetResourceFromQueue()

	// TODO: Maybe reinsert NodeResources depending on the type of error?
	if err != nil {
		log.Printf("Unable to allocate resources for node %v: %s", connectionInfo, err.Error())
		return nil, err
	}

	//TODO: Retrieve path from somewhere else
	// Retrieving kube config
	config := getKubeConfigData("/apate/config")

	// Get connection information and create node
	node := store.NewNode(connectionInfo, nodeResources)

	// Add to apate store
	err = st.AddNode(node)

	if err != nil {
		return nil, err
	}

	log.Printf("Added node to apate store: %v\n", node)

	// TODO: Retrieve proper kube context from somewhere
	return &controlplane.JoinInformation{
		KubeConfig:  config,
		KubeContext: "kind-Apate",
		NodeUuid:    node.UUID.String(),
		Hardware:    &controlplane.NodeHardware{
			Ram:     nodeResources.RAM,
			Cpu:     nodeResources.CPU,
			MaxPods: nodeResources.MaxPods,
		},
	}, nil
}

// LeaveCluster removes the node from the store
// This will maybe also remove it from k8s itself, TBD
func (s *clusterOperationService) LeaveCluster(_ context.Context, leaveInformation *controlplane.LeaveInformation) (*empty.Empty, error) {
	// TODO: Maybe check if the remote address is still the same? idk

	// TODO: Remove node from store and maybe from k8s too?
	log.Printf("Received request to leave apate store from node %s\n", leaveInformation.NodeUuid)
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
