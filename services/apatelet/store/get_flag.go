package store

import (
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"time"
)

type FlagGetter interface {
	// GetNodeFlag returns the value of the given node flag
	GetNodeFlag(events.NodeEventFlag) (interface{}, error)

	// GetPodFlag returns the value of the given pod for a configuration
	GetPodFlag(*corev1.Pod, events.PodEventFlag) (interface{}, error)
}

func (s *store) GetNodeFlag(id events.NodeEventFlag) (interface{}, error) {
	s.nodeFlagLock.RLock()
	defer s.nodeFlagLock.RUnlock()

	if val, ok := s.nodeFlags[id]; ok {
		return val, nil
	}

	if dv, ok := defaultNodeValues[id]; ok {
		return dv, nil
	}

	return nil, errors.New("flag not found in get node flag")
}

func (s *store) GetPodFlag(pod *corev1.Pod, flag events.PodEventFlag) (interface{}, error) {
	s.podFlagLock.Lock()
	defer s.podFlagLock.Unlock()

	label := getPodLabelByPod(pod)
	if label != "" {
		if val, ok := s.podFlags[label][flag]; ok {
			return val, nil
		}

		if val, ok := s.getPodTimeFlag(pod, flag, label); ok {
			return val, nil
		}
	}

	if dv, ok := defaultPodValues[flag]; ok {
		return dv, nil
	}

	return nil, errors.New("flag not found in get pod flag")
}

func (s *store) getPodTimeFlag(pod *corev1.Pod, flag events.PodEventFlag, label string) (interface{}, bool) {
	if _, ok := s.podTimeIndexCache[pod]; !ok {
		s.podTimeIndexCache[pod] = make(map[events.EventFlag]int)
	}

	podTimeIndex := 0
	if val, ok := s.podTimeIndexCache[pod][flag]; ok {
		podTimeIndex = val
	}

	podStartTime := time.Now()
	if pod.Status.StartTime != nil {
		podStartTime = pod.Status.StartTime.Time
	}

	timeFlags := s.podTimeFlags[label]
	lastIndexWithFlag := podTimeIndex
	for i := podTimeIndex; i < len(timeFlags); i++ {
		flags := timeFlags[i]

		if podStartTime.Add(flags.TimeSincePodStart).Before(time.Now()) {
			currentPodFlags := timeFlags[lastIndexWithFlag]
			s.podTimeIndexCache[pod][flag] = lastIndexWithFlag
			return currentPodFlags.Flags[flag], true
		}

		if _, ok := flags.Flags[flag]; ok {
			lastIndexWithFlag = i
		}
	}

	return nil, false
}
