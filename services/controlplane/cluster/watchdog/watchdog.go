// Package watchdog contains the functions for removing an apatelet fully from the cluster and functions to ensure
// all unhealthy nodes are removed from the cluster
package watchdog

import (
	"context"
	"log"
	"time"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario"

	"github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/cluster"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/health"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/store"
)

// StartWatchDog starts the watchdog
// The watchdog checks for unhealthy nodes, and removes them
func StartWatchDog(ctx context.Context, delay time.Duration, st *store.Store, cl *kubernetes.Cluster) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(delay):
				checkUnhealthyApatelets(st, cl)
			}
		}
	}()
}

func checkUnhealthyApatelets(st *store.Store, cl *kubernetes.Cluster) {
	apatelets, err := (*st).GetNodes()

	if err != nil {
		log.Printf("unable to query nodes: %v", err)
	}

	for _, kubelet := range apatelets {
		if kubelet.Status == health.Status_UNHEALTHY {
			apate, _, err := cluster.RemoveNodeWithUUID(kubelet.UUID, st, cl)
			if err != nil {
				log.Printf("error while removing apatelet from cluster: %v", err)
			} else {
				log.Printf("removed apatelet: %v from cluster", kubelet.UUID)
			}

			// If the node was removed from the store, add the resources back to the queue
			if apate {
				err := (*st).AddResourcesToQueue([]scenario.NodeResources{*kubelet.Resources})
				if err != nil {
					log.Printf("error while readding resources to queue: %v", err)
				}
			}
		}
	}
}
