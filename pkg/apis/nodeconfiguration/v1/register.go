package v1

import (
	"io/ioutil"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/nodeconfiguration"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster/kubeconfig"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubectl"
)

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: nodeconfiguration.GroupName, Version: "v1"}

// TODO: this is to register the emulated pod types with the decoder
var schemeGroupVersionInternal = schema.GroupVersion{Group: nodeconfiguration.GroupName, Version: runtime.APIVersionInternal}

var (
	// SchemeBuilder initialises a scheme builder
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	// AddToScheme is a global function that registers this API group & version to a scheme
	AddToScheme = SchemeBuilder.AddToScheme
)

// Adds the list of known types to Scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&NodeConfiguration{},
		&NodeConfigurationList{},
	)
	scheme.AddKnownTypes(schemeGroupVersionInternal,
		&NodeConfiguration{},
		&NodeConfigurationList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}

// CreateInKubernetes registers the generated CRD YAML to Kubernetes
func CreateInKubernetes(config *kubeconfig.KubeConfig) error {
	file, err := ioutil.ReadFile("config/crd/apate.opendc.org_nodeconfigurations.yaml")
	if err != nil {
		return err
	}

	return kubectl.Create(file, config)
}
