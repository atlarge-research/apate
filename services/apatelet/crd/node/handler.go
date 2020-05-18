package node

import (
	"log"

	v1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/nodeconfiguration/v1"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster/kubeconfig"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/crd/node"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
)

// CreateNodeInformer creates a new node informer
func CreateNodeInformer(config *kubeconfig.KubeConfig, st *store.Store, selector string) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	client, err := node.NewForConfig(cfg, "default") //TODO: Change namespace
	if err != nil {
		return err
	}

	client.WatchResources(func(obj interface{}) {
		// Add function
		nodeCfg := obj.(*v1.NodeConfiguration)
		if nodeCfg.Meta.Name == selector {
			err := enqueueNodeTasks(nodeCfg, st)
			if err != nil {
				log.Printf("error while adding node tasks: %v\n", err)
			}
		}
	}, func(_, obj interface{}) {
		// Update function
		nodeCfg := obj.(*v1.NodeConfiguration)
		if nodeCfg.Meta.Name == selector {
			err := enqueueNodeTasks(nodeCfg, st)
			if err != nil {
				log.Printf("error while adding node tasks: %v\n", err)
			}
		}
	}, func(obj interface{}) {
		// Delete function
		// Do nothing here, as control plane will determine which, if any, apatelets should stop
	})

	return nil
}

func enqueueNodeTasks(nodeCfg *v1.NodeConfiguration, st *store.Store) error {
	if nodeCfg.Spec.State != nil {
		SetNodeFlags(st, nodeCfg.Spec.State)
	}

	var tasks []*store.Task
	for _, task := range nodeCfg.Spec.Tasks {
		tasks = append(tasks, store.NewNodeTask(task.Timestamp, &store.NodeTask{State: &task.State}))
	}

	return (*st).SetNodeTasks(tasks)
}
