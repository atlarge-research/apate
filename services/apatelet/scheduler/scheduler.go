package scheduler

import (
	"context"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/apatelet"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/any"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
	"time"
)

type Scheduler struct {
	store *store.Store
}

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
	now := time.Now().Unix()

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

	return
}

func (s Scheduler) taskHandler(ech chan error, t *apatelet.Task) {
	// TODO: Handle revert tasks here or in normalize
	for k, mv := range t.EventFlags {
		v, err := any.Unmarshal(mv)
		if err != nil {
			ech <- err
			continue
		}

		(*s.store).SetFlag(k, v)
	}
}
