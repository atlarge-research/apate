package v1

import (
	"io/ioutil"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/atlarge-research/apate/internal/kubectl"
	"github.com/atlarge-research/apate/pkg/apis/nodeconfiguration"
	"github.com/atlarge-research/apate/pkg/env"
	"github.com/atlarge-research/apate/pkg/kubernetes/kubeconfig"
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

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

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

// UpdateInKubernetes registers or deletes the generated CRD YAML to Kubernetes
func UpdateInKubernetes(config *kubeconfig.KubeConfig, deleteCRD bool) error {
	cpEnv := env.ControlPlaneEnv()

	file, err := ioutil.ReadFile(cpEnv.NodeCRDLocation)
	if err != nil {
		return errors.Wrapf(err, "failed to read crd file at %v", cpEnv.NodeCRDLocation)
	}

	if deleteCRD {
		return errors.Wrap(kubectl.Delete(file, config), "deleting node configuration failed")
	}
	return errors.Wrap(kubectl.Apply(file, config), "applying node configuration failed")
}
