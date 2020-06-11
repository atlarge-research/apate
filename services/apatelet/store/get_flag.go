package store

import (
	"time"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
)

// FlagGetter defines function aiding in retrieving flags
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

	label, ok := getPodLabelByPod(pod)
	if ok {
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

// getPodTimeFlag returns the pod time flag that is currently active for the given pod
// Meaning, given the current time, the pod (from which its start time is retrieved) and the flag, what is the expected state?
// It does this by retrieving the index cache for the flag/pod combination: the last index in the podTimeFlags that is checked for the current pod
// From this index it will continue to check next indices for the flag
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
	previousIndex := podTimeIndex
	for i := podTimeIndex; i < len(timeFlags); i++ {
		flags := timeFlags[i]

		podSinceStart := podStartTime.Add(flags.TimeSincePodStart)

		// The current index contains the expected flag and is still before the podSinceStart
		if _, ok := flags.Flags[flag]; ok && podSinceStart.Before(time.Now()) {
			previousIndex = i
		}

		// If the current flag is set too late or we are in the last iteration
		// We check for last iteration because there are no further flags to test afterwards
		if podSinceStart.After(time.Now()) || i == len(timeFlags)-1 {
			// Look at the previous index
			currentPodFlags := timeFlags[previousIndex]

			// If this index has time flags before now (it might not have if this is the first iteration)
			if podStartTime.Add(currentPodFlags.TimeSincePodStart).Before(time.Now()) {
				if pf, ok := currentPodFlags.Flags[flag]; ok {
					// Set cache and return it
					s.podTimeIndexCache[pod][flag] = previousIndex
					return pf, true
				}
			}

			// Else set the current index, as the next iteration can skip every index thus far
			s.podTimeIndexCache[pod][flag] = i
			break
		}
	}

	return nil, false
}
