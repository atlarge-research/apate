package scheduler

import (
	"context"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/apatelet"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
	"time"
)

type Scheduler struct {
	store *store.Store
}

func (s *Scheduler) StartScheduler(ctx context.Context) chan error {
	ech := make(chan error)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			if err := s.runner(); err != nil {
				ech <- err
			}
			time.Sleep(0)
		}
	}()

	return ech
}

func (s *Scheduler) runner() error {
	now := time.Now().Unix()

	nextT, err := (*s.store).PeekTask()
	if err != nil {
		return err
	}
	if now >= nextT {
		task, err := (*s.store).PopTask()
		if err != nil {
			return err
		}

		go s.taskHandler(task)
	}

	return nil
}

func (s *Scheduler) taskHandler(t *apatelet.Task) {
	// TODO: Implement
}
