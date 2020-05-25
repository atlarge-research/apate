// Package node provides functions and types to deal with the NodeConfiguration on the control plane
package node

import (
	"context"
	"log"
	"sync"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/run"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/clients/apatelet"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario"

	"github.com/google/uuid"

	nodev1 "github.com/atlarge-research/opendc-emulate-kubernetes/internal/crd/node"
	"github.com/atlarge-research/opendc-emulate-kubernetes/internal/service"
	v1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/nodeconfiguration/v1"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes/kubeconfig"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/store"
)

// CreateNodeInformer creates a new node informer
func CreateNodeInformer(ctx context.Context, config *kubeconfig.KubeConfig, st *store.Store, info *service.ConnectionInfo, stopCh <-chan struct{}) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return errors.Wrap(err, "couldn't get kubeconfig for node informer")
	}

	client, err := nodev1.NewForConfig(cfg)
	if err != nil {
		return errors.Wrap(err, "couldn't create node client from config")
	}

	// Create lock for stabilising creation
	var lock sync.Locker = &sync.Mutex{}

	client.WatchResources(func(obj interface{}) {
		go getDesiredApatelets(ctx, obj.(*v1.NodeConfiguration), st, &lock, info)
	}, func(_, obj interface{}) {
		go getDesiredApatelets(ctx, obj.(*v1.NodeConfiguration), st, &lock, info)
	}, func(obj interface{}) {
		cfg := obj.(*v1.NodeConfiguration)
		cfg.Spec.Replicas = 0

		go getDesiredApatelets(ctx, cfg, st, &lock, info)
	}, stopCh)

	return nil
}

func getDesiredApatelets(ctx context.Context, cfg *v1.NodeConfiguration, st *store.Store, lock *sync.Locker, info *service.ConnectionInfo) {
	(*lock).Lock()
	defer (*lock).Unlock()

	res, err := getNodeResources(cfg)
	if err != nil {
		log.Printf("error while retrieving node resources from CRD: %v\n", err)
	}

	selector := nodev1.GetSelector(cfg)
	nodes, err := (*st).GetNodesBySelector(selector)
	if err != nil {
		log.Printf("error while retrieving nodes with selector %s: %v\n", nodev1.GetSelector(cfg), err)
	}

	current := int64(len(nodes))
	desired := cfg.Spec.Replicas

	if current < desired {
		// Not enough apatelets, spawn extra
		err := spawnApatelets(ctx, st, desired, res, info, selector)
		if err != nil {
			log.Printf("error while spawning apatelets: %v\n", err)
			// TODO: Stop, notify, idk?
		}
	} else if current > desired {
		// Too many apatelets, stop a few
		err := stopApatelets(ctx, st, desired, selector)
		if err != nil {
			log.Printf("error while stopping apatelets: %v\n", err)
			// TODO: Stop, notify, idk?
		}
	}
}

func spawnApatelets(ctx context.Context, st *store.Store, desired int64, res scenario.NodeResources, info *service.ConnectionInfo, selector string) error {

	nodes, err := (*st).GetNodesBySelector(selector)
	if err != nil {
		return errors.Wrap(err, "failed getting nodes using selector")
	}

	current := int64(len(nodes))
	diff := desired - current

	log.Printf("Creating %v apatelets", diff)
	resources := createResources(int(diff), res)
	if err = (*st).AddResourcesToQueue(resources); err != nil {
		return errors.Wrap(err, "failed to add Apatalet resources to queue")
	}

	// Create environment for apatelets
	environment := env.DefaultApateletEnvironment()
	environment.AddConnectionInfo(info.Address, info.Port)

	// Start the apatelets
	if err = run.Registry.Run(ctx, int(diff), environment); err != nil {
		log.Print(err)
		return errors.Wrap(err, "error starting apatelets")
	}

	return nil
}

func stopApatelets(ctx context.Context, st *store.Store, desired int64, selector string) error {
	nodes, err := (*st).GetNodesBySelector(selector)
	if err != nil {
		return errors.Wrapf(err, "error while retrieving nodes with selector %s\n", selector)
	}

	current := int64(len(nodes))
	diff := int(current - desired)

	log.Printf("Stopping %v apatelets", diff)

	var wg sync.WaitGroup

	for i, node := range nodes {
		if i >= diff {
			break
		}

		node := node
		wg.Add(diff)

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

func createResources(needed int, base scenario.NodeResources) []scenario.NodeResources {
	var resources []scenario.NodeResources

	for i := 0; i < needed; i++ {
		res := base
		res.UUID = uuid.New()

		resources = append(resources, res)
	}

	return resources
}

func getNodeResources(nodeCfg *v1.NodeConfiguration) (scenario.NodeResources, error) {
	res := nodeCfg.Spec.Resources
	mem, err := scenario.GetInBytes(res.Memory, "memory")
	if err != nil {
		return scenario.NodeResources{}, errors.Wrap(err, "couldn't convert memory to bytes")
	}

	storage, err := scenario.GetInBytes(res.Storage, "storage")
	if err != nil {
		return scenario.NodeResources{}, errors.Wrap(err, "couldn't convert storage to bytes")
	}

	ephemeralStorage, err := scenario.GetInBytes(res.EphemeralStorage, "ephemeral storage")
	if err != nil {
		return scenario.NodeResources{}, errors.Wrap(err, "couldn't convert ephemeral storage to bytes")
	}

	return scenario.NodeResources{
		Memory:           mem,
		CPU:              res.CPU,
		Storage:          storage,
		EphemeralStorage: ephemeralStorage,
		MaxPods:          res.MaxPods,
		Selector:         nodev1.GetSelector(nodeCfg),
	}, nil
}
