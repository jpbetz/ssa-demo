//go:generate ../../../controller-tools/controller-gen paths="." object crd apply

// +groupName=samplecontroller.k8s.io
// +versionName=v1alpha1
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:generate=true
// +kubebuilder:object:root=true

type Foo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   FooSpec   `json:"spec"`
	Status FooStatus `json:"status"`
}

// +kubebuilder:object:generate=true

type FooSpec struct {
	DeploymentName string `json:"deploymentName"`
	Replicas       *int32 `json:"replicas"`
}

// +kubebuilder:object:generate=true

type FooStatus struct {
	AvailableReplicas int32 `json:"availableReplicas"`
}

// +kubebuilder:object:generate=true
// +kubebuilder:object:root=true

type FooList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Foo `json:"items"`
}
