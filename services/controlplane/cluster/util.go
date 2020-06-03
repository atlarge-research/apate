// Package cluster contains utilities to interact with the apate cluster and kubernetes cluster
package cluster

import (
	"log"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/store"
)

// RemoveNodesWithUUID removes the apatelets with the given uuids from the cluster (both apate and k8s)
func RemoveNodesWithUUID(uuids []uuid.UUID, st *store.Store, cl *kubernetes.Cluster) (err error) {
	if apateErr := (*st).RemoveNodes(uuids); apateErr != nil {
		err = errors.Wrapf(err, "removing node with uuids %v failed", uuids)
	}

	if kubernetesErr := cl.RemoveNodesFromCluster(uuids); kubernetesErr != nil {
		if err != nil {
			err = errors.Wrapf(err, "removing node with uuids from cluster %v failed (remove from apate also failed with '%v')", uuids, err)
		} else {
			err = errors.Wrapf(err, "removing node with uuids from cluster %v failed", uuids)
		}
	}

	return
}

// RemoveNodeWithUUID removes the apatelet with the given uuid from the cluster (both apate and k8s)
func RemoveNodeWithUUID(uuid uuid.UUID, st *store.Store, cl *kubernetes.Cluster) (apate bool, k8s bool, err error) {
	log.Printf("Removing %s from the cluster", uuid)
	apate = true
	k8s = true

	if apateErr := (*st).RemoveNode(uuid); apateErr != nil {
		err = errors.Wrapf(err, "removing node with uuid %v failed", uuid)
		apate = false
	}

	if kubernetesErr := cl.RemoveNodeFromCluster("apatelet-" + uuid.String()); kubernetesErr != nil {
		k8s = false
		if err != nil {
			err = errors.Wrapf(err, "removing node with uuid from cluster %v failed (remove from apate also failed with '%v')", uuid, err)
		} else {
			err = errors.Wrapf(err, "removing node with uuid from cluster %v failed", uuid)
		}
	}

	return
}
