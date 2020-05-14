package pod

import (
	"errors"

	v1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/emulatedpod/v1"

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
func (i *Informer) List() (eps []v1.EmulatedPod) {
	for _, ep := range (*i.store).List() {
		eps = append(eps, ep.(v1.EmulatedPod))
	}
	return eps
}

// Find finds an emulated pod which can be identified by the given label
// This label should have the format <namespace>/<name> of the emulated pod.
func (i *Informer) Find(label string) (*v1.EmulatedPod, bool, error) {
	key, exists, err := (*i.store).GetByKey(label)
	if err != nil {
		return nil, false, err
	}

	if exists {
		ep, ok := key.(*v1.EmulatedPod)
		if ok {
			return ep, true, nil
		}

		return nil, false, errors.New("couldn't cast to emulated pod")
	}

	return nil, false, nil
}
