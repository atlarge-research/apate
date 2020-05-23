// Package scheduler handles the scheduling of tasks
package scheduler

import (
	"context"
	"time"

	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/crd/node"

	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/crd/pod"

	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
)

// Scheduler is struct on which all scheduler functionality is implemented.
type Scheduler struct {
	store *store.Store
	ctx   context.Context

	readyCh chan struct{}

	prevT     int64
	startTime int64
}

// New returns a new scheduler
func New(ctx context.Context, st *store.Store) Scheduler {
	return Scheduler{
		store:     st,
		ctx:       ctx,
		readyCh:   make(chan struct{}),
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
			s.runner(ech)
		}
	}()

	return ech
}

// StartScheduler sets the start time and starts the scheduler
func (s *Scheduler) StartScheduler(startTime int64) {
	s.startTime = startTime
	s.readyCh <- struct{}{}
}

func (s *Scheduler) runner(ech chan error) {
	now := time.Now().UnixNano()

	relativeTime, err := (*s.store).PeekTask()
	if err != nil {
		// TODO: Check error type properly @jona, I know this is the bad
		if err.Error() != "no tasks left" {
			ech <- err
		}
		return
	}

	if now >= relativeTime+s.startTime {
		task, err := (*s.store).PopTask()
		if err != nil {
			ech <- err
			return
		}

		if relativeTime >= s.prevT {
			s.prevT = relativeTime
			go s.taskHandler(ech, task)
		}
	}
}

func (s Scheduler) taskHandler(ech chan error, t *store.Task) {
	isPod, err := t.IsPod()
	if err != nil {
		ech <- err
		return
	}

	if isPod {
		err := pod.SetPodFlags(s.store, t.PodTask.Label, t.PodTask.State)
		if err != nil {
			ech <- err
		}
	} else {
		node.SetNodeFlags(s.store, t.NodeTask.State)
	}
}
