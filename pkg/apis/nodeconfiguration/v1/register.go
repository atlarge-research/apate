package v1

import (
	"io/ioutil"

	"github.com/pkg/errors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/nodeconfiguration"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster/kubeconfig"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubectl"
)

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: nodeconfiguration.GroupName, Version: "v1"}
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
	const filename = "config/crd/apate.opendc.org_nodeconfigurations.yaml"
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return errors.Wrapf(err, "reading file: %v failed", filename)
	}

	return kubectl.Create(file, config)
}
