package pod

import (
	"fmt"

	v1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/podconfiguration/v1"

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
func (i *Informer) List() (eps []v1.PodConfiguration) {
	for _, ep := range (*i.store).List() {
		eps = append(eps, ep.(v1.PodConfiguration))
	}
	return eps
}

// Find finds an emulated pod which can be identified by the given label
// This label should have the format <namespace>/<name> of the emulated pod.
func (i *Informer) Find(label string) (*v1.PodConfiguration, bool, error) {
	key, exists, err := (*i.store).GetByKey(label)
	if err != nil {
		return nil, false, err
	}

	if exists {
		ep, ok := key.(*v1.PodConfiguration)
		if ok {
			return ep, true, nil
		}

		return nil, false, fmt.Errorf("couldn't cast %v to emulated pod", key)
	}

	return nil, false, nil
}
