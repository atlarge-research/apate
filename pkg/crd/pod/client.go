// Package pod defines utilities for the PodConfiguration CRD
package pod

import (
	"time"

	v1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/podconfiguration/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

const resource = "podconfigurations"

// ConfigurationClient is the client for the PodConfiguration CRD
type ConfigurationClient struct {
	restClient rest.Interface
	nameSpace  string
}

// NewForConfig creates a new ConfigurationClient based on the given restConfig and namespace
func NewForConfig(c *rest.Config, namespace string) (*ConfigurationClient, error) {
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

	return &ConfigurationClient{restClient: client, nameSpace: namespace}, nil
}

// WatchResources creates an informer which watches for new or updated PodConfigurations and updates the returned store accordingly
func (e *ConfigurationClient) WatchResources(addFunc func(obj interface{}), updateFunc func(oldObj, newObj interface{}), deleteFunc func(obj interface{})) {
	_, podConfigurationController := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo metav1.ListOptions) (result runtime.Object, err error) {
				return e.list(lo)
			},
			WatchFunc: e.watch,
		},
		&v1.PodConfiguration{},
		1*time.Minute,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    addFunc,
			UpdateFunc: updateFunc,
			DeleteFunc: deleteFunc,
		},
	)

	go podConfigurationController.Run(wait.NeverStop)
}

func (e *ConfigurationClient) list(opts metav1.ListOptions) (*v1.PodConfigurationList, error) {
	result := v1.PodConfigurationList{}

	err := e.restClient.Get().
		Namespace(e.nameSpace).
		Resource(resource).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (e *ConfigurationClient) watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true

	wi, err := e.restClient.
		Get().
		Namespace(e.nameSpace).
		Resource(resource).
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()

	if err != nil {
		return nil, err
	}

	return wi, nil
}
