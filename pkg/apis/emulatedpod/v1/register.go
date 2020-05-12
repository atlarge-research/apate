package v1

import (
	"io/ioutil"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster/kubeconfig"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubectl"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/emulatedpod"
)

// SchemeGroupVersion is the group and version of the EmulatedPod CRD
var SchemeGroupVersion = schema.GroupVersion{Group: emulatedpod.GroupName, Version: "v1"}

// Register registers the EmulatedPod structs to the scheme, so it can deserialize responses by Kubernetes
func Register() error {
	schemeBuilder := runtime.NewSchemeBuilder(addKnownTypes)
	return schemeBuilder.AddToScheme(scheme.Scheme)
}

// Adds the list of known types to Scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&EmulatedPod{},
		&EmulatedPodList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}

// CreateInKubernetes registers the generated CRD YAML to Kubernetes
func CreateInKubernetes(config *kubeconfig.KubeConfig) error {
	file, err := ioutil.ReadFile("config/crd/apate.opendc.org_emulatedpods.yaml")
	if err != nil {
		return err
	}

	if err := kubectl.Create(file, config); err != nil {
		return err
	}

	return nil
}
