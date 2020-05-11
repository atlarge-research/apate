package client

import (
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/crd"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

func NewForConfig(c *rest.Config, namespace string) (*emulatedPodClient, error) {
	if err := crd.RegisterCRD(); err != nil {
		return nil, err
	}

	config := *c
	config.ContentConfig.GroupVersion = &crd.GroupVersion
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &emulatedPodClient{restClient: client, nameSpace: namespace}, nil
}
