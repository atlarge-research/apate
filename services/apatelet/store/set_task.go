package store

import (
	"container/heap"
	"github.com/pkg/errors"
)

type TaskSetter interface {
	// SetNodeTasks adds or updates node tasks
	// Existing node tasks will be removed if not in the list of tasks
	SetNodeTasks([]*Task) error

	// SetPodTasks adds or updates pod CRD tasks to the queue based on their label (<namespace>/<name>)
	// Existing pod tasks will be removed if not in the list of tasks
	SetPodTasks(string, []*Task) error
}

func (s *store) setTasksOfType(newTasks []*Task, check TaskTypeCheck) error {
	s.queueLock.Lock()
	defer s.queueLock.Unlock()

	for i, task := range s.queue.tasks {
		typeCheck, err := check(task)
		if err != nil {
			return errors.Wrap(err, "failed to determine task type")
		}

		if typeCheck {
			if len(newTasks) == 0 {
				heap.Remove(s.queue, i)
			} else {
				if newTasks[0] != nil {
					s.queue.tasks[i] = newTasks[0]
					// Replacing and then fixing instead of deleting all and pushing because it's slightly faster, see comments on heap.Fix
					heap.Fix(s.queue, i)
				}
				newTasks = newTasks[1:]
			}
		}
	}

	for _, remainingTask := range newTasks {
		if remainingTask != nil {
			heap.Push(s.queue, remainingTask)
		}
	}

	return nil
}

func (s *store) SetNodeTasks(tasks []*Task) error {
	return s.setTasksOfType(tasks, func(task *Task) (bool, error) {
		isNode, err := task.IsNode()
		if err != nil {
			return false, err
		}

		return isNode, nil
	})
}

func (s *store) SetPodTasks(label string, tasks []*Task) error {
	return s.setTasksOfType(tasks, func(task *Task) (bool, error) {
		isPod, err := task.IsPod()
		if err != nil {
			return false, err
		}

		return isPod && task.PodTask.Label == label, nil
	})
}
