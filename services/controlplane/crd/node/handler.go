// Package node provides functions and types to deal with the NodeConfiguration on the control plane
package node

import (
	"context"
	"log"
	"net"
	"os"
	"sync"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/runner"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/clients/apatelet"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario"

	nodev1 "github.com/atlarge-research/opendc-emulate-kubernetes/internal/crd/node"
	"github.com/atlarge-research/opendc-emulate-kubernetes/internal/service"
	nodeconfigv1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/nodeconfiguration/v1"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/store"
)

// ApateletHandler contains utilities to spawn and stop apatelets, and to update them based on a given node configuration
type ApateletHandler interface {
	// Updates the amount of apatelets based on a given node configuration
	GetDesiredApatelets(context.Context, *nodeconfigv1.NodeConfiguration) error

	// Spawns n apatelets with resources and label
	SpawnApatelets(context.Context, int64, scenario.NodeResources, string) error

	// Stops n apatelets with label
	StopApatelets(context.Context, int64, string) error
}

type apateletHandler struct {
	lock           sync.Mutex
	store          *store.Store
	connectionInfo *service.ConnectionInfo
	runnerRegistry *runner.Registry
}

// NewHandler creates a new ApateletHandler
func NewHandler(st *store.Store, runnerRegistry *runner.Registry, info *service.ConnectionInfo) *ApateletHandler {
	var handler ApateletHandler = &apateletHandler{
		store:          st,
		connectionInfo: info,
		runnerRegistry: runnerRegistry,
	}

	return &handler
}

func (a *apateletHandler) GetDesiredApatelets(ctx context.Context, cfg *nodeconfigv1.NodeConfiguration) error {
	a.lock.Lock()
	defer a.lock.Unlock()

	res, err := getNodeResources(cfg)
	if err != nil {
		return errors.Wrap(err, "error while retrieving node resources from CRD")
	}

	label := nodev1.GetCrdLabel(cfg)
	nodes, err := (*a.store).GetNodesByLabel(label)
	if err != nil {
		return errors.Wrapf(err, "error while retrieving nodes with label %s", nodev1.GetCrdLabel(cfg))
	}

	current := int64(len(nodes))
	desired := cfg.Spec.Replicas

	if current < desired {
		// Not enough apatelets, spawn extra
		err := a.SpawnApatelets(ctx, desired, res, label)
		if err != nil {
			return errors.Wrap(err, "error while spawning apatelets")
		}
	} else if current > desired {
		// Too many apatelets, stop a few
		err := a.StopApatelets(ctx, desired, label)
		if err != nil {
			return errors.Wrap(err, "error while stopping apatelets")
		}
	}

	return nil
}

func (a *apateletHandler) SpawnApatelets(ctx context.Context, desired int64, res scenario.NodeResources, label string) error {
	nodes, err := (*a.store).GetNodesByLabel(label)
	if err != nil {
		return errors.Wrap(err, "failed getting nodes using label")
	}

	current := int64(len(nodes))
	diff := desired - current

	log.Printf("Creating %v apatelets", diff)
	resources := createResources(int(diff), res)
	if err = (*a.store).AddResourcesToQueue(resources); err != nil {
		return errors.Wrap(err, "failed to add Apatalet resources to queue")
	}

	// Create environment for apatelets
	environment, err := env.ApateletEnv()
	if err != nil {
		return errors.Wrap(err, "getting apatelet environment failed")
	}

	// Part of the fixes for DinD CI
	if os.Getenv("CI_COMMIT_REF_SLUG") != "" {
		addr, err := net.ResolveIPAddr("ip", "docker")
		if err != nil {
			log.Fatalf("Resolving ip of docker failed: %v", err)
		}
		environment.CIKubernetesAddress = addr.String()
	}

	environment.AddConnectionInfo(a.connectionInfo.Address, a.connectionInfo.Port)

	// Start the apatelets
	if err = a.runnerRegistry.Run(ctx, int(diff), environment); err != nil {
		log.Print(err)
		return errors.Wrap(err, "error starting apatelets")
	}

	return nil
}

func (a *apateletHandler) StopApatelets(ctx context.Context, desired int64, label string) error {
	nodes, err := (*a.store).GetNodesByLabel(label)
	if err != nil {
		return errors.Wrapf(err, "error while retrieving nodes with label %s\n", label)
	}

	current := int64(len(nodes))
	diff := int(current - desired)

	log.Printf("Stopping %v apatelets", diff)

	var wg sync.WaitGroup
	wg.Add(diff)

	for i, node := range nodes {
		if i >= diff {
			break
		}

		node := node

		go func() {
			defer wg.Done()

			client, err := apatelet.GetApateletClient(&node.ConnectionInfo)
			if err != nil {
				log.Printf("%v", errors.Wrap(err, "failed getting apatelet client"))
			}

			_, err = client.Client.StopApatelet(ctx, new(empty.Empty))
			if err != nil {
				log.Printf("%v", errors.Wrap(err, "failed stopping apatelet"))
			}

			err = client.Conn.Close()
			if err != nil {
				log.Printf("%v", errors.Wrap(err, "failed closing apatelet client connection"))
			}
		}()
	}

	wg.Wait()
	return nil
}
