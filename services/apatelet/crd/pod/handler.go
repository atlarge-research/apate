package pod

// TODO make node equivalent when moving node to CRD

import (
	v1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/podconfiguration/v1"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster/kubeconfig"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/crd/pod"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
)

// CreateCRDInformer creates a new crd informer.
func CreateCRDInformer(config *kubeconfig.KubeConfig, st *store.Store, errch *chan error) *pod.Informer {
	restConfig, err := config.GetConfig()
	if err != nil {
		*errch <- err
		return nil
	}

	podClient, err := pod.NewForConfig(restConfig, "default")
	if err != nil {
		*errch <- err
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
		err := (*st).RemovePodTasks(crdLabel)
		if err != nil {
			*errch <- err
		}
	})

	return inf
}

func enqueueCRD(obj interface{}, st *store.Store) error {
	newCRD, crdLabel := getCRDAndLabel(obj)

	empty := v1.PodConfigurationState{}
	if newCRD.Spec.PodConfigurationState != empty {
		if err := SetPodFlags(st, crdLabel, &newCRD.Spec.PodConfigurationState); err != nil {
			return err
		}
	}

	var tasks []*store.Task
	for _, task := range newCRD.Spec.Tasks {
		state := task.State
		tasks = append(tasks, store.NewPodTask(task.Timestamp, crdLabel, &state))
	}

	return (*st).EnqueuePodTasks(crdLabel, tasks)
}

func getCRDAndLabel(obj interface{}) (*v1.PodConfiguration, string) {
	newCRD := obj.(*v1.PodConfiguration)
	crdLabel := newCRD.Namespace + "/" + newCRD.Name
	return newCRD, crdLabel
}
