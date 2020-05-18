// Package node provides functions and types to deal with the NodeConfiguration on the control plane
package node

import (
	"context"
	"log"
	"sync"

	"github.com/google/uuid"

	"github.com/atlarge-research/opendc-emulate-kubernetes/internal/run"
	v1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/nodeconfiguration/v1"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster/kubeconfig"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/crd/node"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/normalization"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/normalization/translate"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/store"
)

// CreateNodeInformer creates a new node informer
func CreateNodeInformer(ctx context.Context, config *kubeconfig.KubeConfig, st *store.Store, info *service.ConnectionInfo) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	client, err := node.NewForConfig(cfg, "default") // TODO: Change namespace
	if err != nil {
		return err
	}

	// TODO: Decide if we want this. This is the easiest way to ensure the correct amount of apatelets
	// Downside is that updates cannot cancel each other (eg. +200 and -200, will not result in 0 directly, but in
	// two separate 'transactions' of +200 and -200 resp.).
	// 1st (current) option: Lock and check current vs desired once. No cancellation, but easy and possibly fastest
	// 2nd option: No lock and recheck current vs desired after every spawn/despawn. This will result in cancellation,
	// but is possibly slower as it will have to recheck the amount of current nodes every time
	// 3rd option: Don't care about sync and let the watcher stabilise the amount with resyncs (currently every 1 minute).
	// This will result in easy cancellation, but in some cases it can take a few minutes before the desired amount of
	// apatelets is reached.

	// Create lock for stabilising creation
	var lock sync.Locker = &sync.Mutex{}

	client.WatchResources(func(obj interface{}) {
		go getDesiredApatelets(ctx, obj, st, &lock, info)
	}, func(_, obj interface{}) {
		go getDesiredApatelets(ctx, obj, st, &lock, info)
	}, func(obj interface{}) {
		go getDesiredApatelets(ctx, obj, st, &lock, info)
	})

	return nil
}

func getDesiredApatelets(ctx context.Context, obj interface{}, st *store.Store, lock *sync.Locker, info *service.ConnectionInfo) {
	(*lock).Lock()
	defer (*lock).Unlock()

	cfg := obj.(*v1.NodeConfiguration)
	res, err := getNodeResources(cfg)
	if err != nil {
		log.Printf("error while retrieving node resources from CRD: %v\n", err)
	}

	nodes, err := (*st).GetNodesBySelector(getSelector(cfg))
	if err != nil {
		log.Printf("error while retrieving nodes with selector %s: %v\n", getSelector(cfg), err)
	}

	current := int64(len(nodes))
	desired := cfg.Spec.Replicas

	if current < desired {
		// Not enough apatelets, spawn extra
		err := spawnApatelets(ctx, desired-current, st, res, info)
		if err != nil {
			log.Printf("error while spawning apatelets: %v\n", err)
			// TODO: Stop, notify, idk?
		}
	} else if current > desired {
		// Too many apatelets, stop a few
		err := stopApatelets(ctx, current-desired, st)
		if err != nil {
			log.Printf("error while stopping apatelets: %v\n", err)
			// TODO: Stop, notify, idk?
		}
	}
}

func spawnApatelets(ctx context.Context, diff int64, st *store.Store, res normalization.NodeResources, info *service.ConnectionInfo) error {
	log.Printf("Creating %v apatelets", diff)
	resources := createResources(int(diff), res)
	if err := (*st).AddResourcesToQueue(resources); err != nil {
		return err
	}

	// Retrieve pull policy
	pullPolicy := env.RetrieveFromEnvironment(env.ControlPlaneDockerPolicy, env.ControlPlaneDockerPolicyDefault)
	log.Printf("Using pull policy %s to spawn apatelets\n", pullPolicy)

	// Create environment for apatelets
	environment, err := env.DefaultApateletEnvironment()
	if err != nil {
		return err
	}

	environment.AddConnectionInfo(info.Address, info.Port)

	// Start the apatelets
	if err = run.StartApatelets(ctx, int(diff), environment); err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func stopApatelets(_ context.Context, diff int64, _ *store.Store) error {
	log.Printf("Stopping %v apatelets", diff)

	// TODO: Stop apatelets
	return nil
}

func createResources(needed int, base normalization.NodeResources) []normalization.NodeResources {
	var resources []normalization.NodeResources

	for i := 0; i < needed; i++ {
		res := base
		res.UUID = uuid.New()

		resources = append(resources, res)
	}

	return resources
}

func getNodeResources(nodeCfg *v1.NodeConfiguration) (normalization.NodeResources, error) {
	res := nodeCfg.Spec.Resources
	mem, err := translate.GetInBytes(res.Memory, "memory")
	if err != nil {
		return normalization.NodeResources{}, err
	}

	storage, err := translate.GetInBytes(res.Storage, "storage")
	if err != nil {
		return normalization.NodeResources{}, err
	}

	ephemeralStorage, err := translate.GetInBytes(res.EphemeralStorage, "ephemeral storage")
	if err != nil {
		return normalization.NodeResources{}, err
	}

	return normalization.NodeResources{
		Memory:           mem,
		CPU:              res.CPU,
		Storage:          storage,
		EphemeralStorage: ephemeralStorage,
		MaxPods:          res.MaxPods,
		Selector:         getSelector(nodeCfg),
	}, nil
}

func getSelector(cfg *v1.NodeConfiguration) string {
	return cfg.Namespace + "/" + cfg.Name
}
