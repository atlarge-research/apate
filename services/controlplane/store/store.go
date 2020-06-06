// Package store provides state to the apate cluster
package store

import (
	"container/list"
	"sync"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario"

	"github.com/pkg/errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/health"

	"github.com/google/uuid"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/apatelet"
)

//TODO: Multi-master not soon :tm:

// Store represents the store of the control plane
type Store interface {
	// AddNode adds the given Node to the Apate cluster
	AddNode(*Node) error

	// RemoveNode removes the given Node from the Apate cluster by uuid
	RemoveNode(uuid.UUID) error

	// RemoveNodes removes the given Nodes from the Apate cluster by uuids
	RemoveNodes([]uuid.UUID) error

	// GetNode returns the node with the given uuid
	GetNode(uuid.UUID) (Node, error)

	// SetNodeStatus sets the status of the node with the given uuid
	SetNodeStatus(uuid.UUID, health.Status) error

	// GetNodes returns an array containing all nodes in the Apate cluster
	GetNodes() ([]Node, error)

	// GetNodesByLabel returns an array containing all nodes in the Apate cluster with the given label
	GetNodesByLabel(string) ([]Node, error)

	// ClearNodes removes all nodes from the Apate cluster
	ClearNodes() error

	// AddResourcesToQueue adds a node resource to the queue
	AddResourcesToQueue([]scenario.NodeResources) error

	// GetResourceFromQueue returns the first NodeResources struct in the list
	GetResourceFromQueue() (*scenario.NodeResources, error)

	// SetApateletScenario adds the ApateletScenario to the store
	SetApateletScenario(*apatelet.ApateletScenario) error

	// GetApateletScenario gets the ApateletScenario
	GetApateletScenario() (*apatelet.ApateletScenario, error)
}

type store struct {
	nodes        map[uuid.UUID]Node
	nodesByLabel map[string][]Node
	nodeLock     sync.RWMutex

	resourceQueue list.List
	resourceLock  sync.Mutex

	scenario     *apatelet.ApateletScenario
	scenarioLock sync.RWMutex
}

// NewStore creates a new empty cluster
func NewStore() Store {
	return &store{
		nodes:        make(map[uuid.UUID]Node),
		nodesByLabel: make(map[string][]Node),
	}
}

func (s *store) AddNode(node *Node) error {
	s.nodeLock.Lock()
	defer s.nodeLock.Unlock()

	// Check if node already exists (uuid collision)
	if _, ok := s.nodes[node.UUID]; ok {
		return errors.Errorf("node with uuid '%s' already exists", node.UUID.String())
	}

	if len(node.Label) == 0 {
		return errors.Errorf("node %s has no label", node.UUID.String())
	}

	s.nodes[node.UUID] = *node
	s.nodesByLabel[node.Label] = append(s.nodesByLabel[node.Label], *node)

	return nil
}

func (s *store) removeNodeByUUID(uuid uuid.UUID) error {
	node := s.nodes[uuid]
	if node == (Node{}) {
		return nil
	}

	label := node.Label

	if len(label) == 0 {
		return errors.Errorf("node %s has no label", node.UUID.String())
	}

	for i, cur := range s.nodesByLabel[label] {
		if cur.UUID == node.UUID {
			le := len(s.nodesByLabel[label])
			s.nodesByLabel[label][i] = s.nodesByLabel[label][le-1]
			s.nodesByLabel[label] = s.nodesByLabel[label][:le-1]
		}
	}

	delete(s.nodes, node.UUID)
	return nil
}

func (s *store) RemoveNode(uuid uuid.UUID) error {
	s.nodeLock.Lock()
	defer s.nodeLock.Unlock()

	return s.removeNodeByUUID(uuid)
}

func (s *store) RemoveNodes(uuids []uuid.UUID) (err error) {
	s.nodeLock.Lock()
	defer s.nodeLock.Unlock()

	for _, id := range uuids {
		if remErr := s.removeNodeByUUID(id); err == nil && remErr != nil {
			err = remErr
		}
	}

	return errors.Wrap(err, "failed to remove nodes")
}

func (s *store) GetNode(uuid uuid.UUID) (Node, error) {
	s.nodeLock.RLock()
	defer s.nodeLock.RUnlock()

	if node, ok := s.nodes[uuid]; ok {
		return node, nil
	}

	return Node{}, errors.Errorf("node with uuid '%s' not found", uuid.String())
}

func (s *store) SetNodeStatus(uuid uuid.UUID, status health.Status) error {
	s.nodeLock.Lock()
	defer s.nodeLock.Unlock()

	if node, ok := s.nodes[uuid]; ok {
		node.Status = status
		s.nodes[uuid] = node
		return nil
	}

	return errors.Errorf("node with uuid '%s' not found", uuid.String())
}

func (s *store) GetNodes() ([]Node, error) {
	s.nodeLock.RLock()
	defer s.nodeLock.RUnlock()

	nodes := make([]Node, 0, len(s.nodes))

	for _, node := range s.nodes {
		nodes = append(nodes, node)
	}

	return nodes, nil
}

func (s *store) GetNodesByLabel(label string) ([]Node, error) {
	s.nodeLock.RLock()
	defer s.nodeLock.RUnlock()

	return s.nodesByLabel[label], nil
}

func (s *store) ClearNodes() error {
	s.nodeLock.Lock()
	defer s.nodeLock.Unlock()

	s.nodes = make(map[uuid.UUID]Node)
	s.nodesByLabel = make(map[string][]Node)
	return nil
}

func (s *store) AddResourcesToQueue(resources []scenario.NodeResources) error {
	s.resourceLock.Lock()
	defer s.resourceLock.Unlock()

	for _, res := range resources {
		res := res
		s.resourceQueue.PushBack(&res)
	}
	return nil
}

func (s *store) GetResourceFromQueue() (*scenario.NodeResources, error) {
	s.resourceLock.Lock()
	defer s.resourceLock.Unlock()

	if s.resourceQueue.Len() == 0 {
		return nil, errors.New("no NodeResources available for this apatelet")
	}

	res := s.resourceQueue.Front()
	s.resourceQueue.Remove(res)

	return res.Value.(*scenario.NodeResources), nil
}

func (s *store) SetApateletScenario(scenario *apatelet.ApateletScenario) error {
	s.scenarioLock.Lock()
	defer s.scenarioLock.Unlock()

	s.scenario = scenario
	return nil
}

func (s *store) GetApateletScenario() (*apatelet.ApateletScenario, error) {
	s.scenarioLock.RLock()
	defer s.scenarioLock.RUnlock()

	if s.scenario == nil {
		return nil, errors.New("no scenario available yet")
	}

	return s.scenario, nil
}
