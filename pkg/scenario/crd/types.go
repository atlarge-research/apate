package crd

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

//go:generate controller-gen object paths=$GOFILE

type EmulatedPodSpec struct {
	Test int `json:"test"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type EmulatedPod struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec EmulatedPodSpec `json:"spec"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type EmulatedPodList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []EmulatedPod `json:"items"`
}
