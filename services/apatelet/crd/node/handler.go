package node

import (
	"log"
	"time"

	"github.com/pkg/errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/internal/crd/node"
	v1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/nodeconfiguration/v1"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes/kubeconfig"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
)

// CreateNodeInformer creates a new node informer
func CreateNodeInformer(config *kubeconfig.KubeConfig, st *store.Store, selector string, stopCh <-chan struct{}, cb func()) error {
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
		cb()

		nodeCfg := obj.(*v1.NodeConfiguration)
		if node.GetSelector(nodeCfg) == selector {
			err := setNodeTasks(nodeCfg, st)
			if err != nil {
				log.Printf("error while adding node tasks: %v\n", err)
			}
		}
	}, func(_, obj interface{}) {
		// Update function
		cb()

		nodeCfg := obj.(*v1.NodeConfiguration)
		if node.GetSelector(nodeCfg) == selector {
			err := setNodeTasks(nodeCfg, st)
			if err != nil {
				log.Printf("error while adding node tasks: %v\n", err)
			}
		}
	}, func(obj interface{}) {
		// Delete function
		// Do nothing here, as control plane will determine which, if any, apatelets should stop
	}, stopCh)

	return nil
}

func setNodeTasks(nodeCfg *v1.NodeConfiguration, st *store.Store) error {
	if nodeCfg.Spec.NodeConfigurationState != (v1.NodeConfigurationState{}) {
		SetNodeFlags(st, &nodeCfg.Spec.NodeConfigurationState)
	}

	var tasks []*store.Task
	for _, task := range nodeCfg.Spec.Tasks {
		state := task.State
		tasks = append(tasks, store.NewNodeTask(time.Duration(task.Timestamp)*time.Millisecond, &state))
	}

	if err := (*st).SetNodeTasks(tasks); err != nil {
		return errors.Wrap(err, "setting node tasks failed")
	}

	return nil
}
