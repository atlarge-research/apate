package store

import (
	"sort"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
)

// FlagSetter defines function aiding in setting flags
type FlagSetter interface {
	// SetNodeFlag sets the value of the given node flag
	SetNodeFlags(Flags)

	// SetNodeFlag sets the value of the given pod flag for a configuration
	SetPodFlags(string, Flags)

	// SetNodeFlag sets the value of the given pod flag for a configuration
	SetPodTimeFlags(string, []*TimeFlags)
}

func (s *store) SetNodeFlags(flags Flags) {
	s.nodeFlagLock.Lock()
	defer s.nodeFlagLock.Unlock()

	for k, v := range flags {
		s.nodeFlags[k] = v
	}
}

func (s *store) SetPodFlags(label string, flags Flags) {
	s.podFlagLock.Lock()
	if _, ok := s.podFlags[label]; !ok {
		s.podFlags[label] = make(Flags)
	}

	for k, v := range flags {
		s.podFlags[label][k] = v
	}
	s.podFlagLock.Unlock()

	s.podListenersLock.RLock()
	for flag, val := range flags {
		if listeners, ok := s.podListeners[flag]; ok {
			for _, listener := range listeners {
				listener(val)
			}
		}
	}
	s.podListenersLock.RUnlock()
}

func (s *store) SetPodTimeFlags(label string, flags []*TimeFlags) {
	s.podFlagLock.Lock()
	defer s.podFlagLock.Unlock()

	sort.Slice(flags, func(i, j int) bool {
		return flags[i].TimeSincePodStart < flags[j].TimeSincePodStart
	})

	s.podTimeFlags[label] = flags

	for pod := range s.podTimeIndexCache {
		if pl, ok := getPodLabelByPod(pod); ok && pl == label {
			s.podTimeIndexCache[pod] = make(map[events.EventFlag]int)
		}
	}
}
