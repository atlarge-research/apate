package node

import (
	"github.com/pkg/errors"
	"log"

	v1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/nodeconfiguration/v1"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster/kubeconfig"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/crd/node"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
)

// CreateNodeInformer creates a new node informer
func CreateNodeInformer(config *kubeconfig.KubeConfig, st *store.Store, selector string, stopCh chan struct{}) error {
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
		nodeCfg := obj.(*v1.NodeConfiguration)
		if node.GetSelector(nodeCfg) == selector {
			err := enqueueNodeTasks(nodeCfg, st)
			if err != nil {
				log.Printf("error while adding node tasks: %v\n", err)
			}
		}
	}, func(_, obj interface{}) {
		// Update function
		nodeCfg := obj.(*v1.NodeConfiguration)
		if node.GetSelector(nodeCfg) == selector {
			err := enqueueNodeTasks(nodeCfg, st)
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

func enqueueNodeTasks(nodeCfg *v1.NodeConfiguration, st *store.Store) error {
	if nodeCfg.Spec.NodeConfigurationState != (v1.NodeConfigurationState{}) {
		SetNodeFlags(st, &nodeCfg.Spec.NodeConfigurationState)
	}

	var tasks []*store.Task
	for _, task := range nodeCfg.Spec.Tasks {
		tasks = append(tasks, store.NewNodeTask(task.Timestamp, &store.NodeTask{State: &task.State}))
	}

	if err := (*st).SetNodeTasks(tasks); err != nil {
		return errors.Wrap(err, "setting node tasks failed")
	}

	return nil
}
