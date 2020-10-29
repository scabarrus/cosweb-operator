package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CoswebSpec defines the desired state of Cosweb
type CoswebSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Name string `json:"name"`
	Replicas int32 `json:"replicas"`
	Image string `json:"image"`
	Labels map[string]string `json:"labels"`
}

// CoswebStatus defines the observed state of Cosweb
type CoswebStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Cosweb is the Schema for the coswebs API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=coswebs,scope=Namespaced
type Cosweb struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CoswebSpec   `json:"spec,omitempty"`
	Status CoswebStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CoswebList contains a list of Cosweb
type CoswebList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Cosweb `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Cosweb{}, &CoswebList{})
}
