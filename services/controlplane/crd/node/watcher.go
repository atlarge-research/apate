package node

import (
	"context"
	"log"

	"github.com/pkg/errors"

	nodev1 "github.com/atlarge-research/opendc-emulate-kubernetes/internal/crd/node"
	v1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/nodeconfiguration/v1"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes/kubeconfig"
)

// WatchHandler watches a newly created informer for updates and calls functions on apateletHandler on updates
func WatchHandler(ctx context.Context, config *kubeconfig.KubeConfig, handler *ApateletHandler, stopCh <-chan struct{}) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return errors.Wrap(err, "couldn't get kubeconfig for node informer")
	}

	client, err := nodev1.NewForConfig(cfg)
	if err != nil {
		return errors.Wrap(err, "couldn't create node client from config")
	}

	client.WatchResources(func(obj interface{}) {
		go func() {
			if err := (*handler).GetDesiredApatelets(ctx, obj.(*v1.NodeConfiguration)); err != nil {
				log.Printf("error while starting apatelets %v", err)
			}
		}()
	}, func(_, obj interface{}) {
		go func() {
			if err := (*handler).GetDesiredApatelets(ctx, obj.(*v1.NodeConfiguration)); err != nil {
				log.Printf("error while updating apatelets %v", err)
			}
		}()
	}, func(obj interface{}) {
		cfg := obj.(*v1.NodeConfiguration)
		cfg.Spec.Replicas = 0
		go func() {
			if err := (*handler).GetDesiredApatelets(ctx, cfg); err != nil {
				log.Printf("error while stopping apatelets %v", err)
			}
		}()
	}, stopCh)

	return nil
}
