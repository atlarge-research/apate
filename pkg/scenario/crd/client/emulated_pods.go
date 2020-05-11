package client

import (
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/crd"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"time"
)

const resource = "emulatedpods"

type emulatedPodClient struct {
	restClient rest.Interface
	nameSpace  string
}

func (e *emulatedPodClient) WatchResources() cache.Store {
	emulatedPodStore, emulatedPodController := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo metav1.ListOptions) (result runtime.Object, err error) {
				return e.list(lo)
			},
			WatchFunc: func(lo metav1.ListOptions) (watch.Interface, error) {
				return e.watch(lo)
			},
		},
		&crd.EmulatedPod{},
		1*time.Minute,
		cache.ResourceEventHandlerFuncs{},
	)

	go emulatedPodController.Run(wait.NeverStop)
	return emulatedPodStore
}


func (e *emulatedPodClient) list(opts metav1.ListOptions) (*crd.EmulatedPodList, error) {
	result := crd.EmulatedPodList{}

	err := e.restClient.Get().
		Namespace(e.nameSpace).
		Resource(resource).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (e *emulatedPodClient) watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true

	return e.restClient.
		Get().
		Namespace(e.nameSpace).
		Resource(resource).
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}
