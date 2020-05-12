package v1

import (
	"errors"

	"k8s.io/client-go/tools/cache"
)

// Informer is a wrapper around the Kubernetes informer which fills the informer cache
type Informer struct {
	store *cache.Store
}

// NewInformer creates a new informer wrapper
func NewInformer(store *cache.Store) *Informer {
	return &Informer{store: store}
}

// List returns a list of emulated pods.
func (i *Informer) List() (eps []EmulatedPod) {
	for _, ep := range (*i.store).List() {
		eps = append(eps, ep.(EmulatedPod))
	}
	return eps
}

// Find finds an emulated pod which can be identified by the given label
// This label should have the format <namespace>/<name> of the emulated pod.
func (i *Informer) Find(label string) (*EmulatedPod, bool, error) {
	key, exists, err := (*i.store).GetByKey(label)
	if err != nil {
		return nil, false, err
	}

	if exists {
		ep, ok := key.(*EmulatedPod)
		if ok {
			return ep, true, nil
		}

		return nil, false, errors.New("emulated pod not an emulated pod hmmmmmmm")
	}

	return nil, false, nil
}
