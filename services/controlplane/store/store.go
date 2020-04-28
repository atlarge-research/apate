// Package store provides state to the apate cluster
package store

import (
	"fmt"
	"sync"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/health"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/apatelet"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/normalization"

	"github.com/google/uuid"
)

//TODO: Multi-master soon :tm:

// Store represents the store of the control plane
type Store interface {
	// AddNode adds the given Node to the Apate cluster
	AddNode(*Node) error

	// RemoveNode removes the given Node from the Apate cluster
	RemoveNode(*Node) error

	// GetNode returns the node with the given uuid
	GetNode(uuid.UUID) (Node, error)

	SetNodeStatus(uuid.UUID, health.Status) error

	// GetNodes returns an array containing all nodes in the Apate cluster
	GetNodes() ([]Node, error)

	// ClearNodes removes all nodes from the Apate cluster
	ClearNodes() error

	// AddResourcesToQueue adds a node resource to the queue
	AddResourcesToQueue([]normalization.NodeResources) error

	// AddApateletScenario adds the ApateletScenario to the store
	AddApateletScenario(*apatelet.ApateletScenario) error

	// GetApateletScenario gets the ApateletScenario
	GetApateletScenario() (*apatelet.ApateletScenario, error)
}

type store struct {
	nodes    map[uuid.UUID]Node
	nodeLock sync.RWMutex
}

// NewStore creates a new empty cluster
func NewStore() Store {
	return &store{
		nodes: make(map[uuid.UUID]Node),
	}
}

// AddNode adds the given Node to the Apate cluster
func (s *store) AddNode(node *Node) error {
	s.nodeLock.Lock()
	defer s.nodeLock.Unlock()

	// Check if node already exists (uuid collision)
	if _, ok := s.nodes[node.UUID]; ok {
		return fmt.Errorf("node with uuid '%s' already exists", node.UUID.String())
	}

	s.nodes[node.UUID] = *node

	return nil
}

// RemoveNode removes the given Node from the Apate cluster
func (s *store) RemoveNode(node *Node) error {
	s.nodeLock.Lock()
	defer s.nodeLock.Unlock()

	delete(s.nodes, node.UUID)
	return nil
}

// GetNode returns the node with the given uuid
func (s *store) GetNode(uuid uuid.UUID) (Node, error) {
	s.nodeLock.RLock()
	defer s.nodeLock.RUnlock()

	if node, ok := s.nodes[uuid]; ok {
		return node, nil
	}

	return Node{}, fmt.Errorf("node with uuid '%s' not found", uuid.String())
}

func (s *store) SetNodeStatus(uuid uuid.UUID, status health.Status) error {
	s.nodeLock.Lock()
	defer s.nodeLock.Unlock()

	if node, ok := s.nodes[uuid]; ok {
		node.Status = status
		return nil
	}

	return fmt.Errorf("node with uuid '%s' not found", uuid.String())
}

// GetNodes returns an array containing all nodes in the Apate cluster
func (s *store) GetNodes() ([]Node, error) {
	s.nodeLock.RLock()
	defer s.nodeLock.RUnlock()

	nodes := make([]Node, 0, len(s.nodes))

	for _, node := range s.nodes {
		nodes = append(nodes, node)
	}

	return nodes, nil
}

func (s *store) ClearNodes() error {
	s.nodeLock.Lock()
	defer s.nodeLock.Unlock()

	s.nodes = make(map[uuid.UUID]Node)
	return nil
}

func (s *store) AddResourcesToQueue(resources []normalization.NodeResources) error {
	return nil
}

func (s *store) AddApateletScenario(scenario *apatelet.ApateletScenario) error {
	return nil
}

func (s *store) GetApateletScenario() (*apatelet.ApateletScenario, error) {
	return nil, nil
}
