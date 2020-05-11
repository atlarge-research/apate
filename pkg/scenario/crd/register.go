package crd

import (
  "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster/kubeconfig"
  "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubectl"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  "k8s.io/apimachinery/pkg/runtime/schema"

  "k8s.io/client-go/kubernetes/scheme"

  "k8s.io/apimachinery/pkg/runtime"
)

const groupName = "apate.opendc.org"

const emulatedPodCRD = `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
    name: emulatedpods.apate.opendc.org
spec:
    group: apate.opendc.org
    versions:
        -
            name: v1
            served: true
            storage: true
            schema:
                openAPIV3Schema:
                    type: object
                    properties:
                        spec:
                            type: object
                            properties:
                                test:
                                    type: integer
    scope: Namespaced
    names:
        plural: emulatedpods
        singular: emulatedpod
        kind: EmulatedPod
        shortNames: [ep]

`

var GroupVersion = schema.GroupVersion{Group: groupName, Version: "v1"}

func AddCRDToKubernetes(config *kubeconfig.KubeConfig) error {
  if err := kubectl.Apply([]byte(emulatedPodCRD), config); err != nil {
    return err
  }

  return RegisterCRD()
}

func RegisterCRD() error {
  schemeBuilder := runtime.NewSchemeBuilder(addKnownTypes)
  return schemeBuilder.AddToScheme(scheme.Scheme)
}

func addKnownTypes(scheme *runtime.Scheme) error {
  scheme.AddKnownTypes(GroupVersion,
    &EmulatedPod{},
    &EmulatedPodList{},
  )

  metav1.AddToGroupVersion(scheme, GroupVersion)
  return nil
}
