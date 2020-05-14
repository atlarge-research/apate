// Package scheduler handles the scheduling of tasks
package scheduler

import (
	"context"
	"time"

	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/crd"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/any"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
)

// Scheduler is struct on which all scheduler functionality is implemented.
type Scheduler struct {
	store *store.Store
}

// StartScheduler starts running the scheduler
// this will poll the store queue for changes and write errors to a 3-buffered channel
func (s *Scheduler) StartScheduler(ctx context.Context) chan error {
	ech := make(chan error, 3)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			s.runner(ech)
		}
	}()

	return ech
}

func (s *Scheduler) runner(ech chan error) {
	now := time.Now().UnixNano()

	nextT, err := (*s.store).PeekTask()
	if err != nil {
		ech <- err
		return
	}
	if now >= nextT {
		task, err := (*s.store).PopTask()
		if err != nil {
			ech <- err
			return
		}

		go s.taskHandler(ech, task)
	}
}

func (s Scheduler) taskHandler(ech chan error, t *store.Task) {
	if t.IsPod {
		err := crd.SetPodFlags(s.store, t.PodTask)
		if err != nil {
			ech <- err
		}
	} else {
		for k, mv := range (*t.NodeTask).NodeEventFlags {
			v, err := any.Unmarshal(mv)
			if err != nil {
				ech <- err
				continue
			}

			(*s.store).SetNodeFlag(k, v)
		}
	}
}
