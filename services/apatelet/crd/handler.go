package crd

import (
	v1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/emulatedpod/v1"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster/kubeconfig"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/crd/pod"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
)

// CreateCRDInformer creates a new crd informer.
func CreateCRDInformer(config *kubeconfig.KubeConfig, st *store.Store, errch *chan error) *pod.Informer {
	restConfig, err := config.GetConfig()
	if err != nil {
		return nil
	}

	podClient, err := pod.NewForConfig(restConfig, "default")
	if err != nil {
		return nil
	}

	inf := podClient.WatchResources(func(obj interface{}) {
		err := enqueueCRD(obj, st)
		if err != nil {
			*errch <- err
		}
	}, func(oldObj, newObj interface{}) {
		err := enqueueCRD(newObj, st) // just replace all tasks with the <namespace>/<name>
		if err != nil {
			*errch <- err
		}
	}, func(obj interface{}) {
		_, crdLabel := getCRDAndLabel(obj)
		(*st).RemoveCRDTasks(crdLabel)
	})

	return inf
}

func enqueueCRD(obj interface{}, st *store.Store) error {
	newCRD, crdLabel := getCRDAndLabel(obj)

	if err := SetPodFlags(st, &store.PodTask{
		Label: crdLabel,
		State: &newCRD.Spec.DirectState,
	}); err != nil {
		return err
	}

	var tasks []*store.Task
	for _, task := range newCRD.Spec.Tasks {
		tasks = append(tasks, &store.Task{
			RelativeTimestamp: task.Timestamp,
			PodTask: &store.PodTask{
				Label: crdLabel,
				State: &task.State,
			},
		})
	}

	(*st).EnqueueCRDTasks(crdLabel, tasks)

	return nil
}

func getCRDAndLabel(obj interface{}) (*v1.EmulatedPod, string) {
	newCRD := obj.(*v1.EmulatedPod)
	crdLabel := newCRD.Namespace + "/" + newCRD.Name
	return newCRD, crdLabel
}
