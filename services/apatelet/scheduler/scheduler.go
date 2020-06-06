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

const (
	// sleepMargin represents the amount of time the scheduler wakes up before a next task
	sleepMargin = time.Second
)

// Scheduler is struct on which all scheduler functionality is implemented.
type Scheduler struct {
	store *store.Store

	readyCh  chan struct{}
	updateCh chan struct{}

	prevT     time.Time
	startTime time.Time
}

// New returns a new scheduler
func New(st *store.Store) Scheduler {
	return Scheduler{
		store:     st,
		readyCh:   make(chan struct{}),
		updateCh:  make(chan struct{}),
		prevT:     time.Unix(0, 0),
		startTime: time.Unix(0, 0),
	}
}

// EnableScheduler enables the scheduler
// this will wait until StartScheduler() is called, after that it
// will poll the store queue for changes and write errors to a 3-buffered channel
func (s *Scheduler) EnableScheduler(ctx context.Context) <-chan error {
	ech := make(chan error, 3)

	go func() {
		// Wait for start
		select {
		case <-ctx.Done():
			return
		case <-s.readyCh:
			//
		}

		s.scheduleLoop(ctx, ech)
	}()

	return ech
}

func (s *Scheduler) scheduleLoop(ctx context.Context, ech chan<- error) {
	timer := time.NewTimer(time.Millisecond)
	defer timer.Stop()

	for {
		if err := ctx.Err(); err != nil {
			ech <- errors.Wrap(err, "scheduler stopped")
			return
		}

		// Run iteration
		done, delay := s.runner(ech)

		if done {
			select {
			case <-ctx.Done():
				return
			case <-s.updateCh:
				//
			}
		}

		if delay > time.Millisecond {
			timer.Reset(delay)
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
				//
			case <-s.updateCh:
				//
			}
		}
	}
}

// StartScheduler sets the start time and starts the scheduler
func (s *Scheduler) StartScheduler(startTime int64) {
	s.startTime = time.Unix(0, startTime)
	s.readyCh <- struct{}{}
}

// WakeScheduler wakes up the scheduler
func (s *Scheduler) WakeScheduler() {
	select {
	case s.updateCh <- struct{}{}:
	default:
	}
}

func (s *Scheduler) runner(ech chan<- error) (bool, time.Duration) {
	now := time.Now()

	relativeTime, taskFound, err := (*s.store).PeekTask()
	if err != nil {
		ech <- errors.Wrap(err, "failed to peek at the next task in the store")
		return false, 0
	}

	if !taskFound {
		return true, 0
	}

	scheduledTime := s.startTime.Add(relativeTime)

	if now.After(scheduledTime) {
		var task *store.Task
		task, err = (*s.store).PopTask()
		if err != nil {
			ech <- errors.Wrap(err, "failed to pop the next task from the store")
			return false, 0
		}

		if s.prevT.Before(scheduledTime) || s.prevT.Equal(scheduledTime) {
			s.prevT = scheduledTime
			go s.taskHandler(ech, task)
		}
	}

	nextRelativeTime, nextTaskFound, err := (*s.store).PeekTask()
	if err != nil {
		ech <- errors.Wrap(err, "failed to peek the next task from the store")
		return false, 0
	}

	if nextTaskFound {
		nextTime := s.startTime.Add(nextRelativeTime)

		delay := nextTime.Sub(now) - sleepMargin

		if delay < time.Nanosecond {
			return false, time.Nanosecond
		}
		return false, delay
	}

	return false, 0
}

func (s *Scheduler) taskHandler(ech chan<- error, t *store.Task) {
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
