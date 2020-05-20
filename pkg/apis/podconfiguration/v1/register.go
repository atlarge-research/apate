package v1

import (
	"io/ioutil"

	"github.com/pkg/errors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/podconfiguration"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster/kubeconfig"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubectl"
)

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: podconfiguration.GroupName, Version: "v1"}

// this is a hack, to register the emulated pod types with the decoder
var schemeGroupVersionInternal = schema.GroupVersion{Group: podconfiguration.GroupName, Version: runtime.APIVersionInternal}

var (
	// SchemeBuilder initialises a scheme builder
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	// AddToScheme is a global function that registers this API group & version to a scheme
	AddToScheme = SchemeBuilder.AddToScheme
)

// Adds the list of known types to Scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&PodConfiguration{},
		&PodConfigurationList{},
	)
	// TODO find out why this is necessary
	scheme.AddKnownTypes(schemeGroupVersionInternal,
		&PodConfiguration{},
		&PodConfigurationList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)

	return nil
}

// CreateInKubernetes registers the generated CRD YAML to Kubernetes
func CreateInKubernetes(config *kubeconfig.KubeConfig) error {
	const filename = "config/crd/apate.opendc.org_podconfigurations.yaml"
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return errors.Wrapf(err, "failed to read crd file at %v", filename)
	}

	return errors.Wrap(kubectl.Create(file, config), "failed to create crd")
}
