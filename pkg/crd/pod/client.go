// Package pod defines utilities for the EmulatedPod CRD
// TODO make node CRD equivalent
package pod

import (
	"time"

	v1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/emulatedpod/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

const resource = "emulatedpods"

// EmulatedPodClient is the client for the EmulatedPod CRD
type EmulatedPodClient struct {
	restClient rest.Interface
	nameSpace  string
}

// NewForConfig creates a new EmulatedPodClient based on the given restConfig and namespace
func NewForConfig(c *rest.Config, namespace string) (*EmulatedPodClient, error) {
	if err := v1.AddToScheme(scheme.Scheme); err != nil {
		return nil, err
	}

	config := *c
	config.ContentConfig.GroupVersion = &v1.SchemeGroupVersion
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &EmulatedPodClient{restClient: client, nameSpace: namespace}, nil
}

// WatchResources creates an informer which watches for new or updated EmulatedPods and updates the returned store accordingly
func (e *EmulatedPodClient) WatchResources(addFunc func(obj interface{}), updateFunc func(oldObj, newObj interface{}), deleteFunc func(obj interface{})) *Informer {
	emulatedPodStore, emulatedPodController := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo metav1.ListOptions) (result runtime.Object, err error) {
				return e.list(lo)
			},
			WatchFunc: e.watch,
		},
		&v1.EmulatedPod{},
		1*time.Minute,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    addFunc,
			UpdateFunc: updateFunc,
			DeleteFunc: deleteFunc,
		},
	)

	go emulatedPodController.Run(wait.NeverStop)
	return NewInformer(&emulatedPodStore)
}

func (e *EmulatedPodClient) list(opts metav1.ListOptions) (*v1.EmulatedPodList, error) {
	result := v1.EmulatedPodList{}

	err := e.restClient.Get().
		Namespace(e.nameSpace).
		Resource(resource).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (e *EmulatedPodClient) watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true

	return e.restClient.
		Get().
		Namespace(e.nameSpace).
		Resource(resource).
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}
