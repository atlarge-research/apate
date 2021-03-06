package pod

import (
	"log"
	"time"

	"github.com/atlarge-research/apate/internal/crd/pod"

	"github.com/pkg/errors"

	podconfigv1 "github.com/atlarge-research/apate/pkg/apis/podconfiguration/v1"
	"github.com/atlarge-research/apate/pkg/kubernetes/kubeconfig"
	"github.com/atlarge-research/apate/services/apatelet/store"
)

// CreatePodInformer creates a new crd informer.
func CreatePodInformer(config *kubeconfig.KubeConfig, st *store.Store, stopch <-chan struct{}, wakeScheduler func()) error {
	restConfig, err := config.GetConfig()
	if err != nil {
		return errors.Wrap(err, "failed to get restconfig from kubeconfig for the pod informer")
	}

	podClient, err := pod.NewForConfig(restConfig, "default")
	if err != nil {
		return errors.Wrap(err, "failed to get podclient from rest config for pod informer")
	}

	podClient.WatchResources(func(obj interface{}) {
		// Add function
		podCfg := obj.(*podconfigv1.PodConfiguration)
		err := setPodTasks(podCfg, st)
		if err != nil {
			log.Printf("error while adding pod tasks: %v\n", err)
		}

		wakeScheduler()
	}, func(_, obj interface{}) {
		// Update function
		podCfg := obj.(*podconfigv1.PodConfiguration)

		err := setPodTasks(podCfg, st) // just replace all tasks with the <namespace>/<name>
		if err != nil {
			log.Printf("error while adding pod tasks: %v\n", err)
		}

		wakeScheduler()
	}, func(obj interface{}) {
		// Delete function
		podCfg := obj.(*podconfigv1.PodConfiguration)

		crdLabel := getCrdLabel(podCfg)
		err := (*st).RemovePodTasks(crdLabel)
		if err != nil {
			log.Printf("error while removing pod tasks: %v\n", err)
		}
	}, stopch)

	return nil
}

func setPodTasks(podCfg *podconfigv1.PodConfiguration, st *store.Store) error {
	// Validating timestamps before actually doing anything
	var durations = make([]time.Duration, len(podCfg.Spec.Tasks))
	for i, task := range podCfg.Spec.Tasks {
		duration, err := time.ParseDuration(task.Timestamp)
		if err != nil {
			return errors.Wrapf(err, "error while converting timestamp %v to a duration", task.Timestamp)
		}
		durations[i] = duration
	}

	crdLabel := getCrdLabel(podCfg)

	empty := podconfigv1.PodConfigurationState{}
	if podCfg.Spec.PodConfigurationState != empty {
		if err := SetPodFlags(st, crdLabel, &podCfg.Spec.PodConfigurationState); err != nil {
			return errors.Wrap(err, "failed to set pod flags during enqueueing of crd")
		}
	}

	var tasks []*store.Task
	var timeFlags []*store.TimeFlags

	for i, task := range podCfg.Spec.Tasks {
		state := task.State

		if task.RelativeToPod {
			flags, err := TranslatePodFlags(&state)
			if err != nil {
				return errors.Wrap(err, "failed to translate pod state into flags")
			}

			timeFlags = append(timeFlags, &store.TimeFlags{
				TimeSincePodStart: durations[i],
				Flags:             flags,
			})
		} else {
			tasks = append(tasks, store.NewPodTask(durations[i], crdLabel, &state))
		}
	}

	(*st).SetPodTimeFlags(crdLabel, timeFlags)
	return errors.Wrap((*st).SetPodTasks(crdLabel, tasks), "failed to set pod tasks")
}

func getCrdLabel(podCfg *podconfigv1.PodConfiguration) string {
	crdLabel := podCfg.Namespace + "/" + podCfg.Name
	return crdLabel
}
