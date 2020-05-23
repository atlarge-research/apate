// Package watchdog contains the functions for removing an apatelet fully from the cluster and functions to ensure
// all unhealthy nodes are removed from the cluster
package watchdog

import (
	"github.com/pkg/errors"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/health"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/store"
)

// RemoveNodeWithUUID removes the apatelet with the given uuid from the cluster (both apate and k8s)
func RemoveNodeWithUUID(uuid uuid.UUID, st *store.Store, cl *cluster.KubernetesCluster) error {
	log.Printf("Removing %s from the cluster", uuid)

	if err := cl.RemoveNodeFromCluster("apatelet-" + uuid.String()); err != nil {
		return errors.Wrapf(err, "removing node with uuid from cluster: %v failed", uuid)
	}

	if err := (*st).RemoveNode(uuid); err != nil {
		return errors.Wrapf(err, "removing node with uuid: %v failed", uuid)
	}

	return nil
}

// StartWatchDog starts the watchdog
// The watchdog checks for unhealthy nodes, and removes them
func StartWatchDog(delay time.Duration, st *store.Store, cl *cluster.KubernetesCluster) {
	go func() {
		for {
			apatelets, err := (*st).GetNodes()

			if err != nil {
				log.Printf("unable to query nodes: %v", err)
			}

			for _, kubelet := range apatelets {
				if kubelet.Status == health.Status_UNHEALTHY {
					err := RemoveNodeWithUUID(kubelet.UUID, st, cl)
					if err != nil {
						log.Printf("error while removing apatelet from cluster: %v", err)
					} else {
						log.Printf("removed apatelet: %v from cluster", kubelet.UUID)
					}
				}
			}

			time.Sleep(delay)
		}
	}()
}
