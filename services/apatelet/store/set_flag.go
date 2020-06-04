package store

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

	s.nodeFlags = flags
}

func (s *store) SetPodFlags(label string, flags Flags) {
	s.podFlagLock.Lock()
	s.podFlags[label] = flags
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
	s.podTimeFlags[label] = flags
}
