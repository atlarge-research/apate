// Package scheduler handles the scheduling of tasks
package scheduler

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/crd/node"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/crd/pod"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
)

// Scheduler is struct on which all scheduler functionality is implemented.
type Scheduler struct {
	store *store.Store
	ctx   context.Context

	readyCh  chan struct{}
	updateCh chan struct{}

	prevT     time.Duration
	startTime time.Duration
}

// New returns a new scheduler
func New(ctx context.Context, st *store.Store) Scheduler {
	return Scheduler{
		store:     st,
		ctx:       ctx,
		readyCh:   make(chan struct{}),
		updateCh:  make(chan struct{}),
		prevT:     0,
		startTime: 0,
	}
}

// EnableScheduler enables the scheduler
// this will wait until StartScheduler() is called, after that it
// will poll the store queue for changes and write errors to a 3-buffered channel
func (s *Scheduler) EnableScheduler() <-chan error {
	ech := make(chan error, 3)

	go func() {
		// Wait for start
		<-s.readyCh

		for {
			select {
			case <-s.ctx.Done():
				return
			default:
			}

			// Run iteration
			if done, delay := s.runner(ech); done {
				select {
				case <-s.updateCh:
					//
				case <-time.After(delay):
					//
				}
			}
		}
	}()

	return ech
}

// StartScheduler sets the start time and starts the scheduler
func (s *Scheduler) StartScheduler(startTime time.Duration) {
	s.startTime = startTime
	s.readyCh <- struct{}{}
}

// WakeScheduler wakes up the scheduler
func (s *Scheduler) WakeScheduler() {
	select {
	case s.updateCh <- struct{}{}:
	default:
	}
}

func (s *Scheduler) runner(ech chan error) (bool, time.Duration) {
	now := time.Duration(time.Now().UnixNano())

	relativeTime, taskFound, err := (*s.store).PeekTask()
	if err != nil {
		ech <- errors.Wrap(err, "failed to peek at the next task in the store")
		return false, 0
	}

	if !taskFound {
		return true, 0
	}

	if now >= relativeTime+s.startTime {
		task, err := (*s.store).PopTask()
		if err != nil {
			ech <- errors.Wrap(err, "failed to pop the next task from the store")
			return false, 0
		}

		if relativeTime >= s.prevT {
			s.prevT = relativeTime
			go s.taskHandler(ech, task)
		}

		nextRelativeTime, nextTaskFound, err := (*s.store).PeekTask()
		if err != nil {
			ech <- errors.Wrap(err, "failed to peek the next task from the store")
			return false, 0
		}

		if nextTaskFound {
			nextTime := nextRelativeTime + s.startTime

			delay := nextTime - now - time.Second
			if delay < 0 {
				return false, 0
			}
			return false, delay
		}
	}

	return false, 0
}

func (s Scheduler) taskHandler(ech chan error, t *store.Task) {
	isPod, err := t.IsPod()
	if err != nil {
		ech <- errors.Wrap(err, "failed to determine task type")
		return
	}

	if isPod {
		err := pod.SetPodFlags(s.store, t.PodTask.Label, t.PodTask.State)
		if err != nil {
			ech <- errors.Wrap(err, "failed to set pod flags")
		}
	} else {
		node.SetNodeFlags(s.store, t.NodeTask.State)
	}
}
