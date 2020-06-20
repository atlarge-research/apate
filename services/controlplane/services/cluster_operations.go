// Package services contains the implementation of the GRPC services of the control plane
package services

import (
	"context"
	"log"
	"net"

	"github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/cluster"

	"github.com/google/uuid"

	"github.com/pkg/errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/controlplane"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes"

	"google.golang.org/grpc/peer"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/atlarge-research/opendc-emulate-kubernetes/internal/service"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/store"
)

type clusterOperationService struct {
	store             *store.Store
	kubernetesCluster *kubernetes.Cluster
}

// RegisterClusterOperationService registers a new clusterOperationService with the given gRPC server
func RegisterClusterOperationService(server *service.GRPCServer, store *store.Store, kubernetesCluster *kubernetes.Cluster) {
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
	connectionInfo := *service.NewConnectionInfo(addr.IP.String(), int(info.Port))
	log.Printf("Received request to join apate store from %s\n", connectionInfo.Address)

	// Retrieve node resources
	st := *s.store
	nodeResources, err := st.GetResourceFromQueue()

	// TODO: Maybe reinsert NodeResources depending on the type of error?
	if err != nil {
		err = errors.Wrap(err, "failed to get node resources from queue")
		log.Println(err)
		return nil, err
	}

	// Get connection information and create node
	node := store.NewNode(connectionInfo, nodeResources, nodeResources.Label)

	// Add to apate store
	err = st.AddNode(node)

	if err != nil {
		err = errors.Wrap(err, "failed to add node to queue")
		log.Println(err)
		return nil, err
	}

	log.Printf("Added node to apate store: %v\n", node)

	// Check start time for scenario
	time := int64(-1)
	scenario, err := st.GetApateletScenario()
	if err == nil {
		time = scenario.StartTime
	}

	return &controlplane.JoinInformation{
		KubeConfig: s.kubernetesCluster.KubeConfig.Bytes,
		NodeUuid:   node.UUID.String(),
		NodeLabel:  nodeResources.Label,
		StartTime:  time,

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
		return nil, errors.Wrap(err, "failed to remove node from cluster")
	}

	var clusterAPI kubernetes.ClusterAPI = s.kubernetesCluster
	_, _, err = cluster.RemoveNodeWithUUID(id, s.store, &clusterAPI)
	if err != nil {
		return nil, errors.Wrapf(err, "removing node with uuid %v during leave cluster failed", id)
	}

	return &empty.Empty{}, nil
}

func (s *clusterOperationService) GetKubeConfig(_ context.Context, _ *empty.Empty) (*controlplane.KubeConfig, error) {
	return &controlplane.KubeConfig{Config: s.kubernetesCluster.KubeConfig.Bytes}, nil
}
