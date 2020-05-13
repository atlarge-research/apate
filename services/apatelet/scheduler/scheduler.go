// Package scheduler handles the scheduling of tasks
package scheduler

import (
	"context"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario"
	v1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/emulatedpod/v1"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/normalization/translate"
	"time"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/apatelet"
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
		s.setPodFlags(ech, t)
	} else {
		for k, mv := range t.OriginalTask.(*apatelet.Task).NodeEventFlags {
			v, err := any.Unmarshal(mv)
			if err != nil {
				ech <- err
				continue
			}

			(*s.store).SetNodeFlag(k, v)
		}
	}
}

func (s Scheduler) setPodFlags(ech chan error, t *store.Task) {
	task := t.OriginalTask.(*v1.EmulatedPodTask)

	if task.CreatePodResponse != v1.RESPONSE_UNSET {
		(*s.store).SetPodFlag(t.Label, events.PodCreatePodResponse, scenario.Response_value[string(task.CreatePodResponse)])
	}

	if task.UpdatePodResponse != v1.RESPONSE_UNSET {
		(*s.store).SetPodFlag(t.Label, events.PodUpdatePodResponse, scenario.Response_value[string(task.UpdatePodResponse)])
	}

	if task.DeletePodResponse != v1.RESPONSE_UNSET {
		(*s.store).SetPodFlag(t.Label, events.PodDeletePodResponse, scenario.Response_value[string(task.DeletePodResponse)])
	}

	if task.GetPodResponse != v1.RESPONSE_UNSET {
		(*s.store).SetPodFlag(t.Label, events.PodGetPodResponse, scenario.Response_value[string(task.GetPodResponse)])
	}

	if task.GetPodStatusResponse != v1.RESPONSE_UNSET {
		(*s.store).SetPodFlag(t.Label, events.PodGetPodStatusResponse, scenario.Response_value[string(task.GetPodStatusResponse)])
	}

	// Resource usage
	s.setResourceBytes(ech, t, task.ResourceUsage.Memory, events.PodMemoryUsage)
	if task.ResourceUsage.CPU != -1 {
		(*s.store).SetPodFlag(t.Label, events.PodCPUUsage, task.ResourceUsage.CPU)
	}
	s.setResourceBytes(ech, t, task.ResourceUsage.Storage, events.PodStorageUsage)
	s.setResourceBytes(ech, t, task.ResourceUsage.EphemeralStorage, events.PodEphemeralStorageUsage)

	if task.PodStatus != v1.POD_STATUS_UNSET {
		(*s.store).SetPodFlag(t.Label, events.PodStatus, scenario.PodStatus_value[string(task.PodStatus)])
	}
}

func (s Scheduler) setResourceBytes(ech chan error, t *store.Task, unit string, flag events.PodEventFlag) {
	if unit != "-1" {
		bytes, err := translate.GetInBytes(unit, "storage")
		if err != nil {
			ech <- err
		}
		(*s.store).SetPodFlag(t.Label, flag, bytes)
	}
}
