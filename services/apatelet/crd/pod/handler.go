package pod

// TODO make node equivalent when moving node to CRD

import (
	v1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/podconfiguration/v1"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster/kubeconfig"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/crd/pod"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
	"github.com/pkg/errors"
)

// CreateCRDInformer creates a new crd informer.
func CreateCRDInformer(config *kubeconfig.KubeConfig, st *store.Store, errch *chan error) {
	restConfig, err := config.GetConfig()
	if err != nil {
		*errch <- errors.Wrap(err, "failed to get rest config from Kubeconfig")
		return
	}

	podClient, err := pod.NewForConfig(restConfig, "default")
	if err != nil {
		*errch <- errors.Wrap(err, "failed to create new pod client")
		return
	}

	podClient.WatchResources(func(obj interface{}) {
		err := enqueueCRD(obj, st)
		if err != nil {
			*errch <- errors.Wrap(err, "failed to enqueue crd (addfunc)")
		}
	}, func(oldObj, newObj interface{}) {
		err := enqueueCRD(newObj, st) // just replace all tasks with the <namespace>/<name>
		if err != nil {
			*errch <- errors.Wrap(err, "failed to enqueue crd (updatefunc)")
		}
	}, func(obj interface{}) {
		_, crdLabel := getCRDAndLabel(obj)
		err := (*st).RemovePodTasks(crdLabel)
		if err != nil {
			*errch <- errors.Wrap(err, "failed to delete pod tasks (deletefunc)")
		}
	})
}

func enqueueCRD(obj interface{}, st *store.Store) error {
	newCRD, crdLabel := getCRDAndLabel(obj)

	empty := v1.PodConfigurationState{}
	if newCRD.Spec.PodConfigurationState != empty {
		if err := SetPodFlags(st, crdLabel, &newCRD.Spec.PodConfigurationState); err != nil {
			return errors.Wrap(err, "failed to set pod flags")
		}
	}

	var tasks []*store.Task
	for _, task := range newCRD.Spec.Tasks {
		state := task.State
		tasks = append(tasks, store.NewPodTask(task.Timestamp, crdLabel, &state))
	}

	return errors.Wrap((*st).EnqueuePodTasks(crdLabel, tasks), "failed to enqueue new pod tasks")
}

func getCRDAndLabel(obj interface{}) (*v1.PodConfiguration, string) {
	newCRD := obj.(*v1.PodConfiguration)
	crdLabel := newCRD.Namespace + "/" + newCRD.Name
	return newCRD, crdLabel
}
