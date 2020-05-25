// Package pod defines utilities for the PodConfiguration CRD
package pod

import (
	"log"
	"sync"
	"time"

	"github.com/pkg/errors"

	v1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/podconfiguration/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

const resource = "podconfigurations"

var once sync.Once

// ConfigurationClient is the client for the PodConfiguration CRD
type ConfigurationClient struct {
	restClient rest.Interface
	nameSpace  string
}

// NewForConfig creates a new ConfigurationClient based on the given restConfig and namespace
func NewForConfig(c *rest.Config, namespace string) (*ConfigurationClient, error) {
	once.Do(func() {
		if err := v1.AddToScheme(scheme.Scheme); err != nil {
			log.Panicf("%+v", errors.Wrap(err, "failed to add crd information to the scheme"))
		}
	})

	config := *c
	config.ContentConfig.GroupVersion = &v1.SchemeGroupVersion
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new pod crd client for config")
	}

	return &ConfigurationClient{restClient: client, nameSpace: namespace}, nil
}

// WatchResources creates an informer which watches for new or updated PodConfigurations and updates the returned store accordingly
func (e *ConfigurationClient) WatchResources(addFunc func(obj interface{}), updateFunc func(oldObj, newObj interface{}), deleteFunc func(obj interface{}), stopCh <-chan struct{}) {
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

	go podConfigurationController.Run(stopCh)
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
		return nil, errors.Wrap(err, "failed to list pod configurations")
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
		return nil, errors.Wrap(err, "failed to watch pod configurations")
	}

	return wi, nil
}
