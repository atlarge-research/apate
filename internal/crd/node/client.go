// Package node defines utilities for the NodeConfiguration
package node

import (
	"log"
	"sync"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"

	"time"

	nodeconfigv1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/nodeconfiguration/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

const resource = "nodeconfigurations"

// ConfigurationClient is the client for the NodeConfiguration CRD
type ConfigurationClient struct {
	restClient rest.Interface
	namespace  string
}

var once sync.Once

// NewForConfig creates a new ConfigurationClient based on the given restConfig and namespace
func NewForConfig(c *rest.Config, namespace string) (*ConfigurationClient, error) {
	once.Do(func() {
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

	return &ConfigurationClient{restClient: client, namespace: namespace}, nil
}

// WatchResources creates an informer which watches for new or updated NodeConfigurations and updates the store accordingly
// This will also trigger the creation and removal of nodes when applicable
func (e *ConfigurationClient) WatchResources(addFunc func(obj interface{}), updateFunc func(oldObj, newObj interface{}), deleteFunc func(obj interface{}), stopCh <-chan struct{}) {
	_, nodeConfigurationController := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo metav1.ListOptions) (result runtime.Object, err error) {
				return e.list(lo)
			},
			WatchFunc: e.watch,
		},
		&nodeconfigv1.NodeConfiguration{},
		1*time.Minute,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    addFunc,
			UpdateFunc: updateFunc,
			DeleteFunc: deleteFunc,
		},
	)

	go nodeConfigurationController.Run(stopCh)
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

// Delete deletes all node configurations.
func (e *ConfigurationClient) Delete() error {
	return errors.Wrap(e.restClient.
		Delete().
		Namespace(e.namespace).
		Resource(resource).
		VersionedParams(&metav1.ListOptions{}, scheme.ParameterCodec).
		Do().
		Error(),
		"failed to delete node configurations")
}

// GetSelector concatenates the namespace and name to create a unique selector
func GetSelector(cfg *nodeconfigv1.NodeConfiguration) string {
	return cfg.Namespace + "/" + cfg.Name
}
