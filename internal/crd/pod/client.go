// Package pod defines utilities for the PodConfiguration CRD
package pod

import (
	"log"
	"sync"
	"time"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"

	podconfigv1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/podconfiguration/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

const resource = "podconfigurations"

type podClient struct {
	schemeLock         sync.Once
	sharedInformerLock sync.Once
	sharedInformer     *cache.SharedIndexInformer
}

var client podClient

// ConfigurationClient is the client for the PodConfiguration CRD
type ConfigurationClient struct {
	restClient rest.Interface
	namespace  string
}

// Reset will reset the sharedInformerLock, resulting in a new informer being created the next time resources are
// being watched. This is mostly for tests.
// Warning: Calling this during normal runtime will result in unpredictable behaviour, and possibly memory + routine leaks
func Reset() {
	client.sharedInformerLock = sync.Once{}
}

// NewForConfig creates a new ConfigurationClient based on the given restConfig and namespace
func NewForConfig(c *rest.Config, namespace string) (*ConfigurationClient, error) {
	client.schemeLock.Do(func() {
		if err := podconfigv1.AddToScheme(scheme.Scheme); err != nil {
			log.Panicf("%+v", errors.Wrap(err, "failed to add crd information to the scheme"))
		}
	})

	config := *c
	config.ContentConfig.GroupVersion = &podconfigv1.SchemeGroupVersion
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new pod crd client for config")
	}

	return &ConfigurationClient{restClient: client, namespace: namespace}, nil
}

// WatchResources creates an informer which watches for new or updated PodConfigurations and updates the returned store accordingly
func (e *ConfigurationClient) WatchResources(addFunc func(obj interface{}), updateFunc func(oldObj, newObj interface{}), deleteFunc func(obj interface{}), stopCh <-chan struct{}) {
	client.sharedInformerLock.Do(func() {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(lo metav1.ListOptions) (result runtime.Object, err error) {
					return e.list(lo)
				},
				WatchFunc: e.watch,
			},
			&podconfigv1.PodConfiguration{},
			1*time.Minute,
			cache.Indexers{},
		)

		client.sharedInformer = &informer
		go informer.Run(stopCh)
	})

	(*client.sharedInformer).AddEventHandlerWithResyncPeriod(cache.ResourceEventHandlerFuncs{
		AddFunc:    addFunc,
		UpdateFunc: updateFunc,
		DeleteFunc: deleteFunc,
	}, time.Minute)
}

func (e *ConfigurationClient) list(opts metav1.ListOptions) (*podconfigv1.PodConfigurationList, error) {
	result := podconfigv1.PodConfigurationList{}

	err := e.restClient.Get().
		Namespace(e.namespace).
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
		Namespace(e.namespace).
		Resource(resource).
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()

	if err != nil {
		return nil, errors.Wrap(err, "failed to watch pod configurations")
	}

	return wi, nil
}
