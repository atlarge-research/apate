// Package pod provides functions to watch the pod CRD on the control plane
package pod

import (
	"github.com/pkg/errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/internal/crd/pod"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes/kubeconfig"
)

// NoopWatchHandler watches a newly created informer for updates and does nothing with this.
// This is to make sure our goroutines work properly by making this the first call to the WatchResources
func NoopWatchHandler(config *kubeconfig.KubeConfig, stopInformerCh <-chan struct{}) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return errors.Wrap(err, "couldn't get kubeconfig for node informer")
	}

	client, err := pod.NewForConfig(cfg, "default")
	if err != nil {
		return errors.Wrap(err, "couldn't create node client from config")
	}

	client.WatchResources(func(obj interface{}) {
		// Noop
	}, func(_, obj interface{}) {
		// Noop
	}, func(obj interface{}) {
		// Noop
	}, stopInformerCh)

	return nil
}