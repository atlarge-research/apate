// Package node defines utilities for the NodeConfiguration
package node

import (
	"log"
	"sync"
	"time"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/tools/cache"

	nodeconfigv1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/nodeconfiguration/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

const resource = "nodeconfigurations"

// ConfigurationClient is the client for the NodeConfiguration CRD
type ConfigurationClient struct {
	restClient rest.Interface
	restConfig rest.Config
	namespace  string
}

type nodeClient struct {
	schemeLock         sync.Once
	sharedInformerLock sync.Once
	sharedInformer     *cache.SharedIndexInformer
}

var client nodeClient

// Reset will reset the sharedInformerLock, resulting in a new informer being created the next time resources are
// being watched. This is mostly for tests.
// Warning: Calling this during normal runtime will result in unpredictable behaviour, and possibly memory + routine leaks
func Reset() {
	client.sharedInformerLock = sync.Once{}
}

// NewForConfig creates a new ConfigurationClient based on the given restConfig and namespace
func NewForConfig(c *rest.Config, namespace string) (*ConfigurationClient, error) {
	client.schemeLock.Do(func() {
		if err := nodeconfigv1.AddToScheme(scheme.Scheme); err != nil {
			log.Panicf("%+v", errors.Wrap(err, "adding global node scheme failed"))
		}
	})

	config := *c
	config.ContentConfig.GroupVersion = &nodeconfigv1.SchemeGroupVersion
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new node crd client for config")
	}

	return &ConfigurationClient{restClient: client, restConfig: config, namespace: namespace}, nil
}

// WatchResources creates an informer which watches for new or updated NodeConfigurations and updates the store accordingly
// This will also trigger the creation and removal of nodes when applicable
func (e *ConfigurationClient) WatchResources(addFunc func(obj interface{}), updateFunc func(oldObj, newObj interface{}), deleteFunc func(obj interface{}), stopCh <-chan struct{}) {
	client.sharedInformerLock.Do(func() {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(lo metav1.ListOptions) (result runtime.Object, err error) {
					return e.list(lo)
				},
				WatchFunc: e.watch,
			},
			&nodeconfigv1.NodeConfiguration{},
			time.Minute,
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

func (e *ConfigurationClient) list(opts metav1.ListOptions) (*nodeconfigv1.NodeConfigurationList, error) {
	result := nodeconfigv1.NodeConfigurationList{}

	err := e.restClient.Get().
		Resource(resource).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	if err != nil {
		return nil, errors.Wrap(err, "failed to list node configurations")
	}

	return &result, nil
}

func (e *ConfigurationClient) watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true

	wi, err := e.restClient.
		Get().
		Resource(resource).
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()

	if err != nil {
		return nil, errors.Wrap(err, "failed to watch node configurtions")
	}

	return wi, nil
}

// GetCrdLabel concatenates the namespace and name to create a unique label
func GetCrdLabel(cfg *nodeconfigv1.NodeConfiguration) string {
	return cfg.Namespace + "/" + cfg.Name
}
