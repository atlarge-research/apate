package run

import (
	v1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/emulatedpod/v1"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster/kubeconfig"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/crd"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
)

func createCRDInformer(config *kubeconfig.KubeConfig, st *store.Store) (*crd.Informer, error) {
	restConfig, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	podClient, err := crd.NewForConfig(restConfig, "default")
	if err != nil {
		return nil, err
	}

	crdSt := podClient.WatchResources(func(obj interface{}) {
		enqueueCRD(obj, st)
	}, func(oldObj, newObj interface{}) {
		enqueueCRD(newObj, st) // just replace all tasks with the <namespace>/<name>
	}, func(obj interface{}) {
		_, crdLabel := getCRDAndLabel(obj)
		(*st).RemoveCRDTasks(crdLabel)
	})
	return crdSt, nil
}

func enqueueCRD(obj interface{}, st *store.Store) {
	newCRD, crdLabel := getCRDAndLabel(obj)

	var tasks []*store.Task
	for _, task := range newCRD.Spec.Tasks {
		tasks = append(tasks, &store.Task{
			AbsoluteTimestamp: 0, // TODO
			PodTask: store.PodTask{
				Label: crdLabel,
				Task:  &task,
			},
		})
	}

	(*st).EnqueueCRDTasks(crdLabel, tasks)
}

func getCRDAndLabel(obj interface{}) (*v1.EmulatedPod, string) {
	newCRD := obj.(*v1.EmulatedPod)
	crdLabel := newCRD.Namespace + "/" + newCRD.Name
	return newCRD, crdLabel
}
