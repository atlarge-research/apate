package node

import (
	"log"
	"time"

	"github.com/atlarge-research/opendc-emulate-kubernetes/internal/crd/node"

	"github.com/pkg/errors"

	nodeconfigv1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/nodeconfiguration/v1"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes/kubeconfig"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
)

// CreateNodeInformer creates a new node informer
func CreateNodeInformer(config *kubeconfig.KubeConfig, st *store.Store, selector string, stopch <-chan struct{}, wakeScheduler func()) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return errors.Wrap(err, "couldn't get kubeconfig")
	}

	client, err := node.NewForConfig(cfg)
	if err != nil {
		return errors.Wrap(err, "couldn't create client from config for node informer")
	}

	client.WatchResources(func(obj interface{}) {
		// Add function
		nodeCfg := obj.(*nodeconfigv1.NodeConfiguration)
		log.Printf("Received new node CRD on Apatelet\n")

		if node.GetSelector(nodeCfg) == selector {
			err := setNodeTasks(nodeCfg, st)
			if err != nil {
				log.Printf("error while adding node tasks: %v\n", err)
			}
		}

		wakeScheduler()
	}, func(_, obj interface{}) {
		// Update function
		nodeCfg := obj.(*nodeconfigv1.NodeConfiguration)
		log.Printf("Received updated node CRD on Apatelet\n")

		if node.GetSelector(nodeCfg) == selector {
			err := setNodeTasks(nodeCfg, st)
			if err != nil {
				log.Printf("error while adding node tasks: %v\n", err)
			}
		}

		wakeScheduler()
	}, func(obj interface{}) {
		// Delete function
		log.Printf("Received deleted node CRD on Apatelet\n")

		// Do nothing here, as control plane will determine which, if any, apatelets should stop
	}, stopch)

	return nil
}

func setNodeTasks(nodeCfg *nodeconfigv1.NodeConfiguration, st *store.Store) error {
	// Validating timestamps before actually doing anything
	var durations = make([]time.Duration, len(nodeCfg.Spec.Tasks))
	for i, task := range nodeCfg.Spec.Tasks {
		duration, err := time.ParseDuration(task.Timestamp)
		if err != nil {
			return errors.Wrapf(err, "error while converting timestamp %v to a duration", task.Timestamp)
		}
		durations[i] = duration
	}

	if nodeCfg.Spec.NodeConfigurationState != (nodeconfigv1.NodeConfigurationState{}) {
		SetNodeFlags(st, &nodeCfg.Spec.NodeConfigurationState)
	}

	var tasks []*store.Task
	for i, task := range nodeCfg.Spec.Tasks {
		state := task.State
		tasks = append(tasks, store.NewNodeTask(durations[i], &state))
	}

	if err := (*st).SetNodeTasks(tasks); err != nil {
		return errors.Wrap(err, "setting node tasks failed")
	}

	return nil
}
