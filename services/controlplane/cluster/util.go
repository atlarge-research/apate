// Package cluster contains utilities to interact with the apate cluster and kubernetes cluster
package cluster

import (
	"log"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/store"
)

// RemoveNodeWithUUID removes the apatelet with the given uuid from the cluster (both apate and k8s)
func RemoveNodeWithUUID(uuid uuid.UUID, st *store.Store, cl *kubernetes.Cluster) (err error) {
	log.Printf("Removing %s from the cluster", uuid)

	if apateErr := (*st).RemoveNode(uuid); apateErr != nil {
		err = errors.Wrapf(err, "removing node with uuid %v failed", uuid)
	}

	// TODO: Err handling
	go func() {
		if kubernetesErr := cl.RemoveNodeFromCluster("apatelet-" + uuid.String()); kubernetesErr != nil {
			err = errors.Wrapf(err, "removing node with uuid from cluster %v failed", uuid)
		}
	}()

	return
}
