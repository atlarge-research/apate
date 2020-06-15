// Package store provides a way for the apatelet to have state
package store

import (
	"container/heap"
	"sync"
	"time"

	corev1 "k8s.io/api/core/v1"

	podconfigv1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/podconfiguration/v1"

	"github.com/finitum/node-cli/stats"
	"github.com/pkg/errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
)

// Store represents the state of the apatelet
type Store interface {
	TaskSetter
	FlagSetter
	FlagGetter

	// RemovePodTasks removes pod CRD tasks from the queue based on their label (<namespace>/<name>)
	RemovePodTasks(string) error

	// PeekTask returns the start time of the next task in the priority queue, without removing it from the queue
	PeekTask() (time.Duration, bool, error)

	// PopTask returns the first task to be executed and removes it from the queue
	PopTask() (*Task, error)

	// AddPodFlagListener adds a listener which is called when the given flag is updated
	AddPodFlagListener(events.PodEventFlag, func(interface{}))
}

// Flags is a map from event flags to their interface value
type Flags map[events.EventFlag]interface{}

type podFlags map[string]Flags
type podListeners map[events.EventFlag][]func(interface{})

// TimeFlags contains Flags at a certain timestamp relative to the starting time of a pod
type TimeFlags struct {
	TimeSincePodStart time.Duration
	Flags             Flags
}
type podTimeFlags map[string][]*TimeFlags
type podTimeIndexCache map[*corev1.Pod]map[events.EventFlag]int

type store struct {
	queue     *taskQueue
	queueLock sync.RWMutex

	nodeFlags    Flags
	nodeFlagLock sync.RWMutex

	podFlags    podFlags
	podFlagLock sync.RWMutex

	podListeners     podListeners
	podListenersLock sync.RWMutex

	podTimeFlags      podTimeFlags
	podTimeIndexCache podTimeIndexCache
}

// NewStore returns an empty store
func NewStore() Store {
	q := newTaskQueue()
	heap.Init(q)

	return &store{
		queue:        q,
		nodeFlags:    make(Flags),
		podListeners: make(podListeners),
		podFlags:     make(podFlags),

		podTimeFlags:      make(podTimeFlags),
		podTimeIndexCache: make(podTimeIndexCache),
	}
}

func (s *store) RemovePodTasks(label string) error {
	s.queueLock.Lock()
	defer s.queueLock.Unlock()

	for i := len(s.queue.tasks) - 1; i >= 0; i-- {
		task := s.queue.tasks[i]

		isPod, err := task.IsPod()
		if err != nil {
			return errors.Wrap(err, "failed to determine task type")
		}

		if isPod && task.PodTask.Label == label {
			heap.Remove(s.queue, i)
		}
	}

	return nil
}

func (s *store) PeekTask() (time.Duration, bool, error) {
	s.queueLock.RLock()
	defer s.queueLock.RUnlock()

	if s.queue.Len() == 0 {
		return -1, false, nil
	}

	// Make sure the array in the pq didn't magically change to a different type
	if task, ok := s.queue.First().(*Task); ok {
		return task.RelativeTimestamp, true, nil
	}

	return -1, false, errors.New("array in pq magically changed to a different type")
}

func (s *store) PopTask() (*Task, error) {
	s.queueLock.Lock()
	defer s.queueLock.Unlock()

	if s.queue.Len() == 0 {
		return nil, errors.New("no tasks left")
	}

	// Make sure the array in the pq didn't magically change to a different type
	if task, ok := heap.Pop(s.queue).(*Task); ok {
		return task, nil
	}

	return nil, errors.New("array in pq magically changed to a different type")
}

func (s *store) AddPodFlagListener(flag events.PodEventFlag, cb func(interface{})) {
	s.podListenersLock.Lock()
	defer s.podListenersLock.Unlock()

	if listeners, ok := s.podListeners[flag]; ok {
		s.podListeners[flag] = append(listeners, cb)
	} else {
		s.podListeners[flag] = []func(interface{}){cb}
	}
}

func getPodLabelByPod(pod *corev1.Pod) (string, bool) {
	label, ok := pod.Labels[podconfigv1.PodConfigurationLabel]
	if !ok {
		return "", false
	}
	return pod.Namespace + "/" + label, true
}

var defaultNodeValues = map[events.EventFlag]interface{}{
	events.NodeCreatePodResponse:    scenario.ResponseUnset,
	events.NodeUpdatePodResponse:    scenario.ResponseUnset,
	events.NodeDeletePodResponse:    scenario.ResponseUnset,
	events.NodeGetPodResponse:       scenario.ResponseUnset,
	events.NodeGetPodStatusResponse: scenario.ResponseUnset,
	events.NodeGetPodsResponse:      scenario.ResponseUnset,
	events.NodePingResponse:         scenario.ResponseUnset,

	events.NodeAddedLatency: time.Duration(0),
}

var defaultPodValues = map[events.PodEventFlag]interface{}{
	events.PodCreatePodResponse:    scenario.ResponseUnset,
	events.PodUpdatePodResponse:    scenario.ResponseUnset,
	events.PodDeletePodResponse:    scenario.ResponseUnset,
	events.PodGetPodResponse:       scenario.ResponseUnset,
	events.PodGetPodStatusResponse: scenario.ResponseUnset,

	events.PodResources: &stats.PodStats{},

	events.PodStatus: scenario.PodStatusUnset,
}
