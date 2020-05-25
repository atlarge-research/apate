package pod

import (
	"log"
	"time"

	"github.com/pkg/errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/internal/crd/pod"
	podconfigv1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/podconfiguration/v1"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes/kubeconfig"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
)

// CreatePodInformer creates a new crd informer.
func CreatePodInformer(config *kubeconfig.KubeConfig, st *store.Store, stopch <-chan struct{}, cb func()) error {
	restConfig, err := config.GetConfig()
	if err != nil {
		return errors.Wrap(err, "failed to get restconfig from kubeconfig for the pod informer")
	}

	podClient, err := pod.NewForConfig(restConfig, "default")
	if err != nil {
		return errors.Wrap(err, "failed to get podclient from rest config for pod informer")
	}

	podClient.WatchResources(func(obj interface{}) {
		cb()

		err := setPodTasks(obj, st)
		if err != nil {
			log.Printf("error while adding pod tasks: %v\n", err)
		}
	}, func(oldObj, newObj interface{}) {
		cb()

		err := setPodTasks(newObj, st) // just replace all tasks with the <namespace>/<name>
		if err != nil {
			log.Printf("error while adding pod tasks: %v\n", err)
		}
	}, func(obj interface{}) {
		_, crdLabel := getCRDAndLabel(obj)
		err := (*st).RemovePodTasks(crdLabel)
		if err != nil {
			log.Printf("error while removing pod tasks: %v\n", err)
		}
	}, stopch)

	return nil
}

func setPodTasks(obj interface{}, st *store.Store) error {
	newCRD, crdLabel := getCRDAndLabel(obj)

	empty := podconfigv1.PodConfigurationState{}
	if newCRD.Spec.PodConfigurationState != empty {
		if err := SetPodFlags(st, crdLabel, &newCRD.Spec.PodConfigurationState); err != nil {
			return errors.Wrap(err, "failed to set pod flags during enqueueing of crd")
		}
	}

	var tasks []*store.Task
	for _, task := range newCRD.Spec.Tasks {
		state := task.State
		tasks = append(tasks, store.NewPodTask(time.Duration(task.Timestamp), crdLabel, &state))
	}

	return errors.Wrap((*st).SetPodTasks(crdLabel, tasks), "failed to set pod tasks")
}

func getCRDAndLabel(obj interface{}) (*podconfigv1.PodConfiguration, string) {
	newCRD := obj.(*podconfigv1.PodConfiguration)
	crdLabel := newCRD.Namespace + "/" + newCRD.Name
	return newCRD, crdLabel
}
