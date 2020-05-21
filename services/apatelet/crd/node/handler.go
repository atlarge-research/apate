package node

import (
	"log"

	v1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/nodeconfiguration/v1"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster/kubeconfig"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/crd/node"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
)

// CreateNodeInformer creates a new node informer
func CreateNodeInformer(config *kubeconfig.KubeConfig, st *store.Store, cb func(), selector string, stopCh chan struct{}) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	client, err := node.NewForConfig(cfg)
	if err != nil {
		return err
	}

	client.WatchResources(func(obj interface{}) {
		// Add function
		nodeCfg := obj.(*v1.NodeConfiguration)
		if getSelector(nodeCfg) == selector {
			err := enqueueNodeTasks(nodeCfg, st, cb)
			if err != nil {
				log.Printf("error while adding node tasks: %v\n", err)
			}
		}
	}, func(_, obj interface{}) {
		// Update function
		nodeCfg := obj.(*v1.NodeConfiguration)
		if getSelector(nodeCfg) == selector {
			err := enqueueNodeTasks(nodeCfg, st, cb)
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

func enqueueNodeTasks(nodeCfg *v1.NodeConfiguration, st *store.Store, cb func()) error {
	if nodeCfg.Spec.State != nil {
		SetNodeFlags(st, nodeCfg.Spec.State)
	}

	var tasks []*store.Task
	for _, task := range nodeCfg.Spec.Tasks {
		tasks = append(tasks, store.NewNodeTask(task.Timestamp, &store.NodeTask{State: &task.State}))
	}

	if err := (*st).SetNodeTasks(tasks); err != nil {
		return err
	}

	// Notify of update
	cb()
	return nil
}

func getSelector(cfg *v1.NodeConfiguration) string {
	return cfg.Namespace + "/" + cfg.Name
}
