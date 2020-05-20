// Package services contains the implementation of the GRPC services of the control plane
package services

import (
	"context"
	"io/ioutil"
	"log"
	"net"

	"github.com/google/uuid"

	"github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/watchdog"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/controlplane"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster"

	"google.golang.org/grpc/peer"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/store"
)

type clusterOperationService struct {
	store             *store.Store
	kubernetesCluster cluster.KubernetesCluster
}

// RegisterClusterOperationService registers a new clusterOperationService with the given gRPC server
func RegisterClusterOperationService(server *service.GRPCServer, store *store.Store, kubernetesCluster cluster.KubernetesCluster) {
	controlplane.RegisterClusterOperationsServer(server.Server, &clusterOperationService{
		store:             store,
		kubernetesCluster: kubernetesCluster,
	})
}

// JoinCluster accepts an incoming request from an Apatelet to join the store
func (s *clusterOperationService) JoinCluster(ctx context.Context, info *controlplane.ApateletInformation) (*controlplane.JoinInformation, error) {
	// Get connection information
	p, _ := peer.FromContext(ctx)
	addr := p.Addr.(*net.TCPAddr)
	connectionInfo := *service.NewConnectionInfo(addr.IP.String(), int(info.Port), false)
	log.Printf("Received request to join apate store from %s\n", connectionInfo.Address)

	// Retrieve node resources
	st := *s.store
	nodeResources, err := st.GetResourceFromQueue()

	// TODO: Maybe reinsert NodeResources depending on the type of error?
	if err != nil {
		log.Printf("Unable to allocate resources for node %v: %s", connectionInfo, err.Error())
		return nil, err
	}

	// Get connection information and create node
	node := store.NewNode(connectionInfo, nodeResources, nodeResources.Selector)

	// Add to apate store
	err = st.AddNode(node)

	if err != nil {
		return nil, err
	}

	log.Printf("Added node to apate store: %v\n", node)

	return &controlplane.JoinInformation{
		KubeConfig: s.kubernetesCluster.KubeConfig.Bytes,
		NodeUuid:   node.UUID.String(),
		Hardware: &controlplane.NodeHardware{
			Memory:           nodeResources.Memory,
			Cpu:              nodeResources.CPU,
			Storage:          nodeResources.Storage,
			EphemeralStorage: nodeResources.EphemeralStorage,
			MaxPods:          nodeResources.MaxPods,
		},
	}, nil
}

// LeaveCluster removes the node from the store
// This will maybe also remove it from k8s itself, TBD
func (s *clusterOperationService) LeaveCluster(_ context.Context, leaveInformation *controlplane.LeaveInformation) (*empty.Empty, error) {
	log.Printf("Received request to leave apate cluster from node %s\n", leaveInformation.NodeUuid)

	id, err := uuid.Parse(leaveInformation.NodeUuid)
	if err != nil {
		return nil, err
	}

	err = watchdog.RemoveNodeWithUUID(id, s.store, &s.kubernetesCluster)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *clusterOperationService) GetKubeConfig(_ context.Context, _ *empty.Empty) (*controlplane.KubeConfig, error) {
	cfg, err := ioutil.ReadFile("/tmp/apate/config-ext")

	if err != nil {
		return nil, err
	}

	return &controlplane.KubeConfig{Config: cfg}, nil
}
