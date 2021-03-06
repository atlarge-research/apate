package node

import (
	"context"
	"log"

	"github.com/pkg/errors"

	"github.com/atlarge-research/apate/internal/crd/node"
	nodeconfigv1 "github.com/atlarge-research/apate/pkg/apis/nodeconfiguration/v1"
	"github.com/atlarge-research/apate/pkg/kubernetes/kubeconfig"
)

// WatchHandler watches a newly created informer for updates and calls functions on apateletHandler on updates
func WatchHandler(ctx context.Context, config *kubeconfig.KubeConfig, handler *ApateletHandler, stopInformerCh <-chan struct{}) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return errors.Wrap(err, "couldn't get kubeconfig for node informer")
	}

	client, err := node.NewForConfig(cfg, "default")
	if err != nil {
		return errors.Wrap(err, "couldn't create node client from config")
	}

	client.WatchResources(func(obj interface{}) {
		go func() {
			log.Println("Received new node CRD on controlplane")

			if err := (*handler).GetDesiredApatelets(ctx, obj.(*nodeconfigv1.NodeConfiguration)); err != nil {
				log.Printf("error while starting apatelets: %v\n", err)
			}
		}()
	}, func(_, obj interface{}) {
		go func() {
			log.Println("Received updated node CRD on controlplane")

			if err := (*handler).GetDesiredApatelets(ctx, obj.(*nodeconfigv1.NodeConfiguration)); err != nil {
				log.Printf("error while updating apatelets: %v\n", err)
			}
		}()
	}, func(obj interface{}) {
		cfg := obj.(*nodeconfigv1.NodeConfiguration)
		cfg.Spec.Replicas = 0
		go func() {
			log.Println("Received deleted node CRD on controlplane")

			if err := (*handler).GetDesiredApatelets(ctx, cfg); err != nil {
				log.Printf("error while stopping apatelets: %v\n", err)
			}
		}()
	}, stopInformerCh)

	return nil
}
